/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"net/http"

	"github.com/coreos/etcd/clientv3"
	"istio.io/istio/pkg/log"
)

const (
	skyDNSPrefix = "/skydns"
	defaultTTL   = 3600
)

// EtcdHostData to store in etcd directly
type EtcdHostData struct {
	Host string `json:"host"`
	TTL  int    `json:"ttl"`
}

// HostData to store use REST
type HostData struct {
	Zone    string `json:"zone"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address, omitempty"`
	TTL     int    `json:"ttl, omitempty"`
}

// DNSInterface for DNS
type DNSInterface interface {
	Update(domain, ip, suffix string) error
	Delete(domain, ip, suffix string) error
}

type coreDNSEtcd struct {
	Client *clientv3.Client
}

func newCoreDNSEtcd(client *clientv3.Client) *coreDNSEtcd {
	return &coreDNSEtcd{
		Client: client,
	}
}

func newCoreDNSREST(address string) *coreDNSREST {
	log.Infof("new coredns struct use address %s", address)
	return &coreDNSREST{
		dnsAddress: address,
	}
}

func ConvertDomainToKey(domain string) string {
	keys := strings.Split(domain, ".")

	key := skyDNSPrefix
	for i := len(keys) - 1; i >= 0; i-- {
		key += "/" + keys[i]
	}

	return strings.ToLower(key)
}

// Update
func (cd *coreDNSEtcd) Update(domain, ip, suffix string) error {
	key := ConvertDomainToKey(domain + suffix)

	hostData := EtcdHostData{
		Host: ip,
		TTL:  defaultTTL,
	}
	data, _ := json.Marshal(&hostData)
	value := string(data)

	log.Infof("put <%s, %s>", key, value)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := cd.Client.Put(ctx, key, value)
	cancel()
	if err != nil {
		log.Errorf("put %s %s error: %v", key, value, err)
	}
	return err
}

// Delete
func (cd *coreDNSEtcd) Delete(domain, ip, suffix string) error {
	key := ConvertDomainToKey(domain + suffix)
	log.Infof("delete %s", key)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := cd.Client.Delete(ctx, key)
	cancel()
	return err
}

type coreDNSREST struct {
	// address of coredns
	dnsAddress string
}

// Update
func (cd *coreDNSREST) Update(domain, ip, suffix string) error {
	hostData := HostData{
		Zone:    suffix,
		Name:    domain,
		Type:    "A",
		Address: ip,
		TTL:     defaultTTL,
	}

	return cd.doRequest(&hostData, http.MethodPut)
}

func (cd *coreDNSREST) doRequest(hostData *HostData, method string) error {
	if hostData.Zone[0] == '.' {
		hostData.Zone = hostData.Zone[1:]
	}
	if hostData.Zone[len(hostData.Zone)-1] != '.' {
		hostData.Zone += "."
	}
	hostData.Name = strings.ToLower(hostData.Name)

	v, _ := json.Marshal(&hostData)
	data := string(v)

	client := &http.Client{}

	url := cd.dnsAddress + "/dynapi"

	log.Infof("url: %s, method: %s, data: %s", url, method, data)

	reqest, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := client.Do(reqest)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

// Delete
func (cd *coreDNSREST) Delete(domain, ip, suffix string) error {
	hostData := HostData{
		Zone: suffix,
		Name: domain,
		Type: "A",
		Address:ip,
		TTL:  defaultTTL,
	}

	return cd.doRequest(&hostData, http.MethodDelete)
}

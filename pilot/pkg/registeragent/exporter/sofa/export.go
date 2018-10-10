// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sofa

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"io/ioutil"
	"istio.io/fortio/log"
	"encoding/json"
	"strings"
)

type SimpleRpcInfoExporter struct {
}

func NewRpcInfoExporter() *SimpleRpcInfoExporter {
	s := &SimpleRpcInfoExporter{}
	return s
}

type ProviderInterface struct {
	InterfaceName string `json:"interface"`
	Group         string `json:"group"`
	Serialize     string `json:"serialization"`
	Version       string `json:"version"`
}

type InterfacesDTO struct {
	Providers []ProviderInterface `json:"providers"`
	Protocol  string              `json:"protocol"`
}

type HealthInfo struct {
	Status string                   `json:"status"`
	HealthCheckInfo HealthCheckInfo `json:"sofaBootComponentHealthCheckInfo"`
}

type HealthCheckInfo struct {
	Status string         `json:"status"`
	Middleware Middleware `json:"Middleware"`
}

type Middleware struct {
	RuntimeComponent map[string]string `json:"RUNTIME-COMPONENT"`
}

func (r *SimpleRpcInfoExporter) GetRpcServiceInfo(c *gin.Context) {
	localHost := "localhost"
	//default sofa boot actuator port
	serverPort := 8080
	url := fmt.Sprintf("http://%v:%v/health", localHost, serverPort)
	log.Infof("local rpc info url %v", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errf("new local rpc request failed: ", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errf("get local rpc info failed: ", err)
		return
	}
	defer resp.Body.Close()
	log.Infof("status: %v", string(resp.StatusCode))
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Infof("body string: %v", string(body))
		info := HealthInfo{}
		err := json.Unmarshal(body, &info)
		if err != nil {
			log.Errf("Unmarshal rpc info failed: ", err)
			return
		}
		interfacesDTO := InterfacesDTO{}
		interfacesDTO.Protocol = "BOLT"
		for k, v := range info.HealthCheckInfo.Middleware.RuntimeComponent {
			// filter data
			if k != "status" && !strings.Contains(k, "#") && v == "passed" {
				fmt.Println(k)
				uniqueInterfaceName := strings.Split(k, ":")
				if len(uniqueInterfaceName) != 3 {
					log.Errf("interface name: %v format error: %v", k)
				}
				providerInterfaceDto := ProviderInterface{uniqueInterfaceName[1], "", "Protobuf", uniqueInterfaceName[2]}
				interfacesDTO.Providers = append(interfacesDTO.Providers, providerInterfaceDto)
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": interfacesDTO})
	}
}

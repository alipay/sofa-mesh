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

package dubbo

import (
	"encoding/json"
	"fmt"
	"fortio.org/fortio/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
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

type Info map[string]ProviderInterface

func (r *SimpleRpcInfoExporter) GetRpcServiceInfo(c *gin.Context) {
	localHost := "localhost"
	// default dubbo-spring-boot-actuator port
	serverPort := 8081
	url := fmt.Sprintf("http://%v:%v/actuator/dubbo/services", localHost, serverPort)
	log.Infof("local rpc info url %v", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errf("new local rpc request failed: %v", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errf("get local rpc info failed: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Infof("status: %v", string(resp.StatusCode))
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Infof("body string: %v", string(body))
		info := Info{}
		err := json.Unmarshal(body, &info)
		if err != nil {
			log.Errf("Unmarshal rpc info failed: %v", err)
			return
		}
		interfacesDTO := InterfacesDTO{}
		interfacesDTO.Protocol = "DUBBO"
		for _, v := range info {
			// default serialize
			if v.Serialize == "" {
				v.Serialize = "hessian2"
			}
			interfacesDTO.Providers = append(interfacesDTO.Providers, v)
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": interfacesDTO})
	}
}

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

package main

import (
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"istio.io/istio/pilot/pkg/config/registeragent"
	"istio.io/istio/pilot/pkg/registeragent/exporter"
	"istio.io/istio/pkg/log"
)

var (
	flags = struct {
		loggingOptions *log.Options
	}{
		loggingOptions: log.DefaultOptions(),
	}
)

func init() {
	time.LoadLocation("China/BeiJing")
	if err := log.Configure(flags.loggingOptions); err != nil {
		os.Exit(-1)
	}
}

func startExporter(conf *config.Config) {
	runtime.GOMAXPROCS(4)
	router := gin.Default()
	rpcAcutatorExporter := exporter.RpcInfoExporterFactory()
	router.GET("/rpc/interfaces", rpcAcutatorExporter.GetRpcServiceInfo)
	router.Run(":10006")
}

func main() {
	config := config.NewConfig()
	log.Infof("all configs %s", config)
	startExporter(config)
}

package env

const (
	mosnConfigTempJson = `
{
  "node":{},
  "stats_config":{},
  "dynamic_resources":{},
  "tracing": {},
  "admin": {
    "access_log_path": "{{.AccessLogPath}}",
    "address": {
      "socket_address": {
        "address": "127.0.0.1",
        "port_value": "{{.Ports.AdminPort}}"
      }
    }
  },
  "dynamic_resources": {
    "lds_config": {
      "ads": {}
    },
    "cds_config": {
      "ads": {}
    },
    "ads_config": {
      "api_type": "GRPC",
      "grpc_services": [
        {
          "envoy_grpc": {
            "cluster_name": "loop"
          }
        }
      ]
    }
  },
  "static_resources": {
    "clusters": [
      {
        "name": "backend",
        "connect_timeout": "5s",
        "type": "STATIC",
        "hosts": [
          {
            "socket_address": {
              "address": "127.0.0.1",
              "port_value": "{{.Ports.BackendPort}}"
            }
          }
        ]
      },
      {
        "name": "loop",
        "connect_timeout": "5s",
        "type": "STATIC",
        "hosts": [
          {
            "socket_address": {
              "address": "127.0.0.1",
              "port_value": "{{.Ports.ServerProxyPort}}"
            }
          }
        ]
      },
      {
        "name": "mixer_server",
        "http2_protocol_options": {},
        "connect_timeout": "5s",
        "type": "STATIC",
        "hosts": [
          {
            "socket_address": {
              "address": "127.0.0.1",
              "port_value": "{{.Ports.MixerPort}}"
            }
          }
        ],
        "circuit_breakers": {
          "thresholds": [
            {
              "max_connections": 10000,
              "max_pending_requests": 10000,
              "max_requests": 10000,
              "max_retries": 3
            }
          ]
        }
      }
    ],
    "listeners": [
      {
        "name": "server",
        "address": {
          "socket_address": {
            "address": "127.0.0.1",
            "port_value": "{{.Ports.ServerProxyPort}}"
          }
        },
        "filter_chains": [
          {
            "filters": [
              {
                "name": "envoy.http_connection_manager",
                "config": {
                  "codec_type": "AUTO",
                  "stat_prefix": "inbound_http",
                  "access_log": [
                    {
                      "name": "envoy.file_access_log",
                      "config": {
                        "path": "{{.AccessLogPath}}"
                      }
                    }
                  ],
                  "http_filters": [
                    {
                      "name": "mixer",
                      "config": {{.MfConfig.HTTPServerConf | toJSON }}
                    },
                    {
                      "name": "envoy.router"
                    }
                  ],
                  "route_config": {
                    "name": "backend",
                    "virtual_hosts": [
                      {
                        "name": "backend",
                        "domains": ["*"],
                        "routes": [
                          {
                            "match": {
                              "prefix": "/"
                            },
                            "route": {
                              "cluster": "backend",
                              "timeout": "0s"
                            },
                            "per_filter_config": {
                              "mixer": {{.MfConfig.PerRouteConf | toJSON }}
                            }
                          }
                        ]
                      }
                    ]
                  }

                }
              }
            ]
          }
        ]
      }, {
        "name": "client",
        "address": {
          "socket_address": {
            "address": "127.0.0.1",
            "port_value": "{{.Ports.ClientProxyPort}}"
          }
        },
        "filter_chains": [
          {
            "filters": [
              {
                "name": "envoy.http_connection_manager",
                "config": {
                  "codec_type": "AUTO",
                  "stat_prefix": "outbound_http",
                  "access_log": [
                    {
                      "name": "envoy.file_access_log",
                      "config": {
                        "path": "{{.AccessLogPath}}"
                      }
                    }
                  ],
                  "http_filters": [
                    {
                      "name": "mixer",
                      "config": {{.MfConfig.HTTPClientConf | toJSON }}
                    },
                    {
                      "name": "envoy.router"
                    }
                  ],
                  "route_config": {
                    "name": "loop",
                    "virtual_hosts": [
                      {
                        "name": "loop",
                        "domains": ["*"],
                        "routes": [
                          {
                            "match": {
                              "prefix": "/"
                            },
                            "route": {
                              "cluster": "loop",
                              "timeout": "0s"
                            }
                          }
                        ]
                      }
                    ]
                  }

                }
              }
            ]
          }
        ]
      }, {
        "name": "tcp_server",
        "address": {
          "socket_address": {
            "address": "127.0.0.1",
            "port_value": "{{.Ports.TCPProxyPort}}"
          }
        },
        "filter_chains": [
          {
            "filters": [
              {
                "name": "mixer",
                "config": {{.MfConfig.TCPServerConf | toJSON }}
              }, {
                "name": "envoy.tcp_proxy",
                "config": {
                  "stat_prefix": "inbound_tcp",
                  "cluster": "backend"
                }
              }
            ]
          }
        ]
      }
    ]
  }
}
`
)
package main

import (
	"k8s-webhook-example/core/web"
)

func main() {
	config := web.HTTP{
		IP:   "0.0.0.0",
		Port: "3000",
		TLS: web.TLS{
			Switch:   true,
			CertFile: "config/webhook/TLS/server/server.crt",
			KeyFile:  "config/webhook/TLS/server/server.key",
			CaFile:   "config/webhook/TLS/ca/ca.crt",
		},
	}
	config.Start()
}

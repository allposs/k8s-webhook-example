package main

import (
	"flag"
	"fmt"
	"k8s-webhook-example/core/web"
)

var (
	port    = flag.String("port", "3000", "Web server Listen Port")
	address = flag.String("ip", "0.0.0.0", "Web server Address")
	tls     = flag.Bool("switch", true, "webhook TLS switch")
	cert    = flag.String("certFile", "config/webhook/TLS/server/server.crt", "TLS cert file path")
	key     = flag.String("keyFile", "config/webhook/TLS/server/server.key", "TLS key file path")
	ca      = flag.String("caFile", "config/webhook/TLS/ca/ca.crt", "TLS ca file path")
)

func usage() {
	fmt.Println("k8s-webhook-example 1.0.1")
	fmt.Println("usage: k8s-webhook-example [options]")

	fmt.Println("Options:")
	fmt.Println("  --port   	<80>            			Server Listen Port")
	fmt.Println("  --ip      	<0.0.0.0>   				Web server Address")
	fmt.Println("  --switch  	<false>  				Webhook TLS switch")
	fmt.Println("  --certFile  	<config/webhook/TLS/server/server.crt>  TLS cert file path")
	fmt.Println("  --keyFile  	<config/webhook/TLS/server/server.key>  TLS key file path")
	fmt.Println("  --caFile  	<config/webhook/TLS/ca/ca.crt>  	TLS ca file path")

	fmt.Println("Examples:")
	fmt.Println("  Default webhook Server")
	fmt.Println("  $ k8s-webhook-example")

	fmt.Println("  k8s-webhook-example Listen 127.0.0.1:80 and TLS is flase  ")
	fmt.Println("  $ k8s-webhook-example --port 80 --ip 127.0.0.1 --switch flase")
}

func main() {
	flag.Usage = usage
	flag.Parse()
	config := web.HTTP{
		IP:   *address,
		Port: *port,
		TLS: web.TLS{
			Switch:   *tls,
			CertFile: *cert,
			KeyFile:  *key,
			CaFile:   *ca,
		},
	}
	config.Start()
}

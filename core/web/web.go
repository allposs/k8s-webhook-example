package web

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"k8s-webhook-example/core/web/middleware"
	"k8s-webhook-example/core/web/routes"
	"log"
	"net/http"
	"time"
)

// Web Web服务结构体
type Web struct {
	HTTP   HTTP
	Server *http.Server
}

// HTTP web服务HTTP启动参数
type HTTP struct {
	IP   string
	Port string
	TLS  TLS
}

type TLS struct {
	Switch   bool
	CertFile string
	KeyFile  string
	CaFile   string
}

// Mux  web服务的路由日志与中间件
func (web *Web) Mux() {
	mux := http.NewServeMux()
	route := routes.New()
	mux.Handle("/", middleware.Logging(http.HandlerFunc(route.ServeHTTP)))
	web.Server.Handler = mux
}

func (web *Web) RunTLS() {
	addr := fmt.Sprintf("%s:%s", web.HTTP.IP, web.HTTP.Port)
	if web.HTTP.TLS.Switch {
		//curl -v  --cacert ca/ca.crt --key client/client.key --cert client/client.crt
		//这里读取的是CA根证书
		pool := x509.NewCertPool()
		buf, err := ioutil.ReadFile(web.HTTP.TLS.CaFile)
		if err != nil {
			fmt.Println(err)
		}
		pool.AppendCertsFromPEM(buf)
		cfg := &tls.Config{
			//自签名证书不要做验证
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
			ClientCAs:          pool,
			//对客户端验证
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		cfg.BuildNameToCertificate()

		web.Server = &http.Server{
			Addr:           addr,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 16384,
			TLSConfig:      cfg,
		}
	} else {
		web.Server = &http.Server{
			Addr: addr}
	}
}

// Monitor web服务启动监听
func (web *Web) Monitor() {
	log.Printf("Web service startup listener %s:%s", web.HTTP.IP, web.HTTP.Port)
	if web.HTTP.TLS.Switch {
		err := web.Server.ListenAndServeTLS(web.HTTP.TLS.CertFile, web.HTTP.TLS.KeyFile)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := web.Server.ListenAndServe()
		if err != nil {
			log.Printf("ListenAndServe err: %s", err)
		}
	}

}

// Start web服务启动
func (http HTTP) Start() {
	var web Web
	web.HTTP = http
	web.RunTLS()
	web.Mux()
	web.Monitor()
}

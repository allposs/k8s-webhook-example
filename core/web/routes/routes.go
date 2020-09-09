package routes

import (
	v1 "k8s-webhook-example/core/web/api/v1"
	"net/http"
	"strings"
)

type Router struct {
	Route map[string]map[string]http.HandlerFunc
}

// 路由表初始化
func (r *Router) HandleFunc(method, path string, f http.HandlerFunc) {
	method = strings.ToUpper(method)
	if r.Route == nil {
		r.Route = make(map[string]map[string]http.HandlerFunc)
	}
	if r.Route[method] == nil {
		r.Route[method] = make(map[string]http.HandlerFunc)
	}
	r.Route[method][path] = f
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// 此处我做个很多关于http请求header的限制，一级业务版本处理，就不公布出来了
	// ....
	res.Header().Set("Content-Type", "application/json")
	// 正常路由
	if f, ok := r.Route[req.Method][req.URL.Path]; ok {
		f(res, req)
		//404
	} else {
		// 由于此项目目的是提供接口和WEB HTML服务，所以此处本人做了一系列的url处理
	}
}

func New() *Router {
	route := Router{}
	route.HandleFunc("GET", "/api/v1/token", v1.Token)
	route.HandleFunc("POST", "/api/v1/AuthN", v1.AuthN)
	route.HandleFunc("POST", "/api/v1/AuthZ", v1.AuthZ)
	return &route
}

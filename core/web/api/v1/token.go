package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s-webhook-example/util/encrypt"
	"net/http"

	authn "k8s.io/api/authentication/v1beta1"
	authz "k8s.io/api/authorization/v1beta1"
)

// AuthN 认证(Authentication)，决定谁访问了系统.
// AuthN 完成了用户的认证并且获取了用户的相关信息（如 Username，Groups 等）。
/*
curl -vk \
	-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiYWRtaW4iLCJleHAiOjE1NTQyMjQxMzJ9.Uslc7qnugtXEoZsJsI-zPcw9FE45WY8XdAsc7Y1iZkA" \
	https://k8s:8443/api
*/
func AuthN(w http.ResponseWriter, r *http.Request) {
	var data authn.TokenReview
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println(data)
	defer r.Body.Close()

	json.NewEncoder(w).Encode(Request(data.Spec.Token))
}

// Request TokenReview返回值
func Request(token string) *authn.TokenReview {
	var data authn.TokenReview
	data.APIVersion = "authentication.k8s.io/v1beta1"
	data.Kind = "TokenReview"
	user, err := encrypt.ParseToken(token)
	if err == nil {
		groups := []string{"developers", "qa"}
		extrafield1 := []string{"extravalue1", "extravalue2"}
		extra := map[string]authn.ExtraValue{"extrafield1": extrafield1}
		data.Status.Authenticated = true
		data.Status.User.Extra = extra
		data.Status.User.Groups = groups
		data.Status.User.UID = "42"
		data.Status.User.Username = user.User
		return &data
	}
	data.Status.Authenticated = false
	return &data
}

// AuthZ 主要用于授权 （Authorization），决定访问者具有什么样的权限
// AuthZ 则根据这些信息匹配预定义的权限规则，然后决定某个 API 请求是否被允许。
/*
minikub 使用代理启动
minikube start --registry-mirror=https://registry.docker-cn.com \
--docker-env HTTP_PROXY=http://10.226.144.34:8000 \
--docker-env HTTPS_PROXY=http://10.226.144.34:8000 \
修改minikub 里的kube-apiserver.yaml
		- --authorization-mode=Webhook
    - --runtime-config=authentication.k8s.io/v1beta1=true
    - --authorization-webhook-config-file=/data/webhook/AuthZ.yaml
    - --authorization-webhook-cache-authorized-ttl=5m
    - --authorization-webhook-cache-unauthorized-ttl=30s
    - --authentication-token-webhook-config-file=/data/webhook/AuthN.yaml
    - --authentication-token-webhook-cache-ttl=5m

*/
func AuthZ(w http.ResponseWriter, r *http.Request) {
	var data authz.SubjectAccessReview
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//data2, _ := json.MarshalIndent(data, "", "")
	//fmt.Printf("%s\n", data2)
	defer r.Body.Close()
	var request authz.SubjectAccessReview
	request.APIVersion = "authorization.k8s.io/v1beta1"
	request.Kind = "SubjectAccessReview"
	request.Status.Allowed = true
	json.NewEncoder(w).Encode(request)
}

//Token admin用户的token
func Token(w http.ResponseWriter,
	r *http.Request) {
	data, _ := encrypt.GenerateToken("admin", 400000)
	json.NewEncoder(w).Encode(data)
}

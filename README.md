k8s-webhook-example
========

K8s-webhook-example is an implementation template for AuthZ and AuthN for k8s.

[![Badge](https://img.shields.io/badge/link-996.icu-%23FF4D5B.svg)](https://996.icu/#/en_US)
[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)
[![Slack](https://img.shields.io/badge/slack-996icu-green.svg)](https://join.slack.com/t/996icu/shared_invite/enQtNTg4MjA3MzA1MzgxLWQyYzM5M2IyZmIyMTVjMzU5NTE5MGI5Y2Y2YjgwMmJiMWMxMWMzNGU3NDJmOTdhNmRlYjJlNjk5ZWZhNWIwZGM)

Usage
-----
[Download release](https://github.com/jamescun/switcher/releases) or Build:

    $ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build k8s-webhook-example

    $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build k8s-webhook-example

To get help:

    $ ./k8s-webhook-example  --help
        k8s-webhook-example 1.0.1
        usage: k8s-webhook-example [options]
        Options:
          --port        <80>                                    Server Listen Port
          --ip          <0.0.0.0>                               Web server Address
          --switch      <false>                                 Webhook TLS switch
          --certFile    <config/webhook/TLS/server/server.crt>  TLS cert file path
          --keyFile     <config/webhook/TLS/server/server.key>  TLS key file path
          --caFile      <config/webhook/TLS/ca/ca.crt>          TLS ca file path
        Examples:
          Default webhook Server
          $ k8s-webhook-example
          k8s-webhook-example Listen 127.0.0.1:80 and TLS is flase  
          $ k8s-webhook-example --port 80 --ip 127.0.0.1 --switch flase

Example
-------
Run k8s-webhook-example on port 80 of the machine and turn off TLS.

    $ k8s-webhook-example --port 80 --ip 127.0.0.1 --switch flase

Modify AuthN.yaml and AuthZ.yaml file.

AuthN.yaml

Comment out the relevant TLS content:

    Certificate-authority: /data/webhook/TLS/ca/ca.crt
    Client-certificate: /data/webhook/TLS/client/client.crt
    Client-key: /data/webhook/TLS/client/client.key

Modify content

    Server: https://<k8s-webhook-example address>/api/v1/AuthN

AuthZ.yaml

Comment out the relevant TLS content:

    Certificate-authority: /data/webhook/TLS/ca/ca.crt
    Client-certificate: /data/webhook/TLS/client/client.crt
    Client-key: /data/webhook/TLS/client/client.key

Modify content

    Server: https://<k8s-webhook-example address>/api/v1/AuthZ

Start minikube

Do not use a proxy

    $ minikube start

Using a proxy

    $ minikube start --registry-mirror=https://registry.docker-cn.com --docker-env HTTP_PROXY=http://代理地址端口 --docker-env HTTPS_PROXY=http://代理地址端口

Copy the AuthN.yaml and AuthZ.yaml file to the /data/webhook directory in minikube,The minikube/User directory is interoperable with the user directory and can be copied by this method. Of course, it can be other methods.

    $ minikube ssh

Modify the kube-apiserver.yaml file

    $ vi /etc/kubernetes/manifests/kube-apiserver.yaml

modify

    - --authorization-mode=Node,RBAC,Webhook
New

    - --runtime-config=authentication.k8s.io/v1beta1=true
    - --authorization-webhook-config-file=/data/webhook/AuthZ.yaml
    - --authorization-webhook-cache-authorized-ttl=5m
    - --authorization-webhook-cache-unauthorized-ttl=30s
    - --authentication-token-webhook-config-file=/data/webhook/AuthN.yaml
    - --authentication-token-webhook-cache-ttl=5m

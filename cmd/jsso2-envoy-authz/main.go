package main

import (
	envoy_auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/jrockway/jsso2/pkg/envoyauthz"
	"github.com/jrockway/opinionated-server/server"
	"google.golang.org/grpc"
)

func main() {
	server.AppName = "jsso2-envoy-proxy"
	authzCfg := &envoyauthz.Config{}
	server.AddFlagGroup("Authorization Server", authzCfg)
	server.Setup()
	svc := &envoyauthz.Service{}
	server.AddService(func(s *grpc.Server) {
		envoy_auth.RegisterAuthorizationServer(s, svc)
	})
	server.ListenAndServe()
}

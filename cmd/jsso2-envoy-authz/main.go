package main

import (
	envoy_auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/jrockway/opinionated-server/server"
	"google.golang.org/grpc"
)

func main() {
	server.AppName = "jsso2-envoy-proxy"
	server.Setup()
	server.AddService(func(s *grpc.Server) {
		envoy_auth.RegisterAuthorizationServer(s, &envoy_auth.UnimplementedAuthorizationServer{})
	})
	server.ListenAndServe()
}

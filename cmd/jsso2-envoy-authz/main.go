package main

import (
	"context"
	"time"

	envoy_auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/jrockway/jsso2/pkg/client"
	"github.com/jrockway/jsso2/pkg/envoyauthz"
	"github.com/jrockway/opinionated-server/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.AppName = "jsso2-envoy-proxy"
	authzCfg := &envoyauthz.Config{}
	server.AddFlagGroup("Authorization Server", authzCfg)
	server.Setup()

	startupCtx, c := context.WithTimeout(context.Background(), 15*time.Second)
	cli, err := client.Dial(startupCtx, authzCfg.Address, &client.Credentials{}, grpc.FailOnNonTempDialError(true))
	if err != nil {
		c()
		zap.L().Fatal("problem dialing jsso server", zap.Error(err))
	}
	c()

	svc := &envoyauthz.Service{
		SessionClient:  cli.SessionClient,
		UsernameHeader: authzCfg.AddPlaintextUsernameHeader,
	}
	server.AddService(func(s *grpc.Server) {
		envoy_auth.RegisterAuthorizationServer(s, svc)
	})
	server.ListenAndServe()
}

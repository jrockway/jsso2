package main

import (
	"context"
	"net/http"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/jrockway/jsso2/pkg/foopb"
	"github.com/jrockway/opinionated-server/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.AppName = "jsso2"
	server.AddService(func(s *grpc.Server) {
		foopb.RegisterNameServiceService(s, &foopb.NameServiceService{})
	})
	server.SetStartupCallback(func(info server.Info) {
		// This starts up grpcui on the debug port.
		cc, err := grpc.Dial(info.GRPCAddress, grpc.WithInsecure())
		if err != nil {
			zap.L().Fatal("problem connecting to self", zap.Error(err))
		}
		h, err := standalone.HandlerViaReflection(context.Background(), cc, "self")
		if err != nil {
			zap.L().Fatal("problem creating grpcui handler", zap.Error(err))
		}
		http.Handle("/grpcui/", http.StripPrefix("/grpcui", h))
	})
	server.Setup()
	server.ListenAndServe()
}

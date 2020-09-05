package main

import (
	"github.com/jrockway/jsso2/pkg/foopb"
	"github.com/jrockway/opinionated-server/server"
	"google.golang.org/grpc"
)

func main() {
	server.AppName = "jsso2"
	server.AddService(func(s *grpc.Server) {
		foopb.RegisterNameServiceService(s, &foopb.NameServiceService{})
	})
	server.Setup()
	server.ListenAndServe()
}

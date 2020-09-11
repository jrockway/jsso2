package main

import (
	"context"
	"math/rand"
	"net/http"
	"strings"
	"unicode"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/jrockway/jsso2/pkg/foopb"
	"github.com/jrockway/opinionated-server/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type fooServer struct {
}

func (*fooServer) TransformName(ctx context.Context, req *foopb.TransformNameRequest) (*foopb.TransformNameReply, error) {
	result := new(strings.Builder)
	for _, c := range req.GetName() {
		if rand.Float64() < 0.5 {
			result.WriteRune(unicode.ToLower(c))
		} else {
			result.WriteRune(unicode.ToUpper(c))
		}
	}
	return &foopb.TransformNameReply{Result: result.String()}, nil
}

func main() {
	server.AppName = "jsso2"
	foo := new(fooServer)
	server.AddService(func(s *grpc.Server) {
		foopb.RegisterNameServiceService(s, foopb.NewNameServiceService(foo))
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

package main

import (
	"context"
	"net/http"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jsso/cmd"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/opinionated-server/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.AppName = "jsso2"

	appConfig := &cmd.Config{}
	server.AddFlagGroup("application", appConfig)
	dbConfig := &store.Config{}
	server.AddFlagGroup("database", dbConfig)
	authConfig := &internalauth.Config{}
	server.AddFlagGroup("authorization", authConfig)

	server.Setup()
	l := zap.L().Named("startup")

	db, err := cmd.ConnectDB(l, dbConfig)
	if err != nil {
		zap.L().Fatal("problem connecting to database", zap.Error(err))
	}
	app, err := cmd.Setup(appConfig, authConfig, db)
	if err != nil {
		zap.L().Fatal("problem initializing app", zap.Error(err))
	}
	server.AddUnaryInterceptor(app.Permissions.UnaryServerInterceptor())
	server.AddStreamInterceptor(app.Permissions.StreamServerInterceptor())
	server.SetHTTPHandler(app.PublicMux)
	server.AddService(func(s *grpc.Server) {
		jssopb.RegisterEnrollmentService(s, jssopb.NewEnrollmentService(app.EnrollmentService))
		jssopb.RegisterUserService(s, jssopb.NewUserService(app.UserService))
		jssopb.RegisterLoginService(s, jssopb.NewLoginService(app.LoginService))
		jssopb.RegisterSessionService(s, jssopb.NewSessionService(app.SessionService))
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

	server.ListenAndServe()
}

package main

import (
	"context"
	"net/http"
	"time"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jsso/enrollment"
	"github.com/jrockway/jsso2/pkg/jsso/login"
	"github.com/jrockway/jsso2/pkg/jsso/user"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/opinionated-server/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	server.AppName = "jsso2"

	dbConfig := &store.Config{}
	server.AddFlagGroup("database", dbConfig)
	authConfig := &internalauth.Config{}
	server.AddFlagGroup("authorization", authConfig)

	server.Setup()

	startupCtx, c := context.WithTimeout(context.Background(), time.Minute)
	db, err := store.Connect(startupCtx, dbConfig.DatabaseURL)
	if err != nil {
		zap.L().Fatal("failed to connect to database", zap.String("database_url", dbConfig.DatabaseURL), zap.Error(err))
	}
	if dbConfig.RunMigrations {
		zap.L().Info("running database migrations")
		if err := db.MigrateDB(startupCtx); err != nil {
			zap.L().Warn("failed to run database migrations; continuing anyway", zap.Error(err))
		}
	}

	auth := internalauth.NewFromConfig(authConfig, db)
	server.AddUnaryInterceptor(auth.UnaryServerInterceptor())
	server.AddStreamInterceptor(auth.StreamServerInterceptor())

	c()

	server.AddService(func(s *grpc.Server) {
		jssopb.RegisterEnrollmentService(s, jssopb.NewEnrollmentService(&enrollment.Service{}))
		jssopb.RegisterUserService(s, jssopb.NewUserService(&user.Service{DB: db}))
		jssopb.RegisterLoginService(s, jssopb.NewLoginService(&login.Service{}))
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

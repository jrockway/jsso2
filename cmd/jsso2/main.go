package main

import (
	"context"
	"net/http"
	"net/url"
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

type Config struct {
	BaseURL string `long:"base_url" description:"Where the app's public resources are available; used for generating links and cookies." env:"BASE_URL" default:"http://localhost:4000"`
}

func main() {
	server.AppName = "jsso2"

	appConfig := &Config{}
	server.AddFlagGroup("application", appConfig)
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
	c()

	auth := internalauth.NewFromConfig(authConfig, db)
	server.AddUnaryInterceptor(auth.UnaryServerInterceptor())
	server.AddStreamInterceptor(auth.StreamServerInterceptor())

	baseURL, err := url.Parse(appConfig.BaseURL)
	if err != nil {
		zap.L().Fatal("failed to parse base URL", zap.String("base_url", appConfig.BaseURL), zap.Error(err))
	}

	userService := &user.Service{
		DB:          db,
		Permissions: auth,
		BaseURL:     baseURL,
	}

	enrollmentService := &enrollment.Service{
		Permissions: auth,
		Origin:      baseURL.Host,
	}

	server.AddService(func(s *grpc.Server) {
		jssopb.RegisterEnrollmentService(s, jssopb.NewEnrollmentService(enrollmentService))
		jssopb.RegisterUserService(s, jssopb.NewUserService(userService))
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

package main

import (
	"context"
	"net/http"
	"time"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/jrockway/jsso2/pkg/cookies"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jsso/enrollment"
	"github.com/jrockway/jsso2/pkg/jsso/login"
	"github.com/jrockway/jsso2/pkg/jsso/session"
	"github.com/jrockway/jsso2/pkg/jsso/user"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/logout"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/web"
	"github.com/jrockway/jsso2/pkg/webauthn"
	"github.com/jrockway/opinionated-server/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	BaseURL      string `long:"base_url" description:"Where the app's public resources are available; used for generating links and cookies." env:"BASE_URL" default:"http://localhost:4000"`
	SetCookieKey string `long:"set_cookie_key" description:"32 bytes that are used to encrypt and sign set-cookie tokens." env:"SET_COOKIE_KEY"`
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

	linker, err := web.NewLinker(appConfig.BaseURL)
	if err != nil {
		zap.L().Fatal("failed to create linker", zap.String("base_url", appConfig.BaseURL), zap.Error(err))
	}

	cookieConfig := &cookies.Config{
		Name:   "jsso-session-id",
		Domain: linker.Domain(),
		Linker: linker,
	}
	if err := cookieConfig.SetKey([]byte(appConfig.SetCookieKey)); err != nil {
		zap.L().Fatal("failed to set set-cookie encryption key", zap.Error(err))
	}

	auth := internalauth.NewFromConfig(authConfig, db)
	auth.Cookies = cookieConfig
	server.AddUnaryInterceptor(auth.UnaryServerInterceptor())
	server.AddStreamInterceptor(auth.StreamServerInterceptor())

	webauthnConfig := &webauthn.Config{
		RelyingPartyID:   linker.Domain(),
		RelyingPartyName: linker.Domain(),
		Origin:           linker.Origin(),
	}

	userService := &user.Service{
		DB:          db,
		Permissions: auth,
		Linker:      linker,
	}

	enrollmentService := &enrollment.Service{
		DB:          db,
		Permissions: auth,
		Linker:      linker,
		Webauthn:    webauthnConfig,
	}

	loginService := &login.Service{
		DB:          db,
		Permissions: auth,
		Webauthn:    webauthnConfig,
		Cookies:     cookieConfig,
	}

	sessionService := &session.Service{}

	logoutHandler := &logout.Handler{
		Linker:  linker,
		Cookies: cookieConfig,
		DB:      db,
	}

	publicMux := new(http.ServeMux)
	publicMux.HandleFunc("/set-cookie", cookieConfig.HandleSetCookie)
	publicMux.Handle("/logout", logoutHandler)
	server.SetHTTPHandler(publicMux)

	server.AddService(func(s *grpc.Server) {
		jssopb.RegisterEnrollmentService(s, jssopb.NewEnrollmentService(enrollmentService))
		jssopb.RegisterUserService(s, jssopb.NewUserService(userService))
		jssopb.RegisterLoginService(s, jssopb.NewLoginService(loginService))
		jssopb.RegisterSessionService(s, jssopb.NewSessionService(sessionService))
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

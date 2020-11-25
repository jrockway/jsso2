package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jrockway/jsso2/pkg/cookies"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jsso/enrollment"
	"github.com/jrockway/jsso2/pkg/jsso/login"
	"github.com/jrockway/jsso2/pkg/jsso/session"
	"github.com/jrockway/jsso2/pkg/jsso/user"
	"github.com/jrockway/jsso2/pkg/logout"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/web"
	"github.com/jrockway/jsso2/pkg/webauthn"
	"go.uber.org/zap"
)

type Config struct {
	BaseURL      string `long:"base_url" description:"Where the app's public resources are available; used for generating links and cookies." env:"BASE_URL" default:"http://localhost:4000"`
	SetCookieKey string `long:"set_cookie_key" description:"32 bytes that are used to encrypt and sign set-cookie tokens." env:"SET_COOKIE_KEY"`
}

type App struct {
	DB             *store.Connection
	Linker         *web.Linker
	Cookies        *cookies.Config
	WebauthnConfig *webauthn.Config
	Permissions    *internalauth.Permissions

	UserService       *user.Service
	EnrollmentService *enrollment.Service
	LoginService      *login.Service
	SessionService    *session.Service

	PublicMux *http.ServeMux
}

func Setup(appConfig *Config, authConfig *internalauth.Config, db *store.Connection) (*App, error) {
	app := &App{DB: db}
	linker, err := web.NewLinker(appConfig.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("create linker for baseurl %s: %w", appConfig.BaseURL, err)
	}
	app.Linker = linker

	cookieConfig := &cookies.Config{
		Name:   "jsso-session-id",
		Domain: linker.Domain(),
		Linker: linker,
	}
	if err := cookieConfig.SetKey([]byte(appConfig.SetCookieKey)); err != nil {
		return nil, fmt.Errorf("set set-cookie encryption key: %w", err)
	}
	app.Cookies = cookieConfig

	app.Permissions = internalauth.NewFromConfig(authConfig, db)
	app.Permissions.Cookies = cookieConfig

	webauthnConfig := &webauthn.Config{
		RelyingPartyID:   linker.Domain(),
		RelyingPartyName: linker.Domain(),
		Origin:           linker.Origin(),
	}
	app.WebauthnConfig = webauthnConfig

	app.UserService = &user.Service{
		DB:          db,
		Permissions: app.Permissions,
		Linker:      linker,
	}
	app.EnrollmentService = &enrollment.Service{
		DB:          db,
		Permissions: app.Permissions,
		Linker:      linker,
		Webauthn:    webauthnConfig,
	}
	app.LoginService = &login.Service{
		DB:          db,
		Permissions: app.Permissions,
		Webauthn:    webauthnConfig,
		Cookies:     cookieConfig,
	}
	app.SessionService = &session.Service{
		DB:          db,
		Permissions: app.Permissions,
		Cookies:     cookieConfig,
		Linker:      linker,
	}

	logoutHandler := &logout.Handler{
		Linker:  linker,
		Cookies: cookieConfig,
		DB:      db,
	}
	app.PublicMux = new(http.ServeMux)
	app.PublicMux.HandleFunc("/set-cookie", cookieConfig.HandleSetCookie)
	app.PublicMux.Handle("/logout", logoutHandler)

	return app, nil
}

func ConnectDB(l *zap.Logger, dbConfig *store.Config) (*store.Connection, error) {
	startupCtx, c := context.WithTimeout(context.Background(), time.Minute)
	defer c()
	db, err := store.Connect(startupCtx, dbConfig.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect database at %q: %w", dbConfig.DatabaseURL, err)
	}
	if dbConfig.RunMigrations {
		l.Info("running database migrations")
		if err := db.MigrateDB(startupCtx); err != nil {
			l.Warn("failed to run database migrations; continuing anyway", zap.Error(err))
		}
	}
	return db, nil
}

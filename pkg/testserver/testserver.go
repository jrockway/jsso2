package testserver

import (
	"net/url"
	"testing"

	gzap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/jrockway/jsso2/pkg/client"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jsso/enrollment"
	"github.com/jrockway/jsso2/pkg/jsso/login"
	"github.com/jrockway/jsso2/pkg/jsso/user"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/web"
	"github.com/jrockway/jsso2/pkg/webauthn"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type S struct {
	WantRootClient bool
	Credentials    credentials.PerRPCCredentials
	Permissions    *internalauth.Permissions
	Logger         *zap.Logger
}

func New() *S {
	return &S{
		Logger:      zap.NewNop(),
		Permissions: internalauth.NewFromConfig(&internalauth.Config{RootPassword: "root"}, nil),
	}
}

func (s *S) Options(e *jtesting.E) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			s.Permissions.UnaryServerInterceptor(),
			gzap.UnaryServerInterceptor(e.Logger.Named("server")),
		),
		grpc.ChainStreamInterceptor(
			s.Permissions.StreamServerInterceptor(),
			gzap.StreamServerInterceptor(e.Logger.Named("server")),
		),
	}
}

func (s *S) Setup(t *testing.T, e *jtesting.E, server *grpc.Server) {
	db := store.MustGetTestDB(t, e)
	if s.WantRootClient {
		s.Credentials = &client.Credentials{
			Token: "root root",
		}
	}
	s.Permissions.Store = db
	linker := &web.Linker{
		BaseURL: &url.URL{Scheme: "http", Host: "jsso.example.com", Path: "/"},
	}
	webauthnConfig := &webauthn.Config{
		RelyingPartyID:   linker.Domain(),
		RelyingPartyName: linker.Domain(),
		Origin:           linker.Origin(),
	}

	jssopb.RegisterEnrollmentService(server, jssopb.NewEnrollmentService(&enrollment.Service{DB: db, Permissions: s.Permissions, Linker: linker, Webauthn: webauthnConfig}))
	jssopb.RegisterUserService(server, jssopb.NewUserService(&user.Service{DB: db, Permissions: s.Permissions, Linker: linker}))
	jssopb.RegisterLoginService(server, jssopb.NewLoginService(&login.Service{DB: db, Permissions: s.Permissions, Webauthn: webauthnConfig}))
}

// OK, maybe I went overboard with single-letter type names.
func (s *S) ToR(r *jtesting.R) {
	r.GRPCOptions = s.Options
	r.GRPCClientOptions = func(e *jtesting.E) []grpc.DialOption {
		result := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithChainUnaryInterceptor(
				gzap.UnaryClientInterceptor(e.Logger.Named("client")),
			),
			grpc.WithChainStreamInterceptor(
				gzap.StreamClientInterceptor(e.Logger.Named("client")),
			),
		}
		if s.Credentials != nil {
			result = append(result, grpc.WithPerRPCCredentials(s.Credentials))
		}
		return result
	}
	r.GRPC = s.Setup
}

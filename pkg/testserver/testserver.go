package testserver

import (
	"testing"

	gzap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/jrockway/jsso2/pkg/client"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jsso/cmd"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/store"
	"google.golang.org/grpc"
)

type S struct {
	AppConfig      *cmd.Config
	AuthConfig     *internalauth.Config
	WantRootClient bool
	Credentials    *client.Credentials
	App            *cmd.App
}

func New() *S {
	return &S{
		AppConfig: &cmd.Config{
			BaseURL:  "http://jsso.example.com/",
			TokenKey: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		},
		AuthConfig: &internalauth.Config{
			RootPassword: "root",
		},
	}
}

func (s *S) Setup(t *testing.T, e *jtesting.E, server *grpc.Server) {
	jssopb.RegisterEnrollmentService(server, jssopb.NewEnrollmentService(s.App.EnrollmentService))
	jssopb.RegisterUserService(server, jssopb.NewUserService(s.App.UserService))
	jssopb.RegisterLoginService(server, jssopb.NewLoginService(s.App.LoginService))
	jssopb.RegisterSessionService(server, jssopb.NewSessionService(s.App.SessionService))
}

// OK, maybe I went overboard with single-letter type names.
func (s *S) ToR(r *jtesting.R) {
	r.Database = true
	r.Logger = true
	if s.WantRootClient {
		s.Credentials = &client.Credentials{Root: "root"}
	}
	r.DatabaseReady = func(t *testing.T, e *jtesting.E) {
		db := store.MustGetTestDB(t, e)
		var err error
		s.App, err = cmd.Setup(s.AppConfig, s.AuthConfig, db)
		if err != nil {
			t.Fatal(err)
		}
	}
	r.GRPCOptions = func(e *jtesting.E) []grpc.ServerOption {
		return []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				gzap.UnaryServerInterceptor(e.Logger.Named("server")),
				s.App.Permissions.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				gzap.StreamServerInterceptor(e.Logger.Named("server")),
				s.App.Permissions.StreamServerInterceptor(),
			),
		}
	}
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

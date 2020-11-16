// Package internalauth manages authorizing gRPC calls.
package internalauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Config struct {
	RootPassword string `long:"root_password" env:"ROOT_PASSWORD" description:"If set, allow a requestor full privileges if they include this password in their requests.  Should only be used to bootstrap a normal administrative user."`
}

// Permissions manages all authorization in JSSO.
type Permissions struct {
	// If set, a password that can be provided to bypass all access controls.
	RootPassword string

	// RPCs that can be called without any credentials.
	AllowedWithoutAuth map[string]struct{}

	Store *store.Connection
}

// NewFromConfig builds a Permissions object from configuration.
func NewFromConfig(c *Config, s *store.Connection) *Permissions {
	return &Permissions{
		Store:        s,
		RootPassword: c.RootPassword,
		AllowedWithoutAuth: map[string]struct{}{
			"/grpc.health.v1.Health/Check":                                   {},
			"/grpc.health.v1.Health/Watch":                                   {},
			"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {},
			"/jsso.Login/Start":                                              {},
		},
	}
}

// AuthorizeRPC returns whether the credentials provided allow the RPC to be called.
func (p *Permissions) AuthorizeRPC(ctx context.Context, fullMethod string) error {
	if _, ok := p.AllowedWithoutAuth[fullMethod]; ok {
		return nil
	}

	if _, ok := sessions.FromContext(ctx); ok {
		return nil
	}

	return status.Error(codes.Unauthenticated, "no authentication credentials supplied")
}

func (p *Permissions) isRoot(md metadata.MD) bool {
	if p.RootPassword == "" {
		return false
	}
	want := fmt.Sprintf("root %s", p.RootPassword)
	for _, auth := range md.Get("Authorization") {
		if auth == want {
			return true
		}
	}
	return false
}

func (p *Permissions) sessionToContext(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, errors.New("no metadata in incoming context")
	}

	if p.isRoot(md) {
		root := sessions.Root()
		return sessions.NewContext(ctx, root), nil
	}

	sid, err := sessions.FromMetadata(md)
	if err != nil {
		if errors.Is(err, sessions.ErrSessionMissing) {
			// No session is not an error.
			return ctx, nil
		}
		return ctx, fmt.Errorf("load session from metadata: %w", err)
	}
	var session *types.Session
	err = p.Store.DoTx(ctx, ctxzap.Extract(ctx), true, func(tx *sqlx.Tx) error {
		var err error
		session, err = store.LookupSession(ctx, tx, sid.GetId())
		if err != nil {
			return fmt.Errorf("lookup session: %w", err)
		}
		return nil
	})
	if err != nil {
		return ctx, fmt.Errorf("read session from database: %w", err)
	}
	l := ctxzap.Extract(ctx)
	l = l.With(zap.String("user", session.GetUser().GetUsername()))
	ctx = ctxzap.ToContext(ctx, l)
	return sessions.NewContext(ctx, session), nil
}

func (p *Permissions) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		rootCtx := ss.Context()
		ctx, err := p.sessionToContext(rootCtx)
		if err != nil {
			return status.Error(codes.Unauthenticated, fmt.Sprintf("get user from session: %v", err))
		}
		if err := p.AuthorizeRPC(ctx, info.FullMethod); err != nil {
			l := ctxzap.Extract(ctx)
			l.Debug("user not authorized to perform RPC", zap.Error(err))
			return err
		}
		return handler(srv, ss)
	}
}

func (p *Permissions) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(rootCtx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, err := p.sessionToContext(rootCtx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("get user from session: %v", err))
		}
		if err := p.AuthorizeRPC(ctx, info.FullMethod); err != nil {
			l := ctxzap.Extract(ctx)
			l.Debug("user not authorized to perform RPC", zap.Error(err))
			return nil, err
		}
		return handler(ctx, req)
	}
}

// General policy decisions start here.
func (p *Permissions) EnrollmentSessionPrototype(ctx context.Context, target *types.User) (*types.Session, error) {
	id, err := sessions.GenerateID()
	if err != nil {
		return nil, fmt.Errorf("generate session id: %w", err)
	}
	now := time.Now()
	return &types.Session{
		Id:        id,
		User:      target,
		CreatedAt: timestamppb.New(now),
		ExpiresAt: timestamppb.New(now.Add(3 * 24 * time.Hour)),
	}, nil
}

func (p *Permissions) LoginSessionPrototype(ctx context.Context, target *types.User) (*types.Session, error) {
	id, err := sessions.GenerateID()
	if err != nil {
		return nil, fmt.Errorf("generate session id: %w", err)
	}
	now := time.Now()
	return &types.Session{
		Id:        id,
		User:      target,
		CreatedAt: timestamppb.New(now),
		ExpiresAt: timestamppb.New(now.Add(18 * time.Hour)),
	}, nil
}

// The per-operation permissions start here.

func (p *Permissions) AllowUserEdit(ctx context.Context, target *types.User, actor *types.Session) error {
	return nil
}

func (p *Permissions) AllowGenerateEnrollmentLink(ctx context.Context, target *types.User, actor *types.Session) error {
	return nil
}

func (p *Permissions) AllowStartEnrollment(ctx context.Context, target *types.Session) error {
	return nil
}

func (p *Permissions) AllowFinishEnrollment(ctx context.Context, target *types.Session) error {
	return nil
}

func (p *Permissions) AllowStartLogin(ctx context.Context, target *types.User) error {
	return nil
}

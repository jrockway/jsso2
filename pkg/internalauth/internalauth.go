// Package internalauth manages authorizing gRPC calls.
package internalauth

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Config struct {
	RootPassword string `long:"root_password" env:"ROOT_PASSWORD" description:"If set, allow a requestor full privileges if they include this password in their requests.  Should only be used to bootstrap a normal administrative user."`
}

type Interceptor struct {
	rootPassword string
}

func NewFromConfig(c *Config) *Interceptor {
	return &Interceptor{
		rootPassword: c.RootPassword,
	}
}

func (i *Interceptor) IsRoot(md metadata.MD) bool {
	if i.rootPassword == "" {
		return false
	}
	want := fmt.Sprintf("root %s", i.rootPassword)
	for _, auth := range md.Get("Authorization") {
		if auth == want {
			return true
		}
	}
	return false
}

var allowedWithoutAuth = map[string]struct{}{
	"/grpc.health.v1.Health/Check":                                   {},
	"/grpc.health.v1.Health/Watch":                                   {},
	"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {},
}

func (i *Interceptor) Authorize(method string, md metadata.MD) error {
	if i.IsRoot(md) {
		return nil
	}
	if _, ok := allowedWithoutAuth[method]; ok {
		return nil
	}
	return status.Error(codes.PermissionDenied, "user is not root and method is not allowed for non-root user")
}

func (i *Interceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		if err := i.Authorize(info.FullMethod, md); err != nil {
			l := ctxzap.Extract(ctx)
			l.Debug("user not authorized to perform RPC", zap.Error(err))
			return err
		}
		return handler(srv, ss)
	}
}

func (i *Interceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		if err := i.Authorize(info.FullMethod, md); err != nil {
			l := ctxzap.Extract(ctx)
			l.Debug("user not authorized to perform RPC", zap.Error(err))
			return nil, err
		}
		return handler(ctx, req)
	}
}

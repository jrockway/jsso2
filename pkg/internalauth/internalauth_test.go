package internalauth

import (
	"context"
	"testing"

	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func TestAuthorizeRPC(t *testing.T) {
	jtesting.Run(t, "authorizerpc", jtesting.R{Logger: true, Database: true}, func(t *testing.T, e *jtesting.E) {
		c := store.MustGetTestDB(t, e)
		p := NewFromConfig(&Config{RootPassword: "foo"}, c)
		p.AllowedWithoutAuth = map[string]struct{}{
			"ok": {},
		}

		// Unauthenticated.
		if err := p.AuthorizeRPC(e.Context, "bad"); err == nil {
			t.Errorf("rpc 'bad': expected error")
		}
		if err := p.AuthorizeRPC(e.Context, "ok"); err != nil {
			t.Errorf("rpc 'ok': unexpected error: %v", err)
		}

		// With valid session.
		s := store.ValidSession(t, e, c)
		ctx := metadata.NewIncomingContext(e.Context, metadata.Pairs("authorization", sessions.ToHeaderString(s)))
		var err error
		ctx, err = p.sessionToContext(ctx)
		if err != nil {
			t.Fatalf("load session into context: %v", err)
		}
		if err := p.AuthorizeRPC(ctx, "bad"); err != nil {
			t.Errorf("rpc 'bad' with auth: unexpected error: %v", err)
		}

		// With root password.
		ctx = metadata.NewIncomingContext(e.Context, metadata.Pairs("authorization", "root foo"))
		ctx, err = p.sessionToContext(ctx)
		if err != nil {
			t.Fatalf("load root into context: %v", err)
		}
		if err := p.AuthorizeRPC(ctx, "bad"); err != nil {
			t.Errorf("rpc 'bad' with root password: unexpected error: %v", err)
		}
	})
}

func TestInterceptor(t *testing.T) {
	p := NewFromConfig(&Config{RootPassword: "foo"}, nil)
	p.AllowedWithoutAuth = map[string]struct{}{}

	h := health.NewServer()
	setupGRPC := func(t *testing.T, e *jtesting.E, s *grpc.Server) {
		grpc_health_v1.RegisterHealthServer(s, h)
	}
	r := jtesting.R{
		Logger:   true,
		Database: false,
		GRPC:     setupGRPC,
		GRPCOptions: func(e *jtesting.E) []grpc.ServerOption {
			return []grpc.ServerOption{
				grpc.UnaryInterceptor(p.UnaryServerInterceptor()),
				grpc.StreamInterceptor(p.StreamServerInterceptor()),
			}
		},
	}
	jtesting.Run(t, "testinterceptor", r, func(t *testing.T, e *jtesting.E) {
		pctx := metadata.NewOutgoingContext(e.Context, metadata.Pairs("authorization", "root foo"))
		c := grpc_health_v1.NewHealthClient(e.ClientConn)
		testData := []struct {
			name string
			ctx  context.Context
			ok   bool
		}{
			{"without password", e.Context, false}, {"with password", pctx, true},
		}
		for _, test := range testData {
			t.Run(test.name, func(t *testing.T) {
				if _, err := c.Check(test.ctx, &grpc_health_v1.HealthCheckRequest{}); err == nil && !test.ok {
					t.Errorf("check: expected error")
				} else if err != nil && test.ok {
					t.Errorf("check: unexpected error: %v", err)
				}

				w, err := c.Watch(test.ctx, &grpc_health_v1.HealthCheckRequest{})
				if err != nil {
					t.Fatalf("watch: creating server failed: %v", err)
				}
				if _, err := w.Recv(); err == nil && !test.ok {
					t.Errorf("watch: expected error")
				} else if err != nil && test.ok {
					t.Errorf("watch: unexpected error: %v", err)
				}
			})
		}
	})
}

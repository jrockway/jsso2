package user

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/web"
)

type Service struct {
	DB          *store.Connection
	Permissions *internalauth.Permissions
	Linker      *web.Linker
}

// Edit implements jssopb.UserService.
func (s *Service) Edit(ctx context.Context, req *jssopb.EditUserRequest) (*jssopb.EditUserReply, error) {
	reply := new(jssopb.EditUserReply)
	if err := s.Permissions.AllowUserEdit(ctx, req.GetUser(), sessions.MustFromContext(ctx)); err != nil {
		return reply, fmt.Errorf("check permissions: %w", err)
	}

	if err := s.DB.DoTx(ctx, ctxzap.Extract(ctx), false, func(tx *sqlx.Tx) error {
		user := req.GetUser()
		err := store.UpdateUser(ctx, tx, user)
		if err != nil {
			return err
		}
		reply.User = user
		return nil
	}); err != nil {
		return reply, store.AsGRPCError(fmt.Errorf("update user: %w", err))
	}
	return reply, nil
}

func (s *Service) GenerateEnrollmentLink(ctx context.Context, req *jssopb.GenerateEnrollmentLinkRequest) (*jssopb.GenerateEnrollmentLinkReply, error) {
	reply := new(jssopb.GenerateEnrollmentLinkReply)
	if err := s.DB.DoTx(ctx, ctxzap.Extract(ctx), false, func(tx *sqlx.Tx) error {
		if err := store.LookupUser(ctx, tx, req.GetTarget()); err != nil {
			return fmt.Errorf("lookup target user: %w", err)
		}
		if err := s.Permissions.AllowGenerateEnrollmentLink(ctx, req.GetTarget(), sessions.MustFromContext(ctx)); err != nil {
			return fmt.Errorf("check permissions: %w", err)
		}
		session, err := s.Permissions.EnrollmentSessionPrototype(ctx, req.GetTarget())
		if err != nil {
			return fmt.Errorf("generate session prototype: %w", err)
		}
		if err := store.UpdateSession(ctx, tx, session); err != nil {
			return fmt.Errorf("store session: %w", err)
		}
		reply.Token = sessions.ToBase64(session)
		reply.Url = s.Linker.EnrollmentPage(reply.Token)
		return nil
	}); err != nil {
		return reply, store.AsGRPCError(fmt.Errorf("generate enrollment link: %w", err))
	}
	return reply, nil
}

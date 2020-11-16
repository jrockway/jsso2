package enrollment

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/web"
	"github.com/jrockway/jsso2/pkg/webauthn"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	DB          *store.Connection
	Permissions *internalauth.Permissions
	Linker      *web.Linker
}

func (s *Service) Start(ctx context.Context, req *jssopb.StartEnrollmentRequest) (*jssopb.StartEnrollmentReply, error) {
	reply := &jssopb.StartEnrollmentReply{}
	session := sessions.MustFromContext(ctx)
	if err := s.Permissions.AllowStartEnrollment(ctx, session); err != nil {
		return reply, fmt.Errorf("check permissions: %w", err)
	}
	user := session.GetUser()
	reply.User = user

	l := ctxzap.Extract(ctx)
	var creds []*types.Credential
	if err := s.DB.DoTx(ctx, l, true, func(tx *sqlx.Tx) error {
		c, err := store.GetUserCredentials(ctx, tx, user)
		if err != nil {
			return fmt.Errorf("get credentials for user %s: %w", user.GetUsername(), err)
		}
		creds = c
		return nil
	}); err != nil {
		return reply, store.AsGRPCError(err)
	}

	opts, err := webauthn.BeginEnrollment(s.Linker.RPID(), session, creds)
	if err != nil {
		return reply, fmt.Errorf("create challenge: %w", err)
	}
	reply.CredentialCreationOptions = opts
	return reply, nil
}

func (s *Service) Finish(ctx context.Context, req *jssopb.FinishEnrollmentRequest) (*jssopb.FinishEnrollmentReply, error) {
	reply := &jssopb.FinishEnrollmentReply{}
	session := sessions.MustFromContext(ctx)
	if err := s.Permissions.AllowFinishEnrollment(ctx, session); err != nil {
		return reply, fmt.Errorf("check permissions: %w", err)
	}
	credential, err := webauthn.FinishEnrollment(s.Linker.RPID(), s.Linker.Origin(), session, req)
	if err != nil {
		return reply, fmt.Errorf("validate credential: %w", err)
	}
	credential.Id = 0
	credential.Name = "unnamed"
	credential.User = session.GetUser()
	credential.CreatedAt = timestamppb.Now()
	credential.CreatedBySessionId = session.GetId()
	l := ctxzap.Extract(ctx)
	if err := s.DB.DoTx(ctx, l, false, func(tx *sqlx.Tx) error {
		if err := store.AddCredential(ctx, tx, credential); err != nil {
			return fmt.Errorf("add credential: %w", err)
		}
		s, err := store.LookupSession(ctx, tx, session.GetId())
		if err != nil {
			return fmt.Errorf("lookup session: %w", err)
		}
		if sessions.HasTaint(s, sessions.TaintEnrollment) {
			s.ExpiresAt = timestamppb.Now()
			if err := store.UpdateSession(ctx, tx, s); err != nil {
				return fmt.Errorf("expire session: %w", err)
			}
		}
		return nil
	}); err != nil {
		return reply, store.AsGRPCError(err)
	}
	l.Debug("enrolled new credential", zap.Binary("credential_id", credential.GetCredentialId()))
	reply.LoginUrl = s.Linker.LoginPage()
	return reply, nil
}

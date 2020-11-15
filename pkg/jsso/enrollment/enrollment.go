package enrollment

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/webauthn"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	Domain, Origin string
	Permissions    *internalauth.Permissions
}

func (s *Service) Start(ctx context.Context, req *jssopb.StartEnrollmentRequest) (*jssopb.StartEnrollmentReply, error) {
	reply := &jssopb.StartEnrollmentReply{}
	session := sessions.MustFromContext(ctx)
	if err := s.Permissions.AllowEnrollment(ctx, session); err != nil {
		return reply, fmt.Errorf("check permissions: %w", err)
	}
	user := session.GetUser()
	reply.User = user

	opts, err := webauthn.BeginEnrollment(s.Domain, session)
	if err != nil {
		return reply, fmt.Errorf("create challenge: %w", err)
	}
	reply.CredentialCreationOptions = opts
	return reply, nil
}

func (s *Service) Finish(ctx context.Context, req *jssopb.FinishEnrollmentRequest) (*jssopb.FinishEnrollmentReply, error) {
	reply := &jssopb.FinishEnrollmentReply{}
	session := sessions.MustFromContext(ctx)
	credential, err := webauthn.FinishEnrollment(s.Domain, s.Origin, session, req)
	if err != nil {
		return reply, fmt.Errorf("validate credential: %w", err)
	}
	credential.Id = 0
	credential.Name = "not implemented yet"
	credential.User = session.GetUser()
	credential.CreatedAt = timestamppb.Now()
	credential.CreatedBySessionId = session.GetId()
	ctxzap.Extract(ctx).Debug("added credential", zap.Any("credential", credential))
	return reply, nil
}

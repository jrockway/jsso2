package enrollment

import (
	"context"
	"fmt"

	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/webauthn"
)

type Service struct {
	Origin      string
	Permissions *internalauth.Permissions
}

func (s *Service) Start(ctx context.Context, req *jssopb.StartEnrollmentRequest) (*jssopb.StartEnrollmentReply, error) {
	reply := &jssopb.StartEnrollmentReply{}
	session := sessions.MustFromContext(ctx)
	if err := s.Permissions.AllowEnrollment(ctx, session); err != nil {
		return reply, fmt.Errorf("check permissions: %w", err)
	}
	user := session.GetUser()
	reply.User = user

	opts, err := webauthn.BeginEnrollment(s.Origin, session)
	if err != nil {
		return reply, fmt.Errorf("create challenge: %w", err)
	}
	reply.CredentialCreationOptions = opts
	return reply, nil
}

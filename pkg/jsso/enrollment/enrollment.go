package enrollment

import (
	"context"
	"fmt"

	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/sessions"
)

type Service struct {
	Permissions *internalauth.Permissions
}

func (s *Service) Start(ctx context.Context, req *jssopb.StartEnrollmentRequest) (*jssopb.StartEnrollmentReply, error) {
	reply := &jssopb.StartEnrollmentReply{}
	session := sessions.MustFromContext(ctx)
	if err := s.Permissions.AllowEnrollment(ctx, session); err != nil {
		return reply, fmt.Errorf("check permissions: %w", err)
	}
	reply.User = session.GetUser()
	return reply, nil
}

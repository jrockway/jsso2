package login

import (
	"context"

	"github.com/jrockway/jsso2/pkg/jssopb"
)

type Service struct{}

func (s *Service) Start(ctx context.Context, req *jssopb.StartLoginRequest) (*jssopb.StartLoginReply, error) {
	return &jssopb.StartLoginReply{}, nil
}

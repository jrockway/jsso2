package login

import (
	"context"
	"errors"

	"github.com/jrockway/jsso2/pkg/jssopb"
)

type Service struct{}

func (s *Service) Start(ctx context.Context, req *jssopb.StartLoginRequest) (*jssopb.StartLoginReply, error) {
	return nil, errors.New("unimplemented")
}

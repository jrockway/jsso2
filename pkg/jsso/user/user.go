package user

import (
	"context"
	"errors"

	"github.com/jrockway/jsso2/pkg/jssopb"
)

type Service struct{}

func (s *Service) Add(ctx context.Context, req *jssopb.AddUserRequest) (*jssopb.AddUserReply, error) {
	return nil, errors.New("unimplemented")
}

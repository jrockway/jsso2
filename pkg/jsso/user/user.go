package user

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/store"
)

type Service struct {
	DB *store.Connection
}

// Edit implements jssopb.UserService.
func (s *Service) Edit(ctx context.Context, req *jssopb.EditUserRequest) (*jssopb.EditUserReply, error) {
	reply := new(jssopb.EditUserReply)
	err := s.DB.DoTx(ctx, ctxzap.Extract(ctx), false, func(tx *sqlx.Tx) error {
		user := req.GetUser()
		err := store.UpdateUser(ctx, tx, user)
		if err != nil {
			return err
		}
		reply.User = user
		return nil
	})
	return reply, store.AsGRPCError(err)
}

package login

import (
	"context"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/webauthnpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Service struct {
	DB          *store.Connection
	Permissions *internalauth.Permissions
}

func (s *Service) Start(ctx context.Context, req *jssopb.StartLoginRequest) (*jssopb.StartLoginReply, error) {
	reply := &jssopb.StartLoginReply{
		CredentialRequestOptions: &webauthnpb.PublicKeyCredentialRequestOptions{
			Timeout: durationpb.New(60 * time.Second),
		},
	}

	l := ctxzap.Extract(ctx)
	if err := s.DB.DoTx(ctx, l, false, func(tx *sqlx.Tx) error {
		user := &types.User{
			Username: req.GetUsername(),
		}
		if err := store.LookupUser(ctx, tx, user); err != nil {
			return fmt.Errorf("lookup user %q: %w", req.GetUsername(), err)
		}
		if err := s.Permissions.AllowStartLogin(ctx, user); err != nil {
			return fmt.Errorf("authorize user %q to login: %w", user.GetUsername(), err)
		}
		session, err := s.Permissions.LoginSessionPrototype(ctx, user)
		if err != nil {
			return fmt.Errorf("generate session prototype for user %q: %w", user.GetUsername(), err)
		}
		// XXX: We need to taint this session to only allow it to be used to call Finish().
		// The current state is completely insecure; no authentication is required to act as
		// any user.  Obviously, this can't remain :)
		if err := store.AddSession(ctx, tx, session); err != nil {
			return fmt.Errorf("store session: %w", err)
		}
		reply.CredentialRequestOptions.Challenge = session.GetId()

		creds, err := store.GetUserCredentials(ctx, tx, user)
		if err != nil {
			return fmt.Errorf("lookup user credentials: %w", err)
		}
		for _, c := range creds {
			reply.CredentialRequestOptions.AllowedCredentials = append(reply.CredentialRequestOptions.AllowedCredentials, &webauthnpb.PublicKeyCredentialDescriptor{
				Id: c.CredentialId,
				Transports: []webauthnpb.PublicKeyCredentialDescriptor_AuthenticatorTransport{
					webauthnpb.PublicKeyCredentialDescriptor_BLE,
					webauthnpb.PublicKeyCredentialDescriptor_INTERNAL,
					webauthnpb.PublicKeyCredentialDescriptor_NFC,
					webauthnpb.PublicKeyCredentialDescriptor_USB,
				},
				Type: "public-key",
			})
		}
		return nil
	}); err != nil {
		return reply, store.AsGRPCError(fmt.Errorf("generate credential request: %w", err))
	}
	return reply, nil
}

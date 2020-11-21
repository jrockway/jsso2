package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/cookies"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/webauthn"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	DB          *store.Connection
	Permissions *internalauth.Permissions
	Webauthn    *webauthn.Config
	Cookies     *cookies.Config
}

func (s *Service) Start(ctx context.Context, req *jssopb.StartLoginRequest) (*jssopb.StartLoginReply, error) {
	emptyReply := &jssopb.StartLoginReply{}

	l := ctxzap.Extract(ctx)
	user := &types.User{
		Username: req.GetUsername(),
	}

	// Check that it's possible to log in as this user.
	//
	// This is an information leak that we probably want to kill off at some point, but it's
	// difficult.  We can easily say "oh yeah that user definitely exists sure uh huh" but then
	// we have to synthesize credentials, because it will be pretty obvious that valid users
	// have credentials enrolled but invalid users don't.  So for now, we just leak the
	// information in the interest of providing more useful error messages ("that's not your
	// username" and "you forgot to enroll an authenticator").
	if err := s.DB.DoTx(ctx, l, true, func(tx *sqlx.Tx) error {
		return store.LookupUser(ctx, tx, user)
	}); err != nil {
		return emptyReply, store.AsGRPCError(fmt.Errorf("validate username: %w", err))
	}

	// See if they're allowed to log in.
	if err := s.Permissions.AllowStartLogin(ctx, user); err != nil {
		return emptyReply, fmt.Errorf("authorize user %q to login: %w", user.GetUsername(), err)
	}

	// Fetch the credentials that the user has enrolled.  We have to send these back to the
	// browser so it knows what security key (etc.) to try to use.  It is not possible to just
	// send back an empty list of credentials.
	var creds []*types.Credential
	if err := s.DB.DoTx(ctx, l, true, func(tx *sqlx.Tx) error {
		var err error
		creds, err = store.GetUserCredentials(ctx, tx, user)
		if err != nil {
			return fmt.Errorf("lookup user credentials: %w", err)
		}
		return nil
	}); err != nil {
		return emptyReply, store.AsGRPCError(fmt.Errorf("lookup existing credentials: %w", err))
	}

	// Create a session to keep state between Start and Finish.  LoginSessionPrototype
	// taints this session so that it can only be used to call Finish.
	session, err := s.Permissions.LoginSessionPrototype(ctx, user)
	if err != nil {
		return emptyReply, fmt.Errorf("generate session prototype for user %q: %w", user.GetUsername(), err)
	}

	// Fill out the reply with details the browser needs.
	reply, err := s.Webauthn.BeginLogin(session, creds)
	if err != nil {
		if errors.Is(err, webauthn.ErrNoCredentials) {
			return emptyReply, status.Error(codes.FailedPrecondition, fmt.Sprintf("begin login: %s", err.Error()))
		}
		return emptyReply, fmt.Errorf("begin login: %w", err)
	}

	// Send the session ID back as a token that can be used to call Finish.
	reply.Token = sessions.ToBase64(session)

	// We store the session last, so that any errors before this point don't write unneeded
	// sessions to the database.
	if err := s.DB.DoTx(ctx, l, false, func(tx *sqlx.Tx) error {
		return store.UpdateSession(ctx, tx, session)
	}); err != nil {
		return emptyReply, store.AsGRPCError(fmt.Errorf("store session: %w", err))
	}
	return reply, nil
}

func (s *Service) Finish(ctx context.Context, req *jssopb.FinishLoginRequest) (*jssopb.FinishLoginReply, error) {
	reply := &jssopb.FinishLoginReply{}
	l := ctxzap.Extract(ctx)

	session := sessions.MustFromContext(ctx)
	id := session.GetId()
	user := session.GetUser()

	var creds []*types.Credential
	if err := s.DB.DoTx(ctx, l, true, func(tx *sqlx.Tx) error {
		var err error
		creds, err = store.GetUserCredentials(ctx, tx, user)
		if err != nil {
			return fmt.Errorf("lookup user credentials: %w", err)
		}
		return nil
	}); err != nil {
		return reply, store.AsGRPCError(fmt.Errorf("lookup existing credentials: %w", err))
	}

	if err := s.finishLoginAndCheckCounter(ctx, l, session, creds, req); err != nil {
		if revokeErr := revokeSession(ctx, l, s.DB, id); revokeErr != nil {
			l.Warn("failed to revoke session after failed login", zap.Error(err))
			err = fmt.Errorf("%w (additionally: %v)", err, revokeErr)
		}
		return reply, fmt.Errorf("finish login and update counters: %w", err)
	}
	if err := untaintSession(ctx, l, s.DB, id); err != nil {
		return reply, err
	}
	reply.RedirectUrl = s.Cookies.LinkToSetCookie(s.Cookies.SessionToSetCookieToken(session))
	return reply, nil
}

func (s *Service) finishLoginAndCheckCounter(ctx context.Context, l *zap.Logger, session *types.Session, creds []*types.Credential, req *jssopb.FinishLoginRequest) error {
	usedCred, err := s.Webauthn.FinishLogin(session, creds, req)
	if err != nil {
		return fmt.Errorf("finish login: %w", err)
	}
	if err := s.DB.DoTx(ctx, l, false, func(tx *sqlx.Tx) error {
		return store.CheckAndUpdateSignCount(ctx, tx, usedCred)
	}); err != nil {
		return fmt.Errorf("check and update counter: %w", err)
	}
	return nil
}

func untaintSession(ctx context.Context, l *zap.Logger, db *store.Connection, id []byte) error {
	if err := db.DoTx(ctx, l, false, func(tx *sqlx.Tx) error {
		// Refresh the session in a transaction, since we will be editing it.
		session, err := store.LookupSession(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("refresh session: %w", err)
		}
		var newTaints []string
		for _, t := range session.GetTaints() {
			if t != sessions.TaintStartLogin {
				newTaints = append(newTaints, t)
			}
		}
		session.Taints = newTaints
		// TODO(jrockway): Add an "upgraded at" timestamp in the metadata.
		if err := store.UpdateSession(ctx, tx, session); err != nil {
			return fmt.Errorf("store untainted session: %w", err)
		}
		return nil
	}); err != nil {
		return store.AsGRPCError(fmt.Errorf("upgrade session: %w", err))
	}
	return nil
}

func revokeSession(ctx context.Context, l *zap.Logger, db *store.Connection, id []byte) error {
	// There is some question as to whether or not we want to revoke an untainted session here.
	if err := db.DoTx(ctx, l, false, func(tx *sqlx.Tx) error {
		session, err := store.LookupSession(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("refresh session: %w", err)
		}
		// TODO(jrockway): Add a revocation reason into the metadata.
		session.ExpiresAt = timestamppb.Now()
		if err := store.UpdateSession(ctx, tx, session); err != nil {
			return fmt.Errorf("store expired session: %w", err)
		}
		return nil
	}); err != nil {
		return store.AsGRPCError(fmt.Errorf("expire session: %w", err))
	}
	return nil
}

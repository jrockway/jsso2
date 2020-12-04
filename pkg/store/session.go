package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type rawSession struct {
	ID        []byte    `db:"id"`
	UserID    int64     `db:"user_id"`
	Username  string    `db:"username"`
	Metadata  []byte    `db:"metadata"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
	Taints    []byte    `db:"taints"`
}

// UpdateSession writes a session to the database.
func UpdateSession(ctx context.Context, db sqlx.ExtContext, s *types.Session) error {
	if s == nil {
		return &ErrEmpty{Field: "session"}
	}
	if s.GetUser() == nil {
		return &ErrEmpty{Field: "session.user"}
	}
	if s.GetUser().GetId() < 1 {
		return &ErrEmpty{Field: "session.user.id"}
	}
	if sessions.IsZero(s.GetId()) {
		return &ErrEmpty{Field: "session.id"}
	}
	obj, err := fromSession(s)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}
	if _, err := sqlx.NamedExecContext(ctx, db, `insert into session
                  ( id,  user_id,  metadata,  taints,  created_at,  expires_at)
            values(:id, :user_id, :metadata, :taints, :created_at, :expires_at)
            on conflict on constraint session_pkey
            do update set metadata=:metadata, taints=:taints, expires_at=:expires_at
`, obj); err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}

func fromSession(s *types.Session) (*rawSession, error) {
	result := &rawSession{}

	metadataJSON, err := protojson.Marshal(s.GetMetadata())
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	sort.Strings(s.Taints)
	taintsJSON, err := json.Marshal(s.Taints)
	if err != nil {
		return nil, fmt.Errorf("marshal taints: %w", err)
	}
	if bytes.Equal(taintsJSON, []byte("null")) {
		taintsJSON = []byte("[]")
	}

	result.ID = s.GetId()
	result.UserID = s.GetUser().GetId()
	result.Username = s.GetUser().GetUsername()
	result.CreatedAt = s.GetCreatedAt().AsTime()
	result.ExpiresAt = s.GetExpiresAt().AsTime()
	result.Metadata = metadataJSON
	result.Taints = taintsJSON

	return result, nil
}

func (raw *rawSession) toSession() (*types.Session, error) {
	result := &types.Session{User: &types.User{}, Metadata: &types.SessionMetadata{}}
	result.Id = raw.ID
	result.User.Id = raw.UserID
	result.User.Username = raw.Username
	if err := protojson.Unmarshal(raw.Metadata, result.Metadata); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}
	var taints []string
	if len(raw.Taints) > 0 {
		if err := json.Unmarshal(raw.Taints, &taints); err != nil {
			return nil, fmt.Errorf("unmarshal taints: %w", err)
		}
	}
	result.CreatedAt = timestamppb.New(raw.CreatedAt)
	result.ExpiresAt = timestamppb.New(raw.ExpiresAt)
	result.Taints = taints
	return result, nil
}

func getSession(ctx context.Context, db sqlx.ExtContext, id []byte) (*types.Session, error) {
	if len(id) != 64 {
		return nil, fmt.Errorf("session id %s: %w", id, ErrSessionIDInvalid)
	}
	raw := &rawSession{}
	row := db.QueryRowxContext(ctx, `select
            s.id AS id, s.metadata AS metadata, s.taints AS taints, s.created_at AS created_at, s.expires_at AS expires_at,
            u.id AS user_id, u.username as username
            from session s left join "user" u on u.id=s.user_id where s.id=$1`, id)
	if err := row.StructScan(raw); err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	session, err := raw.toSession()
	if err != nil {
		return nil, fmt.Errorf("convert to *types.Session: %w", err)
	}
	return session, nil
}

// LookupSession will return the session object for a provided session ID, if the session is still valid.
func LookupSession(ctx context.Context, db sqlx.ExtContext, id []byte) (*types.Session, error) {
	session, err := getSession(ctx, db, id)
	if err != nil {
		return nil, fmt.Errorf("read session: %w", err)
	}
	if session.GetExpiresAt().AsTime().Before(time.Now()) {
		return nil, ErrSessionExpired
	}
	if session.GetCreatedAt().AsTime().After(time.Now()) {
		return nil, ErrSessionNotYetCreated
	}
	return session, nil
}

// AuthenticateUser checks the database for a valid session in the provided sessions.  The provided
// sessions need only contain a session ID.  Each lookup is done in a separate transaction.
func (c *Connection) AuthenticateUser(ctx context.Context, l *zap.Logger, ss []*types.Session, unusedHeaders []*sessions.UnusedHeader, unusedCookies []*sessions.UnusedCookie) (*types.Session, []error) {
	var errs []error

	// Collect errors about unused authentication material.
	for _, u := range unusedHeaders {
		if u.Err != nil && !errors.Is(u.Err, sessions.ErrUnknownAuthType) {
			errs = append(errs, fmt.Errorf("spurious unparseable authorization header %q: %v", u.Value, u.Err))
		}
	}
	for _, u := range unusedCookies {
		if u.Err != nil {
			errs = append(errs, fmt.Errorf("spurious unparseable session cookie %q: %v", u.Cookie.String(), u.Err))
		}
	}
	if len(ss) == 0 {
		errs = append(errs, errors.New("no sessions provided"))
		return nil, errs
	}

	// Check for all sessions for validity.
	for i, s := range ss {
		if err := c.DoTx(ctx, l, true, func(tx *sqlx.Tx) error {
			var err error
			session, err := LookupSession(ctx, tx, s.GetId())
			if err != nil {
				return fmt.Errorf("lookup session: %w", err)
			}
			ss[i] = session
			return nil
		}); err != nil {
			ss[i] = nil
			errs = append(errs, fmt.Errorf("validate session %d/%d: %v", i+1, len(ss), err))
		}
	}
	// Look for at least one valid session.
	for _, s := range ss {
		if s != nil {
			return s, nil
		}
	}
	return nil, errs
}

// RevokeSession will revoke the provided session.
func RevokeSession(ctx context.Context, tx *sqlx.Tx, id []byte, reason string) error {
	session, err := getSession(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("refresh session: %w", err)
	}
	if time.Until(session.ExpiresAt.AsTime()) < 0 {
		// Already expired.
		return nil
	}

	if session.GetMetadata().GetRevocationReason() == "" {
		session.GetMetadata().RevocationReason = reason
	}
	session.ExpiresAt = timestamppb.Now()
	if err := UpdateSession(ctx, tx, session); err != nil {
		return fmt.Errorf("store expired session: %w", err)
	}
	return nil
}

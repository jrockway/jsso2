package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AddSession writes a session to the database.
func AddSession(ctx context.Context, db sqlx.ExtContext, s *types.Session) error {
	if s == nil {
		return &ErrEmpty{Field: "session"}
	}
	if s.GetUser() == nil {
		return &ErrEmpty{Field: "session.user"}
	}
	if s.GetUser().GetId() < 1 {
		return &ErrEmpty{Field: "session.user.id"}
	}
	if sessions.IsZero(s) {
		return &ErrEmpty{Field: "session.id"}
	}
	metadataJSON, err := protojson.Marshal(s.GetMetadata())
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	obj := map[string]interface{}{
		"id":         s.GetId(),
		"user_id":    s.GetUser().GetId(),
		"metadata":   metadataJSON,
		"created_at": s.GetCreatedAt().AsTime(), // TODO: write a drive.Valuer generator for protos.
		"expires_at": s.GetExpiresAt().AsTime(),
	}
	if _, err := sqlx.NamedExecContext(ctx, db, `insert into session (id, user_id, metadata, created_at, expires_at) values(:id, :user_id, :metadata, :created_at, :expires_at)`, obj); err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}

type rawSession struct {
	ID       []byte
	UserID   int64
	Username string
	Metadata []byte
	// UserCreatedAt    sql.NullTime
	// UserDisabledAt   sql.NullTime
	SessionCreatedAt sql.NullTime
	SessionExpiresAt sql.NullTime
}

func (raw *rawSession) toSession() (*types.Session, error) {
	result := &types.Session{User: &types.User{}, Metadata: &types.SessionMetadata{}}
	result.Id = raw.ID
	result.User.Id = raw.UserID
	result.User.Username = raw.Username
	if err := protojson.Unmarshal(raw.Metadata, result.Metadata); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}
	if t := raw.SessionCreatedAt; t.Valid {
		result.CreatedAt = timestamppb.New(t.Time)
	}
	if t := raw.SessionExpiresAt; t.Valid {
		result.ExpiresAt = timestamppb.New(t.Time)
	}
	return result, nil
}

// LookupSession will return the session object for a provided session ID, if the session is still valid.
func LookupSession(ctx context.Context, db sqlx.ExtContext, id []byte) (*types.Session, error) {
	if len(id) != 64 {
		return nil, fmt.Errorf("session id %s: %w", id, ErrSessionIDInvalid)
	}
	raw := &rawSession{}
	row := db.QueryRowxContext(ctx, `select
            s.id AS id, s.metadata AS metadata, s.created_at AS sessioncreatedat, s.expires_at AS sessionexpiresat,
            u.id AS userid, u.username as username
            from session s left join "user" u on u.id=s.user_id where s.id=$1`, id)
	if err := row.StructScan(raw); err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	session, err := raw.toSession()
	if err != nil {
		return nil, fmt.Errorf("convert to *types.Session: %w", err)
	}
	if session.GetExpiresAt().AsTime().Before(time.Now()) {
		return nil, ErrSessionExpired
	}
	if session.GetCreatedAt().AsTime().After(time.Now()) {
		return nil, ErrSessionNotYetCreated
	}
	return session, nil
}

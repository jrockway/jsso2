package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type rawCredential struct {
	ID                 int64        `db:"id"`
	UserID             int64        `db:"user_id"`
	Username           string       `db:"username"`
	CredentialID       []byte       `db:"credential_id"`
	PublicKey          []byte       `db:"public_key"`
	Name               string       `db:"name"`
	CreatedAt          time.Time    `db:"created_at"`
	DeletedAt          sql.NullTime `db:"deleted_at"`
	CreatedBySessionID []byte       `db:"created_by_session_id"`
}

func (raw *rawCredential) toCredential() *types.Credential {
	u := &types.User{
		Id:       raw.UserID,
		Username: raw.Username,
	}
	c := &types.Credential{User: u}
	c.Id = raw.ID
	c.Name = raw.Name
	c.CredentialId = raw.CredentialID
	c.PublicKey = raw.PublicKey
	c.CreatedAt = timestamppb.New(raw.CreatedAt)
	if t := raw.DeletedAt; t.Valid {
		c.DeletedAt = timestamppb.New(raw.DeletedAt.Time)
	}
	c.CreatedBySessionId = raw.CreatedBySessionID
	return c
}

// AddCredential adds a credential to the database.  The credential object must refer to a valid
// user and session.
func AddCredential(ctx context.Context, db sqlx.ExtContext, c *types.Credential) error {
	if c == nil {
		return &ErrEmpty{Field: "credential"}
	}
	if len(c.GetCredentialId()) == 0 {
		return &ErrEmpty{Field: "credential.credential_id"}
	}
	if len(c.GetPublicKey()) == 0 {
		return &ErrEmpty{Field: "credential.public_key"}
	}
	if c.GetUser() == nil {
		return &ErrEmpty{Field: "credential.user"}
	}
	if c.GetUser().GetId() < 1 {
		return &ErrEmpty{Field: "credential.user.id"}
	}
	if sessions.IsZero(c.GetCreatedBySessionId()) {
		return &ErrEmpty{Field: "credential.created_by_session_id"}
	}
	if c.GetId() != 0 {
		return fmt.Errorf("editing an existing credential is not supported: %w", ErrUnimplemented)
	}
	obj := &rawCredential{
		CredentialID:       c.GetCredentialId(),
		PublicKey:          c.GetPublicKey(),
		UserID:             c.GetUser().GetId(),
		Name:               c.GetName(),
		CreatedAt:          c.GetCreatedAt().AsTime(),
		CreatedBySessionID: c.GetCreatedBySessionId(),
	}
	rows, err := sqlx.NamedQueryContext(ctx, db, `insert into credential (user_id, credential_id, public_key, name, created_at, created_by_session_id) values(:user_id, :credential_id, :public_key, :name, :created_at, :created_by_session_id) returning (id)`, obj)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	defer rows.Close()
	if ok := rows.Next(); !ok {
		return errors.New("insert: no id returned")
	}
	if err := rows.Scan(&c.Id); err != nil {
		return fmt.Errorf("insert: scan id: %w", err)
	}
	return nil
}

// GetUserCredentials returns a list of all currently-valid credentials associated with the provided
// user.
func GetUserCredentials(ctx context.Context, db sqlx.ExtContext, u *types.User) ([]*types.Credential, error) {
	if u == nil {
		return nil, &ErrEmpty{Field: "user"}
	}
	if u.GetId() < 1 {
		return nil, &ErrEmpty{Field: "user.id"}
	}
	var raw []*rawCredential
	if err := sqlx.SelectContext(ctx, db, &raw, `select
            c.id AS id, c.credential_id AS credential_id, c.public_key AS public_key, c.name AS name, c.created_at as created_at, u.id as user_id, u.username as username
            from credential c left join "user" u on u.id=c.user_id
            where deleted_at is null and c.user_id=$1`, u.GetId()); err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	result := make([]*types.Credential, len(raw))
	for i, r := range raw {
		result[i] = r.toCredential()
	}
	return result, nil
}

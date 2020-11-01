package sessions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const sessionSize = 64

var encoder = base64.URLEncoding.WithPadding(base64.NoPadding)

var ErrSessionMissing = errors.New("no session id")

// GenerateID generates a valid session ID.
func GenerateID() ([]byte, error) {
	buf := make([]byte, sessionSize)
	if n, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("read entropy into session ID: %w", err)
	} else if got, want := n, sessionSize; got != want {
		return nil, fmt.Errorf("did not produce the correct amount of session entropy; read %d bytes, want %d bytes", got, want)
	}
	return buf, nil
}

// IsZero returns true if the session ID is all zeros (or is the wrong length).
func IsZero(s *types.Session) bool {
	if len(s.GetId()) != sessionSize {
		return true
	}
	for _, b := range s.Id {
		if b != 0 {
			return false
		}
	}
	return true
}

// FromBase64 extracts a session from a base64-encoded session ID.
func FromBase64(in string) (*types.Session, error) {
	id, err := encoder.DecodeString(in)
	if err != nil {
		return nil, fmt.Errorf("session from base64: %w", err)
	}
	if got, want := len(id), sessionSize; got != want {
		return nil, fmt.Errorf("session size: got %d bytes, want %d bytes", got, want)
	}
	return &types.Session{Id: id}, nil
}

// ToBase64 converts a session to a base64-encoded session ID.
func ToBase64(s *types.Session) string {
	if s == nil {
		s = &types.Session{}
	}
	if len(s.Id) == 0 {
		s.Id = make([]byte, sessionSize)
	}
	return encoder.EncodeToString(s.Id)
}

// FromHeaderString extracts a session from an HTTP header.
func FromHeaderString(header string) (*types.Session, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("header %q did not contain a type and a token; got %d parts, want 2 parts", header, len(parts))
	}
	typ, tok := parts[0], parts[1]
	switch typ {
	case "SessionID":
		return FromBase64(tok)
	default:
		return nil, fmt.Errorf("unknown token type %q", typ)
	}
}

// ToHeaderString formats a session as an Authorization header.
func ToHeaderString(s *types.Session) string {
	return fmt.Sprintf("SessionID %s", ToBase64(s))
}

// FromMetadata extracts a session from gRPC metadata.
func FromMetadata(md metadata.MD) (*types.Session, error) {
	auths := md.Get("Authorization")
	if len(auths) == 0 {
		return nil, fmt.Errorf("no authorization header in metadata: %w", ErrSessionMissing)
	} else if len(auths) > 1 {
		// This will probably be too restrictive in general.
		return nil, errors.New("multiple authorization headers provided")
	}
	return FromHeaderString(auths[0])
}

// ToMetadata adds a session ID to gRPC metadata.
func ToMetadata(dst metadata.MD, s *types.Session) {
	dst.Append("Authorization", ToHeaderString(s))
}

// Root returns a session for the root user.
func Root() *types.Session {
	return &types.Session{
		Id:        make([]byte, sessionSize),
		CreatedAt: timestamppb.Now(),
		ExpiresAt: timestamppb.New(time.Unix(1<<63-1, 0)),
		User: &types.User{
			Id:       -1, // It pains me to make root not 0, but 0 means other things.
			Username: "root",
		},
	}
}

type sessionKey struct{}

var sessionContextKey = sessionKey{}

// NewContext adds the session to the provided context.
func NewContext(ctx context.Context, s *types.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, s)
}

// FromContext gets the session in the context.
func FromContext(ctx context.Context) (*types.Session, bool) {
	val, ok := ctx.Value(sessionContextKey).(*types.Session)
	return val, ok
}

// MustFromContext gets the session in the context, or panics.
func MustFromContext(ctx context.Context) *types.Session {
	if val, ok := FromContext(ctx); ok {
		return val
	}
	panic("no session in context")
}

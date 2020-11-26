// Package tokens generates typed authenticated tokens that allow untrusted third parties to safely carry
// state between applictions.
package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/jrockway/jsso2/pkg/types"
	"github.com/o1egl/paseto/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrInvalidKey = errors.New("key is not 32 bytes")
	ErrEmptyToken = errors.New("provided token is empty")
	ErrTooNew     = errors.New("secure message is too new")
	ErrTooOld     = errors.New("secure message is too old")
)

type GeneratorConfig struct {
	Key []byte // A key with which to sign and encrypt tokens.  Must be exactly 32 bytes.
}

func (c *GeneratorConfig) SetKey(key []byte) error {
	if n := len(key); n != 32 {
		return fmt.Errorf("invalid key length; got %d bytes, want 32 bytes", n)
	}
	c.Key = key
	for _, c := range c.Key {
		if c != 0 {
			return nil
		}
	}
	return errors.New("key is entirely null bytes; probably a configuration problem")
}

// New generates a token from the provided protocol message, encrypting and signing it with the
// provided 32-byte symmetric key.
func New(msg proto.Message, key []byte) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKey
	}
	any, err := anypb.New(msg)
	if err != nil {
		return "", fmt.Errorf("marshal message to Any: %w", err)
	}
	wrapper := &types.SecureToken{
		Message:  any,
		IssuedAt: timestamppb.Now(),
	}
	payload, err := proto.Marshal(wrapper)
	if err != nil {
		return "", fmt.Errorf("marshal SecureToken: %w", err)
	}
	token, err := paseto.Encrypt(key, payload, "")
	if err != nil {
		return "", fmt.Errorf("encrypt payload: %w", err)
	}
	return token, nil
}

// Decrypt returns the decrypted SecureToken.  It's only intended to be used from debugging
// utilities.
func Decrypt(token string, key []byte) (*types.SecureToken, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}
	var payload []byte
	var footer string
	if err := paseto.Decrypt(token, key, &payload, &footer); err != nil {
		return nil, fmt.Errorf("decrypt token: %w", err)
	}
	wrapper := &types.SecureToken{}
	if err := proto.Unmarshal(payload, wrapper); err != nil {
		return nil, fmt.Errorf("unmarshal SecureToken: %w", err)
	}
	return wrapper, nil
}

// VerifyAndUnmarshal unmarshals a token created by NewToken into the provided protocol message.  An
// error is returned if the token is too new, too old, cryptographically invalid, or if the type of
// the destination message and contained message do not match.
func VerifyAndUnmarshal(dst proto.Message, token string, maxAge time.Duration, key []byte) error {
	if token == "" {
		return ErrEmptyToken
	}
	wrapper, err := Decrypt(token, key)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}
	age := time.Since(wrapper.GetIssuedAt().AsTime())
	if age < 0 {
		return fmt.Errorf("%w (issued_at is %s in the future)", ErrTooNew, age.String())
	}
	if age > maxAge {
		return fmt.Errorf("%w (message is %s old)", ErrTooOld, age.String())
	}
	if err := wrapper.GetMessage().UnmarshalTo(dst); err != nil {
		return fmt.Errorf("unmarshal contained message: %w", err)
	}
	return nil
}

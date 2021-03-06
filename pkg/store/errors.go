package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNothingToUpdate      = errors.New("nothing to update")
	ErrSessionExpired       = errors.New("session expired")
	ErrSessionNotYetCreated = errors.New("session not yet created?")
	ErrSessionIDInvalid     = errors.New("session id is not valid")
	ErrSignCountDecreased   = errors.New("authenticator's signature counter is not higher than the stored signature counter; possible cloned authenticator")
)

type ErrEmpty struct {
	Field string
}

func (e *ErrEmpty) Error() string {
	return fmt.Sprintf("required field %q missing", e.Field)
}

func IsErrEmpty(err error) bool {
	target := &ErrEmpty{}
	return errors.As(err, &target)
}

// AsGRPCError converts a store error to one with a gRPC status code.  Is is valid to call with a
// nil error.
func AsGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNothingToUpdate) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, err.Error())
	}
	if IsErrEmpty(err) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	if isRetryable(err) {
		// From codes: "Use Unavailable if the client can retry just the failing call."
		return status.Error(codes.Unavailable, err.Error())
	}
	// SQLSTATE 23XXX is a referential integrity violation; duplicate unique index, null where
	// the schema dictates non-null, etc.
	if strings.Contains(err.Error(), "(SQLSTATE 23") {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	return status.Error(codes.Unknown, err.Error())
}

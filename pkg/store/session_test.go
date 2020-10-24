package store

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAddSession_Validation(t *testing.T) {
	testData := []struct {
		name    string
		session *types.Session
		wantMsg string
	}{
		{
			name:    "nil session",
			wantMsg: `required field "session" missing`,
		},
		{
			name:    "empty session",
			session: &types.Session{},
			wantMsg: `required field "session.user" missing`,
		},
		{
			name:    "empty user",
			session: &types.Session{User: &types.User{}},
			wantMsg: `required field "session.user.id" missing`,
		},
		{
			name:    "invalid user",
			session: &types.Session{User: &types.User{Id: 0}},
			wantMsg: `required field "session.user.id" missing`,
		},
		{
			name:    "invalid session id",
			session: &types.Session{User: &types.User{Id: 1}},
			wantMsg: `required field "session.id" missing`,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			re, err := regexp.Compile(test.wantMsg)
			if err != nil {
				t.Fatalf("parse regexp /%v/: %v", test.wantMsg, err)
			}
			err = AddSession(context.Background(), EmptyDB{}, test.session)
			if err == nil {
				t.Fatalf("expecting error /%v/", test.wantMsg)
			}
			if !re.MatchString(err.Error()) {
				t.Errorf("error does not match:\n  got:  %v \n want: /%v/", err, test.wantMsg)
			}
		})
	}

}

func TestSessions(t *testing.T) {
	jtesting.Run(t, "addsession", jtesting.R{Logger: true, Database: true}, func(t *testing.T, e *jtesting.E) {
		c := MustGetTestDB(t, e)
		id, err := sessions.GenerateID()
		if err != nil {
			t.Fatal(err)
		}
		session := &types.Session{
			Id: id,
			User: &types.User{
				Id:       0,
				Username: "foo",
			},
			Metadata: &types.SessionMetadata{
				UserAgent: "the/tests",
			},
			CreatedAt: timestamppb.New(time.Now().Add(-time.Minute).Round(time.Millisecond)),
			ExpiresAt: timestamppb.New(time.Now().Add(time.Minute).Round(time.Millisecond)),
		}
		if err := UpdateUser(e.Context, c.db, session.User); err != nil {
			t.Fatal(err)
		}
		if err := UpdateUser(e.Context, c.db, &types.User{Username: "bar"}); err != nil {
			t.Fatal(err)
		}
		if err := AddSession(e.Context, c.db, session); err != nil {
			t.Fatal(err)
		}
		got, err := LookupSession(e.Context, c.db, id)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(got, session, protocmp.Transform()); diff != "" {
			t.Errorf("lookup session:\n%s", diff)
		}
		if _, err := LookupSession(e.Context, c.db, make([]byte, 63)); !errors.Is(err, ErrSessionIDInvalid) {
			t.Errorf("expected invalid session; got %v", err)
		}
		if _, err := LookupSession(e.Context, c.db, make([]byte, 64)); !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("expected no session rows; got %v", err)
		}
		expired := &types.Session{
			Id: id[:],
			User: &types.User{
				Id: 1,
			},
			CreatedAt: timestamppb.New(time.Now().Add(-2 * time.Minute).Round(time.Millisecond)),
			ExpiresAt: timestamppb.New(time.Now().Add(-time.Minute).Round(time.Millisecond)),
		}
		expired.Id[0]++
		if err := AddSession(e.Context, c.db, expired); err != nil {
			t.Fatal(err)
		}
		if got, err := LookupSession(e.Context, c.db, id); !errors.Is(err, ErrSessionExpired) {
			t.Errorf("expected expired session; got %v\n  session: %v", err, got)
		}
	})
}

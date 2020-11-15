package store

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCredentials(t *testing.T) {
	jtesting.Run(t, "credentials", jtesting.R{Logger: true, Database: true}, func(t *testing.T, e *jtesting.E) {
		c := MustGetTestDB(t, e)
		user := &types.User{Username: "foo"}
		if err := UpdateUser(e.Context, c.db, user); err != nil {
			t.Fatal(err)
		}
		id, err := sessions.GenerateID()
		if err != nil {
			t.Fatal(err)
		}
		now := timestamppb.New(time.Now().Round(time.Millisecond))
		session := &types.Session{
			Id:        id,
			User:      user,
			CreatedAt: now,
		}
		if err := AddSession(e.Context, c.db, session); err != nil {
			t.Fatal(err)
		}

		baseCredential := &types.Credential{
			Id:                 0,
			User:               session.GetUser(),
			CreatedBySessionId: session.GetId(),
			CreatedAt:          now,
			CredentialId:       []byte("AAAAAAAAAAAAAAAA"),
			PublicKey:          []byte("public key of some sort"),
		}
		for i := byte(0); i < 2; i++ {
			credential := proto.Clone(baseCredential).(*types.Credential)
			credential.Id = 0
			for j := range credential.CredentialId {
				credential.CredentialId[j] += i
			}
			credential.Name = string(credential.CredentialId[0])
			if err := AddCredential(e.Context, c.db, credential); err != nil {
				t.Fatalf("add credential %d: %v", i, err)
			}
			if credential.GetId() == 0 {
				t.Error("unexpected session id: 0")
			}
		}

		testData := []struct {
			user    *types.User
			want    []*types.Credential
			wantErr bool
		}{
			{
				user:    nil,
				wantErr: true,
			},
			{
				user:    &types.User{},
				wantErr: true,
			},
			{
				user: user,
				want: []*types.Credential{
					{
						Id:           1,
						Name:         "A",
						CredentialId: []byte("AAAAAAAAAAAAAAAA"),
						CreatedAt:    now,
						User:         user,
						PublicKey:    baseCredential.PublicKey,
					},
					{
						Id:           2,
						Name:         "B",
						CredentialId: []byte("BBBBBBBBBBBBBBBB"),
						CreatedAt:    now,
						User:         user,
						PublicKey:    baseCredential.PublicKey,
					},
				},
			},
			{
				user: &types.User{Id: 12345},
			},
		}
		for i, test := range testData {
			got, err := GetUserCredentials(e.Context, c.db, test.user)
			if err != nil && !test.wantErr {
				t.Errorf("test %d: %v", i, err)
			}
			if err == nil && test.wantErr {
				t.Errorf("test %d: expected error", i)
			}
			sort.Slice(got, func(i, j int) bool {
				return got[i].Id < got[j].Id
			})
			if diff := cmp.Diff(got, test.want, protocmp.Transform(), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("test %d: returned credentials:\n%s", i, diff)
			}
		}
	})
}

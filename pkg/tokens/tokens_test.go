package tokens

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRoundTrip(t *testing.T) {
	defaultKey := []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	goodMsg := &types.SetCookieRequest{
		SessionId:        []byte("foo"),
		RedirectUrl:      "http://example.com/",
		SessionExpiresAt: timestamppb.New(time.Now().Add(10 * time.Minute)),
	}
	testData := []struct {
		name             string
		encKey           []byte
		decKey           []byte
		input            proto.Message
		wantMarshalErr   string
		maxAge           time.Duration
		unmarshalInto    proto.Message
		wantUnmarshalErr string
		want             proto.Message
	}{
		{
			name:          "everything works",
			input:         goodMsg,
			maxAge:        time.Minute,
			unmarshalInto: &types.SetCookieRequest{},
			want:          goodMsg,
		},
		{
			name:             "mismatched types",
			input:            goodMsg,
			maxAge:           time.Minute,
			unmarshalInto:    &types.Session{},
			wantUnmarshalErr: "mismatched message type",
		},
		{
			name:           "bad encryption key",
			input:          goodMsg,
			encKey:         []byte("foo"),
			wantMarshalErr: ErrInvalidKey.Error(),
		},
		{
			name:             "bad decryption key",
			input:            goodMsg,
			decKey:           []byte("bar"),
			maxAge:           time.Minute,
			unmarshalInto:    &types.SetCookieRequest{},
			wantUnmarshalErr: ErrInvalidKey.Error(),
		},
		{
			name:             "wrong decryption key",
			input:            goodMsg,
			maxAge:           time.Minute,
			decKey:           []byte("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
			unmarshalInto:    &types.SetCookieRequest{},
			wantUnmarshalErr: "invalid token authentication",
		},
		{
			name:             "expired",
			input:            goodMsg,
			maxAge:           -time.Minute,
			unmarshalInto:    &types.SetCookieRequest{},
			wantUnmarshalErr: ErrTooOld.Error(),
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			if len(test.encKey) == 0 {
				test.encKey = defaultKey
			}
			if len(test.decKey) == 0 {
				test.decKey = defaultKey
			}
			token, err := New(test.input, test.encKey)
			if err != nil && test.wantMarshalErr == "" {
				t.Fatalf("marshal: %v", err)
			} else if err != nil && !strings.Contains(err.Error(), test.wantMarshalErr) {
				t.Fatalf("marshal: unexpected error:\n  got: %v\n want: %v", err, test.wantMarshalErr)
			} else if err == nil && test.wantMarshalErr != "" {
				t.Fatal("marshal: expected error")
			}
			if test.wantMarshalErr != "" {
				return
			}

			err = VerifyAndUnmarshal(test.unmarshalInto, token, test.maxAge, test.decKey)
			if err != nil && test.wantUnmarshalErr == "" {
				t.Fatalf("unmarhsal: %v", err)
			} else if err != nil && !strings.Contains(err.Error(), test.wantUnmarshalErr) {
				t.Fatalf("unmarhsal: unexpected error:\n  got: %v\n want: %v", err, test.wantUnmarshalErr)
			} else if err == nil && test.wantMarshalErr != "" {
				t.Fatal("unmarshal: expected error")
			}
			if test.wantUnmarshalErr != "" {
				return
			}

			if diff := cmp.Diff(test.unmarshalInto, test.want, protocmp.Transform()); diff != "" {
				t.Errorf("contained message:\n%s", diff)
			}
		})
	}
}

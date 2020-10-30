package sessions

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGenerate(t *testing.T) {
	for i := 0; i < 10000; i++ {
		id, err := GenerateID()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := len(id), sessionSize; got != want {
			t.Errorf("id length:\n  got: %v\n want: %v", got, want)
		}

		uniqueBytes := map[byte]int{}
		for _, b := range id {
			uniqueBytes[b]++
		}
		if got, want := len(uniqueBytes), 22; got < want {
			t.Fatalf("not enough entropy:\n  got: %v unique symbols\n want: %v", got, want)
		}
	}
}

func TestEncode(t *testing.T) {
	id := make([]byte, 64)
	for i := range id {
		id[i] = 0b11001100
	}
	session := &types.Session{Id: id}
	b64 := "zMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzA"

	if got, want := ToBase64(session), b64; got != want {
		t.Errorf("to base64:\n  got: %v\n want: %v", got, want)
	}

	allZeros := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	if got, want := ToBase64(&types.Session{}), allZeros; got != want {
		t.Errorf("to base64 on nil id:\n  got: %v\n want: %v", got, want)
	}

	if got, want := ToBase64(nil), allZeros; got != want {
		t.Errorf("to base64 on nil session:\n  got: %v\n want: %v", got, want)
	}

	if got, want := ToHeaderString(session), "SessionID "+b64; got != want {
		t.Errorf("to header string:\n  got: %v\n want: %v", got, want)
	}

	gotMD := metadata.MD{}
	wantMD := metadata.MD{"authorization": []string{"SessionID " + b64}}
	ToMetadata(gotMD, session)
	if diff := cmp.Diff(gotMD, wantMD); diff != "" {
		t.Errorf("to metadata:\n%s", diff)
	}
}

func TestDecode(t *testing.T) {
	id := make([]byte, 64)
	for i := range id {
		id[i] = byte(i)
	}
	wantSession := &types.Session{Id: id}
	wantBase64 := "AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0-Pw"

	got, err := FromBase64(wantBase64)
	if err != nil {
		t.Fatalf("from base64: %v", err)
	}
	if diff := cmp.Diff(got, wantSession, protocmp.Transform()); diff != "" {
		t.Errorf("from base64:\n%s", diff)
	}
	if _, err := FromBase64("x"); err == nil {
		t.Error("expected error because base64 is invalid")
	} else if got, want := err.Error(), "session from base64: illegal base64 data at input byte 0"; got != want {
		t.Errorf("invalid base64: error message\n  got: %v\n want: %v", got, want)
	}
	if _, err := FromBase64(wantBase64[:32]); err == nil {
		t.Error("expected error because string is too short")
	} else if got, want := err.Error(), "session size: got 24 bytes, want 64 bytes"; got != want {
		t.Errorf("invalid base64: error message\n  got: %v\n want: %v", got, want)
	}

	got, err = FromHeaderString("SessionID " + wantBase64)
	if err != nil {
		t.Fatalf("from header string: %v", err)
	}
	if diff := cmp.Diff(got, wantSession, protocmp.Transform()); diff != "" {
		t.Errorf("from header string:\n%s", diff)
	}

	got, err = FromHeaderString(wantBase64)
	if err == nil {
		t.Fatal("from header string without type: expected error")
	} else if got, want := err.Error(), `header "AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0-Pw" did not contain a type and a token; got 1 parts, want 2 parts`; got != want {
		t.Errorf("header without type:\n  got: %v\n want: %v", got, want)
	}

	testData := []struct {
		name        string
		md          func(md metadata.MD)
		wantSession *types.Session
		wantErr     bool
	}{
		{
			name: "capital",
			md: func(md metadata.MD) {
				md.Set("Authorization", "SessionID "+wantBase64)
			},
			wantSession: wantSession,
		},
		{
			name: "lowercase",
			md: func(md metadata.MD) {
				md.Set("authorization", "SessionID "+wantBase64)
			},
			wantSession: wantSession,
		},
		{
			name:        "empty md",
			md:          func(md metadata.MD) {},
			wantSession: nil,
			wantErr:     true,
		},
		{
			name: "multiple headers",
			md: func(md metadata.MD) {
				md.Set("authorization", "SessionID "+wantBase64, "Foo bar")
			},
			wantErr: true,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			md := metadata.MD{}
			test.md(md)
			got, err := FromMetadata(md)
			if err != nil && !test.wantErr {
				t.Fatalf("from md: %v", err)
			}
			if err == nil && test.wantErr {
				t.Fatal("from md: expected errror")
			}
			if diff := cmp.Diff(got, test.wantSession, protocmp.Transform()); diff != "" {
				t.Errorf("from md:\n%s", diff)
			}
		})
	}
}

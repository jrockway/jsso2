package webauthn

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/webauthnpb"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
)

func mustMarshalJSON(x map[string]interface{}) []byte {
	result, err := json.Marshal(x)
	if err != nil {
		panic(fmt.Sprintf("marshal json: %v", err))
	}
	return result
}

type attestation struct {
	AttStmt  map[string]interface{} `json:"attStmt"`
	AuthData []byte                 `json:"authData"`
	Fmt      string                 `json:"fmt"`
}

const (
	challenge             = "SZiX3YKzElWKpqKUB+Fj6CgosC3cpz2Ft0XrHgX0un9o0kyTvqhusYC5eLI5T8+vMocyIPz1q7kSGfySD0pGKA=="
	validAttestaionObject = "o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YVikSZYN5YgOjGh0NBcPZHZgW4/krrmihjLHmVzzuoMdl2NFAAAAAQAAAAAAAAAAAAAAAAAAAAAAIK7yHaYKYt1DdeYd5Miqu2U4RJWuwEP+9J7CksTOVYmIpQECAyYgASFYIGaTv32l+u74CTqcanmVMOYTh+/ntnwABMSkQTIastYGIlggY+KXrZRJWAiSQK/UF2cCUqwAmOZPifuUOArhhTj7gkw="
)

func TestFinishEnrollment(t *testing.T) {
	// This test is kind of terrible in that we use a bunch of hard-coded data, captured from a
	// live running server and Chrome instance, and then change it around.  It would be better
	// if we could generate authenticator responses and sign them, so that we can test cases
	// where something invalid has a valid signature.  But, we'd have to implement a WebAuthn
	// client in Go, which I didn't want to do today.  (We really need this someday for fuzz
	// testing, however!)
	sid, err := base64.StdEncoding.DecodeString(challenge)
	if err != nil {
		t.Fatal(err)
	}
	session := &types.Session{
		Id: sid,
		User: &types.User{
			Id:       1,
			Username: "test",
		},
	}
	attObj, err := base64.StdEncoding.DecodeString(validAttestaionObject)
	if err != nil {
		t.Fatal(err)
	}
	var rawAttObj attestation
	if err := cbor.Unmarshal(attObj, &rawAttObj); err != nil {
		t.Fatal(err)
	}
	// Mess up the RPID hash.
	rawAttObj.AuthData[0] = 1 // Should be 0x49.
	badAttObj, err := cbor.Marshal(rawAttObj)
	if err != nil {
		t.Fatal(err)
	}

	testData := []struct {
		name    string
		input   *jssopb.FinishEnrollmentRequest
		wantErr bool
	}{
		{
			name:    "empty request",
			input:   &jssopb.FinishEnrollmentRequest{},
			wantErr: true,
		},
		{
			name: "empty client data",
			input: &jssopb.FinishEnrollmentRequest{
				Credential: &webauthnpb.PublicKeyCredential{
					Response: &webauthnpb.AuthenticatorResponse{
						ClientDataJson: mustMarshalJSON(map[string]interface{}{}),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "almost empty client data",
			input: &jssopb.FinishEnrollmentRequest{
				Credential: &webauthnpb.PublicKeyCredential{
					Response: &webauthnpb.AuthenticatorResponse{
						ClientDataJson: mustMarshalJSON(map[string]interface{}{
							"type": "webauthn.create",
						}),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "correct client data, no attestation object",
			input: &jssopb.FinishEnrollmentRequest{
				Credential: &webauthnpb.PublicKeyCredential{
					Response: &webauthnpb.AuthenticatorResponse{
						ClientDataJson: mustMarshalJSON(map[string]interface{}{
							"type":        "webauthn.create",
							"challenge":   sessions.ToBase64(session),
							"crossOrigin": false,
							"origin":      "http://localhost:4000",
						}),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "correct client data, correct attestation object",
			input: &jssopb.FinishEnrollmentRequest{
				Credential: &webauthnpb.PublicKeyCredential{
					Response: &webauthnpb.AuthenticatorResponse{
						ClientDataJson: mustMarshalJSON(map[string]interface{}{
							"type":        "webauthn.create",
							"challenge":   sessions.ToBase64(session),
							"crossOrigin": false,
							"origin":      "http://localhost:4000",
						}),
						Response: &webauthnpb.AuthenticatorResponse_AttestationResponse{
							AttestationResponse: &webauthnpb.AuthenticatorAttestationResponse{
								AttestationObject: attObj,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "incorrect origin, correct attestation object",
			input: &jssopb.FinishEnrollmentRequest{
				Credential: &webauthnpb.PublicKeyCredential{
					Response: &webauthnpb.AuthenticatorResponse{
						ClientDataJson: mustMarshalJSON(map[string]interface{}{
							"type":        "webauthn.create",
							"challenge":   sessions.ToBase64(session),
							"crossOrigin": false,
							"origin":      "https://example.com",
						}),
						Response: &webauthnpb.AuthenticatorResponse_AttestationResponse{
							AttestationResponse: &webauthnpb.AuthenticatorAttestationResponse{
								AttestationObject: attObj,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "incorrect cross-origin, correct attestation object",
			input: &jssopb.FinishEnrollmentRequest{
				Credential: &webauthnpb.PublicKeyCredential{
					Response: &webauthnpb.AuthenticatorResponse{
						ClientDataJson: mustMarshalJSON(map[string]interface{}{
							"type":        "webauthn.create",
							"challenge":   sessions.ToBase64(session),
							"crossOrigin": true,
							"origin":      "http://localhost:4000",
						}),
						Response: &webauthnpb.AuthenticatorResponse_AttestationResponse{
							AttestationResponse: &webauthnpb.AuthenticatorAttestationResponse{
								AttestationObject: attObj,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "correct client data, incorrect attestation object",
			input: &jssopb.FinishEnrollmentRequest{
				Credential: &webauthnpb.PublicKeyCredential{
					Response: &webauthnpb.AuthenticatorResponse{
						ClientDataJson: mustMarshalJSON(map[string]interface{}{
							"type":        "webauthn.create",
							"challenge":   sessions.ToBase64(session),
							"crossOrigin": false,
							"origin":      "http://localhost:4000",
						}),
						Response: &webauthnpb.AuthenticatorResponse_AttestationResponse{
							AttestationResponse: &webauthnpb.AuthenticatorAttestationResponse{
								AttestationObject: badAttObj,
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			_, err := FinishEnrollment("localhost", "http://localhost:4000", session, test.input)
			if err != nil && !test.wantErr {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && test.wantErr {
				t.Errorf("expected error, but got success")
			}
		})
	}
}

func TestBeginEnrollment(t *testing.T) {
	session := &types.Session{
		Id: []byte("session_id"),
		User: &types.User{
			Id:       1,
			Username: "foo",
		},
	}
	creds := []*types.Credential{
		{
			Id:           1,
			User:         session.User,
			CredentialId: []byte("credential_id"),
			PublicKey:    []byte("public_key"),
			Name:         "test",
		},
	}
	want := &webauthnpb.PublicKeyCredentialCreationOptions{
		Challenge: session.Id,
		Timeout:   durationpb.New(60 * time.Second),
		ExcludeCredentials: []*webauthnpb.PublicKeyCredentialDescriptor{
			{
				Id: []byte("credential_id"),
				Transports: []webauthnpb.PublicKeyCredentialDescriptor_AuthenticatorTransport{
					webauthnpb.PublicKeyCredentialDescriptor_BLE,
					webauthnpb.PublicKeyCredentialDescriptor_INTERNAL,
					webauthnpb.PublicKeyCredentialDescriptor_NFC,
					webauthnpb.PublicKeyCredentialDescriptor_USB,
				},
				Type: "public-key",
			},
		},
		Rp: &webauthnpb.PublicKeyCredentialRpEntity{
			Id:   "localhost",
			Name: "localhost",
		},
		User: &webauthnpb.PublicKeyCredentialUserEntity{
			Id:          []byte{0, 0, 0, 0, 0, 0, 0, 1},
			DisplayName: "foo",
			Name:        "foo",
		},
		PubKeyCredParams: []*webauthnpb.PublicKeyCredentialParameters{
			{Alg: -7, Type: "public-key"},
			{Alg: -35, Type: "public-key"},
			{Alg: -36, Type: "public-key"},
			{Alg: -257, Type: "public-key"},
			{Alg: -258, Type: "public-key"},
			{Alg: -259, Type: "public-key"},
			{Alg: -37, Type: "public-key"},
			{Alg: -38, Type: "public-key"},
			{Alg: -39, Type: "public-key"},
			{Alg: -8, Type: "public-key"},
		},
	}
	got, err := BeginEnrollment("localhost", session, creds)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Fatal(diff)
	}

	if _, err := BeginEnrollment("localhost", &types.Session{}, nil); err == nil {
		t.Error("expected error with empty session")
	}
}

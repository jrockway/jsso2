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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
)

var cfg = &Config{
	RelyingPartyID:   "localhost",
	Origin:           "http://localhost:4000",
	RelyingPartyName: "localhost",
}

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
			_, err := cfg.FinishEnrollment(session, test.input)
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
	got, err := cfg.BeginEnrollment(session, creds)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Fatal(diff)
	}

	if _, err := cfg.BeginEnrollment(&types.Session{}, nil); err == nil {
		t.Error("expected error with empty session")
	}
}

func TestBeginLogin(t *testing.T) {
	want := &jssopb.StartLoginReply{
		CredentialRequestOptions: &webauthnpb.PublicKeyCredentialRequestOptions{
			AllowedCredentials: []*webauthnpb.PublicKeyCredentialDescriptor{
				{
					Id: []byte("cred"),
					Transports: []webauthnpb.PublicKeyCredentialDescriptor_AuthenticatorTransport{
						webauthnpb.PublicKeyCredentialDescriptor_BLE,
						webauthnpb.PublicKeyCredentialDescriptor_INTERNAL,
						webauthnpb.PublicKeyCredentialDescriptor_NFC,
						webauthnpb.PublicKeyCredentialDescriptor_USB,
					},
					Type: "public-key",
				},
			},
			Timeout:   durationpb.New(60 * time.Second),
			Challenge: []byte("session"),
		},
	}
	got, err := cfg.BeginLogin(&types.Session{Id: []byte("session")}, []*types.Credential{{Id: 123, CredentialId: []byte("cred"), PublicKey: []byte("key")}})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Error(diff)
	}
}

func atob(in string) []byte {
	result, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return result
}

func TestFinishLogin(t *testing.T) {
	good := &jssopb.FinishLoginRequest{
		Credential: &webauthnpb.PublicKeyCredential{
			Id:   "xZr5o2CqIFi-CLVtmfxXEK4f26DC2mGErYeaVTffxhc",
			Type: "public-key",
			Response: &webauthnpb.AuthenticatorResponse{
				ClientDataJson: []byte(`{"type":"webauthn.get","challenge":"aWarCSn8VMx_zDKVwJSKaLJSSg6hAV8V2uSnBdTWcR5X3Hqugzn5ljcchtdS6-OU4OxKyxcgPQbNkoSCdE-vBQ","origin":"http://localhost:4000","crossOrigin":false,"other_keys_can_be_added_here":"do not compare clientDataJSON against a template. See https://goo.gl/yabPex"}`),
				Response: &webauthnpb.AuthenticatorResponse_AssertionResponse{
					AssertionResponse: &webauthnpb.AuthenticatorAssertionResponse{
						AuthenticatorData: atob("SZYN5YgOjGh0NBcPZHZgW4/krrmihjLHmVzzuoMdl2MBAAAAAw=="),
						Signature:         atob("MEYCIQDdgGvKRKxoL1UbDMbaXddzLhUdTey0Gz22WflN0gforQIhAKW6dgRujd3C6l/+p/wWXeaq0X8KL9lTW6rGdW9+i2kg"),
						UserHandle:        nil,
					},
				},
			},
		},
	}
	bad := proto.Clone(good).(*jssopb.FinishLoginRequest)
	bad.Credential.Response.ClientDataJson[2] = 'T'

	testData := []struct {
		name string
		req  *jssopb.FinishLoginRequest
		ok   bool
	}{
		{
			name: "nil request",
			req:  nil,
			ok:   false,
		},
		{
			name: "attestation request",
			req: &jssopb.FinishLoginRequest{
				Credential: &webauthnpb.PublicKeyCredential{
					Response: &webauthnpb.AuthenticatorResponse{
						Response: &webauthnpb.AuthenticatorResponse_AttestationResponse{},
					},
				},
			},
			ok: false,
		},
		{
			name: "valid request",
			req:  good,
			ok:   true,
		},
		{
			name: "signature verification failure",
			req:  bad,
			ok:   false,
		},
	}

	session := &types.Session{
		Id:   atob("aWarCSn8VMx/zDKVwJSKaLJSSg6hAV8V2uSnBdTWcR5X3Hqugzn5ljcchtdS6+OU4OxKyxcgPQbNkoSCdE+vBQ=="),
		User: sessions.Anonymous().User,
	}
	creds := []*types.Credential{
		{
			Id:           1,
			CredentialId: atob("xZr5o2CqIFi+CLVtmfxXEK4f26DC2mGErYeaVTffxhc="),
			PublicKey:    atob("pQECAyYgASFYIHKfpxXg/JFHfUyG3zHVgDg91YGp1XVv4SP6IjIwyTZrIlggc2Z8QuuTYcoRj8GYNFLP+pTAI+a2dcx9vvgm1J83OyY="),
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			if _, err := cfg.FinishLogin(session, creds, test.req); (err != nil) == test.ok {
				want := "<error>"
				if test.ok {
					want = "<no error>"
				}
				t.Errorf("result:\n  got: %v\n want: %v", err, want)
			}
		})
	}
}

package webauthn

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/protocol/webauthncose"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/webauthnpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

var optsPrototype = &webauthnpb.PublicKeyCredentialCreationOptions{
	Attestation: webauthnpb.PublicKeyCredentialCreationOptions_NONE,
	Timeout:     durationpb.New(60 * time.Second),
	PubKeyCredParams: func() []*webauthnpb.PublicKeyCredentialParameters {
		var result []*webauthnpb.PublicKeyCredentialParameters
		algs := []webauthncose.COSEAlgorithmIdentifier{
			webauthncose.AlgES256,
			webauthncose.AlgES384,
			webauthncose.AlgES512,
			webauthncose.AlgRS256,
			webauthncose.AlgRS384,
			webauthncose.AlgRS512,
			webauthncose.AlgPS256,
			webauthncose.AlgPS384,
			webauthncose.AlgPS512,
			webauthncose.AlgEdDSA,
		}
		for _, alg := range algs {
			result = append(result, &webauthnpb.PublicKeyCredentialParameters{
				Alg:  int32(alg),
				Type: "public-key",
			})
		}
		return result
	}(),
}

// BeginEnrollment starts the enrollment process, returning a PublicKeyCredentialCreationOptions
// for the browser.
func BeginEnrollment(domain string, session *types.Session, existingCreds []*types.Credential) (*webauthnpb.PublicKeyCredentialCreationOptions, error) {
	opts := proto.Clone(optsPrototype).(*webauthnpb.PublicKeyCredentialCreationOptions)
	opts.Challenge = session.GetId()
	opts.Rp = &webauthnpb.PublicKeyCredentialRpEntity{
		Id:   domain,
		Name: domain,
	}
	user := session.GetUser()
	if user.GetId() < 1 {
		return nil, errors.New("invalid user attempting enrollment")
	}
	opts.User = &webauthnpb.PublicKeyCredentialUserEntity{
		Id:          make([]byte, 8),
		DisplayName: user.GetUsername(),
		Name:        user.GetUsername(),
	}
	binary.PutVarint(opts.User.Id, user.GetId())
	for _, c := range existingCreds {
		opts.ExcludeCredentials = append(opts.ExcludeCredentials, &webauthnpb.PublicKeyCredentialDescriptor{
			Id:   c.GetCredentialId(),
			Type: "public-key",
			Transports: []webauthnpb.PublicKeyCredentialDescriptor_AuthenticatorTransport{
				webauthnpb.PublicKeyCredentialDescriptor_BLE,
				webauthnpb.PublicKeyCredentialDescriptor_INTERNAL,
				webauthnpb.PublicKeyCredentialDescriptor_NFC,
				webauthnpb.PublicKeyCredentialDescriptor_USB,
			},
		})
	}
	return opts, nil
}

type ClientData struct {
	Challenge   string `json:"challenge"`
	CrossOrigin bool   `json:"crossOrigin"`
	Origin      string `json:"origin"`
	Type        string `json:"type"`
}

// Verify the authenticator response generated by the client.  Because we use a slightly different
// RPC format than Duo's webauthn library, we do the non-crypto things here, and delegate to that
// library to verify signations.  The steps below are from:
// https://www.w3.org/TR/webauthn/#registering-a-new-credential
func FinishEnrollment(domain, origin string, session *types.Session, req *jssopb.FinishEnrollmentRequest) (*types.Credential, error) {
	// Step 1: Let JSONtext be the result of running UTF-8 decode on clientDataJSON.  (We
	// actually do this on the client side since gRPC lets use send raw bytes on the wire unlike
	// the JSON that the spec assumes you're going to use; see src/lib/webauthn.ts.)
	//
	// Step 2: Let C be the result of running a JSON parser on the clientDataJSON.  (We call it
	// clientData.)
	clientDataJSON := req.GetCredential().GetClientDataJson()
	var clientData ClientData
	if err := json.Unmarshal(clientDataJSON, &clientData); err != nil {
		return nil, fmt.Errorf("unmarshal client data json: %w", err)
	}

	// Step 3: Verify that C.type is webauthn.create.
	if got, want := clientData.Type, "webauthn.create"; got != want {
		return nil, fmt.Errorf("client data type: got %q, want %q", got, want)
	}

	// Step 4: Verify that the value of C.challenge matches the challenge that was sent to the
	// authenticator.
	//
	// The challenge comes back as a base64url string, which purely by coincidence matches what
	// sessions.ToBase64 does.  We send the challenge as raw bytes, we transport it in the
	// frontend as raw bytes, but when the browser adds it to the clientDataJSON, it uses
	// url-safe base64.
	//
	// https://www.w3.org/TR/webauthn/#dom-collectedclientdata-challenge
	if got, want := clientData.Challenge, sessions.ToBase64(session); got != want {
		return nil, fmt.Errorf("provided challenge does not match the current session: got %q, want %q", got, want)
	}

	// Step 5: Verify that the value of C.origin matches the Relying Party's origin.
	if got, want := clientData.Origin, origin; got != want {
		return nil, fmt.Errorf("credential from invalid origin: got %q, want %q", got, want)
	}
	if clientData.CrossOrigin {
		return nil, errors.New("rejecting cross-origin credential")
	}

	// Step 6: Verify that the value of C.TokenBinding.status matches the state of Token Binding
	// for the TLS connection.
	// TODO(jrockway): Do this?

	// Step 7: Compute the hash of clientDataJSON using SHA-256.
	clientDataHash := sha256.New().Sum(clientDataJSON)

	// Step 8: Perform CBOR decoding on attestationObject.
	attestationResponse := protocol.AuthenticatorAttestationResponse{
		AuthenticatorResponse: protocol.AuthenticatorResponse{
			ClientDataJSON: req.Credential.GetClientDataJson(),
		},
		AttestationObject: req.Credential.GetAttestationObject(),
	}
	attestation, err := attestationResponse.Parse()
	if err != nil {
		return nil, fmt.Errorf("parsing attestation response: %w", err)
	}
	// Step 9: Verify the rpIdHash.  (Handled by Verify.)
	// Step 10: Verify that UserPresent is set.  (Handled by Verify.)
	// Step 11: We skip user verification.
	// Step 12: Verify the client extensions.  Skipped.
	// Step 13: Verify the attestation format.  (Handled by Verify.)
	// Step 14: Verify that the attestation statement is correct.
	if err := attestation.AttestationObject.Verify(domain, clientDataHash, false); err != nil {
		protocolErr := new(protocol.Error)
		if errors.As(err, &protocolErr) {
			info := protocolErr.DevInfo
			if strings.HasPrefix(info, "RP Hash mismatch") {
				// This contains garbage that messes up your terminal.
				info = "RP Hash mismatch"
			}
			return nil, fmt.Errorf("validate attestation object: %w (type: %s, dev info: %s)", err, protocolErr.Type, info)
		}
		return nil, fmt.Errorf("validate attestation object: %w", err)
	}

	// Step 15: Obtain trust anchors.  The Duo library claims this is impossible to do, so we
	// don't.
	//
	// Step 16: Assess the attestation trustworthines.  The Duo library also skips this.
	//
	// Step 17: Check that no other user has this credential ID.  Handled by the caller.
	//
	// Step 18: Associate the credential with the user.  Handled by the caller.
	//
	// Step 19: If the attestation statement is not trustworthy, fail.  (Skipped.)
	attData := attestation.AttestationObject.AuthData.AttData
	return &types.Credential{
		CredentialId: attData.CredentialID,
		PublicKey:    attData.CredentialPublicKey,
	}, nil
}

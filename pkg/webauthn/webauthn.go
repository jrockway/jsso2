package webauthn

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/duo-labs/webauthn/protocol/webauthncose"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/webauthnpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

// type webAuthnUser struct {
// 	*types.User
// }

// func (u *webAuthnUser) WebAuthnID() []byte {
// 	result := make([]byte, 4)
// 	binary.PutVarint(result, u.GetId())
// 	return result
// }
// func (u *webAuthnUser) WebAuthnName() string {
// 	return u.GetUsername()
// }
// func (u *webAuthnUser) WebAuthnDisplayName() string {
// 	return u.GetUsername()
// }
// func (u *webAuthnUser) WebAuthnIcon() string                       { return "" }
// func (u *webAuthnUser) WebAuthnCredentials() []webauthn.Credential { return nil }

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
func BeginEnrollment(origin string, session *types.Session) (*webauthnpb.PublicKeyCredentialCreationOptions, error) {
	opts := proto.Clone(optsPrototype).(*webauthnpb.PublicKeyCredentialCreationOptions)
	opts.Challenge = session.GetId()
	opts.Rp = &webauthnpb.PublicKeyCredentialRpEntity{
		Id:   origin,
		Name: origin,
	}
	user := session.GetUser()
	if user.GetId() < 1 {
		return nil, errors.New("invalid user attempting enrollment")
	}
	opts.User = &webauthnpb.PublicKeyCredentialUserEntity{
		Id:          make([]byte, 4),
		DisplayName: user.GetUsername(),
		Name:        user.GetUsername(),
	}
	binary.PutVarint(opts.User.Id, user.GetId())
	return opts, nil
}

package jsso

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/testserver"
	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/grpc/metadata"
)

func TestEnrollmentHappyPath(t *testing.T) {
	s := testserver.New()
	r := &jtesting.R{Logger: true, Database: true}
	s.ToR(r)
	jtesting.Run(t, "enrollment_happy_path", *r, func(t *testing.T, e *jtesting.E) {
		ctx := metadata.AppendToOutgoingContext(e.Context, "authorization", "root root")
		userClient := jssopb.NewUserClient(e.ClientConn)
		enrollmentClient := jssopb.NewEnrollmentClient(e.ClientConn)
		if _, err := userClient.Edit(ctx, &jssopb.EditUserRequest{
			User: &types.User{
				Username: "happy_enrollee",
			},
		}); err != nil {
			t.Fatalf("create user: %v", err)
		}

		link, err := userClient.GenerateEnrollmentLink(ctx, &jssopb.GenerateEnrollmentLinkRequest{
			Target: &types.User{
				Username: "happy_enrollee",
			},
		})
		if err != nil {
			t.Fatalf("create enrollment link: %v", err)
		}

		userCtx := metadata.AppendToOutgoingContext(e.Context, "authorization", "SessionID "+link.Token)
		opts, err := enrollmentClient.Start(userCtx, &jssopb.StartEnrollmentRequest{})
		if err != nil {
			t.Fatalf("start enrollment: %v", err)
		}

		session, err := sessions.FromBase64(link.Token)
		if err != nil {
			t.Fatalf("parse session id in token %s: %v", link.Token, err)
		}

		if diff := cmp.Diff(session.GetId(), opts.GetCredentialCreationOptions().GetChallenge()); diff != "" {
			t.Errorf("unexpected challenge:\n%s", diff)
		}
	})
}

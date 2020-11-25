package jsso

import (
	"fmt"
	"testing"

	"github.com/jrockway/jsso2/pkg/client"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/testserver"
)

func TestAuthorizeHTTP(t *testing.T) {
	s := testserver.New()
	r := &jtesting.R{Logger: true, Database: true}
	s.ToR(r)
	jtesting.Run(t, "grpc_session", *r, func(t *testing.T, e *jtesting.E) {
		db := store.MustGetTestDB(t, e)
		session := store.ValidSession(t, e, db)
		cs := client.FromCC(e.ClientConn)

		testData := []struct {
			name      string
			req       *jssopb.AuthorizeHTTPRequest
			wantAllow bool
		}{
			{
				name:      "no data",
				req:       &jssopb.AuthorizeHTTPRequest{},
				wantAllow: false,
			},
			{
				name: "invalid authorization header",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						"foobar",
					},
				},
				wantAllow: false,
			},
			{
				name: "invalid session cookie",
				req: &jssopb.AuthorizeHTTPRequest{
					Cookies: []string{
						fmt.Sprintf("jsso_session_id=%s", sessions.ToBase64(sessions.Anonymous())),
					},
				},
				wantAllow: false,
			},
			{
				name: "valid authorization header",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						fmt.Sprintf("SessionID %s", sessions.ToBase64(session)),
					},
				},
				wantAllow: true,
			},
		}

		for _, test := range testData {
			t.Run(test.name, func(t *testing.T) {
				reply, err := cs.SessionClient.AuthorizeHTTP(e.Context, test.req)
				if err != nil {
					t.Fatal(err)
				}
				var gotAllow bool
				if _, ok := reply.GetDecision().(*jssopb.AuthorizeHTTPReply_Allow); ok {
					gotAllow = true
				}
				if got, want := gotAllow, test.wantAllow; got != want {
					t.Errorf("allowed?\n  got: %v\n want: %v", got, want)
				}
			})
		}
	})
}

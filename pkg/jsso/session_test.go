package jsso

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/client"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/testserver"
	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func deny() *jssopb.AuthorizeHTTPReply {
	return &jssopb.AuthorizeHTTPReply{
		Decision: &jssopb.AuthorizeHTTPReply_Deny{
			Deny: &jssopb.Deny{
				Destination: nil,
				Reason:      "generic deny",
			},
		},
	}
}

func allow(msg *jssopb.Allow) *jssopb.AuthorizeHTTPReply {
	return &jssopb.AuthorizeHTTPReply{
		Decision: &jssopb.AuthorizeHTTPReply_Allow{
			Allow: msg,
		},
	}
}

// This makes cmp treat all Deny messages as equal to each other.  This is different from
// protocmp.IgnoreMessages, which suppresses them from the output entirely when producing debugging
// information.
var filterDeny = protocmp.FilterMessage(&jssopb.Deny{}, cmp.Comparer(func(x, y protocmp.Message) bool {
	return true
}))

func TestAuthorizeHTTP(t *testing.T) {
	s := testserver.New()
	r := &jtesting.R{Logger: true, Database: true}
	s.ToR(r)
	jtesting.Run(t, "grpc_session", *r, func(t *testing.T, e *jtesting.E) {
		db := store.MustGetTestDB(t, e)
		cs := client.FromCC(e.ClientConn)
		session := store.ValidSession(t, e, db)
		tainted := proto.Clone(session).(*types.Session)
		tainted.Id[0]++
		tainted.Taints = []string{sessions.TaintStartLogin}
		if err := db.DoTx(e.Context, e.Logger, false, func(tx *sqlx.Tx) error {
			return store.UpdateSession(e.Context, tx, tainted)
		}); err != nil {
			t.Fatal(err)
		}

		testData := []struct {
			name      string
			req       *jssopb.AuthorizeHTTPRequest
			wantAllow bool
			wantReply *jssopb.AuthorizeHTTPReply
		}{
			{
				name:      "no data",
				req:       &jssopb.AuthorizeHTTPRequest{},
				wantReply: deny(),
			},
			{
				name: "invalid authorization header",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						"foobar",
					},
				},
				wantReply: deny(),
			},
			{
				name: "invalid cookie",
				req: &jssopb.AuthorizeHTTPRequest{
					Cookies: []string{
						fmt.Sprintf("jsso-session-id=%s", sessions.ToBase64(tainted)),
					},
				},
				wantReply: deny(),
			},
			{
				name: "invalid cookie and authorization header",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						"foobar",
					},
					Cookies: []string{
						fmt.Sprintf("jsso-session-id=%s", sessions.ToBase64(tainted)),
					},
				},
				wantReply: deny(),
			},
			{
				name: "valid authorization header",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						fmt.Sprintf("SessionID %s", sessions.ToBase64(session)),
					},
				},
				wantReply: allow(&jssopb.Allow{
					AddHeaders: nil,
					Username:   session.GetUser().GetUsername(),
				}),
			},
			{
				name: "valid cookie",
				req: &jssopb.AuthorizeHTTPRequest{
					Cookies: []string{
						fmt.Sprintf("jsso-session-id=%s", sessions.ToBase64(session)),
					},
				},
				wantReply: allow(&jssopb.Allow{
					AddHeaders: nil,
					Username:   session.GetUser().GetUsername(),
				}),
			},
			{
				name: "valid cookie, extra authorization",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						"Token foobar",
						"Bearer barbaz",
					},
					Cookies: []string{
						fmt.Sprintf("jsso-session-id=%s", sessions.ToBase64(session)),
					},
				},
				wantReply: allow(&jssopb.Allow{
					AddHeaders: []*types.Header{
						{
							Key:   "authorization",
							Value: "Token foobar",
						},
						{
							Key:   "authorization",
							Value: "Bearer barbaz",
						},
					},
					Username: session.GetUser().GetUsername(),
				}),
			},
			{
				name: "valid cookie, extra cookie, extra authorization",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						"Token foobar",
						"Bearer barbaz",
					},
					Cookies: []string{
						"super-tracking-cookie=here-i-am",
						"normal-tracking-cookie=you-got-me-too",
						fmt.Sprintf("jsso-session-id=%s", sessions.ToBase64(session)),
					},
				},
				wantReply: allow(&jssopb.Allow{
					AddHeaders: []*types.Header{
						{
							Key:   "authorization",
							Value: "Token foobar",
						},
						{
							Key:   "authorization",
							Value: "Bearer barbaz",
						},
						{
							Key:   "cookie",
							Value: "super-tracking-cookie=here-i-am",
						},
						{
							Key:   "cookie",
							Value: "normal-tracking-cookie=you-got-me-too",
						},
					},
					Username: session.GetUser().GetUsername(),
				}),
			},
			{
				name: "anonymous session",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						fmt.Sprintf("SessionID %s", sessions.ToBase64(sessions.Anonymous())),
					},
				},
				wantReply: deny(),
			},
			{
				name: "valid authorization header, invalid cookie",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						fmt.Sprintf("SessionID %s", sessions.ToBase64(session)),
					},
					Cookies: []string{
						fmt.Sprintf("jsso-session-id=%s", sessions.ToBase64(sessions.Anonymous())),
					},
				},
				wantReply: allow(&jssopb.Allow{
					AddHeaders: []*types.Header{},
					Username:   session.GetUser().GetUsername(),
				}),
			},
			{
				name: "valid authorization header, random cookie",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						fmt.Sprintf("SessionID %s", sessions.ToBase64(session)),
					},
					Cookies: []string{"foo=bar"},
				},
				wantReply: allow(&jssopb.Allow{
					AddHeaders: []*types.Header{
						{
							Key:   "cookie",
							Value: "foo=bar",
						},
					},
					Username: session.GetUser().GetUsername(),
				}),
			},
			{
				name: "valid authorization header, invalid authorization header",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						fmt.Sprintf("SessionID %s", sessions.ToBase64(session)),
						"Bearer foobar",
					},
				},
				wantReply: allow(&jssopb.Allow{
					AddHeaders: []*types.Header{
						{
							Key:   "authorization",
							Value: "Bearer foobar",
						},
					},
					Username: session.GetUser().GetUsername(),
				}),
			},
			{
				name: "invalid authorization header, valid cookie",
				req: &jssopb.AuthorizeHTTPRequest{
					AuthorizationHeaders: []string{
						fmt.Sprintf("SessionID %s", sessions.ToBase64(sessions.Anonymous())),
					},
					Cookies: []string{
						fmt.Sprintf("jsso-session-id=%s", sessions.ToBase64(session)),
					},
				},
				wantReply: allow(&jssopb.Allow{
					AddHeaders: []*types.Header{},
					Username:   session.GetUser().GetUsername(),
				}),
			},
		}

		for _, test := range testData {
			t.Run(test.name, func(t *testing.T) {
				reply, err := cs.SessionClient.AuthorizeHTTP(e.Context, test.req)
				if err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(reply, test.wantReply, protocmp.Transform(), filterDeny); diff != "" {
					t.Error(diff)
				}
			})
		}
	})
}

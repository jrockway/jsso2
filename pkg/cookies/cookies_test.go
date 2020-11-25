package cookies

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/web"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRoundTrip(t *testing.T) {
	cfg := &Config{
		Domain: "localhost",
		Name:   "jsso-session-id",
		Linker: &web.Linker{
			BaseURL: &url.URL{Host: "localhost", Scheme: "http"},
		},
	}
	if err := cfg.SetKey([]byte("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")); err != nil {
		t.Fatal(err)
	}

	session := &types.Session{
		Id:        []byte("SSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS"),
		ExpiresAt: timestamppb.New(time.Now().Add(time.Hour)),
	}
	token, err := cfg.NewSetCookieRequest(session, "")
	if err != nil {
		t.Fatal(err)
	}
	mux := new(http.ServeMux)
	mux.HandleFunc("/set-cookie", cfg.HandleSetCookie)
	req := httptest.NewRequest("GET", cfg.LinkToSetCookie(token), nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	res := rec.Result()
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("non-200 response: %s", res.Status)
	}
	cookies := res.Cookies()
	if got, want := len(cookies), 1; got != want {
		t.Errorf("invalid number of cookies:\n  got: %v\n want: %v", got, want)
	}

	// This code pretends like we're a browser; taking Set-Cookie headers, putting them in the
	// cookie jar, and then sending them to the server.
	creq := httptest.NewRequest("GET", "http://localhost", nil)
	creq.AddCookie(cookies[0])
	md := metadata.Pairs("cookie", creq.Header.Get("cookie"))

	got, unusedAuth, unusedCookies := cfg.SessionsFromMetadata(md)
	if got, want := len(unusedAuth), 0; got != want {
		t.Error("unexpectedly found unused authorization headers")
	}
	if got, want := len(unusedCookies), 0; got != want {
		t.Error("unexpectedly found unused authorization cookies")
	}
	if diff := cmp.Diff(got, []*types.Session{session}, sessions.TransformToID()); diff != "" {
		t.Error(diff)
	}
}

func TestSessionsFromRequest(t *testing.T) {
	cfg := &Config{
		Domain: "localhost",
		Name:   "jsso-session-id",
		Linker: &web.Linker{
			BaseURL: &url.URL{Host: "localhost", Scheme: "http"},
		},
	}
	if err := cfg.SetKey([]byte("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")); err != nil {
		t.Fatal(err)
	}
	id, err := sessions.GenerateID()
	if err != nil {
		t.Fatal(err)
	}
	s := &types.Session{Id: id}
	token, err := cfg.NewSetCookieRequest(s, "")
	if err != nil {
		t.Fatal(err)
	}
	cookie, _, err := cfg.cookieFromToken(token)
	if err != nil {
		t.Fatal(err)
	}
	invalidSession := &types.Session{Id: make([]byte, 64)}
	invalidCookie := &http.Cookie{
		Name:  cfg.Name,
		Value: sessions.ToBase64(invalidSession),
	}

	testData := []struct {
		name              string
		req               *http.Request
		wantSessions      []*types.Session
		wantUnusedAuth    []*UnusedHeader
		wantUnusedCookies []*UnusedCookie
	}{
		{
			name: "empty",
			req:  nil,
		},
		{
			name: "valid auth",
			req: &http.Request{
				Header: http.Header{
					"Authorization": []string{sessions.ToHeaderString(s)},
				},
			},
			wantSessions: []*types.Session{s},
		},
		{
			name: "invalid auth",
			req: &http.Request{
				Header: http.Header{
					"Authorization": []string{sessions.ToHeaderString(invalidSession)},
				},
			},
			wantUnusedAuth: []*UnusedHeader{{Value: sessions.ToHeaderString(invalidSession)}},
		},
		{
			name: "valid cookie",
			req: func() *http.Request {
				req := &http.Request{Header: http.Header{}}
				req.AddCookie(cookie)
				return req
			}(),
			wantSessions: []*types.Session{s},
		},
		{
			name: "invalid cookie",
			req: func() *http.Request {
				req := &http.Request{Header: http.Header{}}
				req.AddCookie(invalidCookie)
				return req
			}(),
			wantUnusedCookies: []*UnusedCookie{{Cookie: invalidCookie}},
		},
		{
			name: "valid auth and valid cookie",
			req: func() *http.Request {
				req := &http.Request{
					Header: http.Header{
						"Authorization": []string{sessions.ToHeaderString(s)},
					},
				}
				req.AddCookie(cookie)
				return req
			}(),
			wantSessions: []*types.Session{s, s},
		},
		{
			name: "valid cookie twice",
			req: func() *http.Request {
				req := &http.Request{Header: http.Header{}}
				req.AddCookie(cookie)
				req.AddCookie(cookie)
				return req
			}(),
			wantSessions: []*types.Session{s, s},
		},
		{
			name: "valid cookie and invalid cookie",
			req: func() *http.Request {
				req := &http.Request{Header: http.Header{}}
				req.AddCookie(cookie)
				req.AddCookie(invalidCookie)
				return req
			}(),
			wantSessions:      []*types.Session{s},
			wantUnusedCookies: []*UnusedCookie{{Cookie: invalidCookie}},
		},
		{
			name: "valid cookie and unrelated cookie",
			req: func() *http.Request {
				req := &http.Request{Header: http.Header{}}
				req.AddCookie(cookie)
				req.AddCookie(&http.Cookie{Name: "super-mega-tracker", Value: "hahaha"})
				return req
			}(),
			wantSessions:      []*types.Session{s},
			wantUnusedCookies: []*UnusedCookie{{Cookie: &http.Cookie{Name: "super-mega-tracker", Value: "hahaha"}}},
		},
		{
			name: "unparseable session id",
			req: func() *http.Request {
				req := &http.Request{Header: http.Header{}}
				req.AddCookie(&http.Cookie{Name: cfg.Name, Value: "foo"})
				return req
			}(),
			wantUnusedCookies: []*UnusedCookie{{Cookie: &http.Cookie{Name: cfg.Name, Value: "foo"}}},
		},
	}
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			got, gotAuth, gotCookies := cfg.SessionsFromRequest(test.req)
			if diff := cmp.Diff(got, test.wantSessions, sessions.TransformToID()); diff != "" {
				t.Errorf("sessions:\n%s", diff)
			}
			if diff := cmp.Diff(gotAuth, test.wantUnusedAuth, cmp.Transformer("UnpackUnusedAuth", func(u *UnusedHeader) string { return u.Value }), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("unused auth:\n%s", diff)
			}
			if diff := cmp.Diff(gotCookies, test.wantUnusedCookies, cmp.Transformer("UnpackUnusedCookie", func(u *UnusedCookie) *http.Cookie { return u.Cookie }), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("unused cookies:\n%s", diff)
			}
		})
	}
}

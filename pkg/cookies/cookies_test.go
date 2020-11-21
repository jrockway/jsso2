package cookies

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

	got, err := cfg.SessionFromMetadata(md)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(got.GetId(), session.GetId()); diff != "" {
		t.Error(diff)
	}
}

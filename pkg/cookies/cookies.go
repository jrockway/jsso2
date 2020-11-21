// Package cookies provides utilities for working with session cookies.
package cookies

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/tokens"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/web"
	"google.golang.org/grpc/metadata"
)

// How long we'll accept a set-cookie token after issuance.  We probably only need it for a few
// milliseconds, but the risk of making this longer is minimal, and a long duration helps with clock
// skew issues.
const SetCookieTokenLifetime = time.Minute

// Config configures the session cookies (and set-cookie tokens) we produce.
type Config struct {
	Name   string      // The name of the cookie (like "jsso-session-id").
	Domain string      // The domain that the cookie should be valid on.  ("sso.example.com" might choose "example.com" here.)
	Linker *web.Linker // A Linker for generating links to the set-cookie handler.
	Key    []byte      // A key with which to sign and encrypt set-cookie tokens.  Must be exactly 32 bytes.
}

func (c *Config) SetKey(key []byte) error {
	if n := len(key); n != 32 {
		return fmt.Errorf("invalid key length; got %d bytes, want 32 bytes", n)
	}
	c.Key = key
	for _, c := range c.Key {
		if c != 0 {
			return nil
		}
	}
	return errors.New("key is entirely null bytes; probably a configuration problem")
}

// NewSetCookieRequest returns a paseto token (a "set-cookie token") that, when provided to the
// HandleSetCookie http Handler below, causes a session cookie to be set for the provided session.
// (It also redirects to the redirectURL after setting the cookie.)  We sign+encrypt the token so
// that random people on the Internet can't induce the handler to set an arbitrary cookie.  We do
// not care about replay attacks -- while one of these tokens can't be revoked, the underlying
// session can be, so a compromised token is not particularly harmful.
func (c *Config) NewSetCookieRequest(s *types.Session, redirectURL string) (string, error) {
	req := &types.SetCookieRequest{
		SessionId:        s.GetId(),
		SessionExpiresAt: s.GetExpiresAt(),
		RedirectUrl:      redirectURL,
	}
	token, err := tokens.New(req, c.Key)
	if err != nil {
		return "", fmt.Errorf("generate set-cookie token: %w", err)
	}
	return token, nil
}

// HandleSetCookie responds to an HTTP GET request with a set-cookie token from NewSetCookieRequest
// in the "set" query parameter with a Set-Cookie header and a redirect to the redirect_url inside
// the token.  If the redirect_url is empty, we just respond with "ok".
func (c *Config) HandleSetCookie(w http.ResponseWriter, req *http.Request) {
	cookie, redirect, err := c.cookieFromToken(req.URL.Query().Get("set"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.SetCookie(w, cookie)
	if redirect == "" {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}
	http.Redirect(w, req, redirect, http.StatusTemporaryRedirect)
}

func (c *Config) cookieFromToken(token string) (*http.Cookie, string, error) {
	req := &types.SetCookieRequest{}
	if err := tokens.VerifyAndUnmarshal(req, token, SetCookieTokenLifetime, c.Key); err != nil {
		return nil, "", fmt.Errorf("verify and unmarshal set-cookie token: %w", err)
	}
	cookie := &http.Cookie{
		Domain:   c.Domain,
		Expires:  req.GetSessionExpiresAt().AsTime(),
		HttpOnly: true,
		Name:     c.Name,
		SameSite: http.SameSiteLaxMode,
		Value:    sessions.ToBase64(&types.Session{Id: req.GetSessionId()}),
	}
	return cookie, req.GetRedirectUrl(), nil
}

// SessionFromMetadata looks for a Cookie header in the provided metadata, and then returns the
// session ID in the first cookie that looks like a session.  An error is returned if any invalid
// cookies are found.  sessions.ErrSessionMissing is returned if no cookie is found.
func (c *Config) SessionFromMetadata(md metadata.MD) (*types.Session, error) {
	req := &http.Request{Header: http.Header{"Cookie": md.Get("cookie")}}
	return c.SessionFromRequest(req)
}

func (c *Config) SessionFromRequest(req *http.Request) (*types.Session, error) {
	for _, cookie := range req.Cookies() {
		if cookie.Name == c.Name {
			s, err := sessions.FromBase64(cookie.Value)
			if err != nil {
				return nil, fmt.Errorf("investigating cookie %q: %w", cookie.Name, err)
			}
			return s, nil
		}
	}
	return nil, fmt.Errorf("no matching cookies found: %w", sessions.ErrSessionMissing)
}

// LinkToSetCookie accepts a token from NewSetCookieRequest and returns the URL that will cause that
// token to actually set a cookie.
func (c *Config) LinkToSetCookie(token string) string {
	return c.Linker.SetCookie(token)
}

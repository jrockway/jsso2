// Package cookies provides utilities for working with session cookies.
package cookies

import (
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
	tokens.GeneratorConfig
	Name   string      // The name of the cookie (like "jsso-session-id").
	Domain string      // The domain that the cookie should be valid on.  ("sso.example.com" might choose "example.com" here.)
	Linker *web.Linker // A Linker for generating links to the set-cookie handler.
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

// Cookies returns the cookie objects in the provided string.
func Cookies(header ...string) []*http.Cookie {
	req := &http.Request{Header: http.Header{"Cookie": header}}
	return req.Cookies()
}

// UnusedHeader is a header we couldn't extract a session from, and the reason why.
type UnusedHeader struct {
	Value string
	Err   error
}

// UnusedCookie is a cookie we couldn't extract a session from, and the reason why.  If Err is null,
// then it simply wasn't a cookie we were looking for.
type UnusedCookie struct {
	Cookie *http.Cookie
	Err    error
}

// SessionFromMetadata extracts authorization headers and cookies from the metadata, returning any
// sessions that were found, a list of unused authorization headers, and a list of unused cookies.
// md must not be nil.
func (c *Config) SessionsFromMetadata(md metadata.MD) ([]*types.Session, []*UnusedHeader, []*UnusedCookie) {
	return c.SessionsFromAny(md.Get("authorization"), md.Get("cookie"))
}

// SessionFromRequest extracts authentication material from the provided request, returning any
// sessions that were found, a list of unused authorization headers, and a list of unused cookies.
func (c *Config) SessionsFromRequest(req *http.Request) ([]*types.Session, []*UnusedHeader, []*UnusedCookie) {
	var result []*types.Session
	if req == nil || req.Header == nil {
		return nil, nil, nil
	}
	ss, unusedAuth := c.SessionsFromAuthorization(req.Header.Get("authorization"))
	result = append(result, ss...)
	ss, unusedCookies := c.SessionsFromCookies(req.Cookies())
	result = append(result, ss...)
	return result, unusedAuth, unusedCookies
}

// SessionFromAny takes a slice of Authorization headers and Cookie headers, and returns valid
// sessions, a list of unused Authorization headers, and a list of unused cookies.
func (c *Config) SessionsFromAny(headers, cookies []string) ([]*types.Session, []*UnusedHeader, []*UnusedCookie) {
	var result []*types.Session
	ss, unusedAuth := c.SessionsFromAuthorization(headers...)
	result = append(result, ss...)
	ss, unusedCookies := c.SessionsFromCookies(Cookies(cookies...))
	result = append(result, ss...)
	return result, unusedAuth, unusedCookies
}

// SessionFromAuthorization extracts sessions from the authorization headers, returning
// unused/invalid authorization headers.
func (c *Config) SessionsFromAuthorization(auths ...string) ([]*types.Session, []*UnusedHeader) {
	var result []*types.Session
	var unused []*UnusedHeader
	for _, a := range auths {
		if a == "" {
			// This doesn't count as unused.  It's mostly an artifact of how
			// http.Request.Header.Get("foo") returns "" when there is no Foo header.
			continue
		}
		s, err := sessions.FromHeaderString(a)
		if err != nil {
			unused = append(unused, &UnusedHeader{Value: a, Err: err})
			continue
		}
		if sessions.IsZero(s.GetId()) {
			unused = append(unused, &UnusedHeader{Value: a, Err: sessions.ErrSessionZero})
			continue
		}
		result = append(result, s)
	}
	return result, unused
}

// SessionFromCookies looks through the provided cookies and returns the sessionID from a cookie
// that looks like a session, and the list of cookies with all matching cookies removed.  An error
// is returned if any invalid cookies are found.  sessions.ErrSessionMissing is returned if no
// cookie is found.
func (c *Config) SessionsFromCookies(cookies []*http.Cookie) ([]*types.Session, []*UnusedCookie) {
	var result []*types.Session
	var unused []*UnusedCookie
	for _, cookie := range cookies {
		if cookie.Name == c.Name {
			s, err := sessions.FromBase64(cookie.Value)
			if err != nil {
				unused = append(unused, &UnusedCookie{Cookie: cookie, Err: err})
				continue
			}
			if sessions.IsZero(s.GetId()) {
				unused = append(unused, &UnusedCookie{Cookie: cookie, Err: sessions.ErrSessionZero})
				continue
			}
			result = append(result, s)
		} else {
			unused = append(unused, &UnusedCookie{Cookie: cookie})
		}
	}
	return result, unused
}

// LinkToSetCookie accepts a token from NewSetCookieRequest and returns the URL that will cause that
// token to actually set a cookie.
func (c *Config) LinkToSetCookie(token string) string {
	return c.Linker.SetCookie(token)
}

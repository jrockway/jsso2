package cookies

import (
	"fmt"
	"net/http"

	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/web"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	Name   string
	Domain string
	Linker *web.Linker
}

func (c *Config) SessionToSetCookieToken(s *types.Session) string {
	// TODO: this needs to be signed, so arbitrary websites can't set cookies.
	cookie := http.Cookie{
		Domain:   c.Domain,
		Expires:  s.GetExpiresAt().AsTime(),
		SameSite: http.SameSiteLaxMode,
		Name:     c.Name,
		Value:    sessions.ToBase64(s),
	}
	return cookie.String()
}

func (c *Config) SessionIDFromMetadata(md metadata.MD) (*types.Session, error) {
	req := &http.Request{Header: http.Header{"Cookie": md.Get("cookie")}}
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

func (c *Config) LinkToSetCookie(token string) string {
	return c.Linker.SetCookie(token)
}

func (c *Config) HandleSetCookie(w http.ResponseWriter, req *http.Request) {
	cookie := req.URL.Query().Get("set")
	w.Header().Add("set-cookie", cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

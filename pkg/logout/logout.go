package logout

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/cookies"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/web"
	"go.uber.org/zap"
)

type Handler struct {
	Linker  *web.Linker
	Cookies *cookies.Config
	DB      *store.Connection
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	l := ctxzap.Extract(ctx)

	ss, _, _ := h.Cookies.SessionsFromRequest(req)
	for _, s := range ss {
		if err := h.DB.DoTx(ctx, l, false, func(tx *sqlx.Tx) error {
			if err := store.RevokeSession(ctx, tx, s.GetId(), "logout"); err != nil {
				return fmt.Errorf("revoke session: %w", err)
			}
			return nil
		}); err != nil {
			l.Info("problem revoking session", zap.Error(err))
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.Cookies.Name,
		Value:    "",
		Expires:  time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).In(time.UTC),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	})
	http.Redirect(w, req, h.Linker.LoginPage(), http.StatusTemporaryRedirect)
}

package middleware

import (
	"net/http"

	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
	"github.com/themoroccandemonlist/tmdl-server/internal/handler"
)

func RequireSetupStatus(h *handler.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := h.Config.Store.Get(r, "session")
			statusVal, _ := session.Values["player_status"].(string)
			actualVal := enum.PlayerStatus(statusVal)

			if actualVal != enum.Setup {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

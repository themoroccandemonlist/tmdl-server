package middleware

import (
	"net/http"

	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
	"github.com/themoroccandemonlist/tmdl-server/internal/handler"
)

func RequireActivePlayer(h *handler.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := h.Config.Store.Get(r, "session")
			statusVal, _ := session.Values["player_status"].(string)
			actual := enum.PlayerStatus(statusVal)

			switch actual {
			case enum.Banned:
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			case enum.Setup:
				http.Redirect(w, r, "/profile-setup", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"net/http"

	"github.com/themoroccandemonlist/tmdl-server/internal/handler"
)

func RequireProfile(h *handler.Handler, requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := h.Config.Store.Get(r, "session")

			username, ok := session.Values["player_username"].(string)
			isUsernameValid := ok && username != ""

			regionID, ok := session.Values["player_region_id"].(string)
			isRegionIDValid := ok && regionID != ""

			if !isUsernameValid || !isRegionIDValid {
				http.Redirect(w, r, "/profile/setup", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

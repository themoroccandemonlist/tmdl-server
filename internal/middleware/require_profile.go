package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/themoroccandemonlist/tmdl-server/internal/handler"
)

func RequireProfile(h *handler.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := h.Config.Store.Get(r, "session")

			username, _ := session.Values["player_username"].(string)
			isUsernameValid := username != ""

			regionID, _ := session.Values["player_region_id"].(uuid.UUID)
			isRegionIDValid := regionID != uuid.Nil

			profileIsValid := isUsernameValid && isRegionIDValid
			isProfileSetupPage := r.URL.Path == "/profile-setup"

			if !profileIsValid && !isProfileSetupPage {
				http.Redirect(w, r, "/profile-setup", http.StatusSeeOther)
				return
			}

			if profileIsValid && isProfileSetupPage {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

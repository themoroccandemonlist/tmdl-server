package middleware

import (
	"net/http"
	"slices"

	"github.com/themoroccandemonlist/tmdl-server/internal/handler"
)

func RequireRole(h *handler.Handler, requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := h.Config.Store.Get(r, "session")
			roles, _ := session.Values["user_roles"].([]string)
			if !hasAnyRole(roles, requiredRoles...) {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func hasAnyRole(userRoles []string, requiredRoles ...string) bool {
	if len(userRoles) == 0 {
		return false
	}

	for _, userRole := range userRoles {
		if slices.Contains(requiredRoles, userRole) {
			return true
		}
	}
	return false
}

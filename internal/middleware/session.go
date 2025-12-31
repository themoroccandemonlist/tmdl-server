package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/themoroccandemonlist/tmdl-server/internal/types"
)

func Session(store sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session")

			sessionData := types.SessionData{}

			if playerID, ok := session.Values["player_id"].(uuid.UUID); ok {
				sessionData.PlayerID = playerID
			}

			if playerUsername, ok := session.Values["player_username"].(string); ok {
				sessionData.PlayerUsername = playerUsername
			}

			if playerAvatar, ok := session.Values["player_avatar"].(string); ok {
				sessionData.PlayerAvatar = playerAvatar
			}

			if userRoles, ok := session.Values["user_roles"].([]string); ok {
				sessionData.UserRoles = userRoles
			}

			ctx := context.WithValue(r.Context(), "session", sessionData)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

package types

import (
	"github.com/google/uuid"
)

type SessionData struct {
	PlayerID       uuid.UUID
	PlayerUsername string
	PlayerAvatar   string
	UserRoles      []string
}

func (s SessionData) HasRole(roleName string) bool {
	for _, role := range s.UserRoles {
		if role == roleName {
			return true
		}
	}
	return false
}

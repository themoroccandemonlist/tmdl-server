package model

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	Email     string
	Sub       string
	Roles     []string
	IsDeleted bool
}

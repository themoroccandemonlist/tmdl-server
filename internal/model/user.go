package model

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Sub       string    `json:"sub"`
	Roles     []string  `json:"roles"`
	IsBanned  bool      `json:"isBanned"`
	IsDeleted bool      `json:"isDeleted"`
}

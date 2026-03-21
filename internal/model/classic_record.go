package model

import (
	"time"

	"github.com/google/uuid"
)

type ClassicRecord struct {
	ID               uuid.UUID
	ClassicLevelID   uuid.UUID `validate:"required"`
	Player           *Player
	PlayerID         uuid.UUID `validate:"required"`
	RecordPercentage int       `validate:"required,min=0,max=100"`
	CompletedAt      time.Time
}

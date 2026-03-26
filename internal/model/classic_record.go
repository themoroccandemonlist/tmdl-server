package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
)

type ClassicRecord struct {
	ID             uuid.UUID
	ClassicLevel   *ClassicLevel
	ClassicLevelID uuid.UUID `validate:"required"`
	Player         *Player
	PlayerID       uuid.UUID   `validate:"required"`
	Progress       int         `validate:"required,min=0,max=100"`
	Device         enum.Device `validate:"required,device"`
	Footage        string      `validate:"required,url"`
	RawFootage     string      `validate:"required,url"`
	CompletedAt    time.Time   `validate:"required"`
}

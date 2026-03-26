package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
)

type AdminClassicRecordRow struct {
	ID           uuid.UUID
	Player       PlayerSubRow
	ClassicLevel ClassicLevelSubRow
	Progress     int
	Date         time.Time
	Device       enum.Device
	Footage      *string
	RawFootage   *string
}

type PlayerSubRow struct {
	ID uuid.UUID
	Username string
}

type ClassicLevelSubRow struct {
	ID uuid.UUID
	Name string
}
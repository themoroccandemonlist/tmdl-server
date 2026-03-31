package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
)

type AdminPlayerRow struct {
	ID                     uuid.UUID
	Username               *string
	Region                 RegionSubRow
	ClassicPoints          decimal.Decimal
	PlatformerPoints       decimal.Decimal
	TotalClassicRecords    int
	// TotalPlatformerRecords int
	Status                 enum.PlayerStatus
}

type RegionSubRow struct {
	ID   uuid.UUID
	Name string
}

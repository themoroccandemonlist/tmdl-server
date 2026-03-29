package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AdminRegionRow struct {
	ID               uuid.UUID
	Name             string
	TotalPlayers     int
	ClassicPoints    decimal.Decimal
	PlatformerPoints decimal.Decimal
}

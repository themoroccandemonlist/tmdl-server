package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Region struct {
	ID               uuid.UUID       `json:"id"`
	Name             string          `json:"name"`
	ClassicPoints    decimal.Decimal `json:"classicPoints"`
	PlatformerPoints decimal.Decimal `json:"platformerPoints"`
}

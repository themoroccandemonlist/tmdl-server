package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Region struct {
	ID               uuid.UUID
	Name             string
	ClassicPoints    decimal.Decimal
	PlatformerPoints decimal.Decimal
}

package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Player struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Username         *string
	ClassicPoints    decimal.Decimal
	PlatformerPoints decimal.Decimal
	Avatar           *string
	RegionID         *uuid.UUID
	Discord          *string
	YouTube          *string
	Twitter          *string
	Twitch           *string
	IsFlagged        bool
}

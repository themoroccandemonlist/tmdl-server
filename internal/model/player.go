package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Player struct {
	ID               uuid.UUID       `json:"id"`
	UserID           uuid.UUID       `json:"userId"`
	Username         *string         `json:"username"`
	ClassicPoints    decimal.Decimal `json:"classicPoints"`
	PlatformerPoints decimal.Decimal `json:"platformerPoints"`
	Avatar           *string         `json:"avatar"`
	RegionID         *uuid.UUID      `json:"regionId"`
	Discord          *string         `json:"discord"`
	YouTube          *string         `json:"youtube"`
	Twitter          *string         `json:"twitter"`
	Twitch           *string         `json:"twitch"`
	IsFlagged        bool            `json:"isFlagged"`
}

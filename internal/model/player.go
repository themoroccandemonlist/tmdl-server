package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
)

type Player struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Username         *string
	ClassicPoints    decimal.Decimal
	PlatformerPoints decimal.Decimal
	Avatar           *string
	Cover            *string
	Region           *Region
	RegionID         *uuid.UUID
	Discord          *string
	YouTube          *string
	Twitter          *string
	Twitch           *string
	Gamemode         *enum.Gamemode
	Device           *enum.Device
	Status           enum.PlayerStatus
}

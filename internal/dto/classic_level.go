package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
)

type AdminClassicLevelRow struct {
	ID               uuid.UUID
	Name             string
	Publisher        string
	Difficulty       enum.Difficulty
	Ranking          int
	Points           decimal.Decimal
}
package model

import (
	"math"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/themoroccandemonlist/tmdl-server/internal/enum"
)

type ClassicLevel struct {
	ID               uuid.UUID
	LevelID          string          `validate:"required"`
	Name             string          `validate:"required"`
	Publisher        string          `validate:"required"`
	Difficulty       enum.Difficulty `validate:"required,difficulty"`
	Duration         enum.Duration   `validate:"required,duration"`
	Ranking          int             `validate:"required,min=1"`
	ListPercentage   int             `validate:"required,min=0,max=100"`
	Points           decimal.Decimal
	MinimumPoints    decimal.Decimal
	OldPoints        decimal.Decimal
	OldMinimumPoints decimal.Decimal
	YoutubeLink      string `validate:"required,url"`
	ThumbnailPath    string
}

func (l *ClassicLevel) CalculatePoints() {
	if l.Ranking == 0 || l.Ranking > 150 {
		l.Points = decimal.Zero
		l.MinimumPoints = decimal.Zero
		return
	}

	logPart := math.Log(float64(l.Ranking)) / math.Log(151)
	logDecimal := decimal.NewFromFloat(logPart)

	one := decimal.NewFromInt(1)
	fiveHundred := decimal.NewFromInt(500)
	three := decimal.NewFromInt(3)

	points := fiveHundred.Mul(one.Sub(logDecimal))

	l.Points = points.Round(2)
	l.MinimumPoints = points.Div(three).Round(2)
}

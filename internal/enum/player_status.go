package enum

import "strings"

type PlayerStatus string

const (
	Setup   PlayerStatus = "SETUP"
	Active  PlayerStatus = "ACTIVE"
	Flagged PlayerStatus = "FLAGGED"
	Banned  PlayerStatus = "BANNED"
)

func (ps PlayerStatus) IsValid() bool {
	switch ps {
	case Setup, Active, Flagged, Banned:
		return true
	}
	return false
}

func FormatPlayerStatus(ps PlayerStatus) string {
	return string(ps[0]) + strings.ToLower(string(ps[1:]))
}

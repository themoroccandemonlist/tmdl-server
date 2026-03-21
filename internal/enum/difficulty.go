package enum

import "strings"

type Difficulty string

const (
	EasyDemon    Difficulty = "EASY_DEMON"
	MediumDemon  Difficulty = "MEDIUM_DEMON"
	HardDemon    Difficulty = "HARD_DEMON"
	InsaneDemon  Difficulty = "INSANE_DEMON"
	ExtremeDemon Difficulty = "EXTREME_DEMON"
)

func (d Difficulty) IsValid() bool {
    switch d {
    case EasyDemon, MediumDemon, HardDemon, InsaneDemon, ExtremeDemon:
        return true
    }
    return false
}

func FormatDifficulty(d Difficulty) string {
	parts := strings.Split(string(d), "_")

	for i, p := range parts {
		p = strings.ToLower(p)
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}

	return strings.Join(parts, " ")
}
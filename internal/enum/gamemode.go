package enum

import "strings"

type Gamemode string

const (
	Classic       Gamemode = "CLASSIC"
	Platformer    Gamemode = "PLATFORMER"
	BothGamemodes Gamemode = "BOTH"
)

func (g Gamemode) IsValid() bool {
	switch g {
	case Classic, Platformer, BothGamemodes:
		return true
	}
	return false
}

func FormatGamemode(g Gamemode) string {
	return string(g[0]) + strings.ToLower(string(g[1:]))
}

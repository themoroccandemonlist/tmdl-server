package enum

import "strings"

type Device string

const (
	PC     Device = "PC"
	Mobile Device = "MOBILE"
)

func (d Device) IsValid() bool {
	switch d {
	case PC, Mobile:
		return true
	}
	return false
}

func FormatDevice(d Device) string {
	s := string(d)

	if s == "PC" {
		return "PC"
	}

	return string(s[0]) + strings.ToLower(string(s[1:]))
}

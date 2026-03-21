package enum

import "strings"

type Duration string

const (
	Tiny   Duration = "TINY"
	Short  Duration = "SHORT"
	Medium Duration = "MEDIUM"
	Long   Duration = "LONG"
	XL     Duration = "XL"
)

func (d Duration) IsValid() bool {
	switch d {
	case Tiny, Short, Medium, Long, XL:
		return true
	}
	return false
}

func FormatDuration(d Duration) string {
	s := string(d)

	if s == "XL" {
		return "XL"
	}

	s = strings.ToLower(s)
	return strings.ToUpper(s[:1]) + s[1:]
}

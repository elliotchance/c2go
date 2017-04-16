package noarch

import (
	"strconv"
	"unicode"
)

func Atoi(a string) int {
	// TODO: It looks like atoi allow other non-digit characters. We need to
	// only pull off the digit characters before we can do the conversion.
	s := ""

	for _, c := range a {
		if !unicode.IsDigit(c) {
			break
		}

		s += string(c)
	}

	// TODO: Does it always return 0 on error?
	v, _ := strconv.Atoi(s)

	return v
}

func Strtol(string, string, int) int {
	panic("strtol is not supported")
}

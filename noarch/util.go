package noarch

import "strings"

func NullTerminatedString(s string) string {
	return strings.TrimRight(s, "\x00")
}

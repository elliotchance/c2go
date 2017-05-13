package noarch

import "strings"

func NullTerminatedString(s string) string {
	return strings.TrimRight(s, "\x00")
}

func NullTerminatedByteSlice(s []byte) string {
	return NullTerminatedString(string(s))
}

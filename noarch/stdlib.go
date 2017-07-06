package noarch

import (
	"strconv"
	"unicode"
)

func Atoi(a []byte) int {
	// TODO: It looks like atoi allows other non-digit characters. We need to
	// only pull off the digit characters before we can do the conversion.
	s := ""

	for _, c := range a {
		if !unicode.IsDigit(rune(c)) {
			break
		}

		s += string(c)
	}

	// TODO: Does it always return 0 on error?
	v, _ := strconv.Atoi(s)

	return v
}

// Strtol - see detail on site:
// http://www.cplusplus.com/reference/cstdlib/strtol/
func Strtol(a []byte, b *[]byte, base int) (value int32) {
	var countByteForInt int
	for i := 0; i < len(a); i++ {
		countByteForInt = len(a) - 1 - i
		s, err := strconv.ParseInt(string(a[0:countByteForInt+1]), base, 32)
		value = int32(s)
		if err != nil {
			continue
		}
		break
	}
	if countByteForInt == 0 {
		panic("function Strtol: Cannot found integer")
	}
	*b = a[countByteForInt+1 : len(a)]
	return value
}

// Free doesn't do anything since memory is managed by the Go garbage collector.
// However, I will leave it here as a placeholder for now.
func Free(anything interface{}) {
}

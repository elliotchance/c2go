package noarch

// NullTerminatedByteSlice returns a string that contains all the bytes in the
// provided C string up until the first NULL character.
func NullTerminatedByteSlice(s []byte) string {
	if s == nil {
		return ""
	}

	end := -1
	for i, b := range s {
		if b == 0 {
			end = i
			break
		}
	}

	if end == -1 {
		end = len(s)
	}

	newSlice := make([]byte, end)
	copy(newSlice, s)

	return string(newSlice)
}

func CStringIsNull(s []byte) bool {
	if s == nil || len(s) < 1 {
		return true
	}

	return s[0] == 0
}

package noarch

// NullTerminatedByteSlice returns a new byte slice that has been truncated to a
// NULL byte. If there is no NULL byte then a new copy of the original byte
// slice is returned.
func NullTerminatedByteSlice(s []byte) []byte {
	if s == nil {
		return nil
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

	return newSlice
}

package noarch

import "time"

// TimeT is the representation of "time_t".
type TimeT int32

// NullToTimeT converts a NULL to an array of TimeT.
func NullToTimeT(i int32) []TimeT {
	return []TimeT{}
}

// Time returns the current time.
func Time(tloc []TimeT) TimeT {
	var t = TimeT(int32(time.Now().Unix()))

	if len(tloc) > 0 {
		tloc[0] = t
	}

	return t
}

// IntToTimeT converts an int32 to a TimeT.
func IntToTimeT(t int32) TimeT {
	return TimeT(t)
}

// Ctime converts TimeT to a string.
func Ctime(tloc []TimeT) []byte {
	if len(tloc) > 0 {
		var t = time.Unix(int64(tloc[0]), 0)
		return []byte(t.Format(time.ANSIC) + "\n")
	}

	return nil
}

// TimeTToFloat64 converts TimeT to a float64. It is used by the tests.
func TimeTToFloat64(t TimeT) float64 {
	return float64(t)
}

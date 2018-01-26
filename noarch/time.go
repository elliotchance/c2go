package noarch

import "time"

// TimeT is the representation of "time_t".
// For historical reasons, it is generally implemented as an integral value
// representing the number of seconds elapsed
// since 00:00 hours, Jan 1, 1970 UTC (i.e., a unix timestamp).
// Although libraries may implement this type using alternative time
// representations.
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

// Tm - base struct in "time.h"
// Structure containing a calendar date and time broken down into its
// components
type Tm struct {
	Tm_sec   int
	Tm_min   int
	Tm_hour  int
	Tm_mday  int
	Tm_mon   int
	Tm_year  int
	Tm_wday  int
	Tm_yday  int
	Tm_isdst int
	// tm_gmtoff int32
	// tm_zone   []byte
}

// Localtime - Convert time_t to tm as local time
// Uses the value pointed by timer to fill a tm structure with the values that
// represent the corresponding time, expressed for the local timezone.
func LocalTime(timer []TimeT) (tm []Tm) {
	t := time.Unix(int64(timer[0]), 0)
	tm = make([]Tm, 1)
	tm[0].Tm_sec = t.Second()
	tm[0].Tm_min = t.Minute()
	tm[0].Tm_hour = t.Hour()
	tm[0].Tm_mday = t.Day()
	tm[0].Tm_mon = int(t.Month())
	tm[0].Tm_year = t.Year()
	tm[0].Tm_wday = int(t.Weekday())
	tm[0].Tm_yday = t.YearDay()
	return
}

// Mktime - Convert tm structure to time_t
// Returns the value of type time_t that represents the local time described
// by the tm structure pointed by timeptr (which may be modified).
func Mktime(tm []Tm) TimeT {
	t := time.Date(tm[0].Tm_year, time.Month(tm[0].Tm_mon), tm[0].Tm_mday,
		tm[0].Tm_hour, tm[0].Tm_min, tm[0].Tm_sec, 0, time.Now().Location())
	return TimeT(int32(t.Unix()))
}

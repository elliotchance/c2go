package noarch

import (
	"math"
)

// Signbitf is provided by the cstdlib.
func Signbitf(x float32) int {
	return BoolToInt(math.Signbit(float64(x)))
}

// Signbitd is provided by the cstdlib.
func Signbitd(x float64) int {
	return BoolToInt(math.Signbit(x))
}

// Signbitl is provided by the cstdlib.
func Signbitl(x float64) int {
	return BoolToInt(math.Signbit(x))
}

// IsNaN is provided by the cstdlib.
func IsNaN(x float64) int {
	return BoolToInt(math.IsNaN(x))
}

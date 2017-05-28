package noarch

import (
	"math"
)

func Signbitf(x float32) int {
	return BoolToInt(math.Signbit(float64(x)))
}

func Signbitd(x float64) int {
	return BoolToInt(math.Signbit(x))
}

func Signbitl(x float64) int {
	return BoolToInt(math.Signbit(x))
}

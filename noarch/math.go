package noarch

import (
	"math"
)

func Signbitf(x float32) int32 {
	return BoolToInt(math.Signbit(float64(x)))
}

func Signbitd(x float64) int32 {
	return BoolToInt(math.Signbit(x))
}

func Signbitl(x float64) int32 {
	return BoolToInt(math.Signbit(x))
}

func IsNaN(x float64) int32 {
	return BoolToInt(math.IsNaN(x))
}

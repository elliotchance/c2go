package linux

import (
	"math"

	"github.com/elliotchance/c2go/noarch"
)

func IsNanf(x float32) int32 {
	return noarch.BoolToInt(math.IsNaN(float64(x)))
}

func IsInff(x float32) int32 {
	return noarch.BoolToInt(math.IsInf(float64(x), 0))
}

func IsInf(x float64) int32 {
	return noarch.BoolToInt(math.IsInf(x, 0))
}

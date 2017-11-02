package linux

import (
	"math"

	"github.com/elliotchance/c2go/noarch"
)

// IsNanf handles __isnanf(float)
func IsNanf(x float32) int {
	return noarch.BoolToInt(math.IsNaN(float64(x)))
}

// IsInff handles __isinff(float)
func IsInff(x float32) int {
	return noarch.BoolToInt(math.IsInf(float64(x), 0))
}

// IsInfd handles __inline_isinfd(double)
func IsInfd(x float64) int {
	return noarch.BoolToInt(math.IsInf(float64(x), 0))
}

// IsInf handles __inline_isinfl(long double)
func IsInf(x float64) int {
	return noarch.BoolToInt(math.IsInf(x, 0))
}

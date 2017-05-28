package linux

import (
	"math"

	"github.com/elliotchance/c2go/noarch"
)

func Isnanf(x float32) int {
	return noarch.BoolToInt(math.IsNaN(float64(x)))
}

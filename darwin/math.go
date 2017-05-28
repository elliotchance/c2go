package darwin

import (
	"math"
)

type Float2 struct {
	Sinval float32
	Cosval float32
}
type Double2 struct {
	Sinval float64
	Cosval float64
}

func Fabsf(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func Fabs(x float64) float64 {
	return math.Abs(x)
}

func Fabsl(x float64) float64 {
	return math.Abs(x)
}

func Inff() float32 {
	return float32(math.Inf(0))
}

func Inf() float64 {
	return math.Inf(0)
}

func Infl() float64 {
	return math.Inf(0)
}

func SincosfStret(x float32) Float2 {
	return Float2{0, 0}
}

func SincosStret(x float64) Double2 {
	return Double2{0, 0}
}

func SincospifStret(x float32) Float2 {
	return Float2{0, 0}
}

func SincospiStret(x float64) Double2 {
	return Double2{0, 0}
}

func NaN(s []byte) float64 {
	return math.NaN()
}

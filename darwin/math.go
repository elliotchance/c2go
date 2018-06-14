package darwin

import (
	"math"
)

// Float2 is used by the functions that end with "Stret". It replaces the
// internal macOS type of __float2.
type Float2 struct {
	Sinval float32
	Cosval float32
}

// Double2 is used by the functions that end with "Stret". It replaces the
// internal macOS type of __double2.
type Double2 struct {
	Sinval float64
	Cosval float64
}

// Fabsf handles __builtin_fabsf().
func Fabsf(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

// Fabs handles __builtin_fabs().
func Fabs(x float64) float64 {
	return math.Abs(x)
}

// Fabsl handles __builtin_fabsl().
func Fabsl(x float64) float64 {
	return math.Abs(x)
}

// Inff handles __builtin_inff().
func Inff() float32 {
	return float32(math.Inf(0))
}

// Inf handles __builtin_inf().
func Inf() float64 {
	return math.Inf(0)
}

// Infl handles __builtin_infl().
func Infl() float64 {
	return math.Inf(0)
}

// SincosfStret handles __sincosf_stret().
func SincosfStret(x float32) Float2 {
	return Float2{0, 0}
}

// SincosStret handles __sincos_stret().
func SincosStret(x float64) Double2 {
	return Double2{0, 0}
}

// SincospifStret handles __sincospif_stret().
func SincospifStret(x float32) Float2 {
	return Float2{0, 0}
}

// SincospiStret handles __sincospi_stret().
func SincospiStret(x float64) Double2 {
	return Double2{0, 0}
}

func NaN(s *byte) float64 {
	return math.NaN()
}

package darwin

import (
	"math"

	"github.com/elliotchance/c2go/noarch"
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

func Signbitf(x float32) int {
	return noarch.BoolToInt(math.Signbit(float64(x)))
}

func Signbitd(x float64) int {
	// if x*-0.0 == x {
	// 	fmt.Println("yes")
	// }
	// x = -1.0 * 0.0
	// fmt.Printf("%f %d\n", x, math.Signbit(x))
	return noarch.BoolToInt(math.Signbit(x))
}

func Signbitl(x float64) int {
	return noarch.BoolToInt(math.Signbit(x))
}

func Nanf(s []byte) float32 {
	return math.Float32frombits(0x7F800000)
}

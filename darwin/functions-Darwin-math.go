package main

import (
    "math"
)

func __builtin_fabsf(x float32) float32 {
    return float32(math.Abs(float64(x)))
}

func __builtin_fabs(x float64) float64 {
    return math.Abs(x)
}

func __builtin_fabsl(x float64) float64 {
    return math.Abs(x)
}

func __builtin_inff() float32 {
    return float32(math.Inf(0))
}

func __builtin_inf() float64 {
    return math.Inf(0)
}

func __builtin_infl() float64 {
    return math.Inf(0)
}

func __sincosf_stret(x float32) __float2 {
    return __float2{0, 0}
}

func __sincos_stret(x float64) __double2 {
    return __double2{0, 0}
}

func __sincospif_stret(x float32) __float2 {
    return __float2{0, 0}
}

func __sincospi_stret(x float64) __double2 {
    return __double2{0, 0}
}

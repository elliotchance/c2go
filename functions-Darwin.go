package main

import (
    "math"
    "unicode"
)

// FIXME: These are wrong.
type __mbstate_t int64
type __builtin_va_list int64
type fpos_t int64

type _RuneLocale struct {
    __runetype [256]uint32
}

var _DefaultRuneLocale _RuneLocale = _RuneLocale{
    __runetype: [256]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, },
}

func __maskrune(_c __darwin_ct_rune_t, _f uint32) uint32 {
    return _DefaultRuneLocale.__runetype[_c & 0xff] & _f;
}

// func __istype(_c __darwin_ct_rune_t, _f uint32) uint32 {
//     if isascii(int(_c)) != 0 {
//         return __not_uint32(__not_uint32((_DefaultRuneLocale.__runetype[_c] & _f)))
//     }

//     return __not_uint32(__not_uint32(__maskrune(_c, _f)))
// }

func __isctype(_c __darwin_ct_rune_t, _f uint32) __darwin_ct_rune_t {
    if _c < 0 || _c >= (1 <<8 ) {
        return 0
    }

    return __darwin_ct_rune_t(__not_uint32(__not_uint32((_DefaultRuneLocale.__runetype[_c] & _f))))
}

func __tolower(c __darwin_ct_rune_t) __darwin_ct_rune_t {
    return __darwin_ct_rune_t(unicode.ToLower(rune(c)))
}

func __toupper(c __darwin_ct_rune_t) __darwin_ct_rune_t {
    return __darwin_ct_rune_t(unicode.ToUpper(rune(c)))
}

// math

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

// stdio

// func printf(format string, ...args interface{}) {

// }

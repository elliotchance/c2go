package noarch

import (
	"strconv"
	"unicode"
)

// Abs returns the absolute value of parameter n.
//
// In C++, this function is also overloaded in header <cmath> for floating-point
// types (see cmath abs), in header <complex> for complex numbers (see complex
// abs), and in header <valarray> for valarrays (see valarray abs).
func Abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

// Atoi parses the C-string str interpreting its content as an integral number,
// which is returned as a value of type int.
//
// The function first discards as many whitespace characters (as in isspace) as
// necessary until the first non-whitespace character is found. Then, starting
// from this character, takes an optional initial plus or minus sign followed by
// as many base-10 digits as possible, and interprets them as a numerical value.
//
// The string can contain additional characters after those that form the
// integral number, which are ignored and have no effect on the behavior of this
// function.
//
// If the first sequence of non-whitespace characters in str is not a valid
// integral number, or if no such sequence exists because either str is empty or
// it contains only whitespace characters, no conversion is performed and zero
// is returned.
func Atoi(a []byte) int {
	// TODO: It looks like atoi allows other non-digit characters. We need to
	// only pull off the digit characters before we can do the conversion.
	s := ""

	for _, c := range a {
		if !unicode.IsDigit(rune(c)) {
			break
		}

		s += string(c)
	}

	// TODO: Does it always return 0 on error?
	v, _ := strconv.Atoi(s)

	return v
}

// Strtol parses the C-string str interpreting its content as an integral number
// of the specified base, which is returned as a long int value. If endptr is
// not a null pointer, the function also sets the value of endptr to point to
// the first character after the number.
//
// The function first discards as many whitespace characters as necessary until
// the first non-whitespace character is found. Then, starting from this
// character, takes as many characters as possible that are valid following a
// syntax that depends on the base parameter, and interprets them as a numerical
// value. Finally, a pointer to the first character following the integer
// representation in str is stored in the object pointed by endptr.
//
// If the value of base is zero, the syntax expected is similar to that of
// integer constants, which is formed by a succession of:
//
// - An optional sign character (+ or -)
// - An optional prefix indicating octal or hexadecimal base ("0" or "0x"/"0X"
//   respectively)
//
// A sequence of decimal digits (if no base prefix was specified) or either
// octal or hexadecimal digits if a specific prefix is present
//
// If the base value is between 2 and 36, the format expected for the integral
// number is a succession of any of the valid digits and/or letters needed to
// represent integers of the specified radix (starting from '0' and up to
// 'z'/'Z' for radix 36). The sequence may optionally be preceded by a sign
// (either + or -) and, if base is 16, an optional "0x" or "0X" prefix.
//
// If the first sequence of non-whitespace characters in str is not a valid
// integral number as defined above, or if no such sequence exists because
// either str is empty or it contains only whitespace characters, no conversion
// is performed.
//
// For locales other than the "C" locale, additional subject sequence forms may
// be accepted.
func Strtol(a, b []byte, c int) int32 {
	// TODO: This is a bad implementation
	return 65535
}

// Free doesn't do anything since memory is managed by the Go garbage collector.
// However, I will leave it here as a placeholder for now.
func Free(anything interface{}) {
}

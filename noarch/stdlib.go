package noarch

import (
	"math"
	"regexp"
	"strconv"
	"strings"
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

// Atof parses the C string str, interpreting its content as a floating point
// number and returns its value as a double.
//
// The function first discards as many whitespace characters (as in isspace) as
// necessary until the first non-whitespace character is found. Then, starting
// from this character, takes as many characters as possible that are valid
// following a syntax resembling that of floating point literals (see below),
// and interprets them as a numerical value. The rest of the string after the
// last valid character is ignored and has no effect on the behavior of this
// function.
//
// C90 (C++98): A valid floating point number for atof using the "C" locale is
// formed by an optional sign character (+ or -), followed by a sequence of
// digits, optionally containing a decimal-point character (.), optionally
// followed by an exponent part (an e or E character followed by an optional
// sign and a sequence of digits).
//
// C99/C11 (C++11): A valid floating point number for atof using the "C" locale
// is formed by an optional sign character (+ or -), followed by one of:
//
//   - A sequence of digits, optionally containing a decimal-point character
//     (.), optionally followed by an exponent part (an e or E character
//     followed by an optional sign and a sequence of digits).
//   - A 0x or 0X prefix, then a sequence of hexadecimal digits (as in isxdigit)
//     optionally containing a period which separates the whole and fractional
//     number parts. Optionally followed by a power of 2 exponent (a p or P
//     character followed by an optional sign and a sequence of hexadecimal
//     digits).
//   - INF or INFINITY (ignoring case).
//   - NAN or NANsequence (ignoring case), where sequence is a sequence of
//     characters, where each character is either an alphanumeric character (as
//     in isalnum) or the underscore character (_).
//
// If the first sequence of non-whitespace characters in str does not form a
// valid floating-point number as just defined, or if no such sequence exists
// because either str is empty or contains only whitespace characters, no
// conversion is performed and the function returns 0.0.
func Atof(str []byte) float64 {
	// First start by removing any trailing whitespace.
	s := strings.TrimSpace(CStringToString(str))

	// Before we get into the more complicated parser below lets just try and
	// interpret the number.
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}

	// Now convert to lowercase, this makes the regexp and comparisons easier
	// and doesn't change the value.
	s = strings.ToLower(s)

	// 1. Hexadecimal integer? This must be checked before floating-point
	// because it starts with a 0.
	r := regexp.MustCompile(`^([+-])?0x([0-9a-f]+)(p[-+]?[0-9a-f]+)?`)
	match := r.FindStringSubmatch(s)
	if match != nil {
		n, err := strconv.ParseUint(match[2], 16, 32)
		if err == nil {
			f := float64(n)

			if match[1] == "-" {
				f *= -1
			}

			if match[3] != "" {
				p, err := strconv.Atoi(match[3][1:])
				if err != nil {
					return 0
				}

				f *= math.Pow(2, float64(p))
			}

			return f
		}

		return 0
	}

	// 2. Floating-point number?
	r = regexp.MustCompile(`^[+-]?\d*(\.\d*)?(e[+-]\d+)?`)
	match = r.FindStringSubmatch(s)
	if match != nil {
		f, err := strconv.ParseFloat(match[0], 64)
		if err == nil {
			return f
		}
	}

	// 3. Infinity?
	if s == "infinity" || s == "+infinity" ||
		s == "inf" || s == "+inf" {
		return math.Inf(1)
	}
	if s == "-infinity" || s == "-inf" {
		return math.Inf(-1)
	}

	// 4. Not a number?
	if len(s) > 2 && s[:3] == "nan" {
		return math.NaN()
	}
	if len(s) > 3 && s[1:4] == "nan" {
		return math.NaN()
	}

	return 0
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

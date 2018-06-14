package noarch

import (
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/elliotchance/c2go/util"
	"math/rand"
	"sync"
	"unsafe"
)

// DivT is the representation of "div_t". It is used by div().
type DivT struct {
	Quot int32 // quotient
	Rem  int32 // remainder
}

// LdivT is the representation of "ldiv_t". It is used by ldiv().
type LdivT struct {
	Quot int32 // quotient
	Rem  int32 // remainder
}

// LldivT is the representation of "lldiv_t". It is used by lldiv().
type LldivT struct {
	Quot int64 // quotient
	Rem  int64 // remainder
}

// Abs returns the absolute value of parameter n.
//
// In C++, this function is also overloaded in header <cmath> for floating-point
// types (see cmath abs), in header <complex> for complex numbers (see complex
// abs), and in header <valarray> for valarrays (see valarray abs).
func Abs(n int32) int32 {
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
func Atof(str *byte) float64 {
	f, _ := atof(str)

	return f
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
func Atoi(str *byte) int32 {
	return int32(Atol(str))
}

// Atol parses the C-string str interpreting its content as an integral number,
// which is returned as a value of C type "long int".
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
func Atol(str *byte) int32 {
	return int32(Atoll(str))
}

// Atoll parses the C-string str interpreting its content as an integral number,
// which is returned as a value of C type long long int.
//
// This function operates like atol to interpret the string, but produces
// numbers of type long long int (see atol for details on the interpretation
// process).
func Atoll(str *byte) int64 {
	x, _ := atoll(str, 10)

	return x
}

func atoll(str *byte, radix int32) (int64, int) {
	// First start by removing any trailing whitespace. We need to record how
	// much whitespace is trimmed off for the correct offset later.
	cStr := CStringToString(str)
	beforeLength := len(cStr)
	s := strings.TrimSpace(cStr)
	whitespaceOffset := beforeLength - len(s)

	// We must convert the input to lowercase so satisfy radix > 10.
	if radix > 10 {
		s = strings.ToLower(s)
	}

	// We must stop consuming characters when we get to a character that is
	// invalid for the radix. Build a regex to satisfy this.
	rx := ""
	var i int32
	for ; i < radix; i++ {
		if i < 10 {
			rx += string(48 + i)
		} else {
			rx += string(87 + i)
		}
	}
	r := util.GetRegex(`^([+-]?[` + rx + `]+)`)
	match := r.FindStringSubmatch(s)
	if match == nil {
		return 0, 0
	}

	// We do not care about the error here because it should be impossible.
	v, _ := strconv.ParseInt(match[1], int(radix), 64)

	return v, whitespaceOffset + len(match[1])
}

// Div returns the integral quotient and remainder of the division of numer by
// denom ( numer/denom ) as a structure of type div_t, ldiv_t or lldiv_t, which
// has two members: quot and rem.
func Div(numer, denom int32) DivT {
	return DivT{
		Quot: numer / denom,
		Rem:  numer % denom,
	}
}

// Exit uses os.Exit to stop program execution.
func Exit(exitCode int32) {
	os.Exit(int(exitCode))
}

// Getenv retrieves a C-string containing the value of the environment variable
// whose name is specified as argument. If the requested variable is not part of
// the environment list, the function returns a null pointer.
//
// The pointer returned points to an internal memory block, whose content or
// validity may be altered by further calls to getenv (but not by other library
// functions).
//
// The string pointed by the pointer returned by this function shall not be
// modified by the program. Some systems and library implementations may allow
// to change environmental variables with specific functions (putenv,
// setenv...), but such functionality is non-portable.
func Getenv(name *byte) *byte {
	key := CStringToString(name)

	if env, found := os.LookupEnv(key); found {
		return StringToCString(env)
	}

	return nil
}

// Labs returns the absolute value of parameter n ( /n/ ).
//
// This is the long int version of abs.
func Labs(n int32) int32 {
	if n < 0 {
		return -n
	}

	return n
}

// Ldiv returns the integral quotient and remainder of the division of numer by
// denom ( numer/denom ) as a structure of type ldiv_t, which has two members:
// quot and rem.
func Ldiv(numer, denom int32) LdivT {
	return LdivT{
		Quot: numer / denom,
		Rem:  numer % denom,
	}
}

// Llabs returns the absolute value of parameter n ( /n/ ).
//
// This is the long long int version of abs.
func Llabs(n int64) int64 {
	if n < 0 {
		return -n
	}

	return n
}

// Lldiv returns the integral quotient and remainder of the division of numer by
// denom ( numer/denom ) as a structure of type lldiv_t, which has two members:
// quot and rem.
func Lldiv(numer, denom int64) LldivT {
	return LldivT{
		Quot: numer / denom,
		Rem:  numer % denom,
	}
}

// Rand returns a random number using math/rand.Int().
func Rand() int32 {
	return int32(rand.Int())
}

// Strtod parses the C-string str interpreting its content as a floating point
// number (according to the current locale) and returns its value as a double.
// If endptr is not a null pointer, the function also sets the value of endptr
// to point to the first character after the number.
//
// The function first discards as many whitespace characters (as in isspace) as
// necessary until the first non-whitespace character is found. Then, starting
// from this character, takes as many characters as possible that are valid
// following a syntax resembling that of floating point literals (see below),
// and interprets them as a numerical value. A pointer to the rest of the string
// after the last valid character is stored in the object pointed by endptr.
func Strtod(str *byte, endptr **byte) float64 {
	f, fLen := atof(str)

	// FIXME: This is actually creating new data for the returned pointer,
	// rather than returning the correct reference. This means that applications
	// that modify the returned pointer will not be manipulating the original
	// str.
	if endptr != nil {
		*endptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(str)) + uintptr(fLen)))
	}

	return f
}

// Strtof works the same way as Strtod but returns a float.
func Strtof(str *byte, endptr **byte) float32 {
	return float32(Strtod(str, endptr))
}

// Strtold works the same way as Strtod but returns a long double.
func Strtold(str *byte, endptr **byte) float64 {
	return Strtod(str, endptr)
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
func Strtol(str *byte, endptr **byte, radix int32) int32 {
	return int32(Strtoll(str, endptr, radix))
}

// Strtoll works the same way as Strtol but returns a long long.
func Strtoll(str *byte, endptr **byte, radix int32) int64 {
	x, xLen := atoll(str, radix)

	// FIXME: This is actually creating new data for the returned pointer,
	// rather than returning the correct reference. This means that applications
	// that modify the returned pointer will not be manipulating the original
	// str.
	if endptr != nil {
		*endptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(str)) + uintptr(xLen)))
	}

	return x
}

// Strtoul works the same way as Strtol but returns a long unsigned int.
func Strtoul(str *byte, endptr **byte, radix int32) uint32 {
	return uint32(Strtoll(str, endptr, radix))
}

// Strtoull works the same way as Strtol but returns a long long unsigned int.
func Strtoull(str *byte, endptr **byte, radix int32) uint64 {
	return uint64(Strtoll(str, endptr, radix))
}

var (
	memMgmt map[uint64]interface{}
	memSync sync.Mutex
)

func init() {
	memMgmt = make(map[uint64]interface{})
}

// Malloc returns a pointer to a memory block of the given length.
//
// To prevent the Go garbage collector from collecting this memory,
// we store the whole block in a map.
func Malloc(numBytes int32) unsafe.Pointer {
	memBlock := make([]byte, numBytes)
	addr := uint64(uintptr(unsafe.Pointer(&memBlock[0])))
	memSync.Lock()
	defer memSync.Unlock()
	memMgmt[addr] = memBlock
	return unsafe.Pointer(&memBlock[0])
}

// Free removes the reference to this memory address,
// so that the Go GC can free it.
func Free(anything unsafe.Pointer) {
	addr := uint64(uintptr(anything))
	memSync.Lock()
	defer memSync.Unlock()
	delete(memMgmt, addr)
}

func atof(str *byte) (float64, int32) {
	// First start by removing any trailing whitespace. We have to record how
	// much whitespace is trimmed off to correct for the final length.
	cStr := CStringToString(str)
	beforeLength := len(cStr)
	s := strings.TrimSpace(cStr)

	whitespaceLength := beforeLength - len(s)

	// Now convert to lowercase, this makes the regexp and comparisons easier
	// and doesn't change the value.
	s = strings.ToLower(s)

	// 1. Hexadecimal integer? This must be checked before floating-point
	// because it starts with a 0.
	r := util.GetRegex(`^([+-])?0x([0-9a-f]+)(p[-+]?[0-9a-f]+)?`)
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
					return 0, 0
				}

				f *= math.Pow(2, float64(p))
			}

			return f, int32(whitespaceLength + len(match[0]))
		}

		return 0, 0
	}

	// 2. Floating-point number?
	r = util.GetRegex(`^[+-]?\d*(\.\d*)?(e[+-]?\d+)?`)
	match = r.FindStringSubmatch(s)
	if match != nil {
		f, err := strconv.ParseFloat(match[0], 64)
		if err == nil {
			return f, int32(whitespaceLength + len(match[0]))
		}
	}

	// 3. Infinity?
	if s == "infinity" || s == "+infinity" ||
		s == "inf" || s == "+inf" {
		return math.Inf(1), int32(len(s))
	}
	if s == "-infinity" || s == "-inf" {
		return math.Inf(-1), int32(len(s))
	}

	// 4. Not a number?
	if len(s) > 2 && s[:3] == "nan" {
		return math.NaN(), 3
	}
	if len(s) > 3 && s[1:4] == "nan" {
		return math.NaN(), 4
	}

	return 0, 0
}

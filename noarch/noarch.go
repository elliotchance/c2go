// Package noarch contains low-level functions that apply to multiple platforms.
package noarch

// BoolToInt converts boolean value to an int, which is a common operation in C.
// 0 and 1 represent false and true respectively.
func BoolToInt(x bool) int32 {
	if x {
		return 1
	}

	return 0
}

// NotInt performs a logical not (!) on an integer and returns an integer.
func NotInt(x int) int {
	if x == 0 {
		return 1
	}

	return 0
}

// NotInt32 works the same as NotInt, but on a int32.
func NotInt32(x int32) int32 {
	if x == 0 {
		return 1
	}

	return 0
}

// NotUint16 works the same as NotInt, but on a uint16.
func NotUint16(x uint16) uint16 {
	if x == 0 {
		return 1
	}

	return 0
}

// NotInt8 works the same as NotInt, but on a int8.
func NotInt8(x int8) int8 {
	if x == 0 {
		return 1
	}

	return 0
}

// Ternary simulates the ternary (also known as the conditional operator). Go
// does not have the equivalent of using if statements as expressions or inline
// if statements. This function takes the true and false parts as closures to be
// sure that only the true or false condition is evaulated - to prevent side
// effects.
func Ternary(a bool, b, c func() interface{}) interface{} {
	if a {
		return b()
	}

	return c()
}

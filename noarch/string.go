package noarch

// Strlen returns the length of a string.
//
// The length of a C string is determined by the terminating null-character: A
// C string is as long as the number of characters between the beginning of the
// string and the terminating null character (without including the terminating
// null character itself).
func Strlen(a []byte) int {
	// TODO: The transpiler should have a syntax that means this proxy function
	// does not need to exist.

	return len(NullTerminatedByteSlice(a))
}

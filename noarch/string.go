package noarch

func Strlen(a []byte) int {
	// TODO: The transpiler should have a syntax that means this proxy function
	// does not need to exist.

	return len(NullTerminatedByteSlice(a))
}

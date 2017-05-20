package noarch

func Strlen(a []byte) int {
	// TODO: The transpiler should have a syntax that means this proxy function
	// does not need to exist.

	// TODO: The real length of the string will only be up to the NULL
	// character.
	return len(a)
}

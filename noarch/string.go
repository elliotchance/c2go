package noarch

func Strlen(a string) int {
	// TODO: The transpiler should have a syntax that means this proxy function
	// does not need to exist.
	return len(a)
}

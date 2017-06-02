package util

// PanicIfNil will panic with the message provided if the check is nil. This is
// a convieniance method to avoid many similar if statements.
func PanicIfNil(check interface{}, message string) {
	if check == nil {
		panic(message)
	}
}

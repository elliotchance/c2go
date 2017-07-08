package util

import "fmt"

// PanicIfNil will panic with the message provided if the check is nil. This is
// a convieniance method to avoid many similar if statements.
func PanicIfNil(check interface{}, message string) {
	if check == nil {
		panic(message)
	}
}

// PanicOnError will panic with the message and error if the error is not nil.
// If the error is nil (no error) then nothing happens.
func PanicOnError(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", message, err.Error()))
	}
}

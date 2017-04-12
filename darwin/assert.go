package darwin

import (
	"fmt"
	"os"
)

func BuiltinExpect(a, b int) bool {
	return a != b
}

func AssertRtn(functionName, filePath string, lineNumber int, expression string) bool {
	fmt.Fprintf(os.Stderr, "Assertion failed: (%s), function %s, file %s, line %d.\n", expression, functionName, filePath, lineNumber)
	os.Exit(134)

	return true
}

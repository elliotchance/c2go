package darwin

import (
	"fmt"
	"os"

	"github.com/elliotchance/c2go/noarch"
)

// BuiltinExpect handles __builtin_expect().
func BuiltinExpect(a, b int32) int32 {
	return noarch.BoolToInt(a != b)
}

// AssertRtn handles __assert_rtn().
func AssertRtn(
	functionName, filePath *byte,
	lineNumber int32,
	expression *byte,
) bool {
	fmt.Fprintf(
		os.Stderr,
		"Assertion failed: (%s), function %s, file %s, line %d.\n",
		noarch.CStringToString(expression),
		noarch.CStringToString(functionName),
		noarch.CStringToString(filePath),
		lineNumber,
	)
	os.Exit(134)

	return true
}

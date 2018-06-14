package linux

import (
	"fmt"
	"os"

	"github.com/elliotchance/c2go/noarch"
)

// AssertFail handles __assert_fail().
func AssertFail(
	expression, filePath *byte,
	lineNumber uint32,
	functionName *byte,
) bool {
	fmt.Fprintf(
		os.Stderr,
		"a.out: %s:%d: %s: Assertion `%s' failed.\n",
		noarch.CStringToString(filePath),
		lineNumber,
		noarch.CStringToString(functionName),
		noarch.CStringToString(expression),
	)
	os.Exit(134)

	return true
}

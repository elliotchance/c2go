package linux

import (
	"fmt"
	"os"

	"github.com/elliotchance/c2go/noarch"
)

// AssertFail handles __assert_fail().
func AssertFail(
	expression, filePath []byte,
	lineNumber uint32,
	functionName []byte,
) bool {
	fmt.Fprintf(
		os.Stderr,
		"a.out: %s:%d: %s: Assertion `%s' failed.\n",
		noarch.NullTerminatedByteSlice(filePath),
		lineNumber,
		noarch.NullTerminatedByteSlice(functionName),
		noarch.NullTerminatedByteSlice(expression),
	)
	os.Exit(134)

	return true
}

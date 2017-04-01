package main

import (
    "fmt"
    "os"
)

// FIXME
type _IO_FILE interface{}

func __assert_fail(functionName, filePath string, lineNumber uint32, expression string) bool {
    fmt.Fprintf(os.Stderr, "%s:%d: %s: Assertion `%s' failed.\n", fileName, lineNumber, functionName, expression)
    os.Exit(134)

    return true
}

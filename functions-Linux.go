package main

import (
    "fmt"
    "os"
)

// FIXME
type _IO_FILE interface{}

func __assert_fail(expression, filePath string, lineNumber uint32, functionName string) bool {
    fmt.Fprintf(os.Stderr, "a.out: %s:%d: %s: Assertion `%s' failed.\n", filePath, lineNumber, functionName, expression)
    os.Exit(134)

    return true
}

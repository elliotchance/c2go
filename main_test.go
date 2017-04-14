package main

import (
	"testing"
	"path/filepath"
)

// This test exists for code coverage and does not actually test anything. The
// real tests are performed with run-tests.h.
//
// In the future it would be nice to combine them so that the files only have to
// be compiled once and we don't need the extra bash script.
func TestIntegrationScripts(t *testing.T) {
	files, err := filepath.Glob("tests/*/*.c")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		Start([]string{"", file})
	}
}

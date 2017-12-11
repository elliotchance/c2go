package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func setupTest(args []string) (*bytes.Buffer, func()) {
	buf := &bytes.Buffer{}
	oldStderr := stderr
	oldArgs := os.Args

	stderr = buf
	os.Args = args

	return buf, func() {
		stderr = oldStderr
		os.Args = oldArgs
	}
}

var cliTests = map[string][]string{
	// Test that help is printed if no files are given
	"TranspileNoFilesHelp": {"test", "transpile"},

	// Test that help is printed if help flag is set, even if file is given
	"TranspileHelpFlag": {"test", "transpile", "-h", "foo.c"},

	// Test that help is printed if no files are given
	"AstNoFilesHelp": {"test", "ast"},

	// Test that help is printed if help flag is set, even if file is given
	"AstHelpFlag": {"test", "ast", "-h", "foo.c"},
}

func TestCLI(t *testing.T) {
	for testName, args := range cliTests {
		t.Run(testName, func(t *testing.T) {
			output, teardown := setupTest(args)
			defer teardown()

			runCommand()

			err := cupaloy.SnapshotMulti(testName, output)
			if err != nil {
				t.Fatalf("error: %s", err)
			}
		})
	}
}

package main

import (
	"bytes"
	"github.com/bradleyjkemp/cupaloy"
	"os"
	"strings"
	"testing"
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
	"TranspileNoFilesHelp": []string{"test", "transpile"},

	// Test that help is printed if help flag is set, even if file is given
	"TranspileHelpFlag": []string{"test", "transpile", "-h", "foo.c"},

	// Test that help is printed if no files are given
	"AstNoFilesHelp": []string{"test", "ast"},

	// Test that help is printed if help flag is set, even if file is given
	"AstHelpFlag": []string{"test", "ast", "-h", "foo.c"},
}

func TestCLI(t *testing.T) {
	for testName, args := range cliTests {
		t.Run(testName, func(t *testing.T) {
			output, teardown := setupTest(args)
			defer teardown()

			runCommand()
			outputLines := strings.Split(output.String(), "\n")
			err := cupaloy.SnapshotMulti(testName, outputLines)
			if err != nil {
				t.Fatalf("error: %s", err)
			}
		})
	}
}

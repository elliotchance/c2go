package main

import (
	"bytes"
	"github.com/bradleyjkemp/cupaloy"
	"os"
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

// Test that help is printed if no files are given
func TestTranspileNoFilesHelp(t *testing.T) {
	output, teardown := setupTest([]string{"test", "transpile"})
	defer teardown()

	runCommand()
	err := cupaloy.Snapshot(output.String())
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}

// Test that help is printed if help flag is set, even if file is given
func TestTranspileHelpFlag(t *testing.T) {
	output, teardown := setupTest([]string{"test", "transpile", "-h", "foo.c"})
	defer teardown()

	runCommand()
	err := cupaloy.Snapshot(output.String())
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}

// Test that help is printed if no files are given
func TestAstNoFilesHelp(t *testing.T) {
	output, teardown := setupTest([]string{"test", "ast"})
	defer teardown()

	runCommand()
	err := cupaloy.Snapshot(output.String())
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}

// Test that help is printed if help flag is set, even if file is given
func TestAstHelpFlag(t *testing.T) {
	output, teardown := setupTest([]string{"test", "ast", "-h", "foo.c"})
	defer teardown()

	runCommand()
	err := cupaloy.Snapshot(output.String())
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}

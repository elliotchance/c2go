// +build integration

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"regexp"

	"github.com/elliotchance/c2go/util"
)

type programOut struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
	isZero bool
}

// TestIntegrationScripts tests all programs in the tests directory.
//
// Integration tests are not run by default (only unit tests). These are
// indicated by the build flags at the top of the file. To include integration
// tests use:
//
//     go test -tags=integration
//
// You can also run a single file with:
//
//     go test -tags=integration -run=TestIntegrationScripts/tests/ctype/isalnum.c
//
func TestIntegrationScripts(t *testing.T) {
	testFiles, err := filepath.Glob("tests/*.c")
	if err != nil {
		t.Fatal(err)
	}

	exampleFiles, err := filepath.Glob("examples/*.c")
	if err != nil {
		t.Fatal(err)
	}

	files := append(testFiles, exampleFiles...)

	isVerbose := flag.CommandLine.Lookup("test.v").Value.String() == "true"

	totalTapTests := 0

	var (
		buildFolder  = "build"
		cFileName    = "a.out"
		goFileName   = "go.out"
		mainFileName = "main.go"
		stdin        = "7"
		args         = []string{"some", "args"}
		separator    = string(os.PathSeparator)
	)

	// Create build folder
	err = os.MkdirAll(buildFolder, os.ModePerm)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	t.Parallel()

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			cProgram := programOut{}
			goProgram := programOut{}

			// create subfolders for test
			subFolder := buildFolder + separator + strings.Split(file, ".")[0] + separator
			cPath := subFolder + cFileName
			goPath := subFolder + goFileName

			// Create build folder
			err := os.MkdirAll(subFolder, os.ModePerm)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			// Compile C.
			out, err := exec.Command("clang", "-lm", "-o", cPath, file).CombinedOutput()
			if err != nil {
				t.Fatalf("error: %s\n%s", err, out)
			}

			// Run C program
			cmd := exec.Command(cPath, args...)
			cmd.Stdin = strings.NewReader(stdin)
			cmd.Stdout = &cProgram.stdout
			cmd.Stderr = &cProgram.stderr
			err = cmd.Run()
			cProgram.isZero = err == nil

			programArgs := ProgramArgs{
				inputFile:   file,
				outputFile:  subFolder + separator + mainFileName,
				packageName: "main",
			}

			// Compile Go
			err = Start(programArgs)
			if err != nil {
				t.Fatalf("error: %s\n%s", err, out)
			}

			fmt.Println("sdsdsdsdsds")
			{
				buildErr := exec.Command("go", "build", "-o", goPath, subFolder+mainFileName)
				var out bytes.Buffer
				var stderr bytes.Buffer
				buildErr.Stdout = &out
				buildErr.Stderr = &stderr
				err = buildErr.Run()
				if err != nil {
					t.Fatalf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
				}
			}

			// Run Go program
			cmd = exec.Command(goPath, args...)
			cmd.Stdin = strings.NewReader(stdin)
			cmd.Stdout = &goProgram.stdout
			cmd.Stderr = &goProgram.stderr
			err = cmd.Run()
			goProgram.isZero = err == nil

			// Check for special exit codes that signal that tests have failed.
			if exitError, ok := err.(*exec.ExitError); ok {
				exitStatus := exitError.Sys().(syscall.WaitStatus).ExitStatus()
				if exitStatus == 101 || exitStatus == 102 {
					t.Fatal(goProgram.stdout.String())
				}
			}

			// Check if both exit codes are zero (or non-zero)
			if cProgram.isZero != goProgram.isZero {
				t.Fatalf("Exit statuses did not match.\n" +
					util.ShowDiff(cProgram.stdout.String(),
						goProgram.stdout.String()),
				)
			}

			// Check stderr
			if cProgram.stderr.String() != goProgram.stderr.String() {
				t.Fatalf("Expected %q, Got: %q",
					cProgram.stderr.String(),
					goProgram.stderr.String())
			}

			// Check stdout
			if cProgram.stdout.String() != goProgram.stdout.String() {
				t.Fatalf(util.ShowDiff(cProgram.stdout.String(),
					goProgram.stdout.String()))
			}

			// If this is not an example we will extact the number of tests run.
			if strings.Index(file, "examples/") == -1 && isVerbose {
				firstLine := strings.Split(goProgram.stdout.String(), "\n")[0]

				matches := regexp.MustCompile(`1\.\.(\d+)`).
					FindStringSubmatch(firstLine)
				if len(matches) == 0 {
					t.Fatalf("Test did not output tap: %s", file)
				}

				fmt.Printf("TAP: # %s: %s tests\n", file, matches[1])
				totalTapTests += util.Atoi(matches[1])
			}
		})
	}

	if isVerbose {
		fmt.Printf("TAP: # Total tests: %d\n", totalTapTests)
	}
}

func TestStartPreprocess(t *testing.T) {
	// temp dir
	tempDir := os.TempDir()

	// create temp file with garantee
	// wrong file body
	tempFile, err := newTempFile(tempDir, "c2go", "preprocess.c")
	if err != nil {
		t.Errorf("Cannot create temp file for execute test")
	}
	defer os.Remove(tempFile.Name())

	fmt.Fprintf(tempFile, "#include <AbsoluteWrongInclude.h>\nint main(void){\nwrong();\n}")

	err = tempFile.Close()
	if err != nil {
		t.Errorf("Cannot close the temp file")
	}

	var args ProgramArgs
	args.inputFile = tempFile.Name()

	err = Start(args)
	if err == nil {
		t.Errorf("Cannot test preprocess of application")
	}
}

func TestGoPath(t *testing.T) {
	gopath := "GOPATH"

	existEnv := os.Getenv(gopath)
	if existEnv == "" {
		t.Errorf("$GOPATH is not set")
	}

	// return env.var.
	defer func() {
		err := os.Setenv(gopath, existEnv)
		if err != nil {
			t.Errorf("Cannot restore the value of $GOPATH")
		}
	}()

	// reset value of env.var.
	err := os.Setenv(gopath, "")
	if err != nil {
		t.Errorf("Cannot set value of $GOPATH")
	}

	// testing
	err = Start(ProgramArgs{})
	if err == nil {
		t.Errorf(err.Error())
	}
}

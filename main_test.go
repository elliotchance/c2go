package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	cPath  = "build/a.out"
	goPath = "build/go.out"
	stdin  = "7"
	args   = []string{"some", "args"}
)

type programOut struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
	isZero bool
}

// TestIntegrationScripts tests all programs in the tests directory
func TestIntegrationScripts(t *testing.T) {
	files, err := filepath.Glob("tests/*/*.c")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		// Create build folder
		os.Mkdir("build/", os.ModePerm)

		t.Run(file, func(t *testing.T) {
			cProgram := programOut{}
			goProgram := programOut{}

			// Compile C
			err := exec.Command("clang", "-lm", "-o", cPath, file).Run()
			if err != nil {
				t.Fatal(err)
			}

			// Run C program
			cmd := exec.Command(cPath, args...)
			cmd.Stdin = strings.NewReader(stdin)
			cmd.Stdout = &cProgram.stdout
			cmd.Stderr = &cProgram.stderr
			err = cmd.Run()
			cProgram.isZero = err == nil

			// Compile Go
			goSrc := Start([]string{file})
			ioutil.WriteFile("build/main.go", []byte(goSrc), os.ModePerm)
			err = exec.Command("go", "build", "-o", goPath, "build/main.go").Run()
			if err != nil {
				t.Fatal(err)
			}

			// Run Go program
			cmd = exec.Command(goPath, args...)
			cmd.Stdin = strings.NewReader(stdin)
			cmd.Stdout = &goProgram.stdout
			cmd.Stderr = &goProgram.stderr
			err = cmd.Run()
			goProgram.isZero = err == nil

			// Check if both exit codes are zero (or non-zero)
			if cProgram.isZero != goProgram.isZero {
				t.Fatalf("Expected: %t, Got: %t", cProgram.isZero, goProgram.isZero)
			}

			// Check stderr
			if cProgram.stderr.String() != goProgram.stderr.String() {
				t.Fatalf("Expected %q, Got: %q", cProgram.stderr.String(), goProgram.stderr.String())
			}

			// Check stdout
			if cProgram.stdout.String() != goProgram.stdout.String() {
				t.Fatalf("Expected %q, Got: %q", cProgram.stdout.String(), goProgram.stdout.String())
			}
		})
	}
}

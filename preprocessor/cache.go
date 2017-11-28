package preprocessor

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Clang - parameters of clang execution
type Clang struct {
	Args  []string
	Files []string
}

// RunClang - run application clang with arguments
func RunClang(c Clang) (out bytes.Buffer, err error) {
	var stderr bytes.Buffer

	var a []string
	a = append(a, c.Args...)
	a = append(a, c.Files...)
	cmd := exec.Command("clang", a...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		err = fmt.Errorf("clang error:\nargs = %v\nfiles = %v\nerror = %v\nstderr = %v",
			c.Args, c.Files, err, stderr.String())
	}
	return
}

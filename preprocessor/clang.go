package preprocessor

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Clang - parameters of clang execution
type Clang struct {
	Args []string
	File string
}

// RunClang - run application clang with arguments
func RunClang(c Clang) (_ []byte, err error) {
	var (
		out    bytes.Buffer
		stderr bytes.Buffer
	)

	var a []string
	a = append(a, c.Args...)
	a = append(a, c.File)
	cmd := exec.Command("clang", a...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		err = fmt.Errorf("clang error:\nargs = %v\nfiles = %v\nerror = %v\nstderr = %v",
			c.Args, c.File, err, stderr.String())
	}
	return out.Bytes(), err
}

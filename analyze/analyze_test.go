package analyze_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Konstantin8105/c2go/analyze"
)

func TestStartPreprocess(t *testing.T) {
	// temp dir
	tempDir := os.TempDir()

	// create temp file with garantee
	// wrong file body
	tempFile, err := New(tempDir, "c2go", "preprocess.c")
	if err != nil {
		t.Errorf("Cannot create temp file for execute test")
	}
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()

	fmt.Fprintf(tempFile, "#include <AbsoluteWrongInclude.h>\nint main(void){\nwrong\n}")

	err = tempFile.Close()
	if err != nil {
		t.Errorf("Cannot close the temp file")
	}

	fmt.Println("tempDir  = ", tempDir)
	fmt.Println("tempFile = ", tempFile)

	var args analyze.ProgramArgs
	args.InputFile = tempFile.Name()

	err = analyze.Start(args)
	fmt.Println("err = ", err)
	if err == nil {
		t.Errorf("Cannot test preprocess of application")
	}
}

// New returns an unused filename for output files.
func New(dir, prefix, suffix string) (*os.File, error) {
	for index := 1; index < 10000; index++ {
		path := filepath.Join(dir, fmt.Sprintf("%s%03d%s", prefix, index, suffix))
		if _, err := os.Stat(path); err != nil {
			return os.Create(path)
		}
	}
	// Give up
	return nil, fmt.Errorf("could not create file of the form %s%03d%s", prefix, 1, suffix)
}

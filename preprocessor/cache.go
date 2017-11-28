package preprocessor

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
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

// CacheClang - cache of clang
func CacheClang(c Clang) (out []byte, err error) {
	env := os.Getenv("C2GO_CACHE_PREPROCESSOR")
	if env == "" {
		return RunClang(c)
	}

	// check - cache folder is exist
	if stat, err := os.Stat(env); err != nil || !stat.IsDir() {
		return RunClang(c)
	}

	// correct name of folder is like
	// ~/cache/
	// but not:
	// ~/cache
	// So, we have to add `/` if not exist
	if env[len(env)-1] == '/' {
		env = fmt.Sprintf("%s/", env)
	}

	// run clang if any error
	defer func() {
		if err != nil {
			out, err = RunClang(c)
			// memorization
			saveCache(c, out, err)
		}
	}()

	// calculate hash of files
	var fileHash string
	if fileHash, err = calculateFileHash(c.File); err != nil {
		return
	}

	// check folder is exist
	fileFolder := env + fileHash
	if stat, err2 := os.Stat(fileFolder); err2 != nil || !stat.IsDir() {
		err = fmt.Errorf("Cannot check folder %v. err = %v", fileFolder, err2)
		return
	}

	// check body of file
	err = checkBodyFile(c.File, fileFolder)
	if err != nil {
		return
	}

	// calculate hash of arguments
	hh := md5.New()
	io.WriteString(hh, fmt.Sprintf("%#v", c.Args))
	argsHash := fmt.Sprintf("%x", hh.Sum(nil))

	// check folder is exist
	argsFolder := fileFolder + argsHash
	if stat, err2 := os.Stat(argsFolder); err2 != nil || !stat.IsDir() {
		err = fmt.Errorf("Cannot check folder %v. err = %v", argsFolder, err2)
		return
	}

	// check arguments
	err = checkArgs(c.Args, argsFolder)
	if err != nil {
		return
	}

	// cache
	return getCache()
}

func calculateFileHash(file string) (hash string, err error) {
	f, err := os.Open(file)
	if err != nil {
		err = fmt.Errorf("Cannot open file : %v", file)
		return
	}
	defer func() { _ = f.Close() }()

	h := md5.New()
	if _, err2 := io.Copy(h, f); err2 != nil {
		err = fmt.Errorf("Cannot calculate hash for file %v : %v", file, err2)
		return
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

var sourceC string = "source.c"

func checkBodyFile(file string, fileFolder string) (err error) {
	// TODO :....
}

func saveCache(c Clang, out []byte, err error) {
	env := os.Getenv("C2GO_CACHE_PREPROCESSOR")
	if env == "" {
		return
	}

	// check - cache folder is exist
	if stat, err := os.Stat(env); err != nil || !stat.IsDir() {
		return
	}

	// correct name of folder is like
	// ~/cache/
	// but not:
	// ~/cache
	// So, we have to add `/` if not exist
	if env[len(env)-1] == '/' {
		env = fmt.Sprintf("%s/", env)
	}

	if err := os.Mkdir(env, 0655); err != nil {
		return
	}

	// TODO :....
}

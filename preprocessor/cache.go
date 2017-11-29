package preprocessor

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
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
	if env[len(env)-1] != '/' {
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
	fileFolder := env + fileHash + "/"
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
	_, err = io.WriteString(hh, fmt.Sprintf("%#v", c.Args))
	if err != nil {
		return
	}
	argsHash := fmt.Sprintf("%x", hh.Sum(nil))

	// check folder is exist
	argsFolder := fileFolder + argsHash + "/"
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
	return getCache(argsFolder)
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

var sourceC = "source.c"

func checkBodyFile(file string, fileFolder string) (err error) {
	var (
		content1 []byte
		err1     error
		content2 []byte
		err2     error
	)
	content1, err1 = ioutil.ReadFile(file)
	content2, err2 = ioutil.ReadFile(fileFolder + sourceC)
	if err1 != nil || err2 != nil {
		err = fmt.Errorf("Error in file reading : \n%v\n%v", err1, err2)
		return
	}
	if !bytes.Equal(content1, content2) {
		err = fmt.Errorf("Content is not same")
		return
	}
	return nil
}

var argsFile = "args.txt"

func checkArgs(args []string, argsFolder string) (err error) {
	var (
		content1 []byte
		content2 []byte
	)
	content1 = []byte(fmt.Sprintf("%#v", args))
	content2, err = ioutil.ReadFile(argsFolder + argsFile)
	if err != nil {
		err = fmt.Errorf("Error in file reading : %v", err)
		return
	}
	if !bytes.Equal(content1, content2) {
		err = fmt.Errorf("Content is not same. [buffer1,buffer2] = {%v,%v}", len(content1), len(content2))
		return
	}
	return nil
}

func saveCache(c Clang, out []byte, errResult error) {
	env := os.Getenv("C2GO_CACHE_PREPROCESSOR")
	if env == "" {
		return
	}

	// check - cache folder is exist
	if stat, err := os.Stat(env); err != nil || !stat.IsDir() {
		fmt.Println("env   err = ", err)
		return
	}

	// correct name of folder is like
	// ~/cache/
	// but not:
	// ~/cache
	// So, we have to add `/` if not exist
	if env[len(env)-1] != '/' {
		env = fmt.Sprintf("%s/", env)
	}

	if err := os.Mkdir(env, 0755); err != nil {
		_ = err //return
	}

	// calculate hash of files
	fileHash, err := calculateFileHash(c.File)
	if err != nil {
		fmt.Println("fileHash   err = ", err)
		return
	}

	// check folder is exist
	fileFolder := env + fileHash + "/"
	if err := os.Mkdir(fileFolder, 0755); err != nil {
		_ = err
	}

	// copy file
	if err := Copy(c.File, fileFolder+sourceC); err != nil {
		return
	}

	// calculate hash of arguments
	hh := md5.New()
	_, err = io.WriteString(hh, fmt.Sprintf("%#v", c.Args))
	if err != nil {
		return
	}
	argsHash := fmt.Sprintf("%x", hh.Sum(nil))

	// check folder is exist
	argsFolder := fileFolder + argsHash + "/"
	if err := os.Mkdir(argsFolder, 0755); err != nil {
		return
	}

	// check arguments
	if err := ioutil.WriteFile(argsFolder+argsFile, []byte(fmt.Sprintf("%#v", c.Args)), 0755); err != nil {
		return
	}

	// cache
	if err := ioutil.WriteFile(argsFolder+outFile, out, 0755); err != nil {
		return
	}

	var sErr string
	if errResult != nil {
		sErr = errResult.Error()
	}
	if err := ioutil.WriteFile(argsFolder+errFile, []byte(sErr), 0755); err != nil {
		return
	}
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return ioutil.WriteFile(dst, []byte(content), 0755)
}

var outFile = "out.txt"
var errFile = "err.txt"

func getCache(folder string) (out []byte, err error) {
	var (
		content1 []byte
		err1     error
		content2 []byte
		err2     error
	)
	content1, err1 = ioutil.ReadFile(folder + "/" + outFile)
	content2, err2 = ioutil.ReadFile(folder + "/" + errFile)
	if err1 != nil || err2 != nil {
		err = fmt.Errorf("Error in file reading : \n%v\n%v", err1, err2)
		return
	}

	if len(content2) != 0 {
		err = fmt.Errorf("%v", content2)
		return
	}
	return content1, nil
}

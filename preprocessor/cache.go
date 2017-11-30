package preprocessor

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type cache Clang

// name of files
var (
	sourceC  = "source.c"
	argsFile = "args.txt"
	outFile  = "out.txt"
	errFile  = "err.txt"
)

// CacheClang - cache of clang
func CacheClang(c Clang) (out []byte, err error) {
	if !isCacheSwitchOn() {
		return RunClang(c)
	}

	// run clang if any error
	defer func() {
		if err != nil {
			out, err = RunClang(c)
			// memorization
			cache(c).saveCache(out, err)
		}
	}()
	return cache(c).getCache()
}

// isCacheSwitchOn - check cache is switch on
func isCacheSwitchOn() bool {
	env := os.Getenv("C2GO_CACHE_PREPROCESSOR")
	if env == "" {
		return false
	}

	// check - cache folder is exist
	if stat, err := os.Stat(env); err != nil || !stat.IsDir() {
		return false
	}

	return true
}

// getCache - return output from cache
func (c cache) getCache() (out []byte, err error) {

	env := os.Getenv("C2GO_CACHE_PREPROCESSOR")
	if env == "" {
		err = fmt.Errorf("Not correct environment variable")
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

	var fileFolder string
	fileFolder, err = checkSourceFile(c.File, env)
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
	return func(folder string) (out []byte, err error) {
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
	}(argsFolder)
}

// checkSourceFile - compare source file
func checkSourceFile(fileBase string, folderEnv string) (folderCacheSource string, err error) {
	// calculate hash of files
	var contentFileBase []byte
	contentFileBase, err = ioutil.ReadFile(fileBase)
	if err != nil {
		return
	}
	h := md5.New()
	_, err = h.Write(contentFileBase)
	if err != nil {
		return
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))

	// check folder is exist
	folderCacheSource = folderEnv + hash + "/"
	if stat, err2 := os.Stat(folderCacheSource); err2 != nil || !stat.IsDir() {
		err = fmt.Errorf("Cannot check folder %v. err = %v", folderCacheSource, err2)
		return
	}

	// check body of file
	var contentCache []byte
	contentCache, err = ioutil.ReadFile(folderCacheSource + sourceC)
	if err != nil {
		return
	}
	if !bytes.Equal(contentFileBase, contentCache) {
		err = fmt.Errorf("Content is not same")
		return
	}
	return
}

// checkArgs - compare arguments
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

// saveCache - save cache if acceptable
func (c cache) saveCache(out []byte, errResult error) {
	if !isCacheSwitchOn() {
		return
	}

	env := os.Getenv("C2GO_CACHE_PREPROCESSOR")
	if env == "" {
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
	contentFileBase, err := ioutil.ReadFile(c.File)
	if err != nil {
		return
	}
	h := md5.New()
	_, err = h.Write(contentFileBase)
	if err != nil {
		return
	}
	fileHash := fmt.Sprintf("%x", h.Sum(nil))

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

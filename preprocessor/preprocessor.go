package preprocessor

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// One simple part of preprocessor code
type entity struct {
	positionInSource int
	include          string
	other            string

	// Zero index of `lines` is look like that:
	// # 11 "/usr/include/x86_64-linux-gnu/gnu/stubs.h" 2 3 4
	// After that 0 or more lines of codes
	lines []*string
}

// isSame - check is Same entities
func (e *entity) isSame(x *entity) bool {
	if e.include != x.include {
		return false
	}
	if e.positionInSource != x.positionInSource {
		return false
	}
	if e.other != x.other {
		return false
	}
	if len(e.lines) != len(x.lines) {
		return false
	}
	for k := range e.lines {
		is := e.lines[k]
		js := x.lines[k]
		if len(*is) != len(*js) {
			return false
		}
		if *is != *js {
			return false
		}
	}
	return true
}

// Analyze - separation preprocessor code to part
func Analyze(inputFiles, clangFlags []string) (pp []byte, err error) {
	var allItems []entity

	allItems, err = analyzeFiles(inputFiles, clangFlags)
	if err != nil {
		return
	}

	// Merge the entities
	var lines []string
	for i := range allItems {
		// If found same part of preprocess code, then
		// don't include in result buffer for transpiling
		// for avoid dublicate of code
		var found bool
		for j := 0; j < i; j++ {
			if allItems[i].isSame(&allItems[j]) {
				found = true
				break
			}
		}
		if found {
			continue
		}

		// Parameter "other" is not included for avoid like:
		// ./tests/multi/head.h:4:28: error: invalid line marker flag '2': cannot pop empty include stack
		// # 2 "./tests/multi/main.c" 2
		//                            ^
		header := fmt.Sprintf("# %d \"%s\"", allItems[i].positionInSource, allItems[i].include)
		lines = append(lines, header)
		if len(allItems[i].lines) > 0 {
			for ii, l := range allItems[i].lines {
				if ii == 0 {
					continue
				}
				lines = append(lines, *l)
			}
		}
	}
	pp = ([]byte)(strings.Join(lines, "\n"))

	return
}

// analyze - analyze single file and separation preprocessor code to part
func analyzeFiles(inputFiles, clangFlags []string) (items []entity, err error) {
	// See : https://clang.llvm.org/docs/CommandGuide/clang.html
	// clang -E <file>    Run the preprocessor stage.
	var out bytes.Buffer
	out, err = getPreprocessSources(inputFiles, clangFlags)
	if err != nil {
		return
	}

	// Parsing preprocessor file
	r := bytes.NewReader(out.Bytes())
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	// counter - get position of line
	var counter int
	// item, items - entity of preprocess file
	var item *entity
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] == '#' &&
			(len(line) >= 7 && line[0:7] != "#pragma") {
			if item != (*entity)(nil) {
				items = append(items, *item)
			}
			item, err = parseIncludePreprocessorLine(line)
			if err != nil {
				err = fmt.Errorf("Cannot parse line : %s with error: %s", line, err)
				return
			}
			if item.positionInSource == 0 {
				// cannot by less 1 for avoid problem with
				// indentification of "0" AST base element
				item.positionInSource = 1
			}
			item.lines = make([]*string, 0)
		}
		counter++
		item.lines = append(item.lines, &line)
	}
	if item != (*entity)(nil) {
		items = append(items, *item)
	}
	return
}

// See : https://clang.llvm.org/docs/CommandGuide/clang.html
// clang -E <file>    Run the preprocessor stage.
func getPreprocessSources(inputFiles, clangFlags []string) (out bytes.Buffer, err error) {
	var stderr bytes.Buffer
	flags := make(map[string]bool)
	for i := range clangFlags {
		flags[clangFlags[i]] = true
	}

	for pos, inputFile := range inputFiles {
		if inputFile[len(inputFile)-1] != 'c' {
			continue
		}

		if pos > 0 {
			var define []string
			define, err = getDefinitionsOfFile(inputFiles[pos-1])
			if err != nil {
				return
			}
			for i := range define {
				fmt.Println(i, "\t", define[i])
				flags[fmt.Sprintf("-D%s", define[i])] = true
			}
		}

		var args []string
		args = append(args, "-E")
		for k := range flags {
			args = append(args, k)
		}
		args = append(args, inputFile)

		var outFile bytes.Buffer
		cmd := exec.Command("clang", args...)
		cmd.Stdout = &outFile
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			err = fmt.Errorf("preprocess for file: %s\nfailed: %v\nStdErr = %v", inputFile, err, stderr.String())
			return
		}
		_, err = out.Write(outFile.Bytes())
		if err != nil {
			return
		}
	}
	return
}

// getIncludeListWithUserSource - Get list of include files
// Example:
// $ clang  -MM -c exit.c
// exit.o: exit.c tests.h
func getIncludeListWithUserSource(inputFile string) (lines []string, err error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("clang", "-MM", "-c", inputFile)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
		return
	}
	return parseIncludeList(out.String())
}

// getIncludeFullList - Get full list of include files
// Example:
// $ clang -M -c triangle.c
// triangle.o: triangle.c /usr/include/stdio.h /usr/include/features.h \
//   /usr/include/stdc-predef.h /usr/include/x86_64-linux-gnu/sys/cdefs.h \
//   /usr/include/x86_64-linux-gnu/bits/wordsize.h \
//   /usr/include/x86_64-linux-gnu/gnu/stubs.h \
//   /usr/include/x86_64-linux-gnu/gnu/stubs-64.h \
//   / ........ and other
func getIncludeFullList(inputFile string) (lines []string, err error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("clang", "-M", "-c", inputFile)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
		return
	}
	return parseIncludeList(out.String())
}

// getDefineList - get list C defines of C file
// Example:
// clang -dM -E  1.c
// Result:
// #define BUFSIZ _IO_BUFSIZ
// #define EOF (-1)
func getDefineList(inputFile string) (define []string, err error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("clang", "-dM", "-E", inputFile)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
		return
	}
	return parseDefineList(out.String())
}

// getDefinitionsOfFile - return #define from user file
func getDefinitionsOfFile(inputFile string) (define []string, err error) {
	// get full list of files ( user + system includes)
	allFiles, err := getIncludeFullList(inputFile)
	if err != nil {
		return
	}

	// get full list definitions from all files
	allDefine, err := getDefineList(inputFile)
	if err != nil {
		return
	}

	// get list of user files
	userFile, err := getIncludeListWithUserSource(inputFile)
	if err != nil {
		return
	}

	// calculate list of system include files, only
	systemFiles := minus(allFiles, userFile)

	// calculate list definitions of system include files, only
	var systemDefine []string
	for _, systemFile := range systemFiles {
		var d []string
		d, err = getDefineList(systemFile)
		if err != nil {
			// some system include files
			// cannot be taked directly
			// so, we use `continue` isteand of `return`
			err = nil
			continue
		}
		systemDefine = append(systemDefine, d...)
	}

	// calculate = (all definitions) minus (system definitions)
	define = minus(allDefine, systemDefine)

	define = minus(define, []string{
		"_FILE_OFFSET_BITS 64",
	})

	var t []string
	for i := range define {
		ss := strings.Split(define[i], " ")
		if len(ss) == 1 {
			t = append(t, ss[0])
		}
	}
	define = t

	return
}

// minus - result of c = a - b
func minus(a, b []string) (c []string) {
	ma := toMap(a)
	mb := toMap(b)

	for kb := range mb {
		if _, ok := ma[kb]; ok {
			delete(ma, kb)
		}
	}
	for ka := range ma {
		c = append(c, ka)
	}
	return
}

func toMap(list []string) (m map[string]bool) {
	m = make(map[string]bool)
	for i := range list {
		m[list[i]] = true
	}
	return
}

func unique(list []string) (res []string) {
	m := toMap(list)
	for k := range m {
		res = append(res, k)
	}
	return res
}

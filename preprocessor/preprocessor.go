package preprocessor

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/scanner"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
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

func (e *entity) parseComments(comments *[]program.Comment) {
	var source bytes.Buffer
	for i := range e.lines {
		if i == 0 {
			continue
		}
		source.Write([]byte(*e.lines[i]))
		source.Write([]byte{'\n'})
	}

	var s scanner.Scanner
	s.Init(strings.NewReader(source.String()))
	s.Mode = scanner.ScanComments
	s.Filename = e.include
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if scanner.TokenString(tok) == "Comment" {
			(*comments) = append(*comments, program.Comment{
				File:    e.include,
				Line:    s.Position.Line + e.positionInSource - 1,
				Comment: s.TokenText(),
			})
		}
	}
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
func Analyze(inputFiles, clangFlags []string) (pp []byte,
	comments []program.Comment, includes []program.IncludeHeader, err error) {

	var allItems []entity

	allItems, err = analyzeFiles(inputFiles, clangFlags)
	if err != nil {
		return
	}

	// Generate list of user files
	userSource := map[string]bool{}
	var us []string
	us, err = GetIncludeListWithUserSource(inputFiles, clangFlags)
	if err != nil {
		return
	}
	var all []string
	all, err = GetIncludeFullList(inputFiles, clangFlags)
	if err != nil {
		return
	}
	// Generate C header list
	includes = generateIncludeList(us, all)

	for j := range us {
		userSource[us[j]] = true
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
		// Parse comments only for user sources
		var isUserSource bool
		if userSource[allItems[i].include] {
			isUserSource = true
		}
		if allItems[i].include[0] == '.' &&
			allItems[i].include[1] == '/' &&
			userSource[allItems[i].include[2:]] {
			isUserSource = true
		}
		if isUserSource {
			allItems[i].parseComments(&comments)
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

// analyzeFiles - analyze single file and separation preprocessor code to part
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

	reg := util.GetRegex("# (\\d+) \".*\".*")

	for scanner.Scan() {
		line := scanner.Text()
		if reg.MatchString(line) {
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
	// get temp dir
	dir, err := ioutil.TempDir("", "c2go-union")
	if err != nil {
		return
	}
	defer func() { _ = os.RemoveAll(dir) }()

	// file name union file
	var unionFileName = dir + "/" + "unionFileName.c"

	// create a body for union file
	var unionBody string
	for i := range inputFiles {
		var absPath string
		absPath, err = filepath.Abs(inputFiles[i])
		if err != nil {
			return
		}
		unionBody += fmt.Sprintf("#include \"%s\"\n", absPath)
	}

	// write a union file
	err = ioutil.WriteFile(unionFileName, []byte(unionBody), 0644)
	if err != nil {
		return
	}

	// Add open source defines
	if runtime.GOOS == "darwin" {
		clangFlags = append(clangFlags, "-D_XOPEN_SOURCE")
	} else {
		clangFlags = append(clangFlags, "-D_GNU_SOURCE")
	}

	// preprocessor clang
	var stderr bytes.Buffer

	var args []string
	args = append(args, "-E", "-C")
	args = append(args, clangFlags...)
	args = append(args, unionFileName) // All inputFiles

	var outFile bytes.Buffer
	cmd := exec.Command("clang", args...)
	cmd.Stdout = &outFile
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("preprocess for file: %v\nfailed: %v\nStdErr = %v", inputFiles, err, stderr.String())
		return
	}
	_, err = out.Write(outFile.Bytes())
	if err != nil {
		return
	}

	return
}

func generateIncludeList(userList, allList []string) (
	includes []program.IncludeHeader) {

	for i := range allList {
		var isUser bool
		for j := range userList {
			if allList[i] == userList[j] {
				isUser = true
				break
			}
		}
		includes = append(includes, program.IncludeHeader{
			HeaderName:   allList[i],
			IsUserSource: isUser,
		})
	}
	return
}

// GetIncludeListWithUserSource - Get list of include files
// Example:
// $ clang  -MM -c exit.c
// exit.o: exit.c tests.h
func GetIncludeListWithUserSource(inputFiles, clangFlags []string) (lines []string, err error) {
	return getIncludeList(inputFiles, clangFlags, "-MM")
}

// GetIncludeFullList - Get full list of include files
// Example:
// $ clang -M -c triangle.c
// triangle.o: triangle.c /usr/include/stdio.h /usr/include/features.h \
//   /usr/include/stdc-predef.h /usr/include/x86_64-linux-gnu/sys/cdefs.h \
//   /usr/include/x86_64-linux-gnu/bits/wordsize.h \
//   /usr/include/x86_64-linux-gnu/gnu/stubs.h \
//   /usr/include/x86_64-linux-gnu/gnu/stubs-64.h \
//   / ........ and other
func GetIncludeFullList(inputFiles, clangFlags []string) (lines []string, err error) {
	return getIncludeList(inputFiles, clangFlags, "-M")
}

func getIncludeList(inputFiles, clangFlags []string, flag string) (lines []string, err error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	var args []string
	for i := range inputFiles {
		inputFiles[i], err = filepath.Abs(inputFiles[i])
		if err != nil {
			return
		}
	}
	args = append(args, flag, "-c")
	args = append(args, inputFiles...)
	args = append(args, clangFlags...)
	cmd := exec.Command("clang", args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
		return
	}
	return parseIncludeList(out.String())
}

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

// Analyze - separation preprocessor code to part
func Analyze(inputFile string) (pp []byte, userPosition int, err error) {
	// See : https://clang.llvm.org/docs/CommandGuide/clang.html
	// clang -E <file>    Run the preprocessor stage.
	out, err := getPreprocessSource(inputFile)
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
	var items []entity
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] == '#' {
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

	// Get list of include files
	includeList, err := getIncludeList(inputFile)
	if err != nil {
		return
	}

	// Renumbering positionInSource in source for user code to unique
	// Let`s call that positionInSource - userPosition
	// for example: if some entity(GenDecl,...) have positionInSource
	// less userPosition, then that is from system library, but not
	// from user source.
	for _, item := range items {
		if userPosition < item.positionInSource {
			userPosition = item.positionInSource
		}
	}
	for _, item := range items {
		userPosition += len(item.lines)
	}
	for i := range items {
		var found bool
		for _, inc := range includeList {
			if inc == items[i].include {
				found = true
			}
		}
		if !found {
			continue
		}
		items[i].positionInSource = userPosition + 1
	}
	// Now, userPosition is unique and more then other

	// Merge the entities
	lines := make([]string, 0, counter)
	for _, item := range items {
		lines = append(lines, fmt.Sprintf("# %d \"%s\" %s", item.positionInSource, item.include, item.other))
		if len(item.lines) > 0 {
			for i, l := range item.lines {
				if i == 0 {
					continue
				}
				lines = append(lines, *l)
			}
		}
	}
	pp = ([]byte)(strings.Join(lines, "\n"))

	// TODO return list of system `includes`
	return
}

// getIncludeList - Get list of include files
// Example:
// $ clang  -MM -c exit.c
// exit.o: exit.c tests.h
func getIncludeList(inputFile string) (lines []string, err error) {
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

// See : https://clang.llvm.org/docs/CommandGuide/clang.html
// clang -E <file>    Run the preprocessor stage.
func getPreprocessSource(inputFile string) (out bytes.Buffer, err error) {
	var stderr bytes.Buffer
	cmd := exec.Command("clang", "-E", inputFile)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
		return
	}
	return
}

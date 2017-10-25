package preprocessor

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// item - part of preprocessor code
type entity struct {
	positionInSource int
	include          string
	other            string

	// Zero index line is look like that:
	// # 11 "/usr/include/x86_64-linux-gnu/gnu/stubs.h" 2 3 4
	// After that 0 or more lines of codes
	lines []string
}

// Analyze - separation preprocessor code to part
func Analyze(inputFile string) (pp []byte, userPosition int, err error) {
	// See : https://clang.llvm.org/docs/CommandGuide/clang.html
	// clang -E <file>    Run the preprocessor stage.
	out, err := preprocessSource(inputFile)
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
			item, err = parseInclude(line)
			if err != nil {
				err = fmt.Errorf("Cannot parse line : %s with error: %s", line, err)
				return
			}
			if item.positionInSource == 0 {
				item.positionInSource = 1 // Hack : cannot by less 1
			}
			item.lines = make([]string, 0)
		}
		counter++
		item.lines = append(item.lines, line)
	}
	if item != (*entity)(nil) {
		items = append(items, *item)
	}

	// Get list of include files
	includeList, err := includesList(inputFile)
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

	// Merge the items
	lines := make([]string, 0, counter)
	for _, item := range items {
		lines = append(lines, fmt.Sprintf("# %d \"%s\" %s", item.positionInSource, item.include, item.other))
		if len(item.lines) > 0 {
			lines = append(lines, item.lines[1:]...)
		}
	}
	pp = ([]byte)(strings.Join(lines, "\n"))

	return
}

// typically parse that line:
// # 11 "/usr/include/x86_64-linux-gnu/gnu/stubs.h" 2 3 4
func parseInclude(line string) (item *entity, err error) {
	if line[0] != '#' {
		err = fmt.Errorf("Cannot parse: first symbol is not # in line %s", line)
		return
	}
	i := strings.Index(line, "\"")
	if i < 0 {
		err = fmt.Errorf("First index is not correct on line %s", line)
		return
	}
	l := strings.LastIndex(line, "\"")
	if l < 0 {
		err = fmt.Errorf("Last index is not correct on line %s", line)
		return
	}
	if i >= l {
		err = fmt.Errorf("Not allowable positions of symbol \" (%d and %d) in line : %s", i, l, line)
		return
	}

	pos, err := strconv.ParseInt(strings.TrimSpace(line[1:i]), 10, 64)
	if err != nil {
		err = fmt.Errorf("Cannot parse position in source : %v", err)
		return
	}

	if l+1 < len(line) {
		item = &entity{
			positionInSource: int(pos),
			include:          line[i+1 : l],
			other:            line[l+1:],
		}
	} else {
		item = &entity{
			positionInSource: int(pos),
			include:          line[i+1 : l],
		}
	}

	return
}

// includesList - Get list of include files
// Example:
// $ clang  -MM -c exit.c
// exit.o: exit.c tests.h
func includesList(inputFile string) (lines []string, err error) {

	// TODO Add test with multilines

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

	// Parse output
	i := strings.Index(out.String(), ":")
	if i < 0 {
		err = fmt.Errorf("First index is not correct on line %s", out.String())
		return
	}

	line := out.String()[i+1:]
	line = line[:len(line)-1] // remove last \n
	lines = strings.Split(line, " ")

	//fmt.Printf("INCLUDE : %#v", lines)
	return
}

// See : https://clang.llvm.org/docs/CommandGuide/clang.html
// clang -E <file>    Run the preprocessor stage.
func preprocessSource(inputFile string) (out bytes.Buffer, err error) {
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

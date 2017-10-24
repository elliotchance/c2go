package preprocessor

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Item - part of preprocessor code
type Item struct {
	Include string
	Lines   []string
}

// Analyze - separation preprocessor code to part
func Analyze(inputFile string) (pp []byte, err error) {
	// See : https://clang.llvm.org/docs/CommandGuide/clang.html
	// clang -E <file>    Run the preprocessor stage.
	var out bytes.Buffer
	{
		var stderr bytes.Buffer
		cmd := exec.Command("clang", "-E", inputFile)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			err = fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
			return
		}
	}

	// Get list of include files
	includeList, err := getIncludeList(inputFile)
	if err != nil {
		return
	}
	_ = includeList

	// Parsing preprocessor file
	r := bytes.NewReader(out.Bytes())
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var positions []int
	var lines []string
	var counter int
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			positions = append(positions, counter)
		}
		counter++
		lines = append(lines, line)
	}
	var item Item
	var items []Item
	for i := range positions {
		item.Include, err = parseInclude(lines[positions[i]])
		if err != nil {
			err = fmt.Errorf("Cannot parse line : %s", lines[positions[i]])
			return
		}

		// Filter of includes
		var found bool
		for _, in := range includeList {
			if in == item.Include {
				found = true
			}
		}
		if !found {
			continue
		}

		var s int
		if i != len(positions)-1 {
			s = positions[i] + 1
		} else {
			if positions[i]+1 < len(lines)-1 {
				s = positions[i] + 1
			} else {
				continue
			}
		}

		var f int
		if i != len(positions)-1 {
			f = positions[i+1]
		} else {
			f = len(lines)
		}
		item.Lines = lines[s:f]

		items = append(items, item)
	}
	_ = items

	// Merge the items
	lines = make([]string, 0)
	for i := range items {
		//	lines = append(lines, "# 1 "+items[i].Include)
		lines = append(lines, items[i].Lines...)
	}
	pp = ([]byte)(strings.Join(lines, "\n"))

	///fmt.Println("pp = ", string(pp))
	return
}

func parseInclude(line string) (inc string, err error) {
	i := strings.Index(line, "\"")
	if i < 0 {
		err = fmt.Errorf("First index is not correct on line %s", line)
	}
	l := strings.LastIndex(line, "\"")
	if l < 0 {
		err = fmt.Errorf("Last index is not correct on line %s", line)
	}

	inc = line[i+1 : l]
	if inc == "" {
		err = fmt.Errorf("Cannot found include in line: %s", line)
		return
	}

	return
}

// getIncludeList - Get list of include files
func getIncludeList(inputFile string) (lines []string, err error) {
	/* Example:
	$ clang  -MM -c exit.c
	exit.o: exit.c tests.h
	*/
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

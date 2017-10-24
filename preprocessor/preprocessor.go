package preprocessor

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Item - part of preprocessor code
type Item struct {
	Include string
	Lines   []string
}

// Analyze - separation preprocessor code to part
func Analyze(pp bytes.Buffer) (items []Item, err error) {
	r := bytes.NewReader(pp.Bytes())
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
			counter++
		}
		lines = append(lines, line)
	}
	var item Item
	for i := range positions {
		item.Include, err = parseInclude(lines[positions[i]])
		if err != nil {
			err = fmt.Errorf("Cannot parse line : %s", lines[positions[i]])
			return
		}
		if i != len(positions)-1 {
			item.Lines = lines[positions[i]:positions[i+1]]
		} else {
			item.Lines = lines[positions[i]:]
		}
		items = append(items, item)
	}
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

package preprocessor

import (
	"bufio"
	"bytes"
)

type Item struct {
	Include string
	Lines   []string
}

func Analyze(pp bytes.Buffer) (items []Item) {
	r := bytes.NewReader(pp.Bytes())
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var item Item
	for scanner.Scan() {
		str := scanner.Text()
		if len(str) == 0 {
			continue
		}
		if str[0] == '#' {
			if item.Include != "" && len(item.Lines) > 0 {
				items = append(items, item)
			}
			item.Include = parseInclude(str)
			continue
		}
		item.Lines = append(item.Lines, str)
	}
	if item.Include != "" {
		items = append(items, item)
	}
	return
}

func parseInclude(line string) (inc string) {
	return line
}

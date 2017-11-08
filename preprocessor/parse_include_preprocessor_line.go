package preprocessor

import (
	"fmt"
	"strconv"
	"strings"
)

// typically parse that line:
// # 11 "/usr/include/x86_64-linux-gnu/gnu/stubs.h" 2 3 4
func parseIncludePreprocessorLine(line string) (item *entity, err error) {
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

package preprocessor

import (
	"fmt"
	"strings"
)

const (
	define = "#define "
)

func parseDefineList(line string) (define []string, err error) {
	split := strings.Split(line, "\n")

	define = make([]string, len(split))

	for i := range split {
		split[i] = strings.Replace(split[i], "\n", " ", -1)
		split[i] = strings.Replace(split[i], "\t", " ", -1)
		split[i] = strings.Replace(split[i], "\r", " ", -1)
		split[i] = strings.Replace(split[i], "\\", " ", -1)
		split[i] = strings.Replace(split[i], "\xFF", " ", -1)
		split[i] = strings.Replace(split[i], "\u0100", " ", -1)

		if len(split[i]) < len(define)+3 {
			err = fmt.Errorf("Not correct length of line : %v", split[i])
			return
		}

		define[i] = split[i][len(define)+3:]
	}

	return
}

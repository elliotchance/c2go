package preprocessor

import (
	"fmt"
	"strings"
)

const (
	defineString = `#define `
)

func parseDefineList(line string) (define []string, err error) {
	split := strings.Split(line, "\n")

	define = make([]string, len(split))

	var counter int
	for i := range split {
		split[i] = strings.Replace(split[i], "\n", " ", -1)
		split[i] = strings.Replace(split[i], "\t", " ", -1)
		split[i] = strings.Replace(split[i], "\r", " ", -1)
		split[i] = strings.Replace(split[i], "\\", " ", -1)
		split[i] = strings.Replace(split[i], "\xFF", " ", -1)
		split[i] = strings.Replace(split[i], "\u0100", " ", -1)

		if len(split[i]) == 0 {
			continue
		}
		if len(split[i]) < len(defineString) {
			err = fmt.Errorf("Not correct length of line : |%v| . define = |%v|", split[i], defineString)
			return
		}

		define[counter] = split[i][len(defineString):]
		counter++
	}

	return
}

package preprocessor

import (
	"strings"
)

// parseIncludeList - parse list of includes
// Example :
// exit.o: exit.c /usr/include/stdlib.h /usr/include/features.h \
//    /usr/include/stdc-predef.h /usr/include/x86_64-linux-gnu/sys/cdefs.h
func parseIncludeList(line string) (lines []string, err error) {
	line = strings.Replace(line, "\n", " ", -1)
	line = strings.Replace(line, "\r", " ", -1)
	line = strings.Replace(line, "\\", " ", -1)
	parts := strings.Split(line, " ")

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if p[len(p)-1] == ':' {
			continue
		}
		lines = append(lines, p)
	}
	return
}

package cc

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/elliotchance/c2go/util"
)

// fileCache contains the previous read files. This is important because
// rereading large preprocessed files is extremely expensive and they should not
// change during the life of the executable.
var fileCache map[string]map[int]string

// ResetCache must be called before the next transpile happens if the executable
// is being used to transpile multiple files.
func ResetCache() {
	fileCache = nil
}

// GetLineFromPreprocessedFile returns a specific line for a file from a
// preprocessed C file.
func GetLineFromPreprocessedFile(inputFilePath, filePath string, lineNumber int) (string, error) {
	// Only load the file once.
	if fileCache == nil {
		fileCache = map[string]map[int]string{}

		inputFile, err := ioutil.ReadFile(inputFilePath)
		if err != nil {
			return "", err
		}

		lines := strings.Split(string(inputFile), "\n")
		currentFile := ""
		currentLine := 1

		// There is also an integer that appears after the path - not sure what this
		// is? I will ignore it for now.
		resetLineRegexp := util.GetRegex(`^# (\d+) "(.+)"`)

		for _, line := range lines {
			if len(line) > 0 && line[0] == '#' {
				matches := resetLineRegexp.FindStringSubmatch(line)

				// Ignore other preprocessor lines like: #pragma pack(4)
				if len(matches) == 0 {
					continue
				}

				currentLine = util.Atoi(matches[1])

				// unescape windows file paths
				currentFile = strings.Replace(matches[2], "\\\\", "\\", -1)

				if _, ok := fileCache[currentFile]; !ok {
					fileCache[currentFile] = map[int]string{}
				}

				continue
			}

			fileCache[currentFile][currentLine] = line

			currentLine++
		}
	}

	if _, ok := fileCache[filePath]; !ok {
		return "", fmt.Errorf("could not find file %s", filePath)
	}

	if line, ok := fileCache[filePath][lineNumber]; ok {
		return line, nil
	}

	return "", fmt.Errorf("could not find %s:%d", filePath, lineNumber)
}

package noarch

import (
	"os"
)

func Fopen(filePath, mode string) int {
	panic("fopen is not supported")
}

func Fclose(int) int {
	panic("fclose is not supported")
}

func Remove(filePath string) int {
	if os.Remove(filePath) != nil {
		return -1
	}

	return 0
}

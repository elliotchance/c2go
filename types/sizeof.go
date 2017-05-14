package types

import (
	"fmt"
	"strings"
)

func removePrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		s = s[len(prefix):]
	}

	return s
}

// SizeOf returns the number of bytes for a type. This the same as using the
// sizeof operator/function in C.
func SizeOf(cType string) int {
	// Remove keywords that do not effect the size.
	cType = removePrefix(cType, "signed ")
	cType = removePrefix(cType, "unsigned ")
	cType = removePrefix(cType, "const ")
	cType = removePrefix(cType, "volatile ")

	// FIXME: The pointer size will be different on different platforms. We
	// should find out the correct size at runtime.
	// pointerSize := 4

	switch cType {
	case "char":
		return 1

	case "short":
		return 2

	case "int", "float":
		return 4

	case "long", "double":
		return 8

	case "long double":
		return 16

	default:
		panic(fmt.Sprintf("cannot determine size of: %s", cType))
	}
}

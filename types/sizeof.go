package types

import (
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/program"
)

func removePrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		s = s[len(prefix):]
	}

	return s
}

// SizeOf returns the number of bytes for a type. This the same as using the
// sizeof operator/function in C.
func SizeOf(p *program.Program, cType string) int {
	// Remove keywords that do not effect the size.
	cType = removePrefix(cType, "signed ")
	cType = removePrefix(cType, "unsigned ")
	cType = removePrefix(cType, "const ")
	cType = removePrefix(cType, "volatile ")

	// FIXME: The pointer size will be different on different platforms. We
	// should find out the correct size at runtime.
	pointerSize := 8

	// A structure will be the sum of its parts.
	if strings.HasPrefix(cType, "struct ") {
		totalBytes := 0

		for _, t := range p.Structs[cType[7:]].Fields {
			totalBytes += SizeOf(p, t)
		}

		// The size of a struct is rounded up to fit the size of the pointer of
		// the OS.
		if totalBytes%pointerSize != 0 {
			totalBytes += pointerSize - (totalBytes % pointerSize)
		}

		return totalBytes
	}

	// Function pointers are one byte?
	if strings.Index(cType, "(") >= 0 {
		return 1
	}

	if strings.HasSuffix(cType, "*") {
		return pointerSize
	}

	switch cType {
	case "char", "void":
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

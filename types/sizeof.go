package types

import (
	"errors"
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
func SizeOf(p *program.Program, cType string) (int, error) {
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

		s := p.Structs[cType]
		if s == nil {
			return 0, errors.New(fmt.Sprintf("could not sizeof: %s", cType))
		}

		for _, t := range s.Fields {
			var bytes int
			var err error

			switch f := t.(type) {
			case string:
				bytes, err = SizeOf(p, f)

			case *program.Struct:
				bytes, err = SizeOf(p, f.Name)
			}

			if err != nil {
				return 0, err
			}
			totalBytes += bytes
		}

		// The size of a struct is rounded up to fit the size of the pointer of
		// the OS.
		if totalBytes%pointerSize != 0 {
			totalBytes += pointerSize - (totalBytes % pointerSize)
		}

		return totalBytes, nil
	}

	// Function pointers are one byte?
	if strings.Index(cType, "(") >= 0 {
		return 1, nil
	}

	if strings.HasSuffix(cType, "*") {
		return pointerSize, nil
	}

	switch cType {
	case "char", "void":
		return 1, nil

	case "short":
		return 2, nil

	case "int", "float":
		return 4, nil

	case "long", "double":
		return 8, nil

	case "long double":
		return 16, nil

	default:
		return pointerSize, errors.New(
			fmt.Sprintf("cannot determine size of: %s", cType))
	}
}

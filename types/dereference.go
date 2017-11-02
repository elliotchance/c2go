package types

import (
	"errors"
	"regexp"
	"strings"
)

// GetDereferenceType returns the C type that would be the result of
// dereferencing (unary "*" operator or accessing a single array element on a
// pointer) a value.
//
// For example if the input type is "char *", then dereferencing or accessing a
// single element would result in a "char".
//
// If the dereferenced type cannot be determined or is impossible ("char" cannot
// be dereferenced, for example) then an error is returned.
func GetDereferenceType(cType string) (string, error) {
	// In the form of: "int [2][3][4]" -> "int [3][4]"
	search := regexp.MustCompile(`([\w ]+)\s*\[\d+\]((\[\d+\])+)`).FindStringSubmatch(cType)
	if len(search) > 0 {
		return search[1] + search[2], nil
	}

	// In the form of: "char [8]" -> "char"
	search = regexp.MustCompile(`([\w ]+)\s*\[\d+\]`).FindStringSubmatch(cType)
	if len(search) > 0 {
		return strings.TrimSpace(search[1]), nil
	}

	// In the form of: "char **" -> "char *"
	search = regexp.MustCompile(`([\w ]+)\s*(\*+)`).FindStringSubmatch(cType)
	if len(search) > 0 {
		return strings.TrimSpace(search[1] + search[2][0:len(search[2])-1]), nil
	}

	return "", errors.New(cType)
}

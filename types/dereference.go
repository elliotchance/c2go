package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/util"
)

// IsDereferenceType - check is that type dereference
func IsDereferenceType(cType string) bool {
	return strings.ContainsAny(cType, "[]*")
}

// GetDereferenceType returns the C type that would be the result of
// dereferencing (unary "*" operator or accessing a single array element on a
// pointer) a value.
//
// For example if the input type is "char *", then dereferencing or accessing a
// single element would result in a "char".
//
// If the dereferenced type cannot be determined or is impossible ("char" cannot
// be dereferenced, for example) then an error is returned.
func GetDereferenceType(cType string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Error in GetDereferenceType : %v", err)
		}
	}()

	// In the form of: "int [2][3][4]" -> "int [3][4]"
	search := util.GetRegex(`([\w\* ]+)\s*\[\d+\]((\[\d+\])+)`).FindStringSubmatch(cType)
	if len(search) > 0 {
		return search[1] + search[2], nil
	}

	// In the form of: "char [8]" -> "char"
	search = util.GetRegex(`([\w\* ]+)\s*\[\d+\]`).FindStringSubmatch(cType)
	if len(search) > 0 {
		return strings.TrimSpace(search[1]), nil
	}

	// In the form of: "char **" -> "char *"
	search = util.GetRegex(`([\w ]+)\s*(\*+)`).FindStringSubmatch(cType)
	if len(search) > 0 {
		return strings.TrimSpace(search[1] + search[2][0:len(search[2])-1]), nil
	}

	return "", errors.New(cType)
}

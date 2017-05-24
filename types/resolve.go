package types

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/elliotchance/c2go/program"
)

// TODO: Some of these are based on assumptions that may not be true for all
// architectures (like the size of an int). At some point in the future we will
// need to find out the sizes of some of there and pick the most compatible
// type.
//
// Please keep them sorted by name.
var simpleResolveTypes = map[string]string{
	"bool":               "bool",
	"char *":             "[]byte",
	"char":               "byte",
	"char*":              "[]byte",
	"double":             "float64",
	"float":              "float32",
	"int":                "int",
	"long double":        "float64",
	"long int":           "int32",
	"long long":          "int64",
	"long unsigned int":  "uint32",
	"long":               "int32",
	"short":              "int16",
	"signed char":        "int8",
	"unsigned char":      "uint8",
	"unsigned int":       "uint32",
	"unsigned long long": "uint64",
	"unsigned long":      "uint32",
	"unsigned short":     "uint16",
	"unsigned short int": "uint16",
	"void":               "",
	"_Bool":              "bool",

	// void* is treated like char*
	"void*":  "[]byte",
	"void *": "[]byte",

	// null is a special case (it should probably have a less ambiguos name)
	// when using the NULL macro.
	"null": "null",

	// Are these built into some compilers?
	"uint32":     "uint32",
	"uint64":     "uint64",
	"__uint16_t": "uint16",
	"__uint32_t": "uint32",
	"__uint64_t": "uint64",

	// Darwin specific
	"__darwin_ct_rune_t": "github.com/elliotchance/c2go/darwin.CtRuneT",
	"fpos_t":             "int",
	"struct __float2":    "github.com/elliotchance/c2go/darwin.Float2",
	"struct __double2":   "github.com/elliotchance/c2go/darwin.Double2",
	"Float2":             "github.com/elliotchance/c2go/darwin.Float2",
	"Double2":            "github.com/elliotchance/c2go/darwin.Double2",

	// These are special cases that almost certainly don't work. I've put
	// them here because for whatever reason there is no suitable type or we
	// don't need these platform specific things to be implemented yet.
	"__builtin_va_list":            "int64",
	"__darwin_pthread_handler_rec": "int64",
	"unsigned __int128":            "uint64",
	"__int128":                     "int64",
	"__mbstate_t":                  "int64",
	"__sbuf":                       "int64",
	"__sFILEX":                     "interface{}",
	"__va_list_tag":                "interface{}",
	"FILE":                         "github.com/elliotchance/c2go/noarch.File",
}

func ResolveType(p *program.Program, s string) (string, error) {
	// Remove any whitespace or attributes that are not relevant to Go.
	s = strings.Replace(s, "const ", "", -1)
	s = strings.Replace(s, "volatile ", "", -1)
	s = strings.Replace(s, "*__restrict", "*", -1)
	s = strings.Replace(s, "*restrict", "*", -1)
	s = strings.Replace(s, "*const", "*", -1)
	s = strings.Trim(s, " \t\n\r")

	// TODO: Unions are not supported.
	// https://github.com/elliotchance/c2go/issues/84
	//
	// For now we will let them be interface{} so that it does not stop the
	// transpilation.
	if strings.HasPrefix(s, "union ") {
		return "interface{}", errors.New("unions are not supported")
	}

	// FIXME: This is a hack to avoid casting in some situations.
	if s == "" {
		return s, errors.New("probably an incorrect type translation 1")
	}

	// FIXME: I have no idea what this is.
	if s == "const" {
		return "interface{}", errors.New("probably an incorrect type translation 4")
	}

	if s == "char *[]" {
		return "interface{}", errors.New("probably an incorrect type translation 2")
	}

	if s == "fpos_t" {
		return "int", nil
	}

	// The simple resolve types are the types that we know there is an exact Go
	// equivalent. For example float, int, etc.
	for k, v := range simpleResolveTypes {
		if k == s {
			return p.ImportType(v), nil
		}
	}

	// If the type is already defined we can proceed with the same name.
	if p.TypeIsAlreadyDefined(s) {
		return p.ImportType(s), nil
	}

	// Structures are by name.
	if strings.HasPrefix(s, "struct ") {
		if s[len(s)-1] == '*' {
			s = s[7 : len(s)-2]

			for _, v := range simpleResolveTypes {
				if v == s {
					return "*" + p.ImportType(simpleResolveTypes[s]), nil
				}
			}

			return "*" + s, nil
		}

		s = s[7:]

		for _, v := range simpleResolveTypes {
			if v == s {
				return p.ImportType(simpleResolveTypes[s]), nil
			}
		}

		return s, nil
	}

	// Enums are by name.
	if strings.HasPrefix(s, "enum ") {
		if s[len(s)-1] == '*' {
			return "*" + s[5:len(s)-2], nil
		}

		return s[5:], nil
	}

	// I have no idea how to handle this yet.
	if strings.Index(s, "anonymous union") != -1 {
		return "interface{}", errors.New("probably an incorrect type translation 3")
	}

	// It may be a pointer of a simple type. For example, float *, int *,
	// etc.
	if regexp.MustCompile("[\\w ]+\\*+$").MatchString(s) {
		// The "-1" is important because there may or may not be a space between
		// the name and the "*". If there is an extra space it will be trimmed
		// off.
		t, err := ResolveType(p, strings.TrimSpace(s[:len(s)-1]))

		// Pointers are always converted into slices, except with some specific
		// entities that are shared in the Go libraries.
		prefix := "*"
		if !strings.Contains(t, "noarch.") {
			prefix = "[]"
		}

		return prefix + t, err
	}

	if regexp.MustCompile(`[\w ]+\*\[\d+\]$`).MatchString(s) {
		return "[]string", nil
	}

	// Function pointers are not yet supported. In the mean time they will be
	// replaced with a type that certainly wont work until we can fix this
	// properly.
	search := regexp.MustCompile("[\\w ]+\\(\\*.*?\\)\\(.*\\)").MatchString(s)
	if search {
		return "interface{}", errors.New("function pointers are not supported")
	}

	search = regexp.MustCompile("[\\w ]+ \\(.*\\)").MatchString(s)
	if search {
		return "interface{}", errors.New("function pointers are not supported")
	}

	// It could be an array of fixed length. These needs to be converted to
	// slices.
	search2 := regexp.MustCompile("([\\w ]+)\\[(\\d+)\\]").FindStringSubmatch(s)
	if len(search2) > 0 {
		t, err := ResolveType(p, search2[1])
		return fmt.Sprintf("[]%s", t), err
	}

	errMsg := fmt.Sprintf(
		"I couldn't find an appropriate Go type for the C type '%s'.", s)
	return "interface{}", errors.New(errMsg)
}

package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// TODO: Some of these are based on assumptions that may not be true for all
// architectures (like the size of an int). At some point in the future we will
// need to find out the sizes of some of there and pick the most compatible
// type.
//
// Please keep them sorted by name.
var simpleResolveTypes = map[string]string{
	"bool":               "bool",
	"char *":             "string",
	"char":               "byte",
	"char*":              "string",
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
	"void *":             "interface{}",
	"void":               "",

	"const char *": "string",

	// Are these built into some compilers?
	"uint32": "uint32",

	// Darwin specific
	"__darwin_ct_rune_t": "github.com/elliotchance/c2go/darwin.Darwin_ct_rune_t",
	"union __mbstate_t":  "__mbstate_t",
	"fpos_t":             "int",
	"struct __float2":    "github.com/elliotchance/c2go/darwin.Float2",
	"struct __double2":   "github.com/elliotchance/c2go/darwin.Double2",

	// These are special cases that almost certainly don"t work. I've put
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
	"FILE":                         "int64",
}

var TypesAlreadyDefined = []string{
	// Linux specific.
	"_LIB_VERSION_TYPE",

	// Darwin specific.
	"__float2",
	"__double2",
}

func typeIsAlreadyDefined(typeName string) bool {
	return inStrings(typeName, TypesAlreadyDefined)
}

func typeIsNowDefined(typeName string) {
	TypesAlreadyDefined = append(TypesAlreadyDefined, typeName)
}

func resolveType(s string) string {
	// Remove any whitespace or attributes that are not relevant to Go.
	s = strings.Replace(s, "const ", "", -1)
	s = strings.Replace(s, "*__restrict", "*", -1)
	s = strings.Replace(s, "*restrict", "*", -1)
	s = strings.Trim(s, " \t\n\r")

	if s == "fpos_t" {
		return "int"
	}

	// The simple resolve types are the types that we know there is an exact Go
	// equivalent. For example float, int, etc.
	for k, v := range simpleResolveTypes {
		if k == s {
			return importType(v)
		}
	}

	// If the type is already defined we can proceed with the same name.
	for _, v := range TypesAlreadyDefined {
		if v == s {
			return importType(s)
		}
	}

	// Structures are by name.
	if strings.HasPrefix(s, "struct ") {
		if s[len(s)-1] == '*' {
			s = s[7 : len(s)-2]

			for _, v := range simpleResolveTypes {
				if v == s {
					return "*" + importType(simpleResolveTypes[s])
				}
			}

			return "*" + s
		} else {
			s = s[7:]

			for _, v := range simpleResolveTypes {
				if v == s {
					return importType(simpleResolveTypes[s])
				}
			}

			return s
		}
	}

	// Enums are by name.
	if strings.HasPrefix(s, "enum ") {
		if s[len(s)-1] == '*' {
			return "*" + s[5:len(s)-2]
		} else {
			return s[5:]
		}
	}

	// I have no idea how to handle this yet.
	if strings.Index(s, "anonymous union") != -1 {
		return "interface{}"
	}

	// It may be a pointer of a simple type. For example, float *, int *,
	// etc.
	if regexp.MustCompile("[\\w ]+\\*+$").MatchString(s) {
		return "*" + resolveType(strings.TrimSpace(s[:len(s)-2]))
	}

	// Function pointers are not yet supported. In th mean time they will be
	// replaced with a type that certainly wont work until we can fix this
	// properly.
	search := regexp.MustCompile("[\\w ]+\\(\\*.*?\\)\\(.*\\)").MatchString(s)
	if search {
		return "interface{}"
	}

	search = regexp.MustCompile("[\\w ]+ \\(.*\\)").MatchString(s)
	if search {
		return "interface{}"
	}

	// It could be an array of fixed length.
	search2 := regexp.MustCompile("([\\w ]+)\\[(\\d+)\\]").FindStringSubmatch(s)
	if len(search2) > 0 {
		return fmt.Sprintf("[%s]%s", search2[2], resolveType(search2[1]))
	}

	panic(fmt.Sprintf("Oh no! I couldn't find an appropriate the Go type for the C type '%s'.", s))
}

func getDereferenceType(cType string) (string, error) {
	// In the form of: "char [8]" -> "char"
	search := regexp.MustCompile(`([\w ]+)\s*\[\d+\]`).FindStringSubmatch(cType)
	if len(search) > 0 {
		return search[1], nil
	}

	// In the form of: "char **" -> "char *"
	search = regexp.MustCompile(`([\w ]+)\s*(\*+)`).FindStringSubmatch(cType)
	if len(search) > 0 {
		return search[1] + search[2][0:len(search[2])-1], nil
	}

	// I give up...
	return "unknown5", errors.New(cType)
}

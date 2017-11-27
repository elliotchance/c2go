package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

// TODO: Some of these are based on assumptions that may not be true for all
// architectures (like the size of an int). At some point in the future we will
// need to find out the sizes of some of there and pick the most compatible
// type.
//
// Please keep them sorted by name.
var simpleResolveTypes = map[string]string{
	"bool":                   "bool",
	"char *":                 "[]byte",
	"char":                   "byte",
	"char*":                  "[]byte",
	"double":                 "float64",
	"float":                  "float32",
	"int":                    "int",
	"long double":            "float64",
	"long int":               "int32",
	"long long":              "int64",
	"long long int":          "int64",
	"long long unsigned int": "uint64",
	"long unsigned int":      "uint32",
	"long":                   "int32",
	"short":                  "int16",
	"signed char":            "int8",
	"unsigned char":          "uint8",
	"unsigned int":           "uint32",
	"unsigned long long":     "uint64",
	"unsigned long":          "uint32",
	"unsigned short":         "uint16",
	"unsigned short int":     "uint16",
	"void":                   "",
	"_Bool":                  "bool",

	// void*
	"void*":  "interface{}",
	"void *": "interface{}",

	// null is a special case (it should probably have a less ambiguos name)
	// when using the NULL macro.
	"null": "null",

	// Non platform-specific types.
	"uint32":     "uint32",
	"uint64":     "uint64",
	"__uint16_t": "uint16",
	"__uint32_t": "uint32",
	"__uint64_t": "uint64",
	"div_t":      "github.com/elliotchance/c2go/noarch.DivT",
	"ldiv_t":     "github.com/elliotchance/c2go/noarch.LdivT",
	"lldiv_t":    "github.com/elliotchance/c2go/noarch.LldivT",
	"time_t":     "github.com/elliotchance/c2go/noarch.TimeT",

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
	"FILE":                         "github.com/elliotchance/c2go/noarch.File",
}

// ResolveType determines the Go type from a C type.
//
// Some basic examples are obvious, such as "float" in C would be "float32" in
// Go. But there are also much more complicated examples, such as compound types
// (structs and unions) and function pointers.
//
// Some general rules:
//
// 1. The Go type must be deterministic. The same C type will ALWAYS return the
//    same Go type, in any condition. This is extremely important since the
//    nature of C is that is may not have certain information available about the
//    rest of the program or libraries when it is being compiled.
//
// 2. Many C type modifiers and properties are lost as they have no sensible or
//    valid translation to Go. Some example of those would be "const" and
//    "volatile". It is left be up to the clang (or other compiler) to warn if
//    types are being abused against the standards in which they are being
//    compiled under. Go will make no assumptions about how you expect it act,
//    only how it is used.
//
// 3. New types are registered (discovered) throughout the transpiling of the
//    program, so not all types are know at any given time. This works exactly
//    the same way in a C compiler that will not let you use a type before it
//    has been defined.
//
// 4. If all else fails an error is returned. However, a type (which is almost
//    certainly incorrect) "interface{}" is also returned. This is to allow the
//    transpiler to step over type errors and put something as a placeholder
//    until a more suitable solution is found for those cases.
func ResolveType(p *program.Program, s string) (string, error) {
	s = CleanCType(s)

	// FIXME: This is a hack to avoid casting in some situations.
	if s == "" {
		return "interface{}", errors.New("probably an incorrect type translation 1")
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

	// Check is it typedef enum
	if _, ok := p.EnumTypedefName[s]; ok {
		return ResolveType(p, "int")
	}

	// If the type is already defined we can proceed with the same name.
	if p.IsTypeAlreadyDefined(s) {
		return p.ImportType(s), nil
	}

	// Structures are by name.
	if strings.HasPrefix(s, "struct ") || strings.HasPrefix(s, "union ") {
		start := 6
		if s[0] == 's' {
			start++
		}

		if s[len(s)-1] == '*' {
			s = s[start : len(s)-2]

			for _, v := range simpleResolveTypes {
				if v == s {
					return "[]" + p.ImportType(simpleResolveTypes[s]), nil
				}
			}

			return "[]" + strings.TrimSpace(s), nil
		}

		s = s[start:]

		for _, v := range simpleResolveTypes {
			if v == s {
				return p.ImportType(simpleResolveTypes[s]), nil
			}
		}

		return ResolveType(p, s)
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
	if util.GetRegex("[\\w ]+\\*+$").MatchString(s) {
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

	// Function pointers are not yet supported. In the mean time they will be
	// replaced with a type that certainly wont work until we can fix this
	// properly.
	search := util.GetRegex("[\\w ]+\\(\\*.*?\\)\\(.*\\)").MatchString(s)
	if search {
		return "interface{}", errors.New("function pointers are not supported")
	}

	search = util.GetRegex("[\\w ]+ \\(.*\\)").MatchString(s)
	if search {
		return "interface{}", errors.New("function pointers are not supported")
	}

	// It could be an array of fixed length. These needs to be converted to
	// slices.
	// int [2][3] -> [][]int
	// int [2][3][4] -> [][][]int
	search2 := util.GetRegex(`([\w\* ]+)((\[\d+\])+)`).FindStringSubmatch(s)
	if len(search2) > 2 {
		t, err := ResolveType(p, search2[1])

		var re = util.GetRegex(`[0-9]+`)
		arraysNoSize := re.ReplaceAllString(search2[2], "")

		return fmt.Sprintf("%s%s", arraysNoSize, t), err
	}

	errMsg := fmt.Sprintf(
		"I couldn't find an appropriate Go type for the C type '%s'.", s)
	return "interface{}", errors.New(errMsg)
}

// ResolveFunction determines the Go type from a C type.
func ResolveFunction(p *program.Program, s string) (fields []string, returns []string, err error) {
	f, r, err := ParseFunction(s)
	if err != nil {
		return
	}
	for i := range f {
		var t string
		t, err = ResolveType(p, f[i])
		if err != nil {
			return
		}
		fields = append(fields, t)
	}
	for i := range r {
		var t string
		t, err = ResolveType(p, r[i])
		if err != nil {
			return
		}
		returns = append(returns, t)
	}
	return
}

// IsFunction - return true if string is function like "void (*)(void)"
func IsFunction(s string) bool {
	parts := strings.Split(s, "(*)")
	if len(parts) != 2 {
		return false
	}
	inside := strings.TrimSpace(parts[1])
	if inside[0] != '(' || inside[len(inside)-1] != ')' {
		return false
	}
	return true
}

// ParseFunction - parsing elements of C function
func ParseFunction(s string) (f []string, r []string, err error) {
	if !IsFunction(s) {
		err = fmt.Errorf("Is not function : %s", s)
		return
	}
	i := strings.Index(s, "(")
	if i == -1 {
		err = fmt.Errorf("Cannot parse (index of function): %v", s)
		return
	}
	r = append(r, s[0:i])
	parts := strings.Split(s, "(*)")
	if len(parts) != 2 {
		err = fmt.Errorf("Cannot parse (separation on parts) : %v", s)
		return
	}
	inside := strings.TrimSpace(parts[1])
	if inside == "" {
		err = fmt.Errorf("Cannot parse (right part is nil) : %v", s)
		return
	}
	f = append(f, strings.Split(inside[1:len(inside)-1], ",")...)

	for i := range r {
		r[i] = strings.TrimSpace(r[i])
	}
	for i := range f {
		f[i] = strings.TrimSpace(f[i])
	}

	return
}

// CleanCType - remove from C type not Go type
func CleanCType(s string) (out string) {
	out = s
	out = strings.Replace(out, "()", "", -1)
	out = strings.Replace(out, "(*)", "", -1)

	// Remove any whitespace or attributes that are not relevant to Go.
	out = strings.Replace(out, "const", "", -1)
	out = strings.Replace(out, "volatile", "", -1)
	out = strings.Replace(out, "__restrict", "", -1)
	out = strings.Replace(out, "restrict", "", -1)
	out = strings.Replace(out, "\t", "", -1)
	out = strings.Replace(out, "\n", "", -1)
	out = strings.Replace(out, "\r", "", -1)

	// remove space from pointer symbols
	out = strings.Replace(out, "* *", "**", -1)

	// remove addition spaces
	out = strings.Replace(out, "  ", " ", -1)

	// remove spaces around
	out = strings.TrimSpace(out)

	if out != s {
		return CleanCType(out)
	}

	return out
}

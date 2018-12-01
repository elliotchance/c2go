package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

// cIntegerType - slice of C integer type
var cIntegerType = []string{
	"int",
	"long long",
	"long long int",
	"long long unsigned int",
	"long unsigned int",
	"long",
	"short",
	"unsigned int",
	"unsigned long long",
	"unsigned long",
	"unsigned short",
	"unsigned short int",
}

// IsCInteger - return true is C type integer
func IsCInteger(p *program.Program, cType string) bool {
	for i := range cIntegerType {
		if cType == cIntegerType[i] {
			return true
		}
	}
	if rt, ok := p.TypedefType[cType]; ok {
		return IsCInteger(p, rt)
	}
	return false
}

// TODO: Some of these are based on assumptions that may not be true for all
// architectures (like the size of an int). At some point in the future we will
// need to find out the sizes of some of there and pick the most compatible
// type.
//
// Please keep them sorted by name.
var simpleResolveTypes = map[string]string{
	"bool":                   "bool",
	"char *":                 "*byte",
	"char":                   "byte",
	"char*":                  "*byte",
	"double":                 "float64",
	"float":                  "float32",
	"int":                    "int32",
	"long double":            "float64",
	"long int":               "int32",
	"long long":              "int64",
	"long long int":          "int64",
	"long long unsigned int": "uint64",
	"long unsigned int":      "uint32",
	"long":                   "int32",
	"ptrdiff_t":              "int32",
	"short":                  "int16",
	"signed char":            "int8",
	"size_t":                 "uint32",
	"uintptr_t":              "uintptr",
	"unsigned char":          "uint8",
	"unsigned int":           "uint32",
	"unsigned long long":     "uint64",
	"unsigned long":          "uint32",
	"unsigned short":         "uint16",
	"unsigned short int":     "uint16",
	"void":                   "",
	"wchar_t":                "rune",
	"_Bool":                  "int8",

	// void*
	"void*":  "unsafe.Pointer",
	"void *": "unsafe.Pointer",

	// null is a special case (it should probably have a less ambiguous name)
	// when using the NULL macro.
	"null": "null",

	// Non platform-specific types.
	"int8":       "int8",
	"int8_t":     "int8",
	"int16":      "int16",
	"int16_t":    "int16",
	"int32":      "int32",
	"int32_t":    "int32",
	"int64":      "int64",
	"int64_t":    "int64",
	"uint8":      "uint8",
	"uint8_t":    "uint8",
	"uint16":     "uint16",
	"uint16_t":   "uint16",
	"uint32":     "uint32",
	"uint32_t":   "uint32",
	"uint64":     "uint64",
	"uint64_t":   "uint64",
	"u_int8_t":   "uint8",
	"u_int16_t":  "uint16",
	"u_int32_t":  "uint32",
	"u_int64_t":  "uint64",
	"__int8_t":   "int8",
	"__int16_t":  "int16",
	"__int32_t":  "int32",
	"__int64_t":  "int64",
	"__uint8_t":  "uint8",
	"__uint16_t": "uint16",
	"__uint32_t": "uint32",
	"__uint64_t": "uint64",

	// These are special cases that almost certainly don't work. I've put
	// them here because for whatever reason there is no suitable type or we
	// don't need these platform specific things to be implemented yet.
	"__builtin_va_list":            "int64",
	"__darwin_pthread_handler_rec": "int64",
	"unsigned __int128":            "uint64",
	"__int128":                     "int64",
	"__mbstate_t":                  "int64",
	"__sbuf":                       "int64",
	"__sFILEX":                     "unsafe.Pointer",
	"FILE":                         "github.com/elliotchance/c2go/noarch.File",
}

var otherStructType = map[string]string{
	"div_t":   "github.com/elliotchance/c2go/noarch.DivT",
	"ldiv_t":  "github.com/elliotchance/c2go/noarch.LdivT",
	"lldiv_t": "github.com/elliotchance/c2go/noarch.LldivT",

	// time.h
	"tm":        "github.com/elliotchance/c2go/noarch.Tm",
	"struct tm": "github.com/elliotchance/c2go/noarch.Tm",
	"time_t":    "github.com/elliotchance/c2go/noarch.TimeT",

	// Darwin specific
	"__darwin_ct_rune_t":     "github.com/elliotchance/c2go/darwin.CtRuneT",
	"fpos_t":                 "int32",
	"struct __float2":        "github.com/elliotchance/c2go/darwin.Float2",
	"struct __double2":       "github.com/elliotchance/c2go/darwin.Double2",
	"Float2":                 "github.com/elliotchance/c2go/darwin.Float2",
	"Double2":                "github.com/elliotchance/c2go/darwin.Double2",
	"struct __va_list_tag *": "github.com/elliotchance/c2go/noarch.VaList",
}

// NullPointer - is look : (double *)(nil) or (FILE *)(nil)
// created only for transpiler.CStyleCastExpr
var NullPointer = "NullPointerType *"

// ToVoid - specific type for ignore the cast
var ToVoid = "ToVoid"

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
func ResolveType(p *program.Program, s string) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot resolve type '%s' : %v", s, err)
		}
	}()
	s = CleanCType(s)

	if s == "_Bool" {
		p.TypedefType[s] = "signed char"
	}

	// FIXME: This is a hack to avoid casting in some situations.
	if s == "" {
		return "unsafe.Pointer", errors.New("probably an incorrect type translation 1")
	}

	// FIXME: I have no idea what this is.
	if s == "const" {
		return "unsafe.Pointer", errors.New("probably an incorrect type translation 4")
	}

	if s == "char *[]" {
		return "unsafe.Pointer", errors.New("probably an incorrect type translation 2")
	}

	if s == "fpos_t" {
		return ResolveType(p, "int")
	}

	// FIXME: I have no idea, how to solve.
	// See : https://github.com/elliotchance/c2go/issues/628
	if strings.Contains(s, "__locale_data") {
		s = strings.Replace(s, "struct __locale_data", "int", -1)
		s = strings.Replace(s, "__locale_data", "int", -1)
	}
	if strings.Contains(s, "__locale_struct") {
		return ResolveType(p, "int")
	}

	// function type is pointer in Go by default
	if len(s) > 2 {
		base := s[:len(s)-2]
		if ff, ok := p.TypedefType[base]; ok {
			if IsFunction(ff) {
				return base, nil
			}
		}
	}

	// The simple resolve types are the types that we know there is an exact Go
	// equivalent. For example float, int, etc.
	if v, ok := simpleResolveTypes[s]; ok {
		return p.ImportType(v), nil
	}

	// No need resolve typedef types
	if ss, ok := p.TypedefType[s]; ok {
		if tt, ok := otherStructType[s]; ok {
			// "div_t":   "github.com/elliotchance/c2go/noarch.DivT",
			ii := p.ImportType(tt)
			return ii, nil
		}
		if s == "_Bool" {
			return ResolveType(p, ss)
		}
		return s, nil
	}

	if tt, ok := otherStructType[s]; ok {
		// "div_t":   "github.com/elliotchance/c2go/noarch.DivT",
		ii := p.ImportType(tt)
		return ii, nil
	}

	// Try resolving correct type (covers anonymous structs/unions)
	correctType := GenerateCorrectType(s)
	if correctType != s {
		if res, err := ResolveType(p, correctType); err == nil {
			return res, nil
		}
	}

	// For function
	if IsFunction(s) {
		g, e := resolveFunction(p, s)
		return g, e
	}

	// Check is it typedef enum
	if _, ok := p.EnumTypedefName[s]; ok {
		return ResolveType(p, "int")
	}

	if v, ok := p.TypedefType[s]; ok {
		if IsFunction(v) {
			// typedef function
			return s, nil
		}
		return ResolveType(p, v)
	}

	// If the type is already defined we can proceed with the same name.
	if p.IsTypeAlreadyDefined(s) {
		return p.ImportType(s), nil
	}

	if _, ok := p.Structs["struct "+s]; ok {
		return p.ImportType(s), nil
	}

	if _, ok := p.Unions["union "+s]; ok {
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

			var t string
			t, err = ResolveType(p, s)
			return "*" + t, err
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
	if strings.Contains(s, "anonymous union") {
		return "unsafe.Pointer", errors.New("probably an incorrect type translation 3")
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

		return prefix + t, err
	}

	// Function pointers are not yet supported. In the mean time they will be
	// replaced with a type that certainly wont work until we can fix this
	// properly.
	search := util.GetRegex("[\\w ]+\\(\\*.*?\\)\\(.*\\)").MatchString(s)
	if search {
		return "interface{}",
			fmt.Errorf("function pointers are not supported [1] : '%s'", s)
	}

	search = util.GetRegex("[\\w ]+ \\(.*\\)").MatchString(s)
	if search {
		return "interface{}",
			fmt.Errorf("function pointers are not supported [2] : '%s'", s)
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
	return "unsafe.Pointer", errors.New(errMsg)
}

// resolveType determines the Go type from a C type.
func resolveFunction(p *program.Program, s string) (goType string, err error) {
	var f, r []string
	f, r, err = SeparateFunction(p, s)
	goType = "func("
	for i := range f {
		goType += fmt.Sprintf("%s", f[i])
		if i < len(f)-1 {
			goType += " , "
		}
	}
	goType += ")("
	for i := range r {
		goType += fmt.Sprintf("%s", r[i])
		if i < len(r)-1 {
			goType += " , "
		}
	}
	goType += ")"
	return
}

// SeparateFunction separate a function C type to Go types parts.
func SeparateFunction(p *program.Program, s string) (
	fields []string, returns []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot separate function '%s' : %v", s, err)
		}
	}()
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
	s = strings.Replace(s, "(*)", "", -1)
	return strings.Contains(s, "(")
}

// IsPointer - check type is pointer
func IsPointer(p *program.Program, s string) bool {
	if strings.ContainsAny(s, "*[]") {
		return true
	}
	if v, ok := p.TypedefType[s]; ok {
		return IsPointer(p, v)
	}
	return false
}

// IsPurePointer - check type is pointer
func IsPurePointer(p *program.Program, s string) bool {
	if strings.ContainsAny(s, "*") {
		return true
	}
	if v, ok := p.TypedefType[s]; ok {
		return IsPurePointer(p, v)
	}
	return false
}

// IsLastArray - check type have array '[]'
func IsLastArray(s string) bool {
	for _, b := range s {
		switch b {
		case '[':
			return true
		case '*':
			break
		}
	}
	return false
}

// IsTypedefFunction - return true if that type is typedef of function.
func IsTypedefFunction(p *program.Program, s string) bool {
	if v, ok := p.TypedefType[s]; ok && IsFunction(v) {
		return true
	}
	s = string(s[0 : len(s)-len(" *")])
	if v, ok := p.TypedefType[s]; ok && IsFunction(v) {
		return true
	}
	return false
}

// ParseFunction - parsing elements of C function
func ParseFunction(s string) (f []string, r []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot parse function '%s' : %v", s, err)
		}
	}()

	s = strings.TrimSpace(s)
	if !IsFunction(s) {
		err = fmt.Errorf("Is not function : %s", s)
		return
	}
	var arguments string
	{
		// Example of function types :
		// int (*)(int, float)
		// int (int, float)
		// int (*)(int (*)(int))
		// void (*(*)(int *, void *, const char *))(void)
		if s[len(s)-1] != ')' {
			err = fmt.Errorf("function type |%s| haven't last symbol ')'", s)
			return
		}
		counter := 1
		var pos int
		for i := len(s) - 2; i >= 0; i-- {
			if i == 0 {
				err = fmt.Errorf("Don't found '(' in type : %s", s)
				return
			}
			if s[i] == ')' {
				counter++
			}
			if s[i] == '(' {
				counter--
			}
			if counter == 0 {
				pos = i
				break
			}
		}
		r = append(r, strings.Replace(s[:pos], "(*)", "", -1))
		arguments = s[pos:]
	}
	arguments = strings.TrimSpace(arguments)
	if arguments == "" {
		err = fmt.Errorf("Cannot parse (right part is nil) : %v", s)
		return
	}
	// separate fields of arguments
	{
		var pos int = 1
		counter := 0
		for i := 1; i < len(arguments)-1; i++ {
			if arguments[i] == '(' {
				counter++
			}
			if arguments[i] == ')' {
				counter--
			}
			if counter == 0 && arguments[i] == ',' {
				f = append(f, strings.TrimSpace(arguments[pos:i]))
				pos = i + 1
			}
		}
		f = append(f, strings.TrimSpace(arguments[pos:len(arguments)-1]))
	}

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

	// remove space from pointer symbols
	out = strings.Replace(out, "* *", "**", -1)

	// add space for simplification redactoring
	out = strings.Replace(out, "*", " *", -1)

	out = strings.Replace(out, "( *)", "(*)", -1)

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
	out = strings.Replace(out, "[", " [", -1)
	out = strings.Replace(out, "] [", "][", -1)

	// remove addition spaces
	out = strings.Replace(out, "  ", " ", -1)

	// remove spaces around
	out = strings.TrimSpace(out)

	if out != s {
		return CleanCType(out)
	}

	return out
}

// GenerateCorrectType - generate correct type
// Example: 'union (anonymous union at tests/union.c:46:3)'
func GenerateCorrectType(name string) string {
	if !strings.Contains(name, "anonymous") {
		return name
	}
	index := strings.Index(name, "(anonymous")
	if index < 0 {
		return name
	}
	name = strings.Replace(name, "anonymous", "", 1)
	var last int
	for last = index; last < len(name); last++ {
		if name[last] == ')' {
			break
		}
	}

	// Create a string, for example:
	// Input (name)   : 'union (anonymous union at tests/union.c:46:3)'
	// Output(inside) : '(anonymous union at tests/union.c:46:3)'
	inside := string(([]byte(name))[index : last+1])

	// change unacceptable C name letters
	inside = strings.Replace(inside, "(", "B", -1)
	inside = strings.Replace(inside, ")", "E", -1)
	inside = strings.Replace(inside, " ", "S", -1)
	inside = strings.Replace(inside, ":", "D", -1)
	inside = strings.Replace(inside, "/", "S", -1)
	inside = strings.Replace(inside, "-", "T", -1)
	inside = strings.Replace(inside, "\\", "S", -1)
	inside = strings.Replace(inside, ".", "P", -1)
	out := string(([]byte(name))[0:index]) + inside + string(([]byte(name))[last+1:])

	// For case:
	// struct siginfo_t::(anonymous at /usr/include/x86_64-linux-gnu/bits/siginfo.h:119:2)
	// we see '::' before 'anonymous' word
	out = strings.Replace(out, ":", "D", -1)

	return CleanCType(out)
}

// GetAmountArraySize - return amount array size
// Example :
// In  : 'char [40]'
// Out : 40
func GetAmountArraySize(cType string) (size int, err error) {
	reg := util.GetRegex("\\[(?P<size>\\d+)\\]")
	match := reg.FindStringSubmatch(cType)

	if reg.NumSubexp() != 1 {
		err = fmt.Errorf("Cannot found size of array in type : %s", cType)
		return
	}

	result := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return strconv.Atoi(result["size"])
}

// GetBaseType - return base type without pointera, array symbols
// Input:
// s =  struct BSstructSatSShomeSlepriconSgoSsrcSgithubPcomSelliotchanceSc2goStestsSstructPcD260D18E [7]
func GetBaseType(s string) string {
	s = strings.TrimSpace(s)
	s = CleanCType(s)
	if s[len(s)-1] == ']' {
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] == '[' {
				s = s[:i]
				return GetBaseType(s)
			}
		}
	}
	if s[len(s)-1] == '*' {
		return GetBaseType(s[:len(s)-1])
	}
	return s
}

// goIntegerType - slice of Go integer type
var goIntegerType = []string{
	"int8",
	"int16",
	"int32",
	"int64",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
}

// IsGoIntegerType - return whether the given Go type is an integer type
func IsGoIntegerType(s string) bool {
	for _, i := range goIntegerType {
		if i == s {
			return true
		}
	}
	return false
}

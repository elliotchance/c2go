package main

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var FunctionSubstitutions = map[string]string{
	// math.h
	"acos":  "math.Acos",
	"asin":  "math.Asin",
	"atan":  "math.Atan",
	"atan2": "math.Atan2",
	"ceil":  "math.Ceil",
	"cos":   "math.Cos",
	"cosh":  "math.Cosh",
	"exp":   "math.Exp",
	"fabs":  "math.Abs",
	"floor": "math.Floor",
	"fmod":  "math.Mod",
	"ldexp": "math.Ldexp",
	"log":   "math.Log",
	"log10": "math.Log10",
	"pow":   "math.Pow",
	"sin":   "math.Sin",
	"sinh":  "math.Sinh",
	"sqrt":  "math.Sqrt",
	"tan":   "math.Tan",
	"tanh":  "math.Tanh",

	// stdio
	"printf": "fmt.Printf",
	"scanf":  "fmt.Scanf",

	// darwin/math.h
	"__builtin_fabs":    "github.com/elliotchance/c2go/darwin.Fabs",
	"__builtin_fabsf":   "github.com/elliotchance/c2go/darwin.Fabsf",
	"__builtin_fabsl":   "github.com/elliotchance/c2go/darwin.Fabsl",
	"__builtin_inf":     "github.com/elliotchance/c2go/darwin.Inf",
	"__builtin_inff":    "github.com/elliotchance/c2go/darwin.Inff",
	"__builtin_infl":    "github.com/elliotchance/c2go/darwin.Infl",
	"__sincospi_stret":  "github.com/elliotchance/c2go/darwin.SincospiStret",
	"__sincospif_stret": "github.com/elliotchance/c2go/darwin.SincospifStret",
	"__sincos_stret":    "github.com/elliotchance/c2go/darwin.SincosStret",
	"__sincosf_stret":   "github.com/elliotchance/c2go/darwin.SincosfStret",

	// darwin/assert.h
	"__builtin_expect": "github.com/elliotchance/c2go/darwin.BuiltinExpect",
	"__assert_rtn":     "github.com/elliotchance/c2go/darwin.AssertRtn",

	// linux/assert.h
	"__assert_fail": "github.com/elliotchance/c2go/linux.AssertFail",
}

// TODO: Some of these are based on assumptions that may not be true for all
// architectures (like the size of an int). At some point in the future we will
// need to find out the sizes of some of there and pick the most compatible type.
//
// Please keep them sorted by name.
var SimpleResolveTypes = map[string]string{
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

	// Darwin specific
	"__darwin_ct_rune_t": "__darwin_ct_rune_t",
	"union __mbstate_t":  "__mbstate_t",
	"fpos_t":             "int",
	"struct __float2":    "github.com/elliotchance/c2go/darwin.Float2",
	"struct __double2":   "github.com/elliotchance/c2go/darwin.Double2",

	// These are special cases that almost certainly don"t work. I've put
	// them here because for whatever reason there is no suitable type or we
	// don't need these platform specific things to be implemented yet.
	"__builtin_va_list":            "int64",
	"__darwin_pthread_handler_rec": "int64",
	"__int128":                     "int64",
	"__mbstate_t":                  "int64",
	"__sbuf":                       "int64",
	"__sFILEX":                     "interface{}",
	"__va_list_tag":                "interface{}",
	"FILE":                         "int64",
}

var TypesAlreadyDefined = []string{
	// Linux specific
	"_LIB_VERSION_TYPE",

	// Darwin specific
	"__float2",
	"__double2",
}

var Imports = []string{"fmt"}

func ucfirst(word string) string {
	return strings.ToUpper(string(word[0])) + word[1:]
}

func getExportedName(field string) string {
	return ucfirst(strings.TrimLeft(field, "_"))
}

func addImport(importName string) {
	for _, i := range Imports {
		if i == importName {
			return
		}
	}

	Imports = append(Imports, importName)
}

func importType(typeName string) string {
	if strings.Index(typeName, ".") != -1 {
		parts := strings.Split(typeName, ".")
		addImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(typeName, "/")
		return parts2[len(parts2)-1]
	}

	return typeName
}

func isKeyword(w string) bool {
	return w == "char" || w == "long" || w == "struct" || w == "void"
}

func isIdentifier(w string) bool {
	return !isKeyword(w) && regexp.MustCompile("[_a-zA-Z][_a-zA-Z0-9]*").
		MatchString(w)
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
	for k, v := range SimpleResolveTypes {
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

			for _, v := range SimpleResolveTypes {
				if v == s {
					return "*" + importType(SimpleResolveTypes[s])
				}
			}

			return "*" + s
		} else {
			s = s[7:]

			for _, v := range SimpleResolveTypes {
				if v == s {
					return importType(SimpleResolveTypes[s])
				}
			}

			return s
		}
	}

	// Enums are by name.
	if s[:5] == "enum " {
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

	panic(fmt.Sprintf("'%s'", s))
}

func inStrings(item string, items []string) bool {
	for _, v := range items {
		if item == v {
			return true
		}
	}

	return false
}

func cast(expr, fromType, toType string) string {
	fromType = resolveType(fromType)
	toType = resolveType(toType)

	if fromType == toType {
		return expr
	}

	types := []string{"int", "int64", "uint32", "__darwin_ct_rune_t",
		"byte", "float32", "float64"}

	for _, v := range types {
		if fromType == v && toType == "bool" {
			return fmt.Sprintf("%s != 0", expr)
		}
	}

	if fromType == "*int" && toType == "bool" {
		return fmt.Sprintf("%s != nil", expr)
	}

	if inStrings(fromType, types) && inStrings(toType, types) {
		return fmt.Sprintf("%s(%s)", toType, expr)
	}

	addImport("github.com/elliotchance/c2go/noarch")
	return fmt.Sprintf("noarch.%sTo%s(%s)", ucfirst(fromType), ucfirst(toType), expr)
}

func printLine(out *bytes.Buffer, line string, indent int) {
	out.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat("\t", indent), line))
}

func renderExpression(node interface{}) []string {
	switch n := node.(type) {
	case *FieldDecl:
		fieldType := resolveType(n.Type)
		name := strings.Replace(n.Name, "used", "", -1)

		// Go does not allow the name of a variable to be called "type".
		// For the moment I will rename this to avoid the error.
		if name == "type" {
			name = "type_"
		}

		suffix := ""
		if len(n.Children) > 0 {
			suffix = fmt.Sprintf(" = %s", renderExpression(n.Children[0])[0])
		}

		if suffix == " = (0)" {
			suffix = " = nil"
		}

		return []string{fmt.Sprintf("%s %s%s", name, fieldType, suffix), "unknown3"}

	case *CallExpr:
		children := n.Children
		func_name := renderExpression(children[0])[0]

		func_def := getFunctionDefinition(func_name)

		if _, ok := FunctionSubstitutions[func_name]; ok {
			parts := strings.Split(FunctionSubstitutions[func_name], ".")
			addImport(strings.Join(parts[:len(parts)-1], "."))

			parts2 := strings.Split(FunctionSubstitutions[func_name], "/")
			func_name = parts2[len(parts2)-1]
		}

		args := []string{}
		i := 0
		for _, arg := range children[1:] {
			e := renderExpression(arg)

			if i > len(func_def.ArgumentTypes)-1 {
				// This means the argument is one of the varargs
				// so we don't know what type it needs to be
				// cast to.
				args = append(args, e[0])
			} else {
				args = append(args, cast(e[0], e[1], func_def.ArgumentTypes[i]))
			}

			i += 1
		}

		parts := []string{}

		for _, v := range args {
			parts = append(parts, v)
		}

		return []string{
			fmt.Sprintf("%s(%s)", func_name, strings.Join(parts, ", ")),
			func_def.ReturnType}

	case *ImplicitCastExpr:
		return renderExpression(n.Children[0])

	case *DeclRefExpr:
		name := n.Name

		if name == "argc" {
			name = "len(os.Args)"
			addImport("os")
		} else if name == "argv" {
			name = "os.Args"
			addImport("os")
		}

		return []string{name, n.Type}

	case *StringLiteral:
		return []string{
			fmt.Sprintf("\"%s\"", strings.Replace(n.Value, "\n", "\\n", -1)),
			"const char *",
		}

	case *VarDecl:
		theType := resolveType(n.Type)
		name := n.Name

		// Go does not allow the name of a variable to be called "type".
		// For the moment I will rename this to avoid the error.
		if name == "type" {
			name = "type_"
		}

		suffix := ""
		if len(n.Children) > 0 {
			children := n.Children
			suffix = fmt.Sprintf(" = %s", renderExpression(children[0])[0])
		}

		if suffix == " = (0)" {
			suffix = " = nil"
		}

		return []string{fmt.Sprintf("var %s %s%s", name, theType, suffix), "unknown3"}

	case *BinaryOperator:
		operator := n.Operator

		left := renderExpression(n.Children[0])
		right := renderExpression(n.Children[1])

		return_type := "bool"
		if inStrings(operator, []string{"|", "&", "+", "-", "*", "/"}) {
			// TODO: The left and right type might be different
			return_type = left[1]
		}

		if operator == "&&" {
			left[0] = cast(left[0], left[1], return_type)
			right[0] = cast(right[0], right[1], return_type)
		}

		if (operator == "!=" || operator == "==") && right[0] == "(0)" {
			right[0] = "nil"
		}

		return []string{fmt.Sprintf("%s %s %s", left[0], operator, right[0]), return_type}

	case *IntegerLiteral:
		literal := n.Value

		// FIXME
		//if str(literal)[-1] == 'L':
		//    literal = '%s(%s)' % (resolveType('long'), literal[:-1])

		return []string{strconv.FormatInt(int64(literal), 10), "int"}

	case *UnaryOperator:
		operator := n.Operator
		expr := renderExpression(n.Children[0])

		if operator == "!" {
			if expr[1] == "bool" {
				return []string{fmt.Sprintf("!(%s)", expr[0]), expr[1]}
			}

			addImport("github.com/elliotchance/c2go/noarch")
			return []string{fmt.Sprintf("%s(%s)", fmt.Sprintf("noarch.Not%s", ucfirst(expr[1])), expr[0]), expr[1]}
		}

		if operator == "*" {
			if expr[1] == "const char *" {
				return []string{fmt.Sprintf("%s[0]", expr[0]), "char"}
			}

			return []string{fmt.Sprintf("*%s", expr[0]), "int"}
		}

		if operator == "++" {
			return []string{fmt.Sprintf("%s += 1", expr[0]), expr[1]}
		}

		if operator == "~" {
			operator = "^"
		}

		return []string{fmt.Sprintf("%s%s", operator, expr[0]), expr[1]}

	case *ArraySubscriptExpr:
		children := n.Children
		return []string{fmt.Sprintf("%s[%s]", renderExpression(children[0])[0],
			renderExpression(children[1])[0]), "unknown1"}

	case *ParenExpr:
		a := renderExpression(n.Children[0])
		return []string{fmt.Sprintf("(%s)", a[0]), a[1]}

	case *ConditionalOperator:
		a := renderExpression(n.Children[0])[0]
		b := renderExpression(n.Children[1])[0]
		c := renderExpression(n.Children[2])[0]

		addImport("github.com/elliotchance/c2go/noarch")
		return []string{
			fmt.Sprintf("noarch.Ternary(%s, func () interface{} { return %s }, func () interface{} { return %s })", a, b, c),
			n.Type,
		}

	case *CStyleCastExpr:
		children := n.Children
		return renderExpression(children[0])

	case *PredefinedExpr:
		if n.Name == "__PRETTY_FUNCTION__" {
			// FIXME
			return []string{"\"void print_number(int *)\"", "const char*"}
		}

		if n.Name == "__func__" {
			// FIXME
			return []string{fmt.Sprintf("\"%s\"", "print_number"), "const char*"}
		}

		panic(fmt.Sprintf("renderExpression: unknown PredefinedExpr: %s", n.Name))

	case *FloatingLiteral:
		return []string{fmt.Sprintf("%f", n.Value), "double"}

	case *MemberExpr:
		children := n.Children

		lhs := renderExpression(children[0])
		lhs_type := resolveType(lhs[1])
		rhs := n.Name

		if inStrings(lhs_type, []string{"darwin.Float2", "darwin.Double2"}) {
			rhs = getExportedName(rhs)
		}

		return []string{
			fmt.Sprintf("%s.%s", lhs[0], rhs),
			children[0].(*DeclRefExpr).Type,
		}

	default:
		panic(fmt.Sprintf("renderExpression: %#v", n))
	}
}

func getFunctionParams(f *FunctionDecl) []*ParmVarDecl {
	r := []*ParmVarDecl{}
	for _, n := range f.Children {
		if v, ok := n.(*ParmVarDecl); ok {
			r = append(r, v)
		}
	}

	return r
}

func getFunctionReturnType(f string) string {
	// The type of the function will be the complete prototype, like:
	//
	//     __inline_isfinitef(float) int
	//
	// will have a type of:
	//
	//     int (float)
	//
	// The arguments will handle themselves, we only care about the
	// return type ('int' in this case)
	return strings.TrimSpace(strings.Split(f, "(")[0])
}

func Render(out *bytes.Buffer, node interface{}, function_name string, indent int, return_type string) {
	switch n := node.(type) {
	case *TranslationUnitDecl:
		for _, c := range n.Children {
			Render(out, c, function_name, indent, return_type)
		}

	case *TypedefDecl:
		name := strings.TrimSpace(n.Name)
		for _, v := range TypesAlreadyDefined {
			if name == v {
				return
			}
		}

		TypesAlreadyDefined = append(TypesAlreadyDefined, name)

		// FIXME: All of the logic here is just to avoid errors, it
		// needs to be fixed up.
		// if ("struct" in node["type"] or "union" in node["type"]) and :
		//     return
		n.Type = strings.Replace(n.Type, "unsigned", "", -1)

		resolved_type := resolveType(n.Type)

		if name == "__mbstate_t" {
			addImport("github.com/elliotchance/c2go/darwin")
			resolved_type = "darwin.C__mbstate_t"
		}

		if name == "__darwin_ct_rune_t" {
			addImport("github.com/elliotchance/c2go/darwin")
			resolved_type = "darwin.C__darwin_ct_rune_t"
		}

		if name == "__builtin_va_list" || name == "__qaddr_t" || name == "definition" || name ==
			"_IO_lock_t" || name == "va_list" || name == "fpos_t" || name == "__NSConstantString" || name ==
			"__darwin_va_list" || name == "__fsid_t" || name == "_G_fpos_t" || name == "_G_fpos64_t" {
			return
		}

		printLine(out, fmt.Sprintf("type %s %s\n", name, resolved_type), indent)

		return

	case *RecordDecl:
		name := strings.TrimSpace(n.Name)
		if inStrings(name, TypesAlreadyDefined) || name == "" {
			return
		}

		TypesAlreadyDefined = append(TypesAlreadyDefined, name)

		if n.Kind == "union" {
			return
		}

		printLine(out, fmt.Sprintf("type %s %s {", name, n.Kind), indent)
		if len(n.Children) > 0 {
			for _, c := range n.Children {
				Render(out, c, function_name, indent+1, "")
			}
		}

		printLine(out, "}\n", indent)
		return

	case *FieldDecl:
		printLine(out, renderExpression(node)[0], indent+1)
		return

	case *FunctionDecl:
		function_name = strings.TrimSpace(n.Name)

		if function_name == "__istype" || function_name == "__isctype" ||
			function_name == "__wcwidth" || function_name == "__sputc" ||
			function_name == "__inline_signbitf" ||
			function_name == "__inline_signbitd" ||
			function_name == "__inline_signbitl" {
			return
		}

		has_body := false
		if len(n.Children) > 0 {
			for _, c := range n.Children {
				if _, ok := c.(*CompoundStmt); ok {
					has_body = true
				}
			}
		}

		args := []string{}
		for _, a := range getFunctionParams(n) {
			args = append(args, fmt.Sprintf("%s %s", a.Name, resolveType(a.Type)))
		}

		if has_body {
			return_type := getFunctionReturnType(n.Type)

			if function_name == "main" {
				printLine(out, "func main() {", indent)
			} else {
				printLine(out, fmt.Sprintf("func %s(%s) %s {",
					function_name, strings.Join(args, ", "),
					resolveType(return_type)), indent)
			}

			for _, c := range n.Children {
				if _, ok := c.(*CompoundStmt); ok {
					Render(out, c, function_name,
						indent+1, n.Type)
				}
			}

			printLine(out, "}\n", indent)

			params := []string{}
			for _, v := range getFunctionParams(n) {
				params = append(params, v.Type)
			}

			addFunctionDefinition(FunctionDefinition{
				Name:          n.Name,
				ReturnType:    getFunctionReturnType(n.Type),
				ArgumentTypes: params,
			})
		}

	case *VarDecl:
	// FIXME?

	case *CompoundStmt:
		for _, c := range n.Children {
			Render(out, c, function_name, indent, return_type)
		}

	case *CallExpr:
		printLine(out, renderExpression(node)[0], indent)

	case *ReturnStmt:
		r := "return"

		if len(n.Children) > 0 && function_name != "main" {
			re := renderExpression(n.Children[0])
			r = "return " + cast(re[0], re[1], "int")
		}

		printLine(out, r, indent)

	case *DeclStmt:
		for _, child := range n.Children {
			printLine(out, renderExpression(child)[0], indent)
		}

	case *ForStmt:
		children := n.Children

		a := renderExpression(children[0])[0]
		b := renderExpression(children[1])[0]
		c := renderExpression(children[2])[0]

		printLine(out, fmt.Sprintf("for %s; %s; %s {", a, b, c), indent)

		Render(out, children[3], function_name, indent+1, return_type)

		printLine(out, "}", indent)

	case *BinaryOperator:
		printLine(out, renderExpression(node)[0], indent)

	case *ParenExpr:
		printLine(out, renderExpression(node)[0], indent)

	case *IfStmt:
		children := n.Children

		e := renderExpression(children[0])
		printLine(out, fmt.Sprintf("if %s {", cast(e[0], e[1], "bool")), indent)

		Render(out, children[1], function_name, indent+1, return_type)

		if len(children) > 2 {
			printLine(out, "} else {", indent)
			Render(out, children[2], function_name, indent+1, return_type)
		}

		printLine(out, "}", indent)

	case *BreakStmt:
		printLine(out, "break", indent)

	case *WhileStmt:
		children := n.Children

		e := renderExpression(children[0])
		printLine(out, fmt.Sprintf("for %s {", cast(e[0], e[1], "bool")), indent)

		// FIXME: Does this do anything?
		Render(out, children[1], function_name, indent+1, return_type)

		printLine(out, "}", indent)

	case *UnaryOperator:
		printLine(out, renderExpression(node)[0], indent)

	case *EnumDecl:
		return

	default:
		panic(reflect.ValueOf(node).Elem().Type())
	}
}

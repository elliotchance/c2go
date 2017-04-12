package c2go

import (
	"bytes"
	"fmt"
	"github.com/elliotchance/c2go/ast"
	"reflect"
	"regexp"
	"strings"
)

type FunctionDefinition struct {
	ReturnType    string
	ArgumentTypes []string
}

var FunctionDefinitions = map[string]FunctionDefinition{
	// darwin/assert.h
	"__builtin_expect": FunctionDefinition{"int", []string{"int", "int"}},
	"__assert_rtn":     FunctionDefinition{"bool", []string{"const char*", "const char*", "int", "const char*"}},

	// darwin/ctype.h
	"__istype":   FunctionDefinition{"uint32", []string{"__darwin_ct_rune_t", "uint32"}},
	"__isctype":  FunctionDefinition{"__darwin_ct_rune_t", []string{"__darwin_ct_rune_t", "uint32"}},
	"__tolower":  FunctionDefinition{"__darwin_ct_rune_t", []string{"__darwin_ct_rune_t"}},
	"__toupper":  FunctionDefinition{"__darwin_ct_rune_t", []string{"__darwin_ct_rune_t"}},
	"__maskrune": FunctionDefinition{"uint32", []string{"__darwin_ct_rune_t", "uint32"}},

	// darwin/math.h
	"__builtin_fabs":    FunctionDefinition{"double", []string{"double"}},
	"__builtin_fabsf":   FunctionDefinition{"float", []string{"float"}},
	"__builtin_fabsl":   FunctionDefinition{"double", []string{"double"}},
	"__builtin_inf":     FunctionDefinition{"double", []string{}},
	"__builtin_inff":    FunctionDefinition{"float", []string{}},
	"__builtin_infl":    FunctionDefinition{"double", []string{}},
	"__sincospi_stret":  FunctionDefinition{"Double2", []string{"float"}},
	"__sincospif_stret": FunctionDefinition{"Float2", []string{"float"}},
	"__sincos_stret":    FunctionDefinition{"Double2", []string{"float"}},
	"__sincosf_stret":   FunctionDefinition{"Float2", []string{"float"}},

	// linux/assert.h
	"__assert_fail": FunctionDefinition{"bool", []string{"const char*", "const char*", "unsigned int", "const char*"}},
}

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

	// Linux specific
	"_IO_FILE": "github.com/elliotchance/c2go/linux.File",

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
	"_IO_FILE",

	// Darwin specific
	"__float2",
	"__double2",
}

var imports = []string{"fmt"}

func ucfirst(word string) string {
	return strings.ToUpper(string(word[0])) + word[1:]
}

func getExportedName(field string) string {
	return ucfirst(strings.TrimLeft(field, "_"))
}

func addImport(importName string) {
	for _, i := range imports {
		if i == importName {
			return
		}
	}

	imports = append(imports, importName)
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

	// It may be a pointer of a simple type. For example, float *, int *, etc.
	//try:
	if regexp.MustCompile("[\\w ]+\\*").MatchString(s) {
		return "*" + resolveType(strings.TrimSpace(s[:len(s)-2]))
	}
	//except NoSuchTypeException:
	//    # Keep trying the next one.
	//    pass

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

	//try:
	// It could be an array of fixed length.
	search2 := regexp.MustCompile("([\\w ]+)\\[(\\d+)\\]").FindStringSubmatch(s)
	if len(search2) > 0 {
		return fmt.Sprintf("[%s]%s", search2[2], resolveType(search2[1]))
	}
	//except NoSuchTypeException as e:
	// Make the nested exception message more contextual.
	//raise NoSuchTypeException(e.message + " (from '%s')" % s)

	//raise NoSuchTypeException("'%s'" % s)
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

	types := []string{"int", "int64", "uint32", "__darwin_ct_rune_t", "byte", "float32",
		"float64"}

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
	case *ast.FieldDecl:
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

	default:
		panic(fmt.Sprintf("renderExpression: %#v", n))
	}
	//    if node['node'] == 'BinaryOperator':
	//        operator = node['operator']
	//
	//        left, left_type = renderExpression(node['children'][0])
	//        right, right_type = renderExpression(node['children'][1])
	//
	//        return_type = 'bool'
	//        if operator in ('|', '&', '+', '-', '*', '/'):
	//            # TODO: The left and right type might be different
	//            return_type = left_type
	//
	//        if operator == '&&':
	//            left = cast(left, left_type, return_type)
	//            right = cast(right, right_type, return_type)
	//
	//        if (operator == '!=' or operator == '==') and right == '(0)':
	//            right = 'nil'
	//
	//        return '%s %s %s' % (left, operator, right), return_type
	//
	//    if node['node'] == 'UnaryOperator':
	//        operator = node['operator']
	//        expr = renderExpression(node['children'][0])
	//
	//        if operator == '!':
	//            if expr[1] == 'bool':
	//                return '!(%s)' % expr[0], expr[1]
	//
	//            return '%s(%s)' % ('__not_%s' % expr[1], expr[0]), expr[1]
	//
	//        if operator == '*':
	//            if expr[1] == 'const char *':
	//                return '%s[0]' % expr[0], 'char'
	//
	//            return '*%s' % expr[0], 'int'
	//
	//        if operator == '++':
	//            return '%s += 1' % expr[0], expr[1]
	//
	//        if operator == '~':
	//            operator = '^'
	//
	//        return '%s%s' % (operator, expr[0]), expr[1]
	//
	//    if node['node'] == 'StringLiteral':
	//        return '"%s"' % node['value'].replace("\n", "\\n"), 'const char *'
	//
	//    if node['node'] == 'FloatingLiteral':
	//        return node['value'], 'double'
	//
	//    if node['node'] == 'IntegerLiteral':
	//        literal = node['value']
	//        if str(literal)[-1] == 'L':
	//            literal = '%s(%s)' % (resolveType('long'), literal[:-1])
	//
	//        return literal, 'int'
	//
	//    if node['node'] == 'DeclRefExpr':
	//        name = node['name']
	//
	//        if name == 'argc':
	//            name = 'len(os.Args)'
	//            addImport("os")
	//        elif name == 'argv':
	//            name = 'os.Args'
	//            addImport("os")
	//
	//        return name, node['type']
	//
	//    if node['node'] == 'ImplicitCastExpr':
	//        return renderExpression(node['children'][0])
	//
	//    if node['node'] == 'CallExpr':
	//        children = node['children']
	//        func_name = renderExpression(children[0])[0]
	//
	//        func_def = FunctionDefinitions[func_name]
	//
	//        if func_name in FunctionSubstitutions:
	//            addImport('.'.join(FunctionSubstitutions[func_name].split('.')[:-1]))
	//            func_name = FunctionSubstitutions[func_name].split('/')[-1]
	//
	//        args = []
	//        i = 0
	//        for arg in children[1:]:
	//            e = renderExpression(arg)
	//
	//            if i > len(func_def[1]) - 1:
	//                # This means the argument is one of the varargs so we don't know
	//                # what type it needs to be cast to.
	//                args.append(e[0])
	//            else:
	//                args.append(cast(e[0], e[1], func_def[1][i]))
	//
	//            i += 1
	//
	//        return '%s(%s)' % (func_name, ', '.join([str(a) for a in args])), func_def[0]
	//
	//    if node['node'] == 'ArraySubscriptExpr':
	//        children = node['children']
	//        return '%s[%s]' % (renderExpression(children[0])[0],
	//            renderExpression(children[1])[0]), 'unknown1'
	//
	//    if node['node'] == 'MemberExpr':
	//        children = node['children']
	//
	//        lhs = renderExpression(children[0])
	//        lhs_type = resolveType(lhs[1])
	//        rhs = node['name']
	//
	//        if lhs_type in ('darwin.Float2', 'darwin.Double2'):
	//            rhs = getExportedName(rhs)
	//
	//        return '%s.%s' % (lhs[0], rhs), children[0]['type']
	//
	//    if node['node'] == 'CStyleCastExpr':
	//        children = node['children']
	//        return renderExpression(children[0])
	//
	//    if node['node'] == 'FieldDecl' or node['node'] == 'VarDecl':
	//        type = resolveType(node['type'])
	//        name = node['name'].replace('used', '')
	//
	//        # Go does not allow the name of a variable to be called "type". For the
	//        # moment I will rename this to avoid the error.
	//        if name == 'type':
	//            name = 'type_'
	//
	//        prefix = ''
	//        if node['node'] == 'VarDecl':
	//            prefix = 'var '
	//
	//        suffix = ''
	//        if 'children' in node:
	//            children = node['children']
	//            suffix = ' = %s' % renderExpression(children[0])[0]
	//
	//            if suffix == ' = (0)':
	//                suffix = ' = nil'
	//
	//        return '%s%s %s%s' % (prefix, name, type, suffix), 'unknown3'
	//
	//    if node['node'] == 'RecordDecl':
	//        return '/* RecordDecl */', 'unknown5'
	//
	//    if node['node'] == 'ParenExpr':
	//        a, b = renderExpression(node['children'][0])
	//        return '(%s)' % a, b
	//
	//    if node['node'] == 'PredefinedExpr':
	//        if node['name'] == '__PRETTY_FUNCTION__':
	//            # FIXME
	//            return '"void print_number(int *)"', 'const char*'
	//
	//        if node['name'] == '__func__':
	//            # FIXME
	//            return '"%s"' % 'print_number', 'const char*'
	//
	//        raise Exception('renderExpression: unknown PredefinedExpr: %s' % node['name'])
	//
	//    if node['node'] == 'ConditionalOperator':
	//        a = renderExpression(node['children'][0])[0]
	//        b = renderExpression(node['children'][1])[0]
	//        c = renderExpression(node['children'][2])[0]
	//
	//        addImport('github.com/elliotchance/c2go/noarch')
	//        return 'noarch.Ternary(%s, func () interface{} { return %s }, func () interface{} { return %s })' % (a, b, c), node['type']
}

func getFunctionParams(f *ast.FunctionDecl) []*ast.ParmVarDecl {
	r := []*ast.ParmVarDecl{}
	for _, n := range f.Children {
		if v, ok := n.(*ast.ParmVarDecl); ok {
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
	case *ast.TranslationUnitDecl:
		for _, c := range n.Children {
			Render(out, c, function_name, indent, return_type)
		}
		panic("nice")

	case *ast.TypedefDecl:
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

	case *ast.RecordDecl:
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

	case *ast.FieldDecl:
		printLine(out, renderExpression(node)[0], indent+1)
		return

	case *ast.FunctionDecl:
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
				if _, ok := c.(*ast.CompoundStmt); ok {
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
				if _, ok := c.(*ast.CompoundStmt); ok {
					Render(out, c, function_name,
						indent+1, n.Type)
				}
			}

			printLine(out, "}\n", indent)

			params := []string{}
			for _, v := range getFunctionParams(n) {
				params = append(params, v.Type)
			}

			FunctionDefinitions[n.Name] = FunctionDefinition{
				getFunctionReturnType(n.Type), params,
			}
		}

	case *ast.VarDecl:
		// FIXME?
		return

	case *ast.CompoundStmt:
		for _, c := range n.Children {
			Render(out, c, function_name, indent, return_type)
		}

	default:
		panic(reflect.ValueOf(node).Elem().Type())
	}
}

//    if node['node'] == 'IfStmt':
//        children = node['children']
//
//        e = renderExpression(children[0])
//        printLine(out, 'if %s {' % cast(e[0], e[1], 'bool'), indent)
//
//        render(out, children[1], function_name, indent + 1, return_type)
//
//        if len(children) > 2:
//            printLine(out, '} else {', indent)
//            render(out, children[2], function_name, indent + 1, return_type)
//
//        printLine(out, '}', indent)
//
//        return
//
//    if node['node'] == 'WhileStmt':
//        children = node['children']
//
//        e = renderExpression(children[0])
//        printLine(out, 'for %s {' % cast(e[0], e[1], 'bool'), indent)
//
//        render(out, children[1], function_name, indent + 1, return_type)
//
//        printLine(out, '}', indent)
//
//        return
//
//    if node['node'] == 'ForStmt':
//        children = node['children']
//
//        a, b, c = [renderExpression(e)[0] for e in children[:3]]
//        printLine(out, 'for %s; %s; %s {' % (a, b, c), indent)
//
//        render(out, children[3], function_name, indent + 1, return_type)
//
//        printLine(out, '}', indent)
//
//        return
//
//    if node['node'] == 'BreakStmt':
//        printLine(out, 'break', indent)
//        return
//
//    if node['node'] == 'UnaryOperator':
//        printLine(out, renderExpression(node)[0], indent)
//        return
//
//    if node['node'] == 'ReturnStmt':
//        r = 'return'
//
//        if 'children' in node and function_name != 'main':
//            expr, type = renderExpression(node['children'][0])
//            r = 'return ' + cast(expr, type, 'int')
//
//        printLine(out, r, indent)
//        return
//
//    if node['node'] in ('BinaryOperator', 'INTEGER_LITERAL', 'CallExpr'):
//        printLine(out, renderExpression(node)[0], indent)
//        return
//
//    if node['node'] == 'EnumDecl':
//        return
//
//
//    if node['node'] == 'DeclStmt':
//        for child in node['children']:
//            printLine(out, renderExpression(child)[0], indent)
//        return
//
//    if node['node'] == 'ParenExpr':
//        printLine(out, renderExpression(node)[0], indent)
//        return
//
//    raise Exception(node['node'])

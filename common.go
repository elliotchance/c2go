package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func printLine(out *bytes.Buffer, line string, indent int) {
	out.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat("\t", indent), line))
}

func renderExpression(node interface{}) []string {
	if n, ok := node.(Renderer); ok {
		return n.Render()
	}

	panic(fmt.Sprintf("renderExpression: %#v", node))
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
		if typeIsAlreadyDefined(name) {
			return
		}

		typeIsNowDefined(name)

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
		if name == "" || typeIsAlreadyDefined(name) {
			return
		}

		typeIsNowDefined(name)

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

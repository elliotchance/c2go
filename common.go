package main

import (
	"bytes"
	"fmt"
	"strings"
)

func printLine(out *bytes.Buffer, line string, indent int) {
	out.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat("\t", indent), line))
}

func renderExpression(node interface{}) []string {
	if n, ok := node.(ExpressionRenderer); ok {
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

func Render(out *bytes.Buffer, node interface{}, functionName string, indent int, returnType string) {
	if n, ok := node.(LineRenderer); ok {
		n.RenderLine(out, functionName, indent, returnType)
		return
	}

	printLine(out, renderExpression(node)[0], indent)
}

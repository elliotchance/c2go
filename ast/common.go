package ast

import (
	"bytes"
	"fmt"
	"strings"
)

func printLine(out *bytes.Buffer, line string, indent int) {
	out.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat("\t", indent), line))
}

func Render(ast *Ast, node Node) string {
	src, _ := node.render(ast)
	return src
}

func renderExpression(ast *Ast, node Node) (string, string) {
	if node == nil {
		return "", "unknown54"
	}

	return node.render(ast)
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

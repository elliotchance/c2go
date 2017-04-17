package ast

import (
	"bytes"
	"fmt"
)

type ForStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseForStmt(line string) *ForStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ForStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *ForStmt) render(ast *Ast) (string, string) {
	out := bytes.NewBuffer([]byte{})

	children := n.Children

	// There are always 5 children in a ForStmt, for example:
	//
	//     for ( c = 0 ; c < n ; c++ ) {
	//         doSomething();
	//     }
	//
	// 1. initExpression = BinaryStmt: c = 0
	// 2. Not sure what this is for, but it's always nil. There is a panic
	//    below in case we discover what it is used for (pun intended).
	// 3. conditionalExpression = BinaryStmt: c < n
	// 4. stepExpression = BinaryStmt: c++
	// 5. body = CompoundStmt: { CallExpr }

	if len(children) != 5 {
		panic(fmt.Sprintf("Expected 5 children in ForStmt, got %#v", children))
	}

	// TODO: The second child of a ForStmt appears to always be null.
	// Are there any cases where it is used?
	if children[1] != nil {
		panic("non-nil child 1 in ForStmt")
	}

	init, _ := renderExpression(ast, children[0])
	conditional, _ := renderExpression(ast, children[2])
	step, _ := renderExpression(ast, children[3])
	body, _ := renderExpression(ast, children[4])

	if init == "" && conditional == "" && step == "" {
		printLine(out, "for {", ast.indent)
	} else {
		printLine(out, fmt.Sprintf("for %s; %s; %s {",
			init, conditional, step), ast.indent)
	}

	printLine(out, body, ast.indent+1)
	printLine(out, "}", ast.indent)

	return out.String(), ""
}

func (n *ForStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

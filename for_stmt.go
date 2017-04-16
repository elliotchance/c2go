package main

import (
	"bytes"
	"fmt"
)

type ForStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseForStmt(line string) *ForStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ForStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *ForStmt) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
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

	init := renderExpression(children[0])[0]
	conditional := renderExpression(children[2])[0]
	step := renderExpression(children[3])[0]
	body := children[4]

	if init == "" && conditional == "" && step == "" {
		printLine(out, "for {", indent)
	} else {
		printLine(out, fmt.Sprintf("for %s; %s; %s {",
			init, conditional, step), indent)
	}

	Render(out, body, functionName, indent+1, returnType)

	printLine(out, "}", indent)
}

// This file contains functions for transpiling common branching and control
// flow, such as "if", "while", "do" and "for". The more complicated control
// flows like "switch" will be put into their own file of the same or sensible
// name.

package transpiler

import (
	"fmt"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"

	goast "go/ast"
)

func transpileIfStmt(n *ast.IfStmt, p *program.Program) (*goast.IfStmt, error) {
	children := n.Children

	// There is always 4 or 5 children in an IfStmt. For example:
	//
	//     if (i == 0) {
	//         return 0;
	//     } else {
	//         return 1;
	//     }
	//
	// 1. Not sure what this is for. This gets removed.
	// 2. Not sure what this is for.
	// 3. conditional = BinaryOperator: i == 0
	// 4. body = CompoundStmt: { return 0; }
	// 5. elseBody = CompoundStmt: { return 1; }
	//
	// elseBody will be nil if there is no else clause.

	// On linux I have seen only 4 children for an IfStmt with the same
	// definitions above, but missing the first argument. Since we don't
	// know what the first argument is for anyway we will just remove it on
	// Mac if necessary.
	if len(children) == 5 && children[0] != nil {
		panic("non-nil child 0 in IfStmt")
	}
	if len(children) == 5 {
		children = children[1:]
	}

	// From here on there must be 4 children.
	if len(children) != 4 {
		panic(fmt.Sprintf("Expected 4 children in IfStmt, got %#v", children))
	}

	// Maybe we will discover what the nil value is?
	if children[0] != nil {
		panic("non-nil child 0 in IfStmt")
	}

	conditional, conditionalType, err := transpileToExpr(children[1], p)
	if err != nil {
		return nil, err
	}

	// The condition in Go must always be a bool.
	boolCondition := types.CastExpr(p, conditional, conditionalType, "bool")

	body, err := transpileToBlockStmt(children[2], p)
	if err != nil {
		return nil, err
	}

	r := &goast.IfStmt{
		If:   token.NoPos,
		Init: nil,
		Cond: boolCondition,
		Body: body,
	}

	if children[3] != nil {
		elseBody, err := transpileToBlockStmt(children[3], p)
		if err != nil {
			return nil, err
		}

		r.Else = elseBody
	}

	return r, nil
}

func transpileForStmt(n *ast.ForStmt, p *program.Program) (*goast.ForStmt, error) {
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

	init, _ := transpileToStmt(children[0], p)
	post, _ := transpileToStmt(children[3], p)
	body, _ := transpileToBlockStmt(children[4], p)

	// The condition can be nil. This means an infinite loop and will be
	// rendered in Go as "for {".
	var condition goast.Expr
	if children[2] != nil {
		var conditionType string
		condition, conditionType, _ = transpileToExpr(children[2], p)
		condition = types.CastExpr(p, condition, conditionType, "bool")
	}

	return &goast.ForStmt{
		Init: init,
		Cond: condition,
		Post: post,
		Body: body,
	}, nil
}

func transpileWhileStmt(n *ast.WhileStmt, p *program.Program) (*goast.ForStmt, error) {
	// TODO: The first child of a WhileStmt appears to always be null.
	// Are there any cases where it is used?
	children := n.Children[1:]

	body, err := transpileToBlockStmt(children[1], p)
	if err != nil {
		return nil, err
	}

	condition, conditionType, err := transpileToExpr(children[0], p)
	if err != nil {
		return nil, err
	}

	return &goast.ForStmt{
		Cond: types.CastExpr(p, condition, conditionType, "bool"),
		Body: body,
	}, nil
}

func transpileDoStmt(n *ast.DoStmt, p *program.Program) (*goast.ForStmt, error) {
	children := n.Children

	body, err := transpileToBlockStmt(children[0], p)
	if err != nil {
		return nil, err
	}

	condition, conditionType, err := transpileToExpr(children[1], p)
	if err != nil {
		return nil, err
	}

	// Add IfStmt to the end of the loop to check the condition
	body.List = append(body.List, &goast.IfStmt{
		Cond: &goast.UnaryExpr{
			Op: token.NOT,
			X:  types.CastExpr(p, condition, conditionType, "bool"),
		},
		Body: &goast.BlockStmt{
			List: []goast.Stmt{&goast.BranchStmt{Tok: token.BREAK}},
		},
	})

	return &goast.ForStmt{
		Body: body,
	}, nil
}

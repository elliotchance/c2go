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
	conditional, _, _ := transpileToExpr(children[2], p)
	post, _ := transpileToStmt(children[3], p)
	body, _ := transpileToBlockStmt(children[4], p)

	return &goast.ForStmt{
		Init: init,
		Cond: conditional,
		Post: post,
		Body: body,
	}, nil
}

func transpileWhileStmt(n *ast.WhileStmt, p *program.Program) (*goast.ForStmt, error) {
	// TODO: The first child of a WhileStmt appears to always be null.
	// Are there any cases where it is used?
	children := n.Children[1:]

	// TODO: Check errors here
	body, _ := transpileToBlockStmt(children[1], p)
	e, _, _ := transpileToExpr(children[0], p)

	return &goast.ForStmt{
		Cond: e,
		Body: body,
	}, nil

	// printLine(out, fmt.Sprintf("for %s {", types.Cast(program, e, eType, "bool")), program.Indent)

	// printLine(out, body, program.Indent+1)

	// printLine(out, "}", program.Indent)

	// return out.String(), ""
}

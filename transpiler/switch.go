package transpiler

import (
	goast "go/ast"
	"go/token"

	"fmt"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

func transpileSwitchStmt(n *ast.SwitchStmt, p *program.Program) (*goast.SwitchStmt, error) {
	// The first two children are nil. I don't know what they are supposed to be
	// for. It looks like the number of children is also not reliable, but we
	// know that we need the last two which represent the condition and body
	// respectively.

	if len(n.Children) < 2 {
		// I don't know what causes this condition. Need to investigate.
		panic(fmt.Sprintf("Less than two children for switch: %#v", n))
	}

	// The condition is the expression to be evaulated against each of the
	// cases.
	condition, _, err := transpileToExpr(n.Children[len(n.Children)-2], p)
	if err != nil {
		return nil, err
	}

	// The body will always be a CompoundStmt because a switch statement is not
	// valid without curly brackets.
	body := n.Children[len(n.Children)-1].(*ast.CompoundStmt)
	cases, err := normalizeSwitchCases(body, p)
	if err != nil {
		return nil, err
	}

	// Convert the normalized cases back into statements so they can be children
	// of goast.SwitchStmt.
	stmts := []goast.Stmt{}
	for _, singleCase := range cases {
		stmts = append(stmts, singleCase)
	}

	return &goast.SwitchStmt{
		Tag: condition,
		Body: &goast.BlockStmt{
			List: stmts,
		},
	}, nil
}

func normalizeSwitchCases(body *ast.CompoundStmt, p *program.Program) ([]*goast.CaseClause, error) {
	// The body of a switch has a non uniform structure. For example:
	//
	//     switch a {
	//     case 1:
	//         foo();
	//         bar();
	//         break;
	//     default:
	//         baz();
	//         qux();
	//     }
	//
	// Looks like:
	//
	//     *ast.CompountStmt
	//         *ast.CaseStmt     // case 1:
	//             *ast.CallExpr //     foo()
	//         *ast.CallExpr     //     bar()
	//         *ast.BreakStmt    //     break
	//         *ast.DefaultStmt  // default:
	//             *ast.CallExpr //     baz()
	//         *ast.CallExpr     //     qux()
	//
	// Each of the cases contains one child that is the first statement, but all
	// the rest are children of the parent CompountStmt. This makes it
	// especially tricky when we want to remove the 'break' or add a
	// 'fallthrough'.
	//
	// To make it easier we normalise the cases. This means that we iterate
	// through all of the statements of the CompountStmt and merge any children
	// that are not 'case' or 'break' with the previous node to give us a
	// structure like:
	//
	//     []*goast.CaseClause
	//         *goast.CaseClause      // case 1:
	//             *goast.CallExpr    //     foo()
	//             *goast.CallExpr    //     bar()
	//             // *ast.BreakStmt  //     break (was removed)
	//         *goast.CaseClause      // default:
	//             *goast.CallExpr    //     baz()
	//             *goast.CallExpr    //     qux()
	//
	// During this translation we also remove 'break' or append a 'fallthrough'.

	cases := []*goast.CaseClause{}
	caseEndedWithBreak := false
	var err error

	for _, x := range body.Children {
		switch c := x.(type) {
		case *ast.CaseStmt, *ast.DefaultStmt:
			cases, err = appendCaseOrDefaultToNormaliziedCases(cases, c, caseEndedWithBreak, p)
			if err != nil {
				return []*goast.CaseClause{}, err
			}
			caseEndedWithBreak = false

		case *ast.BreakStmt:
			caseEndedWithBreak = true

		default:
			stmt, err := transpileToStmt(x, p)
			if err != nil {
				return []*goast.CaseClause{}, err
			}

			cases[len(cases)-1].Body = append(cases[len(cases)-1].Body, stmt)
		}
	}

	return cases, nil
}

func appendCaseOrDefaultToNormaliziedCases(cases []*goast.CaseClause,
	stmt ast.Node, caseEndedWithBreak bool, p *program.Program) ([]*goast.CaseClause, error) {
	if len(cases) > 0 && !caseEndedWithBreak {
		cases[len(cases)-1].Body = append(cases[len(cases)-1].Body, &goast.BranchStmt{
			Tok: token.FALLTHROUGH,
		})
	}
	caseEndedWithBreak = false

	var singleCase *goast.CaseClause
	var err error

	if c, ok := stmt.(*ast.CaseStmt); ok {
		singleCase, err = transpileCaseStmt(c, p)
	}
	if c, ok := stmt.(*ast.DefaultStmt); ok {
		singleCase, err = transpileDefaultStmt(c, p)
	}

	if err != nil {
		return []*goast.CaseClause{}, err
	}

	return append(cases, singleCase), nil
}

func transpileCaseStmt(n *ast.CaseStmt, p *program.Program) (*goast.CaseClause, error) {
	c, _, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, err
	}

	stmts, err := transpileStmts(n.Children[1:], p)
	if err != nil {
		return nil, err
	}

	return &goast.CaseClause{
		List: []goast.Expr{c},
		Body: stmts,
	}, nil
}

func transpileDefaultStmt(n *ast.DefaultStmt, p *program.Program) (*goast.CaseClause, error) {
	stmts, err := transpileStmts(n.Children[0:], p)
	if err != nil {
		return nil, err
	}

	return &goast.CaseClause{
		List: nil,
		Body: stmts,
	}, nil
}

package transpiler

import (
	goast "go/ast"
	"go/token"

	"fmt"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

func normalizeSwitchCases(body *ast.CompoundStmt, p *program.Program) ([]*goast.CaseClause, error) {
	cases := []*goast.CaseClause{}
	caseEndedWithBreak := false

	for _, x := range body.Children {
		switch c := x.(type) {
		case *ast.CaseStmt:
			if len(cases) > 0 && !caseEndedWithBreak {
				cases[len(cases)-1].Body = append(cases[len(cases)-1].Body, &goast.BranchStmt{
					Tok: token.FALLTHROUGH,
				})
			}
			caseEndedWithBreak = false

			singleCase, err := transpileCaseStmt(c, p)
			if err != nil {
				return []*goast.CaseClause{}, err
			}

			cases = append(cases, singleCase)

		case *ast.DefaultStmt:
			if len(cases) > 0 && !caseEndedWithBreak {
				cases[len(cases)-1].Body = append(cases[len(cases)-1].Body, &goast.BranchStmt{
					Tok: token.FALLTHROUGH,
				})
			}
			caseEndedWithBreak = false

			defaultCase, err := transpileDefaultStmt(c, p)
			if err != nil {
				return []*goast.CaseClause{}, err
			}

			cases = append(cases, defaultCase)

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

func transpileSwitchStmt(n *ast.SwitchStmt, p *program.Program) (*goast.SwitchStmt, error) {
	// The first two children are nil. I don't know what they are supposed to be
	// for. It looks like the number of children is also not reliable, but we
	// know that we need the last two which represent the condition and body
	// respectively.

	if len(n.Children) < 2 {
		// I don't know what causes this condition. Need to investigate.
		panic(fmt.Sprintf("Less than two children for switch: %#v", n))
	}

	condition, _, err := transpileToExpr(n.Children[len(n.Children)-2], p)
	if err != nil {
		return nil, err
	}

	// The body will always be a CompoundStmt because a switch statement is not
	// valid without curly brackets. However, the body itself cannot be
	// processed like a normal CompoundStmt because the 'case' and 'break'
	// statements are not grouped together in a way that's easy to tell if we
	// need to drop the 'break' statement or add a 'fallthrough' statement since
	// Go switch statements work very differently to C switch statements.
	body := n.Children[len(n.Children)-1].(*ast.CompoundStmt)

	cases, err := normalizeSwitchCases(body, p)
	if err != nil {
		return nil, err
	}

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

func transpileCaseStmt(n *ast.CaseStmt, p *program.Program) (*goast.CaseClause, error) {
	c, _, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, err
	}

	stmts := []goast.Stmt{}
	for _, s := range n.Children[1:] {
		if s != nil {
			a, err := transpileToStmt(s, p)
			if err != nil {
				return nil, err
			}

			stmts = append(stmts, a)
		}
	}

	return &goast.CaseClause{
		List: []goast.Expr{c},
		Body: stmts,
	}, nil
}

func transpileDefaultStmt(n *ast.DefaultStmt, p *program.Program) (*goast.CaseClause, error) {
	stmts := []goast.Stmt{}
	for _, s := range n.Children[0:] {
		if s != nil {
			a, err := transpileToStmt(s, p)
			if err != nil {
				return nil, err
			}

			stmts = append(stmts, a)
		}
	}

	return &goast.CaseClause{
		List: nil,
		Body: stmts,
	}, nil
}

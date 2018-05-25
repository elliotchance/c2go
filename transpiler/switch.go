// This file contains functions for transpiling a "switch" statement.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
	"golang.org/x/tools/go/ast/astutil"
)

func transpileSwitchStmt(n *ast.SwitchStmt, p *program.Program) (
	_ *goast.SwitchStmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileSwitchStmt : err = %v", err)
		}
	}()

	// The first two children are nil. I don't know what they are supposed to be
	// for. It looks like the number of children is also not reliable, but we
	// know that we need the last two which represent the condition and body
	// respectively.

	if len(n.Children()) < 2 {
		// I don't know what causes this condition. Need to investigate.
		panic(fmt.Sprintf("Less than two children for switch: %#v", n))
	}

	// The condition is the expression to be evaluated against each of the
	// cases.
	condition, _, newPre, newPost, err := transpileToExpr(n.Children()[len(n.Children())-2], p, false)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// separation body of switch on cases
	body := n.Children()[len(n.Children())-1].(*ast.CompoundStmt)

	// solving switch case without body
	// case -1:
	// default: ...
	bodyLen := len(body.Children())
	for i := 0; i < bodyLen; i++ {
		cn := body.ChildNodes[i]
		cs, ok1 := cn.(*ast.CaseStmt)
		ds, ok2 := cn.(*ast.DefaultStmt)
		ls, ok3 := cn.(*ast.LabelStmt)
		if !ok1 && !ok2 && !ok3 || cn == nil || len(cn.Children()) == 0 {
			// Do not consider a node which is not a case, label or default statement here
			continue
		}
		lastCn := cn.Children()[len(cn.Children())-1]
		_, isCase := lastCn.(*ast.CaseStmt)
		_, isDefault := lastCn.(*ast.DefaultStmt)
		_, isLabel := lastCn.(*ast.LabelStmt)
		if isCase || isDefault || isLabel {
			// Insert lastCn before next case in body (https://github.com/golang/go/wiki/SliceTricks)
			body.ChildNodes = append(body.ChildNodes, &ast.CompoundStmt{})
			copy(body.ChildNodes[i+2:], body.ChildNodes[i+1:])
			body.ChildNodes[i+1] = lastCn
			bodyLen++

			if len(cn.Children()) == 1 {
				// If cn child nodes would be empty without lastCn,
				// replace lastCn by an empty CompoundStmt
				cn.Children()[0] = &ast.CompoundStmt{}
			} else {
				// Remove lastCn from cn child nodes
				if ok1 {
					cs.ChildNodes = cs.ChildNodes[:len(cs.ChildNodes)-1]
				}
				if ok2 {
					ds.ChildNodes = ds.ChildNodes[:len(ds.ChildNodes)-1]
				}
				if ok3 {
					ls.ChildNodes = ls.ChildNodes[:len(ls.ChildNodes)-1]
				}
			}
		}
	}

	for i := range body.Children() {
		// For simplification - each CaseStmt will have CompoundStmt
		if v, ok := body.Children()[i].(*ast.CaseStmt); ok {
			if _, ok := v.Children()[len(v.Children())-1].(*ast.CompoundStmt); !ok {
				var compoundStmt ast.CompoundStmt
				compoundStmt.AddChild(v.Children()[len(v.Children())-1])
				v.Children()[len(v.Children())-1] = &compoundStmt
			}
		}
		// For simplification - each DefaultStmt will have CompoundStmt
		if v, ok := body.Children()[i].(*ast.DefaultStmt); ok {
			if _, ok := v.Children()[len(v.Children())-1].(*ast.CompoundStmt); !ok {
				var compoundStmt ast.CompoundStmt
				compoundStmt.AddChild(v.Children()[len(v.Children())-1])
				v.Children()[len(v.Children())-1] = &compoundStmt
			}
		}
		// For simplification - each LabelStmt will have CompoundStmt
		if v, ok := body.Children()[i].(*ast.LabelStmt); ok {
			if _, ok := v.Children()[len(v.Children())-1].(*ast.CompoundStmt); !ok {
				var compoundStmt ast.CompoundStmt
				compoundStmt.AddChild(v.Children()[len(v.Children())-1])
				v.Children()[len(v.Children())-1] = &compoundStmt
			}
		}
	}

	hasLabelCase := false
	// Move element inside CompoundStmt
	for i := 0; i < len(body.Children()); i++ {
		switch body.Children()[i].(type) {
		case *ast.CaseStmt, *ast.DefaultStmt:
			// do nothing
		case *ast.LabelStmt:
			hasLabelCase = true
			// do nothing else
		default:
			if i != 0 {
				lastStmt := body.Children()[i-1].Children()
				if comp, ok := lastStmt[len(lastStmt)-1].(*ast.CompoundStmt); ok {
					// add node in CompoundStmt
					comp.AddChild(body.Children()[i])

					// remove from body
					if i+1 < len(body.Children()) {
						body.ChildNodes = append(body.ChildNodes[:i], body.ChildNodes[i+1:]...)
					} else {
						body.ChildNodes = body.ChildNodes[:i]
					}

					// goto to last iteration
					i--
				} else {
					p.AddMessage(p.GenerateWarningMessage(
						fmt.Errorf("Unexpected element"), n))
				}
			} else {
				p.AddMessage(p.GenerateWarningMessage(
					fmt.Errorf("Unsupport case"), n))
			}

		}
	}

	// The body will always be a CompoundStmt because a switch statement is not
	// valid without curly brackets.
	cases, newPre, newPost, err := normalizeSwitchCases(body, p)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// For simplification switch case:
	// from:
	// case 3:
	// 	{
	// 		var c int
	// 		return
	// 	}
	// 	fallthrough
	// to:
	// case 3:
	// 	var c int
	// 	return
	//
	for i := range cases {
		cs, ok := cases[i].(*goast.CaseClause)
		if !ok {
			continue
		}
		body := cs.Body
		if len(body) != 2 {
			continue
		}
		var isFallThrough bool
		if v, ok := body[1].(*goast.BranchStmt); ok {
			isFallThrough = (v.Tok == token.FALLTHROUGH)
		}
		if !isFallThrough {
			if len(body) > 1 {
				cs.Body = body
			}
			continue
		}
		if v, ok := body[0].(*goast.BlockStmt); ok {
			if len(v.List) > 0 {
				if vv, ok := v.List[len(v.List)-1].(*goast.BranchStmt); ok {
					if vv.Tok == token.BREAK {
						if isFallThrough {
							v.List = v.List[:len(v.List)-1]
							cs.Body = body[:len(body)-1]
							continue
						}
					}
				}
				if _, ok := v.List[len(v.List)-1].(*goast.ReturnStmt); ok {
					cs.Body = body[:len(body)-1]
					continue
				}
			} else {
				cs.Body = []goast.Stmt{body[1]}
			}
		}
	}

	// Convert the normalized cases back into statements so they can be children
	// of goast.SwitchStmt.
	stmts := []goast.Stmt{}
	for _, singleCase := range cases {
		if singleCase == nil {
			panic("nil single case")
		}

		stmts = append(stmts, singleCase)
	}

	if hasLabelCase {
		stmts, newPost = handleLabelCases(cases, p)
		preStmts, postStmts = combinePreAndPostStmts(preStmts, newPost, []goast.Stmt{}, postStmts)
	}

	return &goast.SwitchStmt{
		Tag: condition,
		Body: &goast.BlockStmt{
			List: stmts,
		},
	}, preStmts, postStmts, nil
}

func normalizeSwitchCases(body *ast.CompoundStmt, p *program.Program) (
	_ []goast.Stmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
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

	cases := []goast.Stmt{}
	caseEndedWithBreak := false

	for _, x := range body.Children() {
		switch c := x.(type) {
		case *ast.CaseStmt, *ast.DefaultStmt, *ast.LabelStmt:
			var newPre, newPost []goast.Stmt
			cases, newPre, newPost, err = appendCaseOrDefaultToNormalizedCases(cases, c, caseEndedWithBreak, p)
			if err != nil {
				return []goast.Stmt{}, nil, nil, err
			}
			caseEndedWithBreak = false

			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		case *ast.BreakStmt:
			caseEndedWithBreak = true

		default:
			var stmt goast.Stmt
			var newPre, newPost []goast.Stmt
			stmt, newPre, newPost, err = transpileToStmt(x, p)
			if err != nil {
				return []goast.Stmt{}, nil, nil, err
			}
			preStmts = append(preStmts, newPre...)
			preStmts = append(preStmts, stmt)
			preStmts = append(preStmts, newPost...)
		}
	}

	return cases, preStmts, postStmts, nil
}

func appendCaseOrDefaultToNormalizedCases(cases []goast.Stmt,
	stmt ast.Node, caseEndedWithBreak bool, p *program.Program) (
	[]goast.Stmt, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	if len(cases) > 0 && !caseEndedWithBreak {
		if cs, ok := cases[len(cases)-1].(*goast.CaseClause); ok {
			cs.Body = append(cs.Body, &goast.BranchStmt{
				Tok: token.FALLTHROUGH,
			})
		}
		if ls, ok := cases[len(cases)-1].(*goast.LabeledStmt); ok {
			ft := &goast.BranchStmt{
				Tok: token.FALLTHROUGH,
			}
			if _, ok2 := ls.Stmt.(*goast.EmptyStmt); ok2 {
				ls.Stmt = ft
			} else if bs, ok2 := ls.Stmt.(*goast.BlockStmt); ok2 {
				bs.List = append(bs.List, ft)
			} else {
				ls.Stmt = &goast.BlockStmt{
					List: []goast.Stmt{
						ls.Stmt,
						ft,
					},
				}
			}
		}
	}
	caseEndedWithBreak = false

	var singleCase goast.Stmt
	var err error
	var newPre []goast.Stmt
	var newPost []goast.Stmt

	switch c := stmt.(type) {
	case *ast.CaseStmt:
		singleCase, newPre, newPost, err = transpileCaseStmt(c, p)

	case *ast.DefaultStmt:
		singleCase, err = transpileDefaultStmt(c, p)

	case *ast.LabelStmt:
		singleCase, newPre, newPost, err = transpileLabelStmt(c, p)
		lc, ok := singleCase.(*goast.LabeledStmt)
		if !ok {
			panic("expected *goast.LabeledStmt")
		}
		if len(newPost) == 1 {
			lc.Stmt = newPost[0]
		} else if len(newPost) > 1 {
			lc.Stmt = &goast.BlockStmt{
				List: newPost,
			}
		}
		newPost = []goast.Stmt{}
	}

	if singleCase != nil {
		cases = append(cases, singleCase)
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if err != nil {
		return []goast.Stmt{}, nil, nil, err
	}

	return cases, preStmts, postStmts, nil
}

func transpileCaseStmt(n *ast.CaseStmt, p *program.Program) (
	*goast.CaseClause, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	c, _, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	stmts, err := transpileStmts(n.Children()[1:], p)
	if err != nil {
		return nil, nil, nil, err
	}

	return &goast.CaseClause{
		List: []goast.Expr{c},
		Body: stmts,
	}, preStmts, postStmts, nil
}

func transpileDefaultStmt(n *ast.DefaultStmt, p *program.Program) (*goast.CaseClause, error) {
	stmts, err := transpileStmts(n.Children()[0:], p)
	if err != nil {
		return nil, err
	}

	return &goast.CaseClause{
		List: nil,
		Body: stmts,
	}, nil
}

func handleLabelCases(cases []goast.Stmt, p *program.Program) (newCases []goast.Stmt, postStmts []goast.Stmt) {
	// In C a switch can have labels before a case.
	// Go does not support this.
	// To make it work we translate the switch cases as labels to blocks appended to the switch
	// For example:
	//
	//     switch a {
	//     case 1:
	//         foo();
	//         break;
	//     LABEL:
	//     case 2:
	//         bar();
	//     default:
	//         baz();
	//     }
	//
	// is transpiled as:
	//
	//     switch a {
	//     case 1:
	//         goto SW_1_1
	//     case 2:
	//         goto SW_1_2
	//     default:
	//         goto SW_1_3
	//     }
	//     SW_1_1:
	//         foo()
	//         goto SW_1_END
	//     LABEL:
	//         ;
	//     SW_1_2:
	//         bar()
	//     SW_1_3:
	//         baz()
	//     SW_1_END:
	//         ;
	swEndLabel := p.GetNextIdentifier("SW_GENERATED_LABEL_")
	postStmts = append(postStmts, &goast.BranchStmt{
		Label: util.NewIdent(swEndLabel),
		Tok:   token.GOTO,
	})
	funcTransformBreak := func(cursor *astutil.Cursor) bool {
		if cursor == nil {
			return true
		}
		node := cursor.Node()
		if bs, ok := node.(*goast.BranchStmt); ok {
			if bs.Tok == token.BREAK {
				cursor.Replace(&goast.BranchStmt{
					Label: util.NewIdent(swEndLabel),
					Tok:   token.GOTO,
				})
			}
		}
		if _, ok := node.(*goast.ForStmt); ok {
			return false
		}
		if _, ok := node.(*goast.RangeStmt); ok {
			return false
		}
		if _, ok := node.(*goast.SwitchStmt); ok {
			return false
		}
		if _, ok := node.(*goast.TypeSwitchStmt); ok {
			return false
		}
		if _, ok := node.(*goast.SelectStmt); ok {
			return false
		}
		return true
	}
	for i, x := range cases {
		switch c := x.(type) {
		case *goast.CaseClause:
			caseLabel := p.GetNextIdentifier("SW_GENERATED_LABEL_")

			if len(c.Body) == 0 {
				c.Body = append(c.Body, &goast.BranchStmt{
					Tok: token.BREAK,
				})
			}
			var isFallThrough bool
			// Remove fallthrough
			if v, ok := c.Body[len(c.Body)-1].(*goast.BranchStmt); ok {
				isFallThrough = (v.Tok == token.FALLTHROUGH)
				c.Body = c.Body[:len(c.Body)-1]
			}
			if len(c.Body) == 0 {
				c.Body = append(c.Body, &goast.EmptyStmt{})
			}

			// Replace break's with goto swEndLabel
			astutil.Apply(c, funcTransformBreak, nil)
			body := c.Body

			// append caseLabel label followed by case body
			postStmts = append(postStmts, &goast.LabeledStmt{
				Label: util.NewIdent(caseLabel),
				Stmt:  body[0],
			})
			body = body[1:]
			postStmts = append(postStmts, body...)

			// If not last case && no fallthrough goto swEndLabel
			if i != len(cases)-1 && !isFallThrough {
				postStmts = append(postStmts, &goast.BranchStmt{
					Label: util.NewIdent(swEndLabel),
					Tok:   token.GOTO,
				})
			}

			// In switch case we goto caseLabel
			c.Body = []goast.Stmt{
				&goast.BranchStmt{
					Label: util.NewIdent(caseLabel),
					Tok:   token.GOTO,
				},
			}
			newCases = append(newCases, c)
		case *goast.LabeledStmt:
			var isFallThrough bool
			// Remove fallthrough if it's the only statement
			if v, ok := c.Stmt.(*goast.BranchStmt); ok {
				if v.Tok == token.FALLTHROUGH {
					c.Stmt = &goast.EmptyStmt{}
					isFallThrough = true
				}
			} else if b, ok := c.Stmt.(*goast.BlockStmt); ok {
				// Remove fallthrough if LabeledStmt contains a BlockStmt
				if v, ok := b.List[len(b.List)-1].(*goast.BranchStmt); ok {
					if v.Tok == token.FALLTHROUGH {
						b.List = b.List[:len(b.List)-1]
						isFallThrough = true
					}
				}
			}

			// Replace break's with goto swEndLabel
			astutil.Apply(c, funcTransformBreak, nil)

			// append label followed by label body
			postStmts = append(postStmts, c)

			// If not last case && no fallthrough goto swEndLabel
			if i != len(cases)-1 && !isFallThrough {
				postStmts = append(postStmts, &goast.BranchStmt{
					Label: util.NewIdent(swEndLabel),
					Tok:   token.GOTO,
				})
			}
		}
	}
	postStmts = append(postStmts, &goast.LabeledStmt{
		Label: util.NewIdent(swEndLabel),
		Stmt:  &goast.EmptyStmt{},
	})
	return
}

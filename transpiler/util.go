// This file contains utility and helper methods for the transpiler.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"reflect"
)

func combinePreAndPostStmts(
	pre []goast.Stmt,
	post []goast.Stmt,
	newPre []goast.Stmt,
	newPost []goast.Stmt) ([]goast.Stmt, []goast.Stmt) {
	pre = append(pre, newPre...)
	post = append(post, newPost...)

	return pre, post
}

// combineStmts - combine elements to slice
func combineStmts(stmt goast.Stmt, preStmts, postStmts []goast.Stmt) (stmts []goast.Stmt) {
	if preStmts != nil {
		stmts = append(stmts, preStmts...)
	}
	if stmt != nil {
		stmts = append(stmts, stmt)
	}
	if postStmts != nil {
		stmts = append(stmts, postStmts...)
	}
	return
}

func removeNil(nodes *[]goast.Decl) {
	walkDeclList(nodes)
}

// Helper functions for common node lists. They may be empty.

func walkIdentList(list *[]*goast.Ident) {
	for i := 0; i < len(*list); i++ {
		if reflect.ValueOf((*list)[i]).IsNil() {
			// if value is nil, then remove from slice
			if len(*list) == 1 {
				*list = nil
				break
			} else if i < len(*list)-1 {
				*list = append((*list)[0:i], (*list)[i+1:]...)
			} else {
				// remove last element of slice
				*list = (*list)[:len(*list)-1]
			}
			i--
			continue
		}
		Walk(((*list)[i]))
	}
}

func walkExprList(list *[]goast.Expr) {
	for i := 0; i < len(*list); i++ {
		if reflect.ValueOf((*list)[i]).IsNil() {
			// if value is nil, then remove from slice
			if len(*list) == 1 {
				*list = nil
				break
			} else if i < len(*list)-1 {
				*list = append((*list)[0:i], (*list)[i+1:]...)
			} else {
				// remove last element of slice
				*list = (*list)[:len(*list)-1]
			}
			i--
			continue
		}
		Walk(((*list)[i]))
	}
}

func walkStmtList(list *[]goast.Stmt) {
	for i := 0; i < len(*list); i++ {
		if reflect.ValueOf((*list)[i]).IsNil() {
			// if value is nil, then remove from slice
			if len(*list) == 1 {
				*list = nil
				break
			} else if i < len(*list)-1 {
				*list = append((*list)[0:i], (*list)[i+1:]...)
			} else {
				// remove last element of slice
				*list = (*list)[:len(*list)-1]
			}
			i--
			continue
		}
		Walk(((*list)[i]))
	}
}

func walkDeclList(list *[]goast.Decl) {
	for i := 0; i < len(*list); i++ {
		if reflect.ValueOf((*list)[i]).IsNil() {
			// if value is nil, then remove from slice
			if len(*list) == 1 {
				*list = nil
				break
			} else if i < len(*list)-1 {
				*list = append((*list)[0:i], (*list)[i+1:]...)
			} else {
				// remove last element of slice
				*list = (*list)[:len(*list)-1]
			}
			i--
			continue
		}
		Walk(((*list)[i]))
	}
}

// Walk - walking inside Go AST tree and remove nil's
func Walk(node goast.Node) {
	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *goast.Comment:
		// nothing to do

	case *goast.CommentGroup:
		for _, c := range n.List {
			Walk(c)
		}

	case *goast.Field:
		if n.Doc != nil {
			Walk(n.Doc)
		}
		walkIdentList(&n.Names)
		Walk(n.Type)
		if n.Tag != nil {
			Walk(n.Tag)
		}
		if n.Comment != nil {
			Walk(n.Comment)
		}

	case *goast.FieldList:
		list := &n.List
		for i := 0; i < len(*list); i++ {
			if reflect.ValueOf((*list)[i]).IsNil() {
				// if value is nil, then remove from slice
				if len(*list) == 1 {
					*list = nil
					break
				} else if i < len(*list)-1 {
					*list = append((*list)[0:i], (*list)[i+1:]...)
				} else {
					// remove last element of slice
					*list = (*list)[:len(*list)-1]
				}
				i--
				continue
			}
			Walk(((*list)[i]))
		}

	// Expressions
	case *goast.BadExpr, *goast.Ident, *goast.BasicLit:
		// nothing to do

	case *goast.Ellipsis:
		if n.Elt != nil {
			Walk(n.Elt)
		}

	case *goast.FuncLit:
		Walk(n.Type)
		Walk(n.Body)

	case *goast.CompositeLit:
		if n.Type != nil {
			Walk(n.Type)
		}
		walkExprList(&n.Elts)

	case *goast.ParenExpr:
		Walk(n.X)

	case *goast.SelectorExpr:
		Walk(n.X)
		Walk(n.Sel)

	case *goast.IndexExpr:
		Walk(n.X)
		Walk(n.Index)

	case *goast.SliceExpr:
		Walk(n.X)
		if n.Low != nil {
			Walk(n.Low)
		}
		if n.High != nil {
			Walk(n.High)
		}
		if n.Max != nil {
			Walk(n.Max)
		}

	case *goast.TypeAssertExpr:
		Walk(n.X)
		if n.Type != nil {
			Walk(n.Type)
		}

	case *goast.CallExpr:
		Walk(n.Fun)
		walkExprList(&n.Args)

	case *goast.StarExpr:
		Walk(n.X)

	case *goast.UnaryExpr:
		Walk(n.X)

	case *goast.BinaryExpr:
		Walk(n.X)
		Walk(n.Y)

	case *goast.KeyValueExpr:
		Walk(n.Key)
		Walk(n.Value)

	// Types
	case *goast.ArrayType:
		if n.Len != nil {
			Walk(n.Len)
		}
		Walk(n.Elt)

	case *goast.StructType:
		Walk(n.Fields)

	case *goast.FuncType:
		if n.Params != nil {
			Walk(n.Params)
		}
		if n.Results != nil {
			Walk(n.Results)
		}

	case *goast.InterfaceType:
		Walk(n.Methods)

	case *goast.MapType:
		Walk(n.Key)
		Walk(n.Value)

	case *goast.ChanType:
		Walk(n.Value)

	// Statements
	case *goast.BadStmt:
		// nothing to do

	case *goast.DeclStmt:
		Walk(n.Decl)

	case *goast.EmptyStmt:
		// nothing to do

	case *goast.LabeledStmt:
		Walk(n.Label)
		Walk(n.Stmt)

	case *goast.ExprStmt:
		Walk(n.X)

	case *goast.SendStmt:
		Walk(n.Chan)
		Walk(n.Value)

	case *goast.IncDecStmt:
		Walk(n.X)

	case *goast.AssignStmt:
		walkExprList(&n.Lhs)
		walkExprList(&n.Rhs)

	case *goast.GoStmt:
		Walk(n.Call)

	case *goast.DeferStmt:
		Walk(n.Call)

	case *goast.ReturnStmt:
		walkExprList(&n.Results)

	case *goast.BranchStmt:
		if n.Label != nil {
			Walk(n.Label)
		}

	case *goast.BlockStmt:
		walkStmtList(&n.List)

	case *goast.IfStmt:
		if n.Init != nil {
			Walk(n.Init)
		}
		Walk(n.Cond)
		Walk(n.Body)
		if n.Else != nil {
			Walk(n.Else)
		}

	case *goast.CaseClause:
		walkExprList(&n.List)
		walkStmtList(&n.Body)

	case *goast.SwitchStmt:
		if n.Init != nil {
			Walk(n.Init)
		}
		if n.Tag != nil {
			Walk(n.Tag)
		}
		Walk(n.Body)

	case *goast.TypeSwitchStmt:
		if n.Init != nil {
			Walk(n.Init)
		}
		Walk(n.Assign)
		Walk(n.Body)

	case *goast.CommClause:
		if n.Comm != nil {
			Walk(n.Comm)
		}
		walkStmtList(&n.Body)

	case *goast.SelectStmt:
		Walk(n.Body)

	case *goast.ForStmt:
		if n.Init != nil {
			Walk(n.Init)
		}
		if n.Cond != nil {
			Walk(n.Cond)
		}
		if n.Post != nil {
			Walk(n.Post)
		}
		Walk(n.Body)

	case *goast.RangeStmt:
		if n.Key != nil {
			Walk(n.Key)
		}
		if n.Value != nil {
			Walk(n.Value)
		}
		Walk(n.X)
		Walk(n.Body)

	// Declarations
	case *goast.ImportSpec:
		if n.Doc != nil {
			Walk(n.Doc)
		}
		if n.Name != nil {
			Walk(n.Name)
		}
		Walk(n.Path)
		if n.Comment != nil {
			Walk(n.Comment)
		}

	case *goast.ValueSpec:
		if n.Doc != nil {
			Walk(n.Doc)
		}
		walkIdentList(&n.Names)
		if n.Type != nil {
			Walk(n.Type)
		}
		walkExprList(&n.Values)
		if n.Comment != nil {
			Walk(n.Comment)
		}

	case *goast.TypeSpec:
		if n.Doc != nil {
			Walk(n.Doc)
		}
		Walk(n.Name)
		Walk(n.Type)
		if n.Comment != nil {
			Walk(n.Comment)
		}

	case *goast.BadDecl:
		// nothing to do

	case *goast.GenDecl:
		if n.Doc != nil {
			Walk(n.Doc)
		}
		list := &n.Specs
		for i := 0; i < len(*list); i++ {
			if reflect.ValueOf((*list)[i]).IsNil() {
				// if value is nil, then remove from slice
				if len(*list) == 1 {
					*list = nil
					break
				} else if i < len(*list)-1 {
					*list = append((*list)[0:i], (*list)[i+1:]...)
				} else {
					// remove last element of slice
					*list = (*list)[:len(*list)-1]
				}
				i--
				continue
			}
			Walk(((*list)[i]))
		}

	case *goast.FuncDecl:
		if n.Doc != nil {
			Walk(n.Doc)
		}
		if n.Recv != nil {
			Walk(n.Recv)
		}
		Walk(n.Name)
		Walk(n.Type)
		if n.Body != nil {
			Walk(n.Body)
		}

	// Files and packages
	case *goast.File:
		if n.Doc != nil {
			Walk(n.Doc)
		}
		Walk(n.Name)
		walkDeclList(&n.Decls)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *goast.Package:
		for _, f := range n.Files {
			Walk(f)
		}

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

}

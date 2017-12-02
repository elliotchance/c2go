// This file contains utility and helper methods for the transpiler.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"reflect"

	"github.com/elliotchance/c2go/program"
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

func removeNil(nodes *[]goast.Decl, p *program.Program) {
	_ = walkDeclList(nodes, p)
}

// Helper functions for common node lists. They may be empty.

func walkIdentList(list *[]*goast.Ident, p *program.Program) (err error) {
	if list == nil || *list == nil {
		return fmt.Errorf("Nil in walkIdentList")
	}
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == nil || reflect.ValueOf((*list)[i]).IsNil() {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T", list))
			err = nil
			goto Remove
		}
		if err = Walk((*list)[i], p); err != nil {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T. Addition err = %v", list, err))
			err = nil
			goto Remove
		}
		continue
	Remove:
		// if value is nil, then remove from slice
		if len(*list) == 1 {
			*list = nil
			return fmt.Errorf("List have only 1 element and it is nil")
		} else if i < len(*list)-1 {
			*list = append((*list)[0:i], (*list)[i+1:]...)
		} else {
			// remove last element of slice
			*list = (*list)[:len(*list)-1]
		}
		i--
	}
	return nil
}

func walkExprList(list *[]goast.Expr, p *program.Program) (err error) {
	if list == nil || *list == nil {
		return fmt.Errorf("Nil in walkExprList")
	}
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == nil || reflect.ValueOf((*list)[i]).IsNil() {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T", list))
			err = nil
			goto Remove
		}
		if err = Walk((*list)[i], p); err != nil {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T. Addition err = %v", list, err))
			err = nil
			goto Remove
		}
		continue
	Remove:
		// if value is nil, then remove from slice
		if len(*list) == 1 {
			*list = nil
			return fmt.Errorf("List have only 1 element and it is nil")
		} else if i < len(*list)-1 {
			*list = append((*list)[0:i], (*list)[i+1:]...)
		} else {
			// remove last element of slice
			*list = (*list)[:len(*list)-1]
		}
		i--
	}
	return nil
}

func walkStmtList(list *[]goast.Stmt, p *program.Program) (err error) {
	if list == nil || *list == nil {
		return fmt.Errorf("Nil in walkStmtList")
	}
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == nil || reflect.ValueOf((*list)[i]).IsNil() {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T", list))
			err = nil
			goto Remove
		}
		if err = Walk((*list)[i], p); err != nil {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T. Addition err = %v", list, err))
			err = nil
			goto Remove
		}
		continue
	Remove:
		// if value is nil, then remove from slice
		if len(*list) == 1 {
			*list = nil
			return fmt.Errorf("List have only 1 element and it is nil")
		} else if i < len(*list)-1 {
			*list = append((*list)[0:i], (*list)[i+1:]...)
		} else {
			// remove last element of slice
			*list = (*list)[:len(*list)-1]
		}
		i--
	}
	return nil
}

func walkDeclList(list *[]goast.Decl, p *program.Program) (err error) {
	if list == nil || *list == nil {
		return fmt.Errorf("Nil in walkDeclList")
	}
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == nil || reflect.ValueOf((*list)[i]).IsNil() {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T", list))
			err = nil
			goto Remove
		}
		if err = Walk((*list)[i], p); err != nil {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T. Addition err = %v", list, err))
			err = nil
			goto Remove
		}
		continue
	Remove:
		// if value is nil, then remove from slice
		if len(*list) == 1 {
			*list = nil
			return fmt.Errorf("List have only 1 element and it is nil")
		} else if i < len(*list)-1 {
			*list = append((*list)[0:i], (*list)[i+1:]...)
		} else {
			// remove last element of slice
			*list = (*list)[:len(*list)-1]
		}
		i--
	}
	return nil
}

func walkCommentList(list *[]*goast.Comment, p *program.Program) (err error) {
	if list == nil || *list == nil {
		return fmt.Errorf("Nil in walkCommentList")
	}
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == nil || reflect.ValueOf((*list)[i]).IsNil() {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T", list))
			err = nil
			goto Remove
		}
		if err = Walk((*list)[i], p); err != nil {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T. Addition err = %v", list, err))
			err = nil
			goto Remove
		}
		continue
	Remove:
		// if value is nil, then remove from slice
		if len(*list) == 1 {
			*list = nil
			return fmt.Errorf("List have only 1 element and it is nil")
		} else if i < len(*list)-1 {
			*list = append((*list)[0:i], (*list)[i+1:]...)
		} else {
			// remove last element of slice
			*list = (*list)[:len(*list)-1]
		}
		i--
	}
	return nil
}

func walkFieldList(list *[]*goast.Field, p *program.Program) (err error) {
	if list == nil || *list == nil {
		return fmt.Errorf("Nil in walkFieldList")
	}
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == nil || reflect.ValueOf((*list)[i]).IsNil() {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T", list))
			err = nil
			goto Remove
		}
		if err = Walk((*list)[i], p); err != nil {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T. Addition err = %v", list, err))
			err = nil
			goto Remove
		}
		continue
	Remove:
		// if value is nil, then remove from slice
		if len(*list) == 1 {
			*list = nil
			return fmt.Errorf("List have only 1 element and it is nil")
		} else if i < len(*list)-1 {
			*list = append((*list)[0:i], (*list)[i+1:]...)
		} else {
			// remove last element of slice
			*list = (*list)[:len(*list)-1]
		}
		i--
	}
	return nil
}

func walkSpecList(list *[]goast.Spec, p *program.Program) (err error) {
	if list == nil || *list == nil {
		return fmt.Errorf("Nil in walkSpecList")
	}
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == nil || reflect.ValueOf((*list)[i]).IsNil() {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T", list))
			err = nil
			goto Remove
		}
		if err = Walk((*list)[i], p); err != nil {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST : %T. Addition err = %v", list, err))
			err = nil
			goto Remove
		}
		continue
	Remove:
		// if value is nil, then remove from slice
		if len(*list) == 1 {
			*list = nil
			return fmt.Errorf("List have only 1 element and it is nil")
		} else if i < len(*list)-1 {
			*list = append((*list)[0:i], (*list)[i+1:]...)
		} else {
			// remove last element of slice
			*list = (*list)[:len(*list)-1]
		}
		i--
	}
	return nil
}

// Walk - walking inside Go AST tree and remove nil's
func Walk(node goast.Node, p *program.Program) (err error) {
	if node == nil {
		return fmt.Errorf("Nil node in Walk function")
	}

	var errSecond error

	defer func() {
		if err != nil {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST (Walk): %T. Addition err = %v", node, err))
			err = nil
		}
		if errSecond != nil {
			p.AddMessage(fmt.Sprintf("// Found nil in Go AST (Walk) |Second priority|: %T. Addition err = %v", node, errSecond))
		}
	}()
	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *goast.Comment:
		// nothing to do

	case *goast.CommentGroup:
		err = walkCommentList(&n.List, p)

	case *goast.Field:
		if n.Doc != nil {
			errSecond = Walk(n.Doc, p)
		}
		if err = walkIdentList(&n.Names, p); err != nil {
			return fmt.Errorf("Names")
		}
		if err = Walk(n.Type, p); err != nil {
			return fmt.Errorf("Type")
		}
		if n.Tag != nil {
			errSecond = Walk(n.Tag, p)
		}
		if n.Comment != nil {
			errSecond = Walk(n.Comment, p)
		}

	case *goast.FieldList:
		err = walkFieldList(&n.List, p)

	// Expressions
	case *goast.BadExpr, *goast.Ident, *goast.BasicLit:
		// nothing to do

	case *goast.Ellipsis:
		if n.Elt != nil {
			errSecond = Walk(n.Elt, p)
		}

	case *goast.FuncLit:
		if err = Walk(n.Type, p); err != nil {
			return fmt.Errorf("Type")
		}
		if err = Walk(n.Body, p); err != nil {
			return fmt.Errorf("Body")
		}

	case *goast.CompositeLit:
		if n.Type != nil {
			errSecond = Walk(n.Type, p)
		}
		if err = walkExprList(&n.Elts, p); err != nil {
			return fmt.Errorf("Elts")
		}

	case *goast.ParenExpr:
		err = Walk(n.X, p)

	case *goast.SelectorExpr:
		if err = Walk(n.X, p); err != nil {
			return fmt.Errorf("X")
		}
		if err = Walk(n.Sel, p); err != nil {
			return fmt.Errorf("Sel")
		}

	case *goast.IndexExpr:
		if err = Walk(n.X, p); err != nil {
			return fmt.Errorf("X")
		}
		if err = Walk(n.Index, p); err != nil {
			return fmt.Errorf("Index")
		}

	case *goast.SliceExpr:
		if err = Walk(n.X, p); err != nil {
			return fmt.Errorf("X")
		}
		if n.Low != nil {
			errSecond = Walk(n.Low, p)
		}
		if n.High != nil {
			errSecond = Walk(n.High, p)
		}
		if n.Max != nil {
			errSecond = Walk(n.Max, p)
		}

	case *goast.TypeAssertExpr:
		if err = Walk(n.X, p); err != nil {
			return fmt.Errorf("X")
		}
		if n.Type != nil {
			errSecond = Walk(n.Type, p)
		}

	case *goast.CallExpr:
		if err = Walk(n.Fun, p); err != nil {
			return fmt.Errorf("Fun")
		}
		errSecond = walkExprList(&n.Args, p)

	case *goast.StarExpr:
		err = Walk(n.X, p)

	case *goast.UnaryExpr:
		err = Walk(n.X, p)

	case *goast.BinaryExpr:
		if err = Walk(n.X, p); err != nil {
			return fmt.Errorf("X")
		}
		if err = Walk(n.Y, p); err != nil {
			return fmt.Errorf("Y")
		}

	case *goast.KeyValueExpr:
		if err = Walk(n.Key, p); err != nil {
			return fmt.Errorf("Key")
		}
		if err = Walk(n.Value, p); err != nil {
			return fmt.Errorf("Value")
		}

	// Types
	case *goast.ArrayType:
		if n.Len != nil {
			errSecond = Walk(n.Len, p)
		}
		if err = Walk(n.Elt, p); err != nil {
			return fmt.Errorf("Elt")
		}

	case *goast.StructType:
		err = Walk(n.Fields, p)

	case *goast.FuncType:
		if n.Params != nil {
			errSecond = Walk(n.Params, p)
		}
		if n.Results != nil {
			errSecond = Walk(n.Results, p)
		}

	case *goast.InterfaceType:
		err = Walk(n.Methods, p)

	case *goast.MapType:
		if err = Walk(n.Key, p); err != nil {
			return fmt.Errorf("Key")
		}
		if err = Walk(n.Value, p); err != nil {
			return fmt.Errorf("Value")
		}

	case *goast.ChanType:
		err = Walk(n.Value, p)

	// Statements
	case *goast.BadStmt:
		// nothing to do

	case *goast.DeclStmt:
		err = Walk(n.Decl, p)

	case *goast.EmptyStmt:
		// nothing to do

	case *goast.LabeledStmt:
		if err = Walk(n.Label, p); err != nil {
			return fmt.Errorf("Label")
		}
		if err = Walk(n.Stmt, p); err != nil {
			return fmt.Errorf("Stmt")
		}

	case *goast.ExprStmt:
		err = Walk(n.X, p)

	case *goast.SendStmt:
		if err = Walk(n.Chan, p); err != nil {
			return fmt.Errorf("Chan")
		}
		if err = Walk(n.Value, p); err != nil {
			return fmt.Errorf("Value")
		}

	case *goast.IncDecStmt:
		err = Walk(n.X, p)

	case *goast.AssignStmt:
		if err = walkExprList(&n.Lhs, p); err != nil {
			return fmt.Errorf("Lhs")
		}
		if err = walkExprList(&n.Rhs, p); err != nil {
			return fmt.Errorf("Rhs")
		}

	case *goast.GoStmt:
		err = Walk(n.Call, p)

	case *goast.DeferStmt:
		err = Walk(n.Call, p)

	case *goast.ReturnStmt:
		err = walkExprList(&n.Results, p)

	case *goast.BranchStmt:
		if n.Label != nil {
			errSecond = Walk(n.Label, p)
		}

	case *goast.BlockStmt:
		if err = walkStmtList(&n.List, p); err != nil {
			return fmt.Errorf("List")
		}

	case *goast.IfStmt:
		if n.Init != nil {
			errSecond = Walk(n.Init, p)
		}
		if err = Walk(n.Cond, p); err != nil {
			return fmt.Errorf("Cond")
		}
		if err = Walk(n.Body, p); err != nil {
			return fmt.Errorf("Body")
		}
		if n.Else != nil {
			errSecond = Walk(n.Else, p)
		}

	case *goast.CaseClause:
		if err = walkExprList(&n.List, p); err != nil {
			return fmt.Errorf("List")
		}
		if err = walkStmtList(&n.Body, p); err != nil {
			return fmt.Errorf("Body")
		}

	case *goast.SwitchStmt:
		if n.Init != nil {
			errSecond = Walk(n.Init, p)
		}
		if n.Tag != nil {
			errSecond = Walk(n.Tag, p)
		}
		err = Walk(n.Body, p)

	case *goast.TypeSwitchStmt:
		if n.Init != nil {
			errSecond = Walk(n.Init, p)
		}
		if err = Walk(n.Assign, p); err != nil {
			return fmt.Errorf("Assign")
		}
		if err = Walk(n.Body, p); err != nil {
			return fmt.Errorf("Body")
		}

	case *goast.CommClause:
		if n.Comm != nil {
			errSecond = Walk(n.Comm, p)
		}
		if err = walkStmtList(&n.Body, p); err != nil {
			return fmt.Errorf("Body")
		}

	case *goast.SelectStmt:
		err = Walk(n.Body, p)

	case *goast.ForStmt:
		if n.Init != nil {
			errSecond = Walk(n.Init, p)
		}
		if n.Cond != nil {
			errSecond = Walk(n.Cond, p)
		}
		if n.Post != nil {
			errSecond = Walk(n.Post, p)
		}
		if err = Walk(n.Body, p); err != nil {
			return fmt.Errorf("Body")
		}

	case *goast.RangeStmt:
		if n.Key != nil {
			errSecond = Walk(n.Key, p)
		}
		if n.Value != nil {
			errSecond = Walk(n.Value, p)
		}
		if err = Walk(n.X, p); err != nil {
			return fmt.Errorf("X")
		}
		if err = Walk(n.Body, p); err != nil {
			return fmt.Errorf("Body")
		}

	// Declarations
	case *goast.ImportSpec:
		if n.Doc != nil {
			errSecond = Walk(n.Doc, p)
		}
		if n.Name != nil {
			errSecond = Walk(n.Name, p)
		}
		if err = Walk(n.Path, p); err != nil {
			return fmt.Errorf("Path")
		}
		if n.Comment != nil {
			errSecond = Walk(n.Comment, p)
		}

	case *goast.ValueSpec:
		if n.Doc != nil {
			errSecond = Walk(n.Doc, p)
		}
		if err = walkIdentList(&n.Names, p); err != nil {
			return fmt.Errorf("Names")
		}
		if n.Type != nil {
			errSecond = Walk(n.Type, p)
		}
		if err = walkExprList(&n.Values, p); err != nil {
			return fmt.Errorf("Values")
		}
		if n.Comment != nil {
			errSecond = Walk(n.Comment, p)
		}

	case *goast.TypeSpec:
		if n.Doc != nil {
			errSecond = Walk(n.Doc, p)
		}
		if err = Walk(n.Name, p); err != nil {
			return fmt.Errorf("Name")
		}
		if err = Walk(n.Type, p); err != nil {
			return fmt.Errorf("Type")
		}
		if n.Comment != nil {
			errSecond = Walk(n.Comment, p)
		}

	case *goast.BadDecl:
		// nothing to do

	case *goast.GenDecl:
		if n.Doc != nil {
			errSecond = Walk(n.Doc, p)
		}
		// err = walkSpecList(&n.Specs, p)
		for i := 0; i < len(n.Specs); i++ {
			_ = Walk(n.Specs[i], p)
		}

	case *goast.FuncDecl:
		if n.Doc != nil {
			errSecond = Walk(n.Doc, p)
		}
		if n.Recv != nil {
			errSecond = Walk(n.Recv, p)
		}
		if err = Walk(n.Name, p); err != nil {
			return fmt.Errorf("Name")
		}
		if err = Walk(n.Type, p); err != nil {
			return fmt.Errorf("Type")
		}
		if n.Body != nil {
			errSecond = Walk(n.Body, p)
		}

	// Files and packages
	case *goast.File:
		if n.Doc != nil {
			errSecond = Walk(n.Doc, p)
		}
		if err = Walk(n.Name, p); err != nil {
			return fmt.Errorf("Name")
		}
		if err = walkDeclList(&n.Decls, p); err != nil {
			return fmt.Errorf("Decls")
		}
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *goast.Package:
		for _, f := range n.Files {
			_ = Walk(f, p)
		}

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	return nil
}

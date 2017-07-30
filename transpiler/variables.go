package transpiler

import (
	"errors"
	"fmt"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
)

func transpileDeclRefExpr(n *ast.DeclRefExpr, p *program.Program) (
	*goast.Ident, string, error) {
	theType := n.Type

	// FIXME: This is for linux to make sure the globals have the right type.
	if n.Name == "stdout" || n.Name == "stdin" || n.Name == "stderr" {
		theType = "FILE *"
	}

	return util.NewIdent(n.Name), theType, nil
}

func getDefaultValueForVar(p *program.Program, a *ast.VarDecl) (
	[]goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
	if len(a.Children) == 0 {
		return nil, "", nil, nil, nil
	}

	defaultValue, defaultValueType, newPre, newPost, err := transpileToExpr(a.Children[0], p)
	if err != nil {
		return nil, defaultValueType, newPre, newPost, err
	}

	var values []goast.Expr
	if !types.IsNullExpr(defaultValue) {
		t, err := types.CastExpr(p, defaultValue, defaultValueType, a.Type)
		if !p.AddMessage(ast.GenerateWarningMessage(err, a)) {
			values = []goast.Expr{t}
		}
	}

	return values, defaultValueType, newPre, newPost, nil
}

func newDeclStmt(a *ast.VarDecl, p *program.Program) (
	*goast.DeclStmt, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	defaultValue, _, newPre, newPost, err := getDefaultValueForVar(p, a)
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// Allocate slice so that it operates like a fixed size array.
	arrayType, arraySize := types.GetArrayTypeAndSize(a.Type)
	if arraySize != -1 && defaultValue == nil {
		goArrayType, err := types.ResolveType(p, arrayType)
		p.AddMessage(ast.GenerateWarningMessage(err, a))

		defaultValue = []goast.Expr{
			util.NewCallExpr(
				"make",
				&goast.ArrayType{
					Elt: util.NewTypeIdent(goArrayType),
				},
				util.NewIntLit(arraySize),
				util.NewIntLit(arraySize),
			),
		}
	}

	t, err := types.ResolveType(p, a.Type)
	p.AddMessage(ast.GenerateWarningMessage(err, a))

	return &goast.DeclStmt{
		Decl: &goast.GenDecl{
			Tok: token.VAR,
			Specs: []goast.Spec{
				&goast.ValueSpec{
					Names:  []*goast.Ident{util.NewIdent(a.Name)},
					Type:   util.NewTypeIdent(t),
					Values: defaultValue,
				},
			},
		},
	}, preStmts, postStmts, nil
}

func transpileDeclStmt(n *ast.DeclStmt, p *program.Program) (
	[]goast.Stmt, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	// There may be more than one variable defined on the same line. With C it
	// is possible for them to have similar but different types, whereas in Go
	// this is not possible. The easiest way around this is to split the
	// variables up into multiple declarations. That is why this function
	// returns one or more DeclStmts.
	decls := []goast.Stmt{}

	for _, c := range n.Children {
		switch a := c.(type) {
		case *ast.RecordDecl:
			// I'm not sure why this is ignored. Maybe we haven't found a
			// situation where this is needed yet?

		case *ast.VarDecl:
			e, newPre, newPost, err := newDeclStmt(a, p)
			if err != nil {
				return nil, nil, nil, err
			}

			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

			decls = append(decls, e)

		case *ast.TypedefDecl:
			p.AddMessage(ast.GenerateWarningMessage(errors.New("cannot use TypedefDecl for DeclStmt"), c))

		default:
			panic(a)
		}
	}

	return decls, preStmts, postStmts, nil
}

func transpileArraySubscriptExpr(n *ast.ArraySubscriptExpr, p *program.Program) (
	*goast.IndexExpr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	children := n.Children
	expression, expressionType, newPre, newPost, err := transpileToExpr(children[0], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	index, _, newPre, newPost, err := transpileToExpr(children[1], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	newType, err := types.GetDereferenceType(expressionType)
	if err != nil {
		message := fmt.Sprintf(
			"Cannot dereference type '%s' for the expression '%s'",
			expressionType, expression)
		return nil, newType, nil, nil, errors.New(message)
	}

	return &goast.IndexExpr{
		X:     expression,
		Index: index,
	}, newType, preStmts, postStmts, nil
}

func transpileMemberExpr(n *ast.MemberExpr, p *program.Program) (
	goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	lhs, lhsType, newPre, newPost, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	lhsResolvedType, err := types.ResolveType(p, lhsType)
	p.AddMessage(ast.GenerateWarningMessage(err, n))

	// lhsType will be something like "struct foo"
	structType := p.GetStruct(lhsType)
	rhs := n.Name
	rhsType := "void *"
	if structType == nil {
		// This case should not happen in the future. Any structs should be
		// either parsed correctly from the source or be manually setup when the
		// parser starts if the struct if hidden or shared between libraries.
		//
		// Some other things to keep in mind:
		//   1. Types need to be stripped of their pointer, 'FILE *' -> 'FILE'.
		//   2. Types may refer to one or more other types in a chain that have
		//      to be resolved before the real field type can be determined.
		err = errors.New("cannot determine type for LHS '" + lhsType +
			"', will use 'void *' for all fields")
		p.AddMessage(ast.GenerateWarningMessage(err, n))
	} else {
		if s, ok := structType.Fields[rhs].(string); ok {
			rhsType = s
		} else {
			err = errors.New("cannot determine type for RHS, will use" +
				" 'void *' for all fields")
			p.AddMessage(ast.GenerateWarningMessage(err, n))
		}
	}

	// FIXME: This is just a hack
	if util.InStrings(lhsResolvedType, []string{"darwin.Float2", "darwin.Double2"}) {
		rhs = util.GetExportedName(rhs)
		rhsType = "int"
	}

	// Construct code for getting value to an union field
	if structType != nil && structType.IsUnion {
		ident := lhs.(*goast.Ident)
		funcName := getFunctionNameForUnionGetter(ident.Name, lhsResolvedType, n.Name)
		resExpr := util.NewCallExpr(funcName)

		return resExpr, rhsType, preStmts, postStmts, nil
	}

	x := lhs
	if n.IsPointer {
		x = &goast.IndexExpr{X: x, Index: util.NewIntLit(0)}
	}

	return &goast.SelectorExpr{
		X:   x,
		Sel: util.NewIdent(rhs),
	}, rhsType, preStmts, postStmts, nil
}

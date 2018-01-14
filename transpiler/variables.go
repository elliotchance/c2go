package transpiler

import (
	"errors"
	"fmt"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
	"go/parser"
	"go/token"
)

// This map is used to rename struct member names.
var structFieldTranslations = map[string]map[string]string{
	"div_t": {
		"quot": "Quot",
		"rem":  "Rem",
	},
	"ldiv_t": {
		"quot": "Quot",
		"rem":  "Rem",
	},
	"lldiv_t": {
		"quot": "Quot",
		"rem":  "Rem",
	},
}

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
	if len(a.Children()) == 0 {
		return nil, "", nil, nil, nil
	}

	// Memory allocation is translated into the Go-style.
	if allocSize := getAllocationSizeNode(a.Children()[0]); allocSize != nil {
		// type
		var t string
		if v, ok := a.Children()[0].(*ast.ImplicitCastExpr); ok {
			t = v.Type
		}
		if v, ok := a.Children()[0].(*ast.CStyleCastExpr); ok {
			t = v.Type
		}
		if t != "" {
			right, newPre, newPost, err := generateAlloc(p, allocSize, t)
			if err != nil {
				p.AddMessage(p.GenerateWarningMessage(err, a))
				return nil, "", nil, nil, err
			}

			return []goast.Expr{right}, t, newPre, newPost, nil
		}
	}

	if va, ok := a.Children()[0].(*ast.VAArgExpr); ok {
		outType, err := types.ResolveType(p, va.Type)
		if err != nil {
			return nil, "", nil, nil, err
		}
		var argsName string
		if a, ok := va.Children()[0].(*ast.ImplicitCastExpr); ok {
			if a, ok := a.Children()[0].(*ast.DeclRefExpr); ok {
				argsName = a.Name
			} else {
				return nil, "", nil, nil, fmt.Errorf("Expect DeclRefExpr for vaar, but we have %T", a)
			}
		} else {
			return nil, "", nil, nil, fmt.Errorf("Expect ImplicitCastExpr for vaar, but we have %T", a)
		}
		src := fmt.Sprintf(`package main
var temp = func() %s {
	var ret %s
	if v, ok := %s[c2goVaListPosition].(int32); ok{
		// for 'rune' type
		ret = %s(v)
	} else {
		ret = %s[c2goVaListPosition].(%s)
	}
	c2goVaListPosition++
	return ret
}()`, outType,
			outType,
			argsName,
			outType,
			argsName, outType)

		// Create the AST by parsing src.
		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, "", src, 0)
		if err != nil {
			return nil, "", nil, nil, err
		}

		expr := f.Decls[0].(*goast.GenDecl).Specs[0].(*goast.ValueSpec).Values
		return expr, va.Type, nil, nil, nil
	}

	defaultValue, defaultValueType, newPre, newPost, err := transpileToExpr(a.Children()[0], p, false)
	if err != nil {
		return nil, defaultValueType, newPre, newPost, err
	}

	var values []goast.Expr
	if !types.IsNullExpr(defaultValue) {
		t, err := types.CastExpr(p, defaultValue, defaultValueType, a.Type)
		if !p.AddMessage(p.GenerateWarningMessage(err, a)) {
			values = append(values, t)
			defaultValueType = a.Type
		}
	}

	return values, defaultValueType, newPre, newPost, nil
}

// GenerateFuncType in according to types
/*
Type: *ast.FuncType {
.  Func: 13:7
.  Params: *ast.FieldList {
.  .  Opening: 13:12
.  .  List: []*ast.Field (len = 2) {
.  .  .  0: *ast.Field {
.  .  .  .  Type: *ast.Ident {
.  .  .  .  .  NamePos: 13:13
.  .  .  .  .  Name: "int"
.  .  .  .  }
.  .  .  }
.  .  .  1: *ast.Field {
.  .  .  .  Type: *ast.Ident {
.  .  .  .  .  NamePos: 13:17
.  .  .  .  .  Name: "int"
.  .  .  .  }
.  .  .  }
.  .  }
.  .  Closing: 13:20
.  }
.  Results: *ast.FieldList {
.  .  Opening: -
.  .  List: []*ast.Field (len = 1) {
.  .  .  0: *ast.Field {
.  .  .  .  Type: *ast.Ident {
.  .  .  .  .  NamePos: 13:21
.  .  .  .  .  Name: "string"
.  .  .  .  }
.  .  .  }
.  .  }
.  .  Closing: -
.  }
}
*/
func GenerateFuncType(fields, returns []string) *goast.FuncType {
	var ft goast.FuncType
	{
		var fieldList goast.FieldList
		fieldList.Opening = 1
		fieldList.Closing = 2
		for i := range fields {
			fieldList.List = append(fieldList.List, &goast.Field{Type: &goast.Ident{Name: fields[i]}})
		}
		ft.Params = &fieldList
	}
	{
		var fieldList goast.FieldList
		for i := range returns {
			fieldList.List = append(fieldList.List, &goast.Field{Type: &goast.Ident{Name: returns[i]}})
		}
		ft.Results = &fieldList
	}
	return &ft
}

func transpileInitListExpr(e *ast.InitListExpr, p *program.Program) (goast.Expr, string, error) {
	resp := []goast.Expr{}
	var hasArrayFiller = false

	for _, node := range e.Children() {
		// Skip ArrayFiller
		if _, ok := node.(*ast.ArrayFiller); ok {
			hasArrayFiller = true
			continue
		}

		var expr goast.Expr
		var err error
		expr, _, _, _, err = transpileToExpr(node, p, true)
		if err != nil {
			return nil, "", err
		}

		resp = append(resp, expr)
	}

	var t goast.Expr
	var cTypeString string

	arrayType, arraySize := types.GetArrayTypeAndSize(e.Type1)
	if arraySize != -1 {
		goArrayType, err := types.ResolveType(p, arrayType)
		p.AddMessage(p.GenerateWarningMessage(err, e))

		cTypeString = fmt.Sprintf("%s[%d]", arrayType, arraySize)

		if hasArrayFiller {
			t = &goast.ArrayType{
				Elt: &goast.Ident{
					Name: goArrayType,
				},
				Len: util.NewIntLit(arraySize),
			}

			// Array fillers do not work with slices.
			// We initialize the array first, then convert to a slice.
			// For example: (&[4]int{1,2})[:]
			return &goast.SliceExpr{
				X: &goast.ParenExpr{
					X: &goast.UnaryExpr{
						Op: token.AND,
						X: &goast.CompositeLit{
							Type: t,
							Elts: resp,
						},
					},
				},
			}, cTypeString, nil
		}

		t = &goast.ArrayType{
			Elt: &goast.Ident{
				Name: goArrayType,
			},
		}
	} else {
		goType, err := types.ResolveType(p, e.Type1)
		if err != nil {
			return nil, "", err
		}

		t = &goast.Ident{
			Name: goType,
		}

		cTypeString = e.Type1
	}

	return &goast.CompositeLit{
		Type: t,
		Elts: resp,
	}, cTypeString, nil
}

func transpileDeclStmt(n *ast.DeclStmt, p *program.Program) (stmts []goast.Stmt, err error) {
	if len(n.Children()) == 0 {
		return
	}
	var tud ast.TranslationUnitDecl
	tud.ChildNodes = n.Children()
	var decls []goast.Decl
	decls, err = transpileToNode(&tud, p)
	if err != nil {
		p.AddMessage(p.GenerateErrorMessage(err, n))
		err = nil
	}
	stmts = convertDeclToStmt(decls)

	return
}

func transpileArraySubscriptExpr(n *ast.ArraySubscriptExpr, p *program.Program) (
	_ *goast.IndexExpr, theType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile ArraySubscriptExpr. err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	children := n.Children()
	expression, expressionType, newPre, newPost, err := transpileToExpr(children[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	index, _, newPre, newPost, err := transpileToExpr(children[1], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	theType, err = types.GetDereferenceType(expressionType)
	if err != nil {
		message := fmt.Sprintf(
			"Cannot dereference type '%s' for the expression '%#v'. err = %v",
			expressionType, expression, err)
		return nil, theType, nil, nil, errors.New(message)
	}

	return &goast.IndexExpr{
		X:     expression,
		Index: index,
	}, n.Type, preStmts, postStmts, nil
}

func transpileMemberExpr(n *ast.MemberExpr, p *program.Program) (
	_ goast.Expr, _ string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {

	n.Type = types.GenerateCorrectType(n.Type)
	n.Type2 = types.GenerateCorrectType(n.Type2)

	lhs, lhsType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	lhsType = types.GenerateCorrectType(lhsType)

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	lhsResolvedType, err := types.ResolveType(p, lhsType)
	p.AddMessage(p.GenerateWarningMessage(err, n))

	// lhsType will be something like "struct foo"
	structType := p.GetStruct(lhsType)
	// added for support "struct typedef"
	if structType == nil {
		structType = p.GetStruct("struct " + lhsType)
	}
	// added for support "union typedef"
	if structType == nil {
		structType = p.GetStruct("union " + lhsType)
	}
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
		err = fmt.Errorf("cannot determine type for LHS '%v'"+
			", will use 'void *' for all fields. Is lvalue = %v", lhsType, n.IsLvalue)
		p.AddMessage(p.GenerateWarningMessage(err, n))
	} else {
		if s, ok := structType.Fields[rhs].(string); ok {
			rhsType = s
		} else {
			err = fmt.Errorf("cannot determine type for RHS '%v', will use"+
				" 'void *' for all fields. Is lvalue = %v", rhs, n.IsLvalue)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}

	// FIXME: This is just a hack
	if util.InStrings(lhsResolvedType, []string{"darwin.Float2", "darwin.Double2"}) {
		rhs = util.GetExportedName(rhs)
		rhsType = "int"
	}

	// Construct code for getting value to an union field
	if structType != nil && structType.IsUnion {
		var resExpr goast.Expr

		switch t := lhs.(type) {
		case *goast.Ident:
			funcName := getFunctionNameForUnionGetter(t.Name, lhsResolvedType, n.Name)
			resExpr = util.NewCallExpr(funcName)
		case *goast.SelectorExpr:
			funcName := getFunctionNameForUnionGetter("", lhsResolvedType, n.Name)
			if id, ok := t.X.(*goast.Ident); ok {
				funcName = id.Name + "." + t.Sel.Name + funcName
			}
			resExpr = &goast.CallExpr{
				Fun:  goast.NewIdent(funcName),
				Args: nil,
			}
		}

		return resExpr, rhsType, preStmts, postStmts, nil
	}

	x := lhs
	if n.IsPointer {
		x = &goast.IndexExpr{X: x, Index: util.NewIntLit(0)}
	}

	// Check for member name translation.
	if member, ok := structFieldTranslations[lhsType]; ok {
		if alias, ok := member[rhs]; ok {
			rhs = alias
		}
	}

	// anonymous struct member?
	if rhs == "" {
		rhs = "anon"
	}

	return &goast.SelectorExpr{
		X:   x,
		Sel: util.NewIdent(rhs),
	}, rhsType, preStmts, postStmts, nil
}

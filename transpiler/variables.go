package transpiler

import (
	"fmt"
	"strings"

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
	"struct tm": {
		"tm_sec":   "Tm_sec",
		"tm_min":   "Tm_min",
		"tm_hour":  "Tm_hour",
		"tm_mday":  "Tm_mday",
		"tm_mon":   "Tm_mon",
		"tm_year":  "Tm_year",
		"tm_wday":  "Tm_wday",
		"tm_yday":  "Tm_yday",
		"tm_isdst": "Tm_isdst",
	},
}

func transpileDeclRefExpr(n *ast.DeclRefExpr, p *program.Program) (
	expr *goast.Ident, exprType string, err error) {

	if n.For == "EnumConstant" {
		// clang don`t show enum constant with enum type,
		// so we have to use hack for repair the type
		if v, ok := p.EnumConstantToEnum[n.Name]; ok {
			expr, exprType, err = util.NewIdent(n.Name), v, nil
			return
		}
	}

	theType := n.Type

	// FIXME: This is for linux to make sure the globals have the right type.
	if n.Name == "stdout" || n.Name == "stdin" || n.Name == "stderr" {
		theType = "FILE *"
	}

	return util.NewIdent(n.Name), theType, nil
}

func getDefaultValueForVar(p *program.Program, a *ast.VarDecl) (
	_ []goast.Expr, _ string, _ []goast.Stmt, _ []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot getDefaultValueForVar : err = %v", err)
		}
	}()
	if len(a.Children()) == 0 {
		return nil, "", nil, nil, nil
	}

	// Memory allocation is translated into the Go-style.
	if allocSize := getAllocationSizeNode(p, a.Children()[0]); allocSize != nil {
		// type
		var t string
		if v, ok := a.Children()[0].(*ast.ImplicitCastExpr); ok {
			t = v.Type
		}
		if v, ok := a.Children()[0].(*ast.CStyleCastExpr); ok {
			t = v.Type
		}
		if v, ok := a.Children()[0].(*ast.CallExpr); ok {
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
		if a, ok := va.Children()[0].(*ast.ImplicitCastExpr); ok {
		} else {
			return nil, "", nil, nil, fmt.Errorf("Expect ImplicitCastExpr for vaar, but we have %T", a)
		}
		src := fmt.Sprintf(`package main
var temp = func() %s {
	var ret %s
	if v, ok := c2goVaList.Args[c2goVaList.Pos].(int32); ok{
		// for 'rune' type
		ret = %s(v)
	} else {
		ret = c2goVaList.Args[c2goVaList.Pos].(%s)
	}
	c2goVaList.Pos++
	return ret
}()`, outType,
			outType,
			outType,
			outType)

		// Create the AST by parsing src.
		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, "", src, 0)
		if err != nil {
			return nil, "", nil, nil, err
		}

		expr := f.Decls[0].(*goast.GenDecl).Specs[0].(*goast.ValueSpec).Values
		return expr, va.Type, nil, nil, nil
	}

	defaultValue, defaultValueType, newPre, newPost, err := atomicOperation(a.Children()[0], p)
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
	e.Type1 = types.GenerateCorrectType(e.Type1)
	e.Type2 = types.GenerateCorrectType(e.Type2)

	var goType string
	arrayType, arraySize := types.GetArrayTypeAndSize(e.Type1)
	if arraySize != -1 {
		goArrayType, err := types.ResolveType(p, arrayType)
		if err == nil {
			goType = goArrayType
		}
	} else {
		goType2, err := types.ResolveType(p, e.Type1)
		if err == nil {
			goType = goType2
		}
	}
	var goStruct *program.Struct
	if e.Type1 == e.Type2 {
		goStruct = p.GetStruct(goType)
		if goStruct == nil {
			goStruct = p.GetStruct("struct " + goType)
		}
	}
	fieldIndex := 0

	for _, node := range e.Children() {
		// Skip ArrayFiller
		if _, ok := node.(*ast.ArrayFiller); ok {
			hasArrayFiller = true
			continue
		}

		var expr goast.Expr
		var exprType string
		var err error
		expr, exprType, _, _, err = transpileToExpr(node, p, true)
		if err != nil {
			return nil, "", err
		}
		if goStruct != nil {
			if fieldIndex >= len(goStruct.FieldNames) {
				// index out of range
				goto CONTINUE_INIT
			}
			fn := goStruct.FieldNames[fieldIndex]
			if _, ok := goStruct.Fields[fn]; !ok {
				// field name not in map
				goto CONTINUE_INIT
			}
			if field, ok := goStruct.Fields[goStruct.FieldNames[fieldIndex]].(string); ok {
				expr2, err := types.CastExpr(p, expr, exprType, field)
				if err == nil {
					expr = expr2
				}
			}
			fieldIndex++
		}
	CONTINUE_INIT:
		resp = append(resp, expr)
	}

	var t goast.Expr
	var cTypeString string

	arrayType, arraySize = types.GetArrayTypeAndSize(e.Type1)
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

func transpileArraySubscriptExpr(n *ast.ArraySubscriptExpr, p *program.Program, exprIsStmt bool) (
	_ goast.Expr, theType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile ArraySubscriptExpr. err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	children := n.Children()

	expression, leftType, newPre, newPost, err := transpileToExpr(children[0], p, exprIsStmt)
	if err != nil {
		return nil, "", nil, nil, err
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	index, indexType, newPre, newPost, err := atomicOperation(children[1], p)
	if err != nil {
		return nil, "", nil, nil, err
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if se, ok := expression.(*goast.SliceExpr); ok && se.High == nil && se.Low == nil && se.Max == nil {
		// simplify the expression
		expression = se.X
	}

	isConst, indexInt := util.EvaluateConstExpr(index)
	if isConst && indexInt < 0 {
		indexInt = -indexInt
		expression, leftType, newPre, newPost, err =
			pointerArithmetic(p, expression, leftType, util.NewIntLit(int(indexInt)), "int", token.SUB)
		return &goast.StarExpr{
			X: expression,
		}, n.Type, newPre, newPost, err
	} else {
		resolvedLeftType, err := types.ResolveType(p, leftType)
		if err != nil {
			return nil, "", nil, nil, err
		}
		if types.IsPurePointer(p, resolvedLeftType) {
			if !isConst || indexInt != 0 {
				expression, leftType, newPre, newPost, err =
					pointerArithmetic(p, expression, leftType, index, indexType, token.ADD)
			}
			return &goast.StarExpr{
				X: expression,
			}, n.Type, newPre, newPost, err
		}
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
	lhsType = types.CleanCType(lhsType)

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

	x := lhs
	if n.IsPointer {
		x = &goast.ParenExpr{
			X: &goast.StarExpr{X: x},
		}
	}

	// Check for member name translation.
	lhsType = strings.TrimSpace(lhsType)
	if lhsType[len(lhsType)-1] == '*' {
		lhsType = lhsType[:len(lhsType)-len(" *")]
	}
	if member, ok := structFieldTranslations[lhsType]; ok {
		if alias, ok := member[rhs]; ok {
			rhs = alias
		}
	}

	// anonymous struct member?
	if rhs == "" {
		rhs = "anon"
	}

	if isUnionMemberExpr(p, n) {
		return &goast.ParenExpr{
			Lparen: 1,
			X: &goast.StarExpr{
				Star: 1,
				X: &goast.CallExpr{
					Fun: &goast.SelectorExpr{
						X:   x,
						Sel: util.NewIdent(rhs),
					},
					Lparen: 1,
				},
			},
		}, n.Type, preStmts, postStmts, nil
	}

	_ = rhsType

	return &goast.SelectorExpr{
		X:   x,
		Sel: util.NewIdent(rhs),
	}, n.Type, preStmts, postStmts, nil
}

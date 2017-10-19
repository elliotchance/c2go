// This file contains transpiling for enums.

package transpiler

import (
	"go/token"

	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

// ctypeEnumValue generates a specific expression for values used by some
// constants in ctype.h. This is to get around an issue that the real values
// need to be evaulated by the compiler; which c2go does not yet do.
//
// TOOD: Ability to evaluate constant expressions at compile time
// https://github.com/elliotchance/c2go/issues/77
func ctypeEnumValue(value int, t token.Token) goast.Expr {
	// Produces an expression like: ((1 << (0)) << 8)
	return &goast.ParenExpr{
		X: util.NewBinaryExpr(
			&goast.ParenExpr{
				X: util.NewBinaryExpr(
					util.NewIntLit(1),
					token.SHL,
					util.NewIntLit(value),
					"int",
					false,
				),
			},
			t,
			util.NewIntLit(8),
			"int",
			false,
		),
	}
}

func transpileEnumConstantDecl(p *program.Program, n *ast.EnumConstantDecl) (
	*goast.ValueSpec, []goast.Stmt, []goast.Stmt) {
	var value goast.Expr = util.NewIdent("iota")
	valueType := "int"
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	// Special cases for linux ctype.h. See the description for the
	// ctypeEnumValue() function.
	switch n.Name {
	case "_ISupper":
		value = ctypeEnumValue(0, token.SHL) // "((1 << (0)) << 8)"
		valueType = "uint16"
	case "_ISlower":
		value = ctypeEnumValue(1, token.SHL) // "((1 << (1)) << 8)"
		valueType = "uint16"
	case "_ISalpha":
		value = ctypeEnumValue(2, token.SHL) // "((1 << (2)) << 8)"
		valueType = "uint16"
	case "_ISdigit":
		value = ctypeEnumValue(3, token.SHL) // "((1 << (3)) << 8)"
		valueType = "uint16"
	case "_ISxdigit":
		value = ctypeEnumValue(4, token.SHL) // "((1 << (4)) << 8)"
		valueType = "uint16"
	case "_ISspace":
		value = ctypeEnumValue(5, token.SHL) // "((1 << (5)) << 8)"
		valueType = "uint16"
	case "_ISprint":
		value = ctypeEnumValue(6, token.SHL) // "((1 << (6)) << 8)"
		valueType = "uint16"
	case "_ISgraph":
		value = ctypeEnumValue(7, token.SHL) // "((1 << (7)) << 8)"
		valueType = "uint16"
	case "_ISblank":
		value = ctypeEnumValue(8, token.SHR) // "((1 << (8)) >> 8)"
		valueType = "uint16"
	case "_IScntrl":
		value = ctypeEnumValue(9, token.SHR) // "((1 << (9)) >> 8)"
		valueType = "uint16"
	case "_ISpunct":
		value = ctypeEnumValue(10, token.SHR) // "((1 << (10)) >> 8)"
		valueType = "uint16"
	case "_ISalnum":
		value = ctypeEnumValue(11, token.SHR) // "((1 << (11)) >> 8)"
		valueType = "uint16"
	default:
		if len(n.Children()) > 0 {
			var err error
			value, _, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, false)
			if err != nil {
				panic(err)
			}
		}
	}

	return &goast.ValueSpec{
		Names:  []*goast.Ident{util.NewIdent(n.Name)},
		Type:   util.NewTypeIdent(valueType),
		Values: []goast.Expr{value},
	}, preStmts, postStmts
}

func transpileEnumDecl(p *program.Program, n *ast.EnumDecl) error {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	if n.Name == "" {
		for _, c := range n.Children() {
			e, newPre, newPost := transpileEnumConstantDecl(p, c.(*ast.EnumConstantDecl))
			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

			p.File.Decls = append(p.File.Decls, &goast.GenDecl{
				Tok: token.CONST,
				Specs: []goast.Spec{
					e,
				},
			})
		}
		return nil
	}

	theType, err := types.ResolveType(p, "int")
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}
	baseType := util.NewTypeIdent(theType)

	bt := goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: &goast.Ident{
					Name: n.Name,
					Obj:  goast.NewObj(goast.Typ, n.Name),
				},
				Type: baseType,
			},
		},
	}
	// Registration new type in program.Program
	if !p.IsTypeAlreadyDefined(n.Name) {
		p.DefineType(n.Name)
	}

	p.File.Decls = append(p.File.Decls, &bt)

	baseType2 := util.NewTypeIdent(n.Name)

	decl := &goast.GenDecl{
		Tok: token.CONST,
	}

	for i, c := range n.Children() {
		e, newPre, newPost := transpileEnumConstantDecl(p, c.(*ast.EnumConstantDecl))
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		e.Names[0].Obj = goast.NewObj(goast.Con, e.Names[0].Name)

		if i > 0 {
			e.Type = nil
			e.Values = nil
		}

		if i == 0 {
			e.Type = baseType2
			if t, ok := e.Type.(*goast.Ident); ok {
				t.Obj = &goast.Object{
					Name: n.Name,
					Kind: goast.Typ,
					Decl: &goast.TypeSpec{
						Name: &goast.Ident{
							Name: n.Name,
						},
						Type: &goast.Ident{
							Name: "int",
						},
					},
				}
			}
		}

		e.Names[0].Obj.Data = i
		decl.Specs = append(decl.Specs, e)
	}

	decl.Specs[0].(*goast.ValueSpec).Names[0].Obj.Decl = nil
	decl.Specs[1].(*goast.ValueSpec).Names[0].Obj = nil
	decl.Specs[2].(*goast.ValueSpec).Names[0].Obj = nil

	decl.Lparen = 1

	p.File.Decls = append(p.File.Decls, decl)

	return nil
}

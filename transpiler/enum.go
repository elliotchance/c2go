// This file contains transpiling for enums.

package transpiler

import (
	"go/token"

	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

// ctypeEnumValue generates a specific expression for values used by some
// constants in ctype.h. This is to get around an issue that the real values
// need to be evaulated by the compiler; which c2go does not yet do.
//
// TOOD: Ability to evaluate constant expressions at compile time
// https://github.com/elliotchance/c2go/issues/77
func ctypeEnumValue(value string, t token.Token) goast.Expr {
	// Produces an expression like: ((1 << (0)) << 8)
	return &goast.ParenExpr{
		X: &goast.BinaryExpr{
			X: &goast.ParenExpr{
				X: &goast.BinaryExpr{
					X: &goast.BasicLit{
						Kind:  token.INT,
						Value: "1",
					},
					Op: token.SHL,
					Y: &goast.BasicLit{
						Kind:  token.INT,
						Value: value,
					},
				},
			},
			Op: t,
			Y: &goast.BasicLit{
				Kind:  token.INT,
				Value: "8",
			},
		},
	}
}

func transpileEnumConstantDecl(p *program.Program, n *ast.EnumConstantDecl) *goast.ValueSpec {
	var value goast.Expr = goast.NewIdent("iota")
	valueType := "int"

	// Special cases for linux ctype.h. See the description for the
	// ctypeEnumValue() function.
	switch n.Name {
	case "_ISupper":
		value = ctypeEnumValue("0", token.SHL) // "((1 << (0)) << 8)"
		valueType = "uint16"
	case "_ISlower":
		value = ctypeEnumValue("1", token.SHL) // "((1 << (1)) << 8)"
		valueType = "uint16"
	case "_ISalpha":
		value = ctypeEnumValue("2", token.SHL) // "((1 << (2)) << 8)"
		valueType = "uint16"
	case "_ISdigit":
		value = ctypeEnumValue("3", token.SHL) // "((1 << (3)) << 8)"
		valueType = "uint16"
	case "_ISxdigit":
		value = ctypeEnumValue("4", token.SHL) // "((1 << (4)) << 8)"
		valueType = "uint16"
	case "_ISspace":
		value = ctypeEnumValue("5", token.SHL) // "((1 << (5)) << 8)"
		valueType = "uint16"
	case "_ISprint":
		value = ctypeEnumValue("6", token.SHL) // "((1 << (6)) << 8)"
		valueType = "uint16"
	case "_ISgraph":
		value = ctypeEnumValue("7", token.SHL) // "((1 << (7)) << 8)"
		valueType = "uint16"
	case "_ISblank":
		value = ctypeEnumValue("8", token.SHR) // "((1 << (8)) >> 8)"
		valueType = "uint16"
	case "_IScntrl":
		value = ctypeEnumValue("9", token.SHR) // "((1 << (9)) >> 8)"
		valueType = "uint16"
	case "_ISpunct":
		value = ctypeEnumValue("10", token.SHR) // "((1 << (10)) >> 8)"
		valueType = "uint16"
	case "_ISalnum":
		value = ctypeEnumValue("11", token.SHR) // "((1 << (11)) >> 8)"
		valueType = "uint16"
	default:
		if len(n.Children) > 0 {
			var err error
			value, _, err = transpileToExpr(n.Children[0], p)
			if err != nil {
				panic(err)
			}
		}
	}

	return &goast.ValueSpec{
		Names:  []*goast.Ident{goast.NewIdent(n.Name)},
		Type:   goast.NewIdent(valueType),
		Values: []goast.Expr{value},
	}
}

func transpileEnumDecl(p *program.Program, n *ast.EnumDecl) error {
	for _, c := range n.Children {
		p.File.Decls = append(p.File.Decls, &goast.GenDecl{
			Tok: token.CONST,
			Specs: []goast.Spec{
				transpileEnumConstantDecl(p, c.(*ast.EnumConstantDecl)),
			},
		})
	}

	return nil
}

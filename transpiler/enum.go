// This file contains transpiling for enums.

package transpiler

import (
	"fmt"
	"go/token"
	"strconv"

	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

// ctypeEnumValue generates a specific expression for values used by some
// constants in ctype.h. This is to get around an issue that the real values
// need to be evaluated by the compiler; which c2go does not yet do.
//
// TODO: Ability to evaluate constant expressions at compile time
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
	valueType := "int32"
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
		Doc:    p.GetMessageComments(),
	}, preStmts, postStmts
}

func transpileEnumDecl(p *program.Program, n *ast.EnumDecl) (decls []goast.Decl, err error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	// For case `enum` without name
	if n.Name == "" {
		// create all EnumConstant like just constants
		var counter int
		for _, child := range n.Children() {
			if c, ok := child.(*ast.EnumConstantDecl); ok {
				var (
					e       goast.Spec
					newPre  []goast.Stmt
					newPost []goast.Stmt
					val     *goast.ValueSpec
				)
				val, newPre, newPost = transpileEnumConstantDecl(p, c)

				if len(newPre) > 0 || len(newPost) > 0 {
					p.AddMessage(p.GenerateWarningMessage(fmt.Errorf("Check - added in code : (%d)(%d)", len(newPre), len(newPost)), n))
				}

				preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

				parseEnumBasicLit := func(b *goast.BasicLit) (_ goast.Spec, counter int, err error) {
					value, err := strconv.Atoi(b.Value)
					if err != nil {
						err = fmt.Errorf("Cannot parse '%s' in BasicLit", b.Value)
						return
					}
					return &goast.ValueSpec{
						Names:  []*goast.Ident{{Name: c.Name}},
						Values: []goast.Expr{&goast.BasicLit{Kind: token.INT, Value: b.Value}},
						Type:   val.Type,
						Doc:    p.GetMessageComments(),
					}, value, nil
				}

				switch v := val.Values[0].(type) {
				case *goast.Ident:
					e = &goast.ValueSpec{
						Names:  []*goast.Ident{{Name: c.Name}},
						Values: []goast.Expr{&goast.BasicLit{Kind: token.INT, Value: strconv.Itoa(counter)}},
						Type:   val.Type,
						Doc:    p.GetMessageComments(),
					}
					counter++

				case *goast.BasicLit:
					var value int
					e, value, err = parseEnumBasicLit(v)
					if err != nil {
						e = val
						counter++
						p.AddMessage(p.GenerateWarningMessage(
							fmt.Errorf("Cannot parse '%s' in BasicLit", v.Value), n))
						break
					}
					counter = value
					counter++

				case *goast.CallExpr:
					e = val
					if id, ok := v.Fun.(*goast.Ident); !ok || len(v.Args) != 1 ||
						!types.IsGoIntegerType(id.Name) {
						p.AddMessage(p.GenerateWarningMessage(fmt.Errorf("Add support of continues counter for type : *goast.CallExpr != integer cast"), n))
						break
					}
					if lit, ok := v.Args[0].(*goast.BasicLit); ok {
						var value int
						e, value, err = parseEnumBasicLit(lit)
						if err != nil {
							e = val
							counter++
							p.AddMessage(p.GenerateWarningMessage(
								fmt.Errorf("Cannot parse '%s' in BasicLit", lit.Value), n))
							break
						}
						counter = value
						counter++
					} else {
						p.AddMessage(p.GenerateWarningMessage(fmt.Errorf("Add support of continues counter for type : *goast.CallExpr (integer cast) with argument type : %T", v), n))
					}
				default:
					e = val
					p.AddMessage(p.GenerateWarningMessage(fmt.Errorf("Add support of continues counter for type : %T", v), n))
				}

				decls = append(decls, &goast.GenDecl{
					Tok: token.CONST,
					Specs: []goast.Spec{
						e,
					},
				})
			}
		}
		err = nil
		return
	}

	// For case `enum` with name
	theType, err := types.ResolveType(p, "int")
	if err != nil {
		// by defaults enum in C is INT
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}

	// Create alias of enum for int
	decls = append(decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: &goast.Ident{
					Name: n.Name,
					Obj:  goast.NewObj(goast.Typ, n.Name),
				},
				Type: util.NewTypeIdent(theType),
			},
		},
	})

	// Registration new type in program.Program
	if !p.IsTypeAlreadyDefined(n.Name) {
		p.DefineType(n.Name)
	}

	baseType := util.NewTypeIdent(n.Name)

	decl := &goast.GenDecl{
		Tok: token.CONST,
	}

	// counter for replace iota
	var counter int
	var i int
	for _, c := range n.Children() {
		if _, ok := c.(*ast.EnumConstantDecl); !ok {
			// add for avoid comments elements
			continue
		}
		e, newPre, newPost := transpileEnumConstantDecl(p, c.(*ast.EnumConstantDecl))
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		e.Names[0].Obj = goast.NewObj(goast.Con, e.Names[0].Name)

		if i > 0 {
			e.Type = nil
			e.Values = nil
		}

		if i == 0 {
			e.Type = baseType
			if t, ok := e.Type.(*goast.Ident); ok {
				t.Obj = &goast.Object{
					Name: n.Name,
					Kind: goast.Typ,
					Decl: &goast.TypeSpec{
						Name: &goast.Ident{
							Name: n.Name,
						},
						Type: &goast.Ident{
							Name: "int", // enum in C is "INT" by default
						},
					},
				}
			}
		}

		if len(c.(*ast.EnumConstantDecl).ChildNodes) > 0 {
			if integr, ok := c.(*ast.EnumConstantDecl).ChildNodes[0].(*ast.IntegerLiteral); ok {
				is, err := strconv.ParseInt(integr.Value, 10, 64)
				if err != nil {
					p.AddMessage(p.GenerateWarningMessage(err, n))
				}
				counter = int(is)
			}
		}

		// Insert value of constants
		e.Values = []goast.Expr{
			&goast.BasicLit{
				Kind:  token.INT,
				Value: strconv.Itoa(counter),
			},
		}

		// Position inside (....), it is
		// not value of constants
		e.Names[0].Obj.Data = i
		counter++

		decl.Specs = append(decl.Specs, e)

		// registration of enum constants
		p.EnumConstantToEnum[e.Names[0].Name] = "enum " + n.Name

		// calculate next position without comments
		i++
	}

	// important value for creating (.....)
	// with constants inside
	decl.Lparen = 1
	decl.Rparen = 2

	decls = append(decls, decl)

	err = nil
	return
}

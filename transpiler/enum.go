// This file contains transpiling for enums.

package transpiler

import (
	"fmt"
	"go/parser"
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

/*
Example of AST tree:
|-EnumDecl 0x32fb5a0 <enum.c:3:1, col:45> col:6 week
| |-EnumConstantDecl 0x32fb650 <col:11> col:11 Mon 'int'
| |-EnumConstantDecl 0x32fb6a0 <col:16> col:16 Tue 'int'
*/
/*
type w int

const (
 A w = iota
 B
)
   23  .  .  1: *ast.GenDecl {
   24  .  .  .  TokPos: 7:1
   25  .  .  .  Tok: type
   26  .  .  .  Lparen: -
   27  .  .  .  Specs: []ast.Spec (len = 1) {
   28  .  .  .  .  0: *ast.TypeSpec {
   29  .  .  .  .  .  Name: *ast.Ident {
   30  .  .  .  .  .  .  NamePos: 7:6
   31  .  .  .  .  .  .  Name: "w"
   32  .  .  .  .  .  .  Obj: *ast.Object {
   33  .  .  .  .  .  .  .  Kind: type
   34  .  .  .  .  .  .  .  Name: "w"
   35  .  .  .  .  .  .  .  Decl: *(obj @ 28)
   36  .  .  .  .  .  .  }
   37  .  .  .  .  .  }
   38  .  .  .  .  .  Type: *ast.Ident {
   39  .  .  .  .  .  .  NamePos: 7:8
   40  .  .  .  .  .  .  Name: "int"
   41  .  .  .  .  .  }
   42  .  .  .  .  }
   43  .  .  .  }
   44  .  .  .  Rparen: -
   45  .  .  }

   46  .  .  2: *ast.GenDecl {
   47  .  .  .  TokPos: 9:1
   48  .  .  .  Tok: const
   49  .  .  .  Lparen: 9:7
   50  .  .  .  Specs: []ast.Spec (len = 2) {
   51  .  .  .  .  0: *ast.ValueSpec {
   52  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
   53  .  .  .  .  .  .  0: *ast.Ident {
   54  .  .  .  .  .  .  .  NamePos: 10:2
   55  .  .  .  .  .  .  .  Name: "A"
   56  .  .  .  .  .  .  .  Obj: *ast.Object {
   57  .  .  .  .  .  .  .  .  Kind: const
   58  .  .  .  .  .  .  .  .  Name: "A"
   59  .  .  .  .  .  .  .  .  Decl: *(obj @ 51)
   60  .  .  .  .  .  .  .  .  Data: 0
   61  .  .  .  .  .  .  .  }
   62  .  .  .  .  .  .  }
   63  .  .  .  .  .  }
   64  .  .  .  .  .  Type: *ast.Ident {
   65  .  .  .  .  .  .  NamePos: 10:4
   66  .  .  .  .  .  .  Name: "w"
   67  .  .  .  .  .  .  Obj: *(obj @ 32)
   68  .  .  .  .  .  }
   69  .  .  .  .  .  Values: []ast.Expr (len = 1) {
   70  .  .  .  .  .  .  0: *ast.Ident {
   71  .  .  .  .  .  .  .  NamePos: 10:8
   72  .  .  .  .  .  .  .  Name: "iota"
   73  .  .  .  .  .  .  }
   74  .  .  .  .  .  }
   75  .  .  .  .  }
   76  .  .  .  .  1: *ast.ValueSpec {
   77  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
   78  .  .  .  .  .  .  0: *ast.Ident {
   79  .  .  .  .  .  .  .  NamePos: 11:2
   80  .  .  .  .  .  .  .  Name: "B"
   81  .  .  .  .  .  .  .  Obj: *ast.Object {
   82  .  .  .  .  .  .  .  .  Kind: const
   83  .  .  .  .  .  .  .  .  Name: "B"
   84  .  .  .  .  .  .  .  .  Decl: *(obj @ 76)
   85  .  .  .  .  .  .  .  .  Data: 1
   86  .  .  .  .  .  .  .  }
   87  .  .  .  .  .  .  }
   88  .  .  .  .  .  }
   89  .  .  .  .  }
   90  .  .  .  }
   91  .  .  .  Rparen: 12:1
   92  .  .  }
*/
func transpileEnumDecl(p *program.Program, n *ast.EnumDecl) error {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	theType, err := types.ResolveType(p, "int")
	//p.AddMessage(p.GenerateWarningMessage(err, n))
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
	/*
		goast.Fprint(os.Stdout, token.NewFileSet(), bt, func(name string, value reflect.Value) bool {
			return true
		})
	*/
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
			(e.Type.(*goast.Ident)).Obj = &goast.Object{
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

		e.Names[0].Obj.Data = i
		decl.Specs = append(decl.Specs, e)
	}
	/*
	   54  .  .  1: *ast.ValueSpec {
	   55  .  .  .  Doc: nil
	   56  .  .  .  Names: []*ast.Ident (len = 1) {
	   57  .  .  .  .  0: *ast.Ident {
	   58  .  .  .  .  .  NamePos: -
	   59  .  .  .  .  .  Name: "wwwb"
	   60  .  .  .  .  .  Obj: *ast.Object {
	   61  .  .  .  .  .  .  Kind: const
	   62  .  .  .  .  .  .  Name: "wwwb"
	   63  .  .  .  .  .  .  Decl: nil
	   64  .  .  .  .  .  .  Data: 1
	   65  .  .  .  .  .  .  Type: nil
	   66  .  .  .  .  .  }
	   67  .  .  .  .  }
	   68  .  .  .  }
	   69  .  .  .  Type: nil
	   70  .  .  .  Values: nil
	   71  .  .  .  Comment: nil
	   72  .  .  }
	*/
	decl.Specs = append(decl.Specs, &goast.ValueSpec{
		Names: []*goast.Ident{
			&goast.Ident{
				Name: "asdasdasd",
				Obj: &goast.Object{
					Kind: goast.Con,
					Name: "asdasdasd",
				},
			},
		},
	})
	// fmt.Println("ABDSD")

	/*
		goast.Fprint(os.Stdout, token.NewFileSet(), decl, func(name string, value reflect.Value) bool {
			return true
		})
	*/
	fmt.Println("\n\n\n\n=================")

	src := `package main

type WWW2 int

const (
	wwwa2 WWW2 = iota
	wwwb2
	wwwc2
)
`

	// Create the AST by parsing src.
	fset2 := token.NewFileSet() // positions are relative to fset
	f3, err := parser.ParseFile(fset2, "", src, 0)
	if err != nil {
		panic(err)
	}

	f3.Decls = []goast.Decl{f3.Decls[1]}
	f3.Scope = nil
	f3.Unresolved = nil
	f3.Name = nil
	fset2 = token.NewFileSet()

	decl2 := f3.Decls[0].(*goast.GenDecl)

	decl2.Specs[0].(*goast.ValueSpec).Names[0].Obj.Decl = nil
	decl2.Specs[0].(*goast.ValueSpec).Type.(*goast.Ident).Obj = nil
	decl2.Specs[1].(*goast.ValueSpec).Names[0].Obj = nil
	decl2.Specs[2].(*goast.ValueSpec).Names[0].Obj = nil

	decl.Specs[0].(*goast.ValueSpec).Names[0].Obj.Decl = nil
	decl.Specs[0].(*goast.ValueSpec).Type.(*goast.Ident).Obj = nil
	decl.Specs[1].(*goast.ValueSpec).Names[0].Obj = nil
	decl.Specs[2].(*goast.ValueSpec).Names[0].Obj = nil

	//decl.Specs = append(decl.Specs, decl2.Specs...)
	//decl2.Specs = append(decl2.Specs, decl.Specs...)

	//decl2.Specs[2].(*goast.ValueSpec).Names[0] = nil
	//decl2.Rparen = 0
	//decl2.Lparen = 0
	decl.Lparen = 1

	p.File.Decls = append(p.File.Decls, []goast.Decl{decl}...)
	p.File.Decls = append(p.File.Decls, decl2)
	p.File.Decls = append(p.File.Decls, decl)
	p.File.Decls = append(p.File.Decls, decl2)
	p.File.Decls = append(p.File.Decls, decl2)

	// Print the AST.
	/*
		fmt.Printf("fset2 = %#v\n", fset2)

		fmt.Printf("decl  =%#v\n", decl)
		fmt.Printf("decl2 = %#v\n", decl2)

		fmt.Printf("f3-decl = %#v\n", f3.Decls)
	*/
	/*
		goast.Fprint(os.Stdout, token.NewFileSet(), decl, func(name string, value reflect.Value) bool {
			return true
		})
	*/
	/*
		goast.Fprint(os.Stdout, token.NewFileSet(), decl2, func(name string, value reflect.Value) bool {
			return true
		})
	*/

	fmt.Printf("%+v\n", decl)
	fmt.Printf("%+v\n", decl2)

	return nil
}

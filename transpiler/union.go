package transpiler

import (
	"strings"

	goast "go/ast"
	"go/token"

	"github.com/elliotchance/c2go/util"
)

func transpileUnion(name string, size int, fields []*goast.Field) []goast.Decl {
	res := []goast.Decl{
		// Type declaration (array: [x]byte with x the size of union)
		&goast.GenDecl{
			Tok: token.TYPE,
			Specs: []goast.Spec{
				&goast.TypeSpec{
					Name: util.NewIdent(name),
					Type: &goast.ArrayType{
						Elt: util.NewIdent("byte"),
						Len: util.NewIntLit(size), // Size of the union
					},
				},
			},
		},

		// cast() method
		&goast.FuncDecl{
			Name: util.NewIdent("cast"),
			Recv: &goast.FieldList{
				List: []*goast.Field{
					{
						Names: []*goast.Ident{util.NewIdent("self")},
						Type: &goast.StarExpr{
							X: util.NewIdent(name),
						},
					},
				},
			},
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						{
							Names: []*goast.Ident{util.NewIdent("t")},
							Type: &goast.SelectorExpr{
								X:   util.NewIdent("reflect"),
								Sel: util.NewIdent("Type"),
							},
						},
					},
				},
				Results: &goast.FieldList{
					List: []*goast.Field{
						{
							Type: &goast.SelectorExpr{
								X:   util.NewIdent("reflect"),
								Sel: util.NewIdent("Value"),
							},
						},
					},
				},
			},
			Body: &goast.BlockStmt{
				List: []goast.Stmt{
					&goast.ReturnStmt{
						Results: []goast.Expr{
							util.NewCallExpr("reflect.NewAt",
								util.NewIdent("t"),
								&goast.CallExpr{
									Fun: &goast.SelectorExpr{
										X:   util.NewIdent("unsafe"),
										Sel: util.NewIdent("Pointer"),
									},
									Args: []goast.Expr{
										&goast.UnaryExpr{
											Op: token.AND,
											X: &goast.IndexExpr{
												X:     util.NewIdent("self"),
												Index: util.NewIntLit(0),
											},
										},
									},
								},
							),
						},
					},
				},
			},
		},

		// assign() method
		&goast.FuncDecl{
			Name: util.NewIdent("assign"),
			Recv: &goast.FieldList{
				List: []*goast.Field{
					{
						Names: []*goast.Ident{util.NewIdent("self")},
						Type: &goast.StarExpr{
							X: util.NewIdent(name),
						},
					},
				},
			},
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						{
							Names: []*goast.Ident{util.NewIdent("v")},
							Type: &goast.InterfaceType{
								Methods: new(goast.FieldList),
							},
						},
					},
				},
			},
			Body: &goast.BlockStmt{
				List: []goast.Stmt{
					&goast.AssignStmt{
						Lhs: []goast.Expr{
							util.NewIdent("value"),
						},
						Tok: token.DEFINE,
						Rhs: []goast.Expr{
							&goast.CallExpr{
								Fun: &goast.SelectorExpr{
									X:   util.NewCallExpr("reflect.ValueOf", util.NewIdent("v")),
									Sel: util.NewIdent("Elem"),
								},
							},
						},
					},
					util.NewExprStmt(
						util.NewCallExpr("value.Set",
							&goast.CallExpr{
								Fun: &goast.SelectorExpr{
									X:   util.NewCallExpr("self.cast", util.NewCallExpr("value.Type")),
									Sel: util.NewIdent("Elem"),
								},
							},
						),
					),
				},
			},
		},

		// UntypedSet() method
		&goast.FuncDecl{
			Name: util.NewIdent("UntypedSet"),
			Recv: &goast.FieldList{
				List: []*goast.Field{
					{
						Names: []*goast.Ident{util.NewIdent("self")},
						Type: &goast.StarExpr{
							X: util.NewIdent(name),
						},
					},
				},
			},
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						{
							Names: []*goast.Ident{util.NewIdent("v")},
							Type: &goast.InterfaceType{
								Methods: new(goast.FieldList),
							},
						},
					},
				},
			},
			Body: &goast.BlockStmt{
				List: []goast.Stmt{
					&goast.AssignStmt{
						Lhs: []goast.Expr{
							util.NewIdent("value"),
						},
						Tok: token.DEFINE,
						Rhs: []goast.Expr{
							util.NewCallExpr("reflect.ValueOf", util.NewIdent("v")),
						},
					},
					util.NewExprStmt(
						&goast.CallExpr{
							Fun: &goast.SelectorExpr{
								X: &goast.CallExpr{
									Fun: &goast.SelectorExpr{
										X:   util.NewCallExpr("self.cast", util.NewCallExpr("value.Type")),
										Sel: util.NewIdent("Elem"),
									},
								},
								Sel: util.NewIdent("Set"),
							},
							Args: []goast.Expr{
								util.NewIdent("value"),
							},
						},
					),
				},
			},
		},
	}

	// Methods for each union field
	for _, f := range fields {
		fieldID := strings.Title(f.Names[0].Name)

		res = append(res,
			// Setter method (SetXX)
			&goast.FuncDecl{
				Name: util.NewIdent("Set" + fieldID),
				Recv: &goast.FieldList{
					List: []*goast.Field{
						{
							Names: []*goast.Ident{util.NewIdent("self")},
							Type: &goast.StarExpr{
								X: util.NewIdent(name),
							},
						},
					},
				},
				Type: &goast.FuncType{
					Params: &goast.FieldList{
						List: []*goast.Field{
							{
								Names: []*goast.Ident{util.NewIdent("v")},
								Type:  f.Type,
							},
						},
					},
					Results: &goast.FieldList{
						List: []*goast.Field{
							{
								Type: f.Type,
							},
						},
					},
				},
				Body: &goast.BlockStmt{
					List: []goast.Stmt{
						util.NewExprStmt(
							util.NewCallExpr("self.UntypedSet", util.NewIdent("v")),
						),
						&goast.ReturnStmt{
							Results: []goast.Expr{
								util.NewIdent("v"),
							},
						},
					},
				},
			},

			// Getter method (GetXX)
			&goast.FuncDecl{
				Name: util.NewIdent("Get" + fieldID),
				Recv: &goast.FieldList{
					List: []*goast.Field{
						{
							Names: []*goast.Ident{util.NewIdent("self")},
							Type: &goast.StarExpr{
								X: util.NewIdent(name),
							},
						},
					},
				},
				Type: &goast.FuncType{
					Results: &goast.FieldList{
						List: []*goast.Field{
							{
								Names: []*goast.Ident{util.NewIdent("res")},
								Type:  f.Type,
							},
						},
					},
				},
				Body: &goast.BlockStmt{
					List: []goast.Stmt{
						util.NewExprStmt(
							util.NewCallExpr("self.assign", util.NewUnaryExpr(token.AND, util.NewIdent("res"))),
						),
						new(goast.ReturnStmt),
					},
				},
			},
		)
	}

	return res
}

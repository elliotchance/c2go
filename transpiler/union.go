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
					Name: goast.NewIdent(name),
					Type: &goast.ArrayType{
						Elt: goast.NewIdent("byte"),
						Len: util.NewIntLit(size), // Size of the union
					},
				},
			},
		},

		// cast() method
		&goast.FuncDecl{
			Name: goast.NewIdent("cast"),
			Recv: &goast.FieldList{
				List: []*goast.Field{
					&goast.Field{
						Names: []*goast.Ident{goast.NewIdent("self")},
						Type: &goast.StarExpr{
							X: goast.NewIdent(name),
						},
					},
				},
			},
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						&goast.Field{
							Names: []*goast.Ident{goast.NewIdent("t")},
							Type: &goast.SelectorExpr{
								X:   goast.NewIdent("reflect"),
								Sel: goast.NewIdent("Type"),
							},
						},
					},
				},
				Results: &goast.FieldList{
					List: []*goast.Field{
						&goast.Field{
							Type: &goast.SelectorExpr{
								X:   goast.NewIdent("reflect"),
								Sel: goast.NewIdent("Value"),
							},
						},
					},
				},
			},
			Body: &goast.BlockStmt{
				List: []goast.Stmt{
					&goast.ReturnStmt{
						Results: []goast.Expr{
							&goast.CallExpr{
								Fun: &goast.SelectorExpr{
									X:   goast.NewIdent("reflect"),
									Sel: goast.NewIdent("NewAt"),
								},
								Args: []goast.Expr{
									goast.NewIdent("t"),
									&goast.CallExpr{
										Fun: &goast.SelectorExpr{
											X:   goast.NewIdent("unsafe"),
											Sel: goast.NewIdent("Pointer"),
										},
										Args: []goast.Expr{
											&goast.UnaryExpr{
												Op: token.AND,
												X: &goast.IndexExpr{
													X:     goast.NewIdent("self"),
													Index: util.NewIntLit(0),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		// assign() method
		&goast.FuncDecl{
			Name: goast.NewIdent("assign"),
			Recv: &goast.FieldList{
				List: []*goast.Field{
					&goast.Field{
						Names: []*goast.Ident{goast.NewIdent("self")},
						Type: &goast.StarExpr{
							X: goast.NewIdent(name),
						},
					},
				},
			},
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						&goast.Field{
							Names: []*goast.Ident{goast.NewIdent("v")},
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
							goast.NewIdent("value"),
						},
						Tok: token.DEFINE,
						Rhs: []goast.Expr{
							&goast.CallExpr{
								Fun: &goast.SelectorExpr{
									X: &goast.CallExpr{
										Fun: &goast.SelectorExpr{
											X:   goast.NewIdent("reflect"),
											Sel: goast.NewIdent("ValueOf"),
										},
										Args: []goast.Expr{
											goast.NewIdent("v"),
										},
									},
									Sel: goast.NewIdent("Elem"),
								},
							},
						},
					},
					&goast.ExprStmt{
						X: &goast.CallExpr{
							Fun: &goast.SelectorExpr{
								X:   goast.NewIdent("value"),
								Sel: goast.NewIdent("Set"),
							},
							Args: []goast.Expr{
								&goast.CallExpr{
									Fun: &goast.SelectorExpr{
										X: &goast.CallExpr{
											Fun: &goast.SelectorExpr{
												X:   goast.NewIdent("self"),
												Sel: goast.NewIdent("cast"),
											},
											Args: []goast.Expr{
												&goast.CallExpr{
													Fun: &goast.SelectorExpr{
														X:   goast.NewIdent("value"),
														Sel: goast.NewIdent("Type"),
													},
												},
											},
										},
										Sel: goast.NewIdent("Elem"),
									},
								},
							},
						},
					},
				},
			},
		},

		// UntypedSet() method
		&goast.FuncDecl{
			Name: goast.NewIdent("UntypedSet"),
			Recv: &goast.FieldList{
				List: []*goast.Field{
					&goast.Field{
						Names: []*goast.Ident{goast.NewIdent("self")},
						Type: &goast.StarExpr{
							X: goast.NewIdent(name),
						},
					},
				},
			},
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						&goast.Field{
							Names: []*goast.Ident{goast.NewIdent("v")},
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
							goast.NewIdent("value"),
						},
						Tok: token.DEFINE,
						Rhs: []goast.Expr{
							&goast.CallExpr{
								Fun: &goast.SelectorExpr{
									X:   goast.NewIdent("reflect"),
									Sel: goast.NewIdent("ValueOf"),
								},
								Args: []goast.Expr{
									goast.NewIdent("v"),
								},
							},
						},
					},
					&goast.ExprStmt{
						X: &goast.CallExpr{
							Fun: &goast.SelectorExpr{
								X: &goast.CallExpr{
									Fun: &goast.SelectorExpr{
										X: &goast.CallExpr{
											Fun: &goast.SelectorExpr{
												X:   goast.NewIdent("self"),
												Sel: goast.NewIdent("cast"),
											},
											Args: []goast.Expr{
												&goast.CallExpr{
													Fun: &goast.SelectorExpr{
														X:   goast.NewIdent("value"),
														Sel: goast.NewIdent("Type"),
													},
												},
											},
										},
										Sel: goast.NewIdent("Elem"),
									},
								},
								Sel: goast.NewIdent("Set"),
							},
							Args: []goast.Expr{
								goast.NewIdent("value"),
							},
						},
					},
				},
			},
		},
	}

	// Methods for each union field
	for _, f := range fields {
		field_id := strings.Title(f.Names[0].Name)

		res = append(res,
			// Setter method (SetXX)
			&goast.FuncDecl{
				Name: goast.NewIdent("Set" + field_id),
				Recv: &goast.FieldList{
					List: []*goast.Field{
						&goast.Field{
							Names: []*goast.Ident{goast.NewIdent("self")},
							Type: &goast.StarExpr{
								X: goast.NewIdent(name),
							},
						},
					},
				},
				Type: &goast.FuncType{
					Params: &goast.FieldList{
						List: []*goast.Field{
							&goast.Field{
								Names: []*goast.Ident{goast.NewIdent("v")},
								Type:  f.Type,
							},
						},
					},
					Results: &goast.FieldList{
						List: []*goast.Field{
							&goast.Field{
								Type: f.Type,
							},
						},
					},
				},
				Body: &goast.BlockStmt{
					List: []goast.Stmt{
						&goast.ExprStmt{
							&goast.CallExpr{
								Fun: &goast.SelectorExpr{
									X:   goast.NewIdent("self"),
									Sel: goast.NewIdent("UntypedSet"),
								},
								Args: []goast.Expr{
									goast.NewIdent("v"),
								},
							},
						},
						&goast.ReturnStmt{
							Results: []goast.Expr{
								goast.NewIdent("v"),
							},
						},
					},
				},
			},

			// Getter method (GetXX)
			&goast.FuncDecl{
				Name: goast.NewIdent("Get" + field_id),
				Recv: &goast.FieldList{
					List: []*goast.Field{
						&goast.Field{
							Names: []*goast.Ident{goast.NewIdent("self")},
							Type: &goast.StarExpr{
								X: goast.NewIdent(name),
							},
						},
					},
				},
				Type: &goast.FuncType{
					Results: &goast.FieldList{
						List: []*goast.Field{
							&goast.Field{
								Names: []*goast.Ident{goast.NewIdent("res")},
								Type:  f.Type,
							},
						},
					},
				},
				Body: &goast.BlockStmt{
					List: []goast.Stmt{
						&goast.ExprStmt{
							&goast.CallExpr{
								Fun: &goast.SelectorExpr{
									X:   goast.NewIdent("self"),
									Sel: goast.NewIdent("assign"),
								},
								Args: []goast.Expr{
									&goast.UnaryExpr{
										Op: token.AND,
										X:  goast.NewIdent("res"),
									},
								},
							},
						},
						new(goast.ReturnStmt),
					},
				},
			},
		)
	}

	return res
}

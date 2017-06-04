package transpiler

import (
    "strconv"
    "strings"

    goast "go/ast"
    "go/token"
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
                        Len: &goast.BasicLit{
                            Kind:  token.INT,
                            Value: strconv.Itoa(size), // Size of the union
                        },
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
                                                    X:  goast.NewIdent("self"),
                                                    Index: &goast.BasicLit{
                                                        Kind:  token.INT,
                                                        Value: "0",
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
        },
    }

    for _, f := range fields {
        field_id := strings.Title(f.Names[0].Name)

        res = append(res,
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

            // Get pointer method (PtrXX)
            &goast.FuncDecl{
                Name: goast.NewIdent("Ptr" + field_id),
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
        )
    }

    return res
}

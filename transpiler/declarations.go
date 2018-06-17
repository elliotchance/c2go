// This file contains functions for transpiling declarations of variables and
// types. The usage of variables is handled in variables.go.

package transpiler

import (
	"errors"
	"fmt"
	goast "go/ast"
	"go/token"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

func newFunctionField(p *program.Program, name, cType string) (_ *goast.Field, err error) {
	if name == "" {
		err = fmt.Errorf("Name of function field cannot be empty")
		return
	}
	if !types.IsFunction(cType) {
		err = fmt.Errorf("Cannot create function field for type : %s", cType)
		return
	}

	field := &goast.Field{
		Names: []*goast.Ident{
			util.NewIdent(name),
		},
	}
	var arg, ret []string
	arg, ret, err = types.SeparateFunction(p, cType)
	if err != nil {
		return
	}
	funcType := &goast.FuncType{}
	argFieldList := []*goast.Field{}
	for _, aa := range arg {
		argFieldList = append(argFieldList, &goast.Field{
			Type: goast.NewIdent(aa),
		})
	}
	funcType.Params = &goast.FieldList{
		List: argFieldList,
	}
	funcType.Results = &goast.FieldList{
		List: []*goast.Field{
			&goast.Field{
				Type: goast.NewIdent(ret[0]),
			},
		},
	}
	field.Type = funcType

	return field, nil
}

func transpileFieldDecl(p *program.Program, n *ast.FieldDecl) (field *goast.Field, err error) {
	if types.IsFunction(n.Type) {
		field, err = newFunctionField(p, n.Name, n.Type)
		if err == nil {
			return
		}
	}

	name := n.Name

	// FIXME: What causes this? See __darwin_fp_control for example.
	if name == "" {
		return nil, fmt.Errorf("Error : name of FieldDecl is empty")
	}

	// Add for fix bug in "stdlib.h"
	// build/tests/exit/main_test.go:90:11: undefined: wait
	// it is "union" with some anonymous struct
	if n.Type == "union wait *" {
		return nil, fmt.Errorf("Avoid struct `union wait *` in FieldDecl")
	}

	fieldType, err := types.ResolveType(p, n.Type)
	p.AddMessage(p.GenerateWarningMessage(err, n))

	// TODO: The name of a variable or field cannot be a reserved word
	// https://github.com/elliotchance/c2go/issues/83
	// Search for this issue in other areas of the codebase.
	if util.IsGoKeyword(name) {
		name += "_"
	}

	arrayType, arraySize := types.GetArrayTypeAndSize(n.Type)
	if arraySize != -1 {
		fieldType, err = types.ResolveType(p, arrayType)
		p.AddMessage(p.GenerateWarningMessage(err, n))
		fieldType = fmt.Sprintf("[%d]%s", arraySize, fieldType)
		err = nil
	}

	return &goast.Field{
		Names: []*goast.Ident{util.NewIdent(name)},
		Type:  util.NewTypeIdent(fieldType),
	}, nil
}

func transpileRecordDecl(p *program.Program, n *ast.RecordDecl) (decls []goast.Decl, err error) {
	name := n.Name

	if name == "" || p.IsTypeAlreadyDefined(name) {
		err = nil
		return
	}

	name = types.GenerateCorrectType(name)
	p.DefineType(name)

	// TODO: Some platform structs are ignored.
	// https://github.com/elliotchance/c2go/issues/85
	if name == "__locale_struct" ||
		name == "__sigaction" ||
		name == "sigaction" {
		err = nil
		return
	}

	var fields []*goast.Field

	for pos := range n.Children() {
		switch field := n.Children()[pos].(type) {
		case *ast.FieldDecl:
			field.Type = types.GenerateCorrectType(field.Type)
			field.Type2 = types.GenerateCorrectType(field.Type2)
			f, err := transpileFieldDecl(p, field)
			if err != nil {
				p.AddMessage(p.GenerateWarningMessage(err, field))
			} else {
				fields = append(fields, f)
			}

		case *ast.RecordDecl:
			if field.Kind == "union" && pos+2 <= len(n.Children()) {
				if inField, ok := n.Children()[pos+1].(*ast.FieldDecl); ok {
					inField.Type = types.GenerateCorrectType(inField.Type)
					inField.Type2 = types.GenerateCorrectType(inField.Type2)
					field.Name = string(([]byte(inField.Type))[len("union "):])
					declUnion, err := transpileRecordDecl(p, field)
					if err != nil {
						p.AddMessage(p.GenerateWarningMessage(err, field))
					}
					pos++
					decls = append(decls, declUnion...)
				}
			} else if field.Kind == "struct" &&
				pos+2 <= len(n.Children()) {
				if inField, ok := n.Children()[pos+1].(*ast.FieldDecl); ok {
					inField.Type = types.GenerateCorrectType(inField.Type)
					inField.Type2 = types.GenerateCorrectType(inField.Type2)

					field.Name = inField.Type

					if strings.HasPrefix(field.Name, "const ") {
						field.Name = field.Name[len("const "):]
					}

					if strings.HasPrefix(field.Name, "struct ") {
						field.Name = field.Name[len("struct "):]
					}

					if field.Name[len(field.Name)-1] == '*' {
						// star in struct
						field.Name = field.Name[:len(field.Name)-len(" *")]
					}

					if strings.Contains(field.Name, "[") {
						p.AddMessage(p.GenerateWarningMessage(
							fmt.Errorf("Not acceptable name of struct %s", field.Name), n))
						continue
					}
					declStruct, err := transpileRecordDecl(p, field)
					if err != nil {
						p.AddMessage(p.GenerateWarningMessage(err, field))
					}
					pos++
					decls = append(decls, declStruct...)
				}
			} else {
				decls, err = transpileRecordDecl(p, field)
				if err != nil {
					message := fmt.Sprintf("could not parse %v", field)
					p.AddMessage(p.GenerateWarningMessage(errors.New(message), field))
				}
			}

		case *ast.FullComment:
			// We haven't Go ast struct for easy inject a comments.
			// All comments are added like CommentsGroup.
			// So, we can ignore that comment, because all comments
			// will be added by another way.

		default:
			message := fmt.Sprintf("could not parse %v", field)
			p.AddMessage(p.GenerateWarningMessage(errors.New(message), field))
		}
	}

	s := program.NewStruct(n)
	if s.IsUnion {
		p.Unions["union "+s.Name] = s
	} else {
		p.Structs["struct "+s.Name] = s
	}
	if s.IsUnion {
		// Union size
		var size int
		size, err = types.SizeOf(p, "union "+name)

		// In normal case no error is returned,
		if err != nil {
			// but if we catch one, send it as a warning
			message := fmt.Sprintf("could not determine the size of type `union %s` for that reason: %s", name, err)
			p.AddMessage(p.GenerateWarningMessage(errors.New(message), n))
		} else {
			// So, we got size, then
			// Add imports needed
			p.AddImports("unsafe")

			// Declaration for implementing union type
			d, err2 := transpileUnion(name, size, fields)
			if err2 != nil {
				return nil, err2
			}
			decls = append(decls, d...)
		}
		return
	}

	decls = append(decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: util.NewIdent(name),
				Type: &goast.StructType{
					Fields: &goast.FieldList{
						List: fields,
					},
				},
			},
		},
	})

	return
}

func transpileTypedefDecl(p *program.Program, n *ast.TypedefDecl) (decls []goast.Decl, err error) {
	// implicit code from clang at the head of each clang AST tree
	if n.IsImplicit && n.Pos.File == ast.PositionBuiltIn {
		return
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile Typedef Decl : err = %v", err)
		}
	}()
	name := n.Name
	n.Type = types.CleanCType(n.Type)
	n.Type2 = types.CleanCType(n.Type2)

	if types.IsFunction(n.Type) {
		var field *goast.Field
		field, err = newFunctionField(p, n.Name, n.Type)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		} else {
			// registration type
			p.TypedefType[n.Name] = n.Type

			decls = append(decls, &goast.GenDecl{
				Tok: token.TYPE,
				Specs: []goast.Spec{
					&goast.TypeSpec{
						Name: util.NewIdent(name),
						Type: field.Type,
					},
				},
			})
			err = nil
			return
		}
	}

	// added for support "typedef enum {...} dd" with empty name of struct
	// Result in Go: "type dd int32"
	if strings.Contains(n.Type, "enum") {
		// Registration new type in program.Program
		if !p.IsTypeAlreadyDefined(n.Name) {
			p.DefineType(n.Name)
			p.EnumTypedefName[n.Name] = true
		}
		decls = append(decls, &goast.GenDecl{
			Tok: token.TYPE,
			Specs: []goast.Spec{
				&goast.TypeSpec{
					Name: util.NewIdent(name),
					Type: util.NewTypeIdent("int32"),
				},
			},
		})
		err = nil
		return
	}

	if p.IsTypeAlreadyDefined(name) {
		err = nil
		return
	}

	p.DefineType(name)

	resolvedType, err := types.ResolveType(p, n.Type)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}

	// There is a case where the name of the type is also the definition,
	// like:
	//
	//     type _RuneEntry _RuneEntry
	//
	// This of course is impossible and will cause the Go not to compile.
	// It itself is caused by lack of understanding (at this time) about
	// certain scenarios that types are defined as. The above example comes
	// from:
	//
	//     typedef struct {
	//        // ... some fields
	//     } _RuneEntry;
	//
	// Until which time that we actually need this to work I am going to
	// suppress these.
	if name == resolvedType {
		err = nil
		return
	}

	if name == "__darwin_ct_rune_t" {
		resolvedType = p.ImportType("github.com/elliotchance/c2go/darwin.CtRuneT")
	}

	if name == "div_t" || name == "ldiv_t" || name == "lldiv_t" {
		intType := "int"
		if name == "ldiv_t" {
			intType = "long int"
		} else if name == "lldiv_t" {
			intType = "long long int"
		}

		// I don't know to extract the correct fields from the typedef to create
		// the internal definition. This is used in the noarch package
		// (stdio.go).
		//
		// The name of the struct is not prefixed with "struct " because it is a
		// typedef.
		p.Structs[name] = &program.Struct{
			Name:    name,
			IsUnion: false,
			Fields: map[string]interface{}{
				"quot": intType,
				"rem":  intType,
			},
		}
	}

	err = nil
	if resolvedType == "" {
		resolvedType = "interface{}"
	}
	decls = append(decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: util.NewIdent(name),
				Type: util.NewTypeIdent(resolvedType),
			},
		},
	})

	if v, ok := p.Structs["struct "+resolvedType]; ok {
		// Registration "typedef struct" with non-empty name of struct
		p.Structs["struct "+name] = v
	} else if v, ok := p.EnumConstantToEnum["enum "+resolvedType]; ok {
		// Registration "enum constants"
		p.EnumConstantToEnum["enum "+resolvedType] = v
	} else {
		// Registration "typedef type type2"
		p.TypedefType[n.Name] = n.Type
	}

	return
}

func transpileVarDecl(p *program.Program, n *ast.VarDecl) (
	decls []goast.Decl, theType string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileVarDecl : err = %v", err)
		}
	}()
	// There may be some startup code for this global variable.
	if p.Function == nil {
		name := n.Name
		switch name {
		// Below are for macOS.
		case "__stdinp", "__stdoutp":
			theType = "*noarch.File"
			p.AddImport("github.com/elliotchance/c2go/noarch")
			p.AppendStartupExpr(
				util.NewBinaryExpr(
					goast.NewIdent(name),
					token.ASSIGN,
					util.NewTypeIdent("noarch."+util.Ucfirst(name[2:len(name)-1])),
					"*noarch.File",
					true,
				),
			)
			return []goast.Decl{&goast.GenDecl{
				Tok: token.VAR,
				Specs: []goast.Spec{&goast.ValueSpec{
					Names: []*goast.Ident{{Name: name}},
					Type:  util.NewTypeIdent(theType),
					Doc:   p.GetMessageComments(),
				}},
			}}, "", nil

		// Below are for linux.
		case "stdout", "stdin", "stderr":
			theType = "*noarch.File"
			p.AddImport("github.com/elliotchance/c2go/noarch")
			p.AppendStartupExpr(
				util.NewBinaryExpr(
					goast.NewIdent(name),
					token.ASSIGN,
					util.NewTypeIdent("noarch."+util.Ucfirst(name)),
					theType,
					true,
				),
			)
			return []goast.Decl{&goast.GenDecl{
				Tok: token.VAR,
				Specs: []goast.Spec{&goast.ValueSpec{
					Names: []*goast.Ident{{Name: name}},
					Type:  util.NewTypeIdent(theType),
				}},
				Doc: p.GetMessageComments(),
			}}, "", nil

		default:
			// No init needed.
		}
	}

	// Ignore extern as there is no analogy for Go right now.
	if n.IsExtern && len(n.ChildNodes) == 0 {
		return
	}

	if strings.Contains(n.Type, "va_list") && strings.Contains(n.Type2, "va_list_tag") {
		// variable for va_list. see "variadic function"
		// header : <stdarg.h>
		// Example :
		// DeclStmt 0x2fd87e0 <line:442:2, col:14>
		// `-VarDecl 0x2fd8780 <col:2, col:10> col:10 used args 'va_list':'struct __va_list_tag [1]'
		// Result:
		// ... - convert to - c2goArgs ...interface{}
		// var args = c2goArgs
		return []goast.Decl{&goast.GenDecl{
			Tok: token.VAR,
			Specs: []goast.Spec{
				&goast.ValueSpec{
					Names:  []*goast.Ident{util.NewIdent("c2goVaList")},
					Values: []goast.Expr{util.NewVaListTag()},
				},
			},
		}}, "", nil
	}

	/*
		Example of DeclStmt for C code:
		void * a = NULL;
		void(*t)(void) = a;
		Example of AST:
		`-VarDecl 0x365fea8 <col:3, col:20> col:9 used t 'void (*)(void)' cinit
		  `-ImplicitCastExpr 0x365ff48 <col:20> 'void (*)(void)' <BitCast>
		    `-ImplicitCastExpr 0x365ff30 <col:20> 'void *' <LValueToRValue>
		      `-DeclRefExpr 0x365ff08 <col:20> 'void *' lvalue Var 0x365f8c8 'r' 'void *'
	*/

	if len(n.Children()) > 0 {
		if v, ok := (n.Children()[0]).(*ast.ImplicitCastExpr); ok {
			if len(v.Type) > 0 {
				// Is it function ?
				if types.IsFunction(v.Type) {
					var fields, returns []string
					fields, returns, err = types.SeparateFunction(p, v.Type)
					if err != nil {
						err = fmt.Errorf("Cannot resolve function : %v", err)
						return
					}
					functionType := GenerateFuncType(fields, returns)
					nameVar1 := n.Name

					if vv, ok := v.Children()[0].(*ast.ImplicitCastExpr); ok {
						if decl, ok := vv.Children()[0].(*ast.DeclRefExpr); ok {
							nameVar2 := decl.Name

							return []goast.Decl{&goast.GenDecl{
								Tok: token.VAR,
								Specs: []goast.Spec{&goast.ValueSpec{
									Names: []*goast.Ident{{Name: nameVar1}},
									Type:  functionType,
									Values: []goast.Expr{&goast.TypeAssertExpr{
										X: &goast.SelectorExpr{
											X: util.NewCallExpr("*noarch.InterfaceWrapper",
												&goast.Ident{Name: nameVar2}),
											Sel: util.NewIdent("X"),
										},
										Type: functionType,
									}},
									Doc: p.GetMessageComments(),
								},
								}}}, "", nil
						}
					}
				}
			}
		}
	}

	if types.IsFunction(n.Type) {
		var fields, returns []string
		fields, returns, err = types.SeparateFunction(p, n.Type)
		if err != nil {
			p.AddMessage(p.GenerateErrorMessage(fmt.Errorf("Cannot resolve function : %v", err), n))
			err = nil // Error is ignored
			return
		}
		functionType := GenerateFuncType(fields, returns)
		nameVar1 := n.Name
		decls = append(decls, &goast.GenDecl{
			Tok: token.VAR,
			Specs: []goast.Spec{&goast.ValueSpec{
				Names: []*goast.Ident{{Name: nameVar1}},
				Type:  functionType,
				Doc:   p.GetMessageComments(),
			},
			}})
		err = nil
		return
	}

	theType = n.Type
	_, isTypedefType := p.TypedefType[theType]

	theType, err = types.ResolveType(p, n.Type)
	if err != nil {
		p.AddMessage(p.GenerateErrorMessage(fmt.Errorf("Cannot resolve type %s : %v", theType, err), n))
		err = nil // Error is ignored
	}

	p.GlobalVariables[n.Name] = theType

	name := n.Name
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	// TODO: Some platform structs are ignored.
	// https://github.com/elliotchance/c2go/issues/85
	if name == "_LIB_VERSION" ||
		name == "_IO_2_1_stdin_" ||
		name == "_IO_2_1_stdout_" ||
		name == "_IO_2_1_stderr_" ||
		name == "_DefaultRuneLocale" ||
		name == "_CurrentRuneLocale" {
		theType = "unknown10"
		return
	}

	defaultValue, _, newPre, newPost, err := getDefaultValueForVar(p, n)
	if err != nil {
		p.AddMessage(p.GenerateErrorMessage(err, n))
		err = nil // Error is ignored
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// Allocate slice so that it operates like a fixed size array.
	arrayType, arraySize := types.GetArrayTypeAndSize(n.Type)

	if arraySize != -1 && defaultValue == nil {
		if len(n.Children()) == 0 {
			var goArrayType string
			goArrayType, err = types.ResolveType(p, arrayType)
			if err != nil {
				p.AddMessage(p.GenerateErrorMessage(err, n))
				err = nil // Error is ignored
			}

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
		} else {
			if iniList, ok := n.Children()[0].(*ast.InitListExpr); ok {
				var list goast.Expr
				list, _, err = transpileInitListExpr(iniList, p)
				if err != nil {
					p.AddMessage(p.GenerateErrorMessage(err, n))
					err = nil // Error is ignored
				} else {
					defaultValue = []goast.Expr{list}
				}
			}
		}
	}

	if len(preStmts) != 0 || len(postStmts) != 0 {
		p.AddMessage(p.GenerateErrorMessage(fmt.Errorf("Not acceptable length of Stmt : pre(%d), post(%d)", len(preStmts), len(postStmts)), n))
	}

	var typeResult goast.Expr
	if isTypedefType {
		typeResult = goast.NewIdent(theType)
	} else {
		typeResult = util.NewTypeIdent(theType)
	}

	return []goast.Decl{&goast.GenDecl{
		Tok: token.VAR,
		Specs: []goast.Spec{
			&goast.ValueSpec{
				Names:  []*goast.Ident{util.NewIdent(n.Name)},
				Type:   typeResult,
				Values: defaultValue,
				Doc:    p.GetMessageComments(),
			},
		},
	}}, "", nil
}

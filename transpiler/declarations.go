// This file contains functions for transpiling declarations of variables and
// types. The usage of variables is handled in variables.go.

package transpiler

import (
	"errors"
	"fmt"
	goast "go/ast"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

func transpileFieldDecl(p *program.Program, n *ast.FieldDecl) (*goast.Field, string) {
	name := n.Name

	// FIXME: What causes this? See __darwin_fp_control for example.
	if name == "" {
		return nil, ""
	}

	fieldType, err := types.ResolveType(p, n.Type)
	p.AddMessage(ast.GenerateWarningMessage(err, n))

	// TODO: The name of a variable or field cannot be "type"
	// https://github.com/elliotchance/c2go/issues/83
	if name == "type" {
		name = "type_"
	}

	return &goast.Field{
		Names: []*goast.Ident{goast.NewIdent(name)},
		Type:  goast.NewIdent(fieldType),
	}, "unknown3"
}

func transpileRecordDecl(p *program.Program, n *ast.RecordDecl) error {
	name := n.Name
	if name == "" || p.TypeIsAlreadyDefined(name) {
		return nil
	}

	p.TypeIsNowDefined(name)

	s := program.NewStruct(n)
	p.Structs["struct "+s.Name] = s

	// TODO: Unions are not supported.
	// https://github.com/elliotchance/c2go/issues/84
	if n.Kind == "union" {
		return nil
	}

	// TODO: Some platform structs are ignored.
	// https://github.com/elliotchance/c2go/issues/85
	if name == "__locale_struct" ||
		name == "__sigaction" ||
		name == "sigaction" {
		return nil
	}

	var fields []*goast.Field
	for _, c := range n.Children {
		if field, ok := c.(*ast.FieldDecl); ok {
			f, _ := transpileFieldDecl(p, field)

			if f != nil {
				fields = append(fields, f)
			}
		} else {
			message := fmt.Sprintf("could not parse %v", c)
			p.AddMessage(ast.GenerateWarningMessage(errors.New(message), c))
		}
	}

	p.File.Decls = append(p.File.Decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: goast.NewIdent(name),
				Type: &goast.StructType{
					Fields: &goast.FieldList{
						List: fields,
					},
				},
			},
		},
	})

	return nil
}

func transpileTypedefDecl(p *program.Program, n *ast.TypedefDecl) error {
	name := n.Name

	if p.TypeIsAlreadyDefined(name) {
		return nil
	}

	p.TypeIsNowDefined(name)

	resolvedType, err := types.ResolveType(p, n.Type)
	p.AddMessage(ast.GenerateWarningMessage(err, n))

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
		return nil
	}

	if name == "__darwin_ct_rune_t" {
		resolvedType = p.ImportType("github.com/elliotchance/c2go/darwin.CtRuneT")
	}

	// TODO: Some platform structs are ignored.
	// https://github.com/elliotchance/c2go/issues/85
	if name == "__builtin_va_list" ||
		name == "__qaddr_t" ||
		name == "definition" ||
		name == "_IO_lock_t" ||
		name == "va_list" ||
		name == "fpos_t" ||
		name == "__NSConstantString" ||
		name == "__darwin_va_list" ||
		name == "__fsid_t" ||
		name == "_G_fpos_t" ||
		name == "_G_fpos64_t" ||
		name == "__locale_t" ||
		name == "locale_t" ||
		name == "fsid_t" ||
		name == "sigset_t" {
		return nil
	}

	p.File.Decls = append(p.File.Decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: goast.NewIdent(name),
				Type: goast.NewIdent(resolvedType),
			},
		},
	})

	return nil
}

func transpileVarDecl(p *program.Program, n *ast.VarDecl) (
	[]goast.Stmt, []goast.Stmt, string) {
	theType, err := types.ResolveType(p, n.Type)
	p.AddMessage(ast.GenerateWarningMessage(err, n))

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
		return nil, nil, "unknown10"
	}

	// TODO: The name of a variable or field cannot be "type"
	// https://github.com/elliotchance/c2go/issues/83
	if name == "type" {
		name = "type_"
	}

	// There may be some startup code for this global variable.
	if p.Function == nil {
		switch name {
		// Below are for macOS.
		case "__stdinp", "__stdoutp":
			p.AddImports("github.com/elliotchance/c2go/noarch", "os")
			p.AppendStartupExpr(
				util.NewBinaryExpr(
					util.NewIdent(name),
					token.ASSIGN,
					util.NewCallExpr(
						"noarch.NewFile",
						goast.NewIdent("os."+util.Ucfirst(name[2:len(name)-1])),
					),
				),
			)

		// Below are for linux.
		case "stdout", "stdin", "stderr":
			theType = "*noarch.File"
			p.AddImports("github.com/elliotchance/c2go/noarch", "os")
			p.AppendStartupExpr(
				util.NewBinaryExpr(
					util.NewIdent(name),
					token.ASSIGN,
					util.NewCallExpr(
						"noarch.NewFile",
						goast.NewIdent("os."+util.Ucfirst(name)),
					),
				),
			)

		default:
			// No init needed.
		}
	}

	defaultValue, _, newPre, newPost, err := getDefaultValueForVar(p, n)
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	p.File.Decls = append(p.File.Decls, &goast.GenDecl{
		Tok: token.VAR,
		Specs: []goast.Spec{
			&goast.ValueSpec{
				Names: []*goast.Ident{
					goast.NewIdent(name),
				},
				Type:   goast.NewIdent(theType),
				Values: defaultValue,
			},
		},
	})

	return nil, nil, theType
}

// This file contains functions for transpiling declarations of variables and
// types. The usage of variables is handled in variables.go.

package transpiler

import (
	goast "go/ast"
	"go/token"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

func transpileFieldDecl(p *program.Program, n *ast.FieldDecl) (*goast.Field, string) {
	fieldType := types.ResolveType(p, n.Type)
	name := n.Name

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
		f, _ := transpileFieldDecl(p, c.(*ast.FieldDecl))
		fields = append(fields, f)
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

	resolvedType := types.ResolveType(p, n.Type)

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

	if name == "__mbstate_t" {
		resolvedType = p.ImportType("github.com/elliotchance/c2go/darwin.C__mbstate_t")
	}

	if name == "__darwin_ct_rune_t" {
		resolvedType = p.ImportType("github.com/elliotchance/c2go/darwin.Darwin_ct_rune_t")
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

func transpileVarDecl(p *program.Program, n *ast.VarDecl) string {
	theType := types.ResolveType(p, n.Type)
	name := n.Name

	// TODO: Some platform structs are ignored.
	// https://github.com/elliotchance/c2go/issues/85
	if name == "_LIB_VERSION" ||
		name == "_IO_2_1_stdin_" ||
		name == "_IO_2_1_stdout_" ||
		name == "_IO_2_1_stderr_" ||
		name == "stdin" ||
		name == "stdout" ||
		name == "stderr" ||
		name == "_DefaultRuneLocale" ||
		name == "_CurrentRuneLocale" {
		return "unknown10"
	}

	// TODO: The name of a variable or field cannot be "type"
	// https://github.com/elliotchance/c2go/issues/83
	if name == "type" {
		name = "type_"
	}

	var defaultValues []goast.Expr
	if len(n.Children) > 0 {
		defaultValue, defaultValueType, err := transpileToExpr(n.Children[0], p)
		if err != nil {
			panic(err)
		}

		defaultValues = []goast.Expr{
			types.CastExpr(p, defaultValue, defaultValueType, n.Type),
		}
	}

	p.File.Decls = append(p.File.Decls, &goast.GenDecl{
		Tok: token.VAR,
		Specs: []goast.Spec{
			&goast.ValueSpec{
				Names: []*goast.Ident{
					goast.NewIdent(name),
				},
				Type:   goast.NewIdent(theType),
				Values: defaultValues,
			},
		},
	})

	return n.Type
}

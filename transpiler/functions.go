package transpiler

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"

	goast "go/ast"
)

func transpileCallExpr(n *ast.CallExpr, p *program.Program) (*goast.CallExpr, string, error) {
	functionName := n.Children[0].(*ast.ImplicitCastExpr).Children[0].(*ast.DeclRefExpr).Name
	functionDef := program.GetFunctionDefinition(functionName)

	if functionDef == nil {
		panic(fmt.Sprintf("unknown function: %s", functionName))
	}

	if functionDef.Substitution != "" {
		parts := strings.Split(functionDef.Substitution, ".")
		importName := strings.Join(parts[:len(parts)-1], ".")
		p.AddImport(importName)

		parts2 := strings.Split(functionDef.Substitution, "/")
		functionName = parts2[len(parts2)-1]
	}

	args := []goast.Expr{}
	i := 0
	for _, arg := range n.Children[1:] {
		e, eType, err := transpileToExpr(arg, p)
		if err != nil {
			return nil, "", err
		}

		if i > len(functionDef.ArgumentTypes)-1 {
			// This means the argument is one of the varargs
			// so we don't know what type it needs to be
			// cast to.
			args = append(args, e)
		} else {
			args = append(args, types.CastExpr(p, e, eType, functionDef.ArgumentTypes[i]))
		}

		i++
	}

	return &goast.CallExpr{
		Fun:      goast.NewIdent(functionName),
		Lparen:   token.NoPos,
		Args:     args,
		Ellipsis: token.NoPos,
		Rparen:   token.NoPos,
	}, "", nil

	// src := fmt.Sprintf("%s(%s)", functionName, strings.Join(parts, ", "))
	// return src, functionDef.ReturnType
}

func transpileFunctionDecl(n *ast.FunctionDecl, p *program.Program) error {
	var body *goast.BlockStmt

	// Always register the new function. Only from this point onwards will
	// we be allowed to refer to the function.
	if program.GetFunctionDefinition(n.Name) == nil {
		program.AddFunctionDefinition(program.FunctionDefinition{
			Name:       n.Name,
			ReturnType: "int",
			// FIXME
			ArgumentTypes: []string{},
			Substitution:  "",
		})
	}

	// If the function has a direct substitute in Go we do not want to
	// output the C definition of it.
	if f := program.GetFunctionDefinition(n.Name); f != nil &&
		f.Substitution != "" {
		return nil
	}

	hasBody := false
	for _, c := range n.Children {
		if b, ok := c.(*ast.CompoundStmt); ok {
			var err error
			body, err = transpileToBlockStmt(b, p)
			if err != nil {
				return err
			}

			hasBody = true
			break
		}
	}

	if n.Name == "__istype" ||
		n.Name == "__isctype" ||
		n.Name == "__wcwidth" ||
		n.Name == "__sputc" ||
		n.Name == "__inline_signbitf" ||
		n.Name == "__inline_signbitd" ||
		n.Name == "__inline_signbitl" {
		return nil
	}

	if hasBody {
		fieldList, err := getFieldList(n, p)
		if err != nil {
			return err
		}

		p.File.Decls = append(p.File.Decls, &goast.FuncDecl{
			Doc:  nil,
			Recv: nil,
			Name: goast.NewIdent(n.Name),
			Type: &goast.FuncType{
				Params:  fieldList,
				Results: nil,
			},
			Body: body,
		})
	}

	return nil
}

func getFieldList(f *ast.FunctionDecl, p *program.Program) (*goast.FieldList, error) {
	// The main() function does not have arguments or a return value.
	if f.Name == "main" {
		return &goast.FieldList{}, nil
	}

	r := []*goast.Field{}
	for _, n := range f.Children {
		if v, ok := n.(*ast.ParmVarDecl); ok {
			r = append(r, &goast.Field{
				Doc:     nil,
				Names:   []*goast.Ident{goast.NewIdent(v.Name)},
				Type:    goast.NewIdent(types.ResolveType(p, v.Type)),
				Tag:     nil,
				Comment: nil,
			})
		}
	}

	return &goast.FieldList{
		List: r,
	}, nil
}

func transpileReturnStmt(n *ast.ReturnStmt, p *program.Program) (*goast.ReturnStmt, error) {
	return &goast.ReturnStmt{
		Return:  token.NoPos,
		Results: nil,
	}, nil
}

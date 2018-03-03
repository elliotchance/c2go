package transpiler

import (
	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

func transpileTranslationUnitDecl(p *program.Program, n *ast.TranslationUnitDecl) (
	decls []goast.Decl, err error) {

	for i := 0; i < len(n.Children()); i++ {
		presentNode := n.Children()[i]
		var runAfter func()

		if rec, ok := presentNode.(*ast.RecordDecl); ok && rec.Name == "" {
			if i+1 < len(n.Children()) {
				switch recNode := n.Children()[i+1].(type) {
				case *ast.VarDecl:
					rec.Name = types.GetBaseType(recNode.Type)
				case *ast.TypedefDecl:
					rec.Name = types.GetBaseType(recNode.Type)
				}
			}
		}
		if rec, ok := presentNode.(*ast.RecordDecl); ok {
			// ignore RecordDecl if haven`t definition
			if rec.Name == "" && !rec.Definition {
				continue
			}
		}

		var d []goast.Decl
		d, err = transpileToNode(presentNode, p)
		if err != nil {
			p.AddMessage(p.GenerateErrorMessage(err, n))
			err = nil
		} else {
			decls = append(decls, d...)
			if runAfter != nil {
				runAfter()
			}
		}

	}
	return
}

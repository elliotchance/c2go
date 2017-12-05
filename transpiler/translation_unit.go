package transpiler

import (
	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

func transpileTranslationUnitDecl(p *program.Program, n *ast.TranslationUnitDecl) (decls []goast.Decl, err error) {
	for i := 0; i < len(n.Children()); i++ {
		switch v := n.Children()[i].(type) {
		case *ast.RecordDecl:
			// specific for `typedef struct` without name
			if v.Name != "" || i == len(n.Children())-1 {
				var d []goast.Decl
				d, err = transpileRecordDecl(p, v)
				if err != nil {
					return
				}
				decls = append(decls, d...)
			}
			for counter := 1; i+counter < len(n.Children()); counter++ {
				if vv, ok := n.Children()[i+counter].(*ast.TypedefDecl); ok {
					nameTypedefStruct := vv.Name
					fields := v.Children()
					// create a struct in according to
					// name and fields
					var recordDecl ast.RecordDecl
					recordDecl.Name = nameTypedefStruct
					recordDecl.ChildNodes = fields
					var d []goast.Decl
					d, err = transpileRecordDecl(p, &recordDecl)
					if err != nil {
						p.AddMessage(p.GenerateErrorMessage(err, n))
						err = nil
					} else {
						decls = append(decls, d...)
					}
				} else {
					counter--
					i = i + counter
					break
				}
			}
		default:
			var d []goast.Decl
			d, err = transpileToNode(n.Children()[i], p)
			if err != nil {
				p.AddMessage(p.GenerateErrorMessage(err, n))
				err = nil
			} else {
				decls = append(decls, d...)
			}
		}
	}
	return
}

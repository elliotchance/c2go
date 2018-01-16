package transpiler

import (
	goast "go/ast"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

func transpileTranslationUnitDecl(p *program.Program, n *ast.TranslationUnitDecl) (decls []goast.Decl, err error) {
	for i := 0; i < len(n.Children()); i++ {
		switch v := n.Children()[i].(type) {
		case *ast.RecordDecl:
			// for case :
			// typedef struct C C;
			// typedef union  C C;
			if len(v.Children()) == 0 {
				if i+1 < len(n.Children()) {
					if vv, ok := n.Children()[i+1].(*ast.TypedefDecl); ok {
						if isSameTypedefNames(vv) {
							i++
							continue
						}
					}
				}
			}
			// specific for `typedef struct`, `typedef union` without name
			if v.Name != "" || i == len(n.Children())-1 {
				var d []goast.Decl
				d, err = transpileRecordDecl(p, v)
				if err != nil {
					return
				}
				decls = append(decls, d...)
			}
			for counter := 1; i+counter < len(n.Children()); counter++ {
				if vv, ok := n.Children()[i+counter].(*ast.TypedefDecl); ok && !types.IsFunction(vv.Type) {
					nameTypedefStruct := vv.Name
					fields := v.Children()
					// create a struct in according to
					// name and fields
					var recordDecl ast.RecordDecl
					recordDecl.Name = nameTypedefStruct
					recordDecl.Kind = "struct"
					if strings.Contains(vv.Type, "union ") {
						recordDecl.Kind = "union"
					}
					recordDecl.ChildNodes = fields
					var d []goast.Decl
					d, err = transpileRecordDecl(p, &recordDecl)
					if err != nil {
						p.AddMessage(p.GenerateErrorMessage(err, n))
						err = nil
					} else {
						decls = append(decls, d...)
					}
					break
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

func isSameTypedefNames(v *ast.TypedefDecl) bool {
	// for structs :
	/*
	   TypedefDecl 0x33da010 <col:3, col:21> col:21 referenced Uq 'struct Uq':'struct Uq'
	   `-ElaboratedType 0x33d9fc0 'struct Uq' sugar
	     `-RecordType 0x33d9fa0 'struct Uq'
	       `-Record 0x33da090 'Uq'
	*/
	// for unions:
	/*
		TypedefDecl 0x38bc070 <col:1, col:23> col:23 referenced myunion 'union myunion':'union myunion'
		`-ElaboratedType 0x38bc020 'union myunion' sugar
		  `-RecordType 0x38bc000 'union myunion'
		    `-Record 0x38bc0d8 'myunion'
	*/
	if ("struct "+v.Name == v.Type2 || "union "+v.Name == v.Type2) && v.Type == v.Type2 {
		if vv, ok := v.Children()[0].(*ast.ElaboratedType); ok && vv.Type == v.Type {
			if vvv, ok := vv.Children()[0].(*ast.RecordType); ok && vvv.Type == v.Type2 {
				if vvvv, ok := vvv.Children()[0].(*ast.Record); ok && vvvv.Type == v.Name {
					return true
				}
			}
		}
	}
	return false
}

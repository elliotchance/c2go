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
		if rec, ok := n.Children()[i].(*ast.RecordDecl); ok {
			if i+1 < len(n.Children()) {
				switch recNode := n.Children()[i+1].(type) {
				case *ast.VarDecl:
					name := types.GenerateCorrectType(recNode.Type)
					if strings.HasPrefix(name, "union ") {
						rec.Name = name[len("union "):]
						recNode.Type = "union " + name
					}
					if strings.HasPrefix(name, "struct ") {
						rec.Name = name[len("struct "):]
						recNode.Type = "struct " + name
					}
				case *ast.TypedefDecl:
					if strings.Contains(recNode.Type, "__locale_struct") {
						i++
						continue
					}
					if isSameTypedefNames(recNode) {
						i++
						continue
					}
					name := types.GenerateCorrectType(recNode.Type)
					if strings.HasPrefix(name, "union ") {
						rec.Name = name[len("union "):]
						recNode.Type = "union " + name
					}
					if strings.HasPrefix(name, "struct ") {
						rec.Name = name[len("struct "):]
						recNode.Type = "struct " + name
					}
				}
			}
		}
		var d []goast.Decl
		d, err = transpileToNode(n.Children()[i], p)
		if err != nil {
			p.AddMessage(p.GenerateErrorMessage(err, n))
			err = nil
		} else {
			decls = append(decls, d...)
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

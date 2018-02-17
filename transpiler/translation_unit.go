package transpiler

import (
	"fmt"
	goast "go/ast"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

func transpileTranslationUnitDecl(p *program.Program, n *ast.TranslationUnitDecl) (
	decls []goast.Decl, err error) {

	for i := 0; i < len(n.Children()); i++ {
		presentNode := n.Children()[i]
		var runAfter func()

		if rec, ok := presentNode.(*ast.RecordDecl); ok {
			fmt.Println("Record = ", rec.Name, "\t", rec.Kind, "\t", rec.Pos)
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
					if isSameTypedefNames(recNode) {
						i++
						continue
					}
					name := types.GenerateCorrectType(types.CleanCType(recNode.Type))
					// fmt.Println("Typedef ", recNode.Name, " -> ", recNode.Type, "||", recNode.Type2)
					if strings.HasPrefix(name, "union ") {
						if recNode.Type == "union "+rec.Name {
							names := []string{rec.Name, recNode.Name}
							for _, name := range names {
								rec.Name = name
								var d []goast.Decl
								d, err = transpileToNode(rec, p)
								if err != nil {
									p.AddMessage(p.GenerateErrorMessage(err, n))
									err = nil
								} else {
									decls = append(decls, d...)
								}
							}

							i++
							continue
						} else {
							rec.Name = name[len("union "):]
							recNode.Type = "union " + name
						}
					}
					if strings.HasPrefix(name, "struct ") {

						// fmt.Println(">>>>  Typedef ", recNode.Name, " -> ", recNode.Type, "||", recNode.Type2)
						// fmt.Println(">>>> ", name, "\t")
						// rec.name = name[len("struct "):]
						// recnode.type = "struct " + name

						// From :
						// TypedefDecl __locale_t 'struct __locale_struct *'
						// To   :
						// VarDecl st7a 'struct st4':'struct st4'
						if rec.Name != "" {
							runAfter = func() {
								// var v ast.VarDecl
								// v.Name = recNode.Name
								// v.Type = recNode.Type[len("struct "):]
								//
								// v.Type = types.GenerateCorrectType(v.Type)
								// fmt.Println("VarDecl : ", v.Name, " > ", v.Type)

								var d []goast.Decl
								d, err = transpileToNode(recNode, p)
								if err != nil {
									p.AddMessage(p.GenerateErrorMessage(err, n))
									err = nil
								} else {
									decls = append(decls, d...)
								}
							}

							i++
						} else {
							rec.Name = name[len("struct "):]
							recNode.Type = "struct " + name
						}
					}
				}
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

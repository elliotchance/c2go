package transpiler

import (
	goast "go/ast"
	"strings"

	"fmt"
	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"reflect"
)

func transpileTranslationUnitDecl(p *program.Program, n *ast.TranslationUnitDecl) (
	decls []goast.Decl, err error) {

	for i := 0; i < len(n.Children()); i++ {
		presentNode := n.Children()[i]
		var runAfter func()

		if rec, ok := presentNode.(*ast.RecordDecl); ok {
			if i+1 < len(n.Children()) {
				switch recNode := n.Children()[i+1].(type) {
				case *ast.VarDecl:
					name := types.GenerateCorrectType(types.CleanCType(recNode.Type))
					if rec.Name == "" {
						recNode.Type = types.GenerateCorrectType(recNode.Type)
						recNode.Type2 = types.GenerateCorrectType(recNode.Type2)
						if strings.HasPrefix(name, "union ") {
							rec.Name = name[len("union "):]
							recNode.Type = types.CleanCType("union " + name)
						}
						if strings.HasPrefix(name, "struct ") {
							name = types.GetBaseType(name)
							rec.Name = name[len("struct "):]
						}
					}
				case *ast.TypedefDecl:
					if isSameTypedefNames(recNode) && !rec.Definition {
						// this is just the declaration of a type, the implementation comes later
						i++
						continue
					}
					name := types.GenerateCorrectType(types.CleanCType(recNode.Type))
					if isSameTypedefNames(recNode) {
						i++
						// continue on, so that this type is defined
					} else if strings.HasPrefix(name, "union ") {
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
							recNode.Type = types.CleanCType("union " + name)
						}
					} else if strings.HasPrefix(name, "struct ") {
						if rec.Name != "" {
							runAfter = func() {
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
							recNode.Type = types.CleanCType("struct " + name)
						}
					}
					if !rec.Definition {
						// This was not the real definition of the type,
						// we have to go and look it up
						var typeToDeclare *ast.RecordDecl
						records := ast.GetAllNodesOfType(recNode, reflect.TypeOf(&ast.Record{}))
						if len(records) > 0 {
							record := records[0].(*ast.Record)
							if n, ok := p.NodeMap[record.Addr]; ok {
								if toDeclare, ok2 := n.(*ast.RecordDecl); ok2 {
									typeToDeclare = toDeclare
								}
							}
						}
						if typeToDeclare == nil {
							p.AddMessage(p.GenerateWarningMessage(fmt.Errorf("could not lookup type definition for : %v", rec.Name), rec))
							typeToDeclare = rec
						}
						p.DeclareType(typeToDeclare, types.GenerateCorrectType(rec.Name))
						if runAfter != nil {
							runAfter()
						}
						continue
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

package indicator

import (
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/ast"
)

var functionDeclPos []ast.Position

func InjectInC(tree []ast.Node, inputFiles []string) (err error) {
	for j := range inputFiles {
		fmt.Printf("# File : %s\n", inputFiles[j])
		for i := range tree[0].Children() {
			n := tree[0].Children()[i]
			if fd, ok := n.(*ast.FunctionDecl); ok {
				if strings.Contains(inputFiles[j], fd.Pos.File) {
					fmt.Printf("# Found function : %s\n", fd.Name)
					functionDeclPos = append(functionDeclPos, fd.Pos)
				}
			}
		}
	}
	fmt.Printf("%#v\n", functionDeclPos)
	return nil
}

package transpiler

import (
	"fmt"
	"go/token"
	"strings"

	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
)

func transpileStringLiteral(n *ast.StringLiteral) *goast.BasicLit {
	// TODO: There are other escape characters.
	value := fmt.Sprintf("\"%s\"", strings.Replace(n.Value, "\n", "\\n", -1))

	return &goast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}
}

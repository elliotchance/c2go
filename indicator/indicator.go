package indicator

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/elliotchance/c2go/ast"
)

var uniqIterator int

func getInjectCode() string {
	uniqIterator++
	return fmt.Sprintf("printf(\"%d \\n\");", uniqIterator)
}

func InjectInC(tree []ast.Node, inputFiles []string) (err error) {
	for j := range inputFiles {
		// find functions
		fmt.Printf("# File : %s\n", inputFiles[j])
		functionDeclPos := make([]ast.Position, 0, 10)
		for i := range tree[0].Children() {
			n := tree[0].Children()[i]
			if fd, ok := n.(*ast.FunctionDecl); ok {
				if co, ok := fd.Children()[0].(*ast.CompoundStmt); ok {
					if strings.Contains(inputFiles[j], fd.Pos.File) {
						fmt.Printf("# Found function : %s\n", fd.Name)
						functionDeclPos = append(functionDeclPos, co.Pos)
					}
				}
			}
		}
		// inject
		var base []byte
		base, err = ioutil.ReadFile(inputFiles[j])
		if err != nil {
			return
		}
		r := bytes.NewReader(base)
		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanLines)
		var buf bytes.Buffer
		var linePos int
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if linePos == 0 {
				buf.WriteString("#include <stdio.h>")
			}
			buf.WriteString(line)
			if len(line) > 0 {
				if line[len(line)-1] == ';' {
					var isInsideFunc bool
					for i := range functionDeclPos {
						if functionDeclPos[i].Line <= linePos && linePos < functionDeclPos[i].LineEnd {
							isInsideFunc = true
							break
						}
					}
					if isInsideFunc {
						buf.WriteString(getInjectCode())
					}
				}
			}
			buf.WriteByte('\n')
			linePos++
		}
		// Write data to dst
		err = ioutil.WriteFile(inputFiles[j], buf.Bytes(), 0644)
		if err != nil {
			return
		}
		fmt.Println(buf.String())
	}
	return nil
}

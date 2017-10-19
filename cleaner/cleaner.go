package cleaner

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"golang.org/x/tools/imports"
	"honnef.co/go/tools/lint/lintutil"
	"honnef.co/go/tools/unused"
)

type un int

const (
	unusedConstans un = iota
	unusedFunction
	unusedType
	unusedVariable
)

const (
	suffix = " is unused (U1000)"
)

var unusedMap = map[string]un{
	"const": unusedConstans,
	"func":  unusedFunction,
	"type":  unusedType,
	"var":   unusedVariable,
}

// Go - clean single Go code from unused variables, ...
func Go(inFile, outFile string, verbose bool) (err error) {
	if verbose {
		fmt.Println("Start cleaning ...")
	}

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("Error : %v", r)
			}
		}
	}()

	stderr := os.Stderr
	os.Stderr, err = ioutil.TempFile("", "temp")
	if err != nil {
		return err
	}
	defer func() {
		os.Stderr = stderr
	}()
	defer func() {
		var buf []byte
		_, err2 := os.Stderr.Read(buf)
		if err2 != nil {
			err = err2
			return
		}
		if len(buf) > 0 {
			err = fmt.Errorf(string(buf))
		}
	}()

	if verbose {
		fmt.Println("\tLinter working ...")
	}
	// prepare linter
	var mode unused.CheckMode
	mode |= unused.CheckConstants
	//mode |= unused.CheckFields
	mode |= unused.CheckFunctions
	mode |= unused.CheckTypes
	mode |= unused.CheckVariables
	checker := unused.NewChecker(mode)
	l := unused.NewLintChecker(checker)

	fs := lintutil.FlagSet("unused")
	err = fs.Parse([]string{inFile})
	if err != nil {
		return fmt.Errorf("Error in flag parsing : %v", err)
	}

	// take result of linter work
	ps, _, err := lintutil.Lint(l, fs.Args(), &lintutil.Options{})
	if err != nil {
		return fmt.Errorf("Error in linter : %v", err)
	}

	// checking stdErr from linter is not empty
	{
		var buf []byte
		_, err := os.Stderr.Read(buf)
		if err != nil {
			return err
		}
		if len(buf) > 0 {
			return fmt.Errorf(string(buf))
		}
	}

	// linter is not found any unused elements
	if len(ps) == 0 {
		return nil
	}

	if verbose {
		fmt.Println("\tCreating AST tree...")
	}
	// create ast tree
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, inFile, nil, 0)
	if err != nil {
		return fmt.Errorf("Error: Cannot parse : %v", err)
	}

	// parse unused strings
	if verbose {
		fmt.Println("\tParsing lines to unused entities...")
	}
	type unusedParameter struct {
		u        un
		name     string
		position token.Pos
	}
	var unusedParameters []unusedParameter
	for _, p := range ps {
		p.Text = p.Text[0 : len(p.Text)-len(suffix)]
		for k, v := range unusedMap {
			if strings.Contains(p.Text, k) {
				unusedParameters = append(unusedParameters, unusedParameter{
					u:        v,
					name:     strings.TrimSpace(p.Text[len(k):len(p.Text)]),
					position: p.Position,
				})
			}
		}
	}

	var removeItems []int

	// remove unused parts of AST tree
	if verbose {
		fmt.Println("\tRemoving unused entities...")
	}
	for _, param := range unusedParameters {
		switch param.u {
		// remove unused constants
		case unusedConstans:
			{
				for i := 0; i < len(tree.Decls); i++ {
					gen, ok := tree.Decls[i].(*ast.GenDecl)
					if !ok || gen == (*ast.GenDecl)(nil) || gen.Tok != token.CONST {
						goto nextConstDecl
					}
					if s, ok := gen.Specs[0].(*ast.ValueSpec); ok {
						for _, n := range s.Names {
							if n.String() == param.name && param.position == s.Pos() {
								removeItems = append(removeItems, i)
								continue
							}
						}
					}
				nextConstDecl:
				}
			}

		// remove unused functions
		case unusedFunction:
			{
				for i := 0; i < len(tree.Decls); i++ {
					gen, ok := tree.Decls[i].(*ast.FuncDecl)
					if !ok || gen == (*ast.FuncDecl)(nil) {
						continue
					}
					if gen.Name.String() == param.name && gen.Pos() <= param.position && param.position <= gen.End() {
						removeItems = append(removeItems, i)
						continue
					}
				}
			}

		// remove unused types
		case unusedType:
			{
				for i := 0; i < len(tree.Decls); i++ {
					gen, ok := tree.Decls[i].(*ast.GenDecl)
					if !ok || gen == (*ast.GenDecl)(nil) || gen.Tok != token.TYPE {
						continue
					}
					if s, ok := gen.Specs[0].(*ast.TypeSpec); ok {
						if s.Name.String() == param.name && param.position == s.Pos() {
							removeItems = append(removeItems, i)
							continue
						}
					}
				}
			}

		// remove unused variables
		case unusedVariable:
			{
				for i := 0; i < len(tree.Decls); i++ {
					gen, ok := tree.Decls[i].(*ast.GenDecl)
					if !ok || gen == (*ast.GenDecl)(nil) || gen.Tok != token.VAR {
						continue
					}
					if s, ok := gen.Specs[0].(*ast.ValueSpec); ok {
						for _, n := range s.Names {
							if n.String() == param.name && param.position == s.Pos() {
								removeItems = append(removeItems, i)
								continue
							}
						}
					}
				}
			}

		}
	}
	// removing methods of types
	// example of type:
	// 		type pthread_attr_t [56]byte
	// example of unused result:
	// 		exit.go:135:6: type pthread_attr_t is unused (U1000)
	// example of that type method:
	// 		func (self *pthread_attr_t) cast(t reflect.Type) reflect.Value {
	for _, param := range unusedParameters {
		if param.u != unusedType {
			continue
		}
		for i := 0; i < len(tree.Decls); i++ {
			gen, ok := tree.Decls[i].(*ast.FuncDecl)
			if !ok || gen == (*ast.FuncDecl)(nil) {
				continue
			}
			if gen.Recv == (*ast.FieldList)(nil) {
				continue
			}
			if s, ok := gen.Recv.List[0].Type.(*ast.StarExpr); ok {
				if t, ok := s.X.(*ast.Ident); ok {
					if t.Name == param.name {
						removeItems = append(removeItems, i)
					}
				}
			}
		}
	}

	// Sorting slice with remove index elements
	sort.Ints(removeItems)

	// remove Decl element from tree
	if verbose {
		fmt.Println("\tAST tree corrections...")
	}
	tempTree := make([]ast.Decl, len(tree.Decls)-len(removeItems))
	var counter int
	for i := range tree.Decls {
		var found bool
		for _, remove := range removeItems {
			if i == remove {
				found = true
			}
		}
		if !found {
			tempTree[counter] = tree.Decls[i]
			counter++
		}
	}
	tree.Decls = tempTree

	// convert AST tree to Go code
	if verbose {
		fmt.Println("\tConverting AST tree to Go code...")
	}
	var buf bytes.Buffer
	err = format.Node(&buf, fset, tree)
	if err != nil {
		return fmt.Errorf("Error: convert AST tree to Go code : %v", err)
	}

	// Remove imports
	b, err := imports.Process(inFile, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("Error: Cannot modify imports : %v", err)
	}
	buf = *bytes.NewBuffer(b)

	// write buffer with Go code to file
	if verbose {
		fmt.Println("\tWriting Go code...")
	}
	err = ioutil.WriteFile(outFile, buf.Bytes(), 0755)
	if err != nil {
		return fmt.Errorf("writing Go output file failed: %v", err)
	}

	return nil
}

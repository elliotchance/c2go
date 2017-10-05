// Package c2go contains the main function for running the executable.
//
// Installation
//
//     go get -u github.com/elliotchance/c2go
//
// Usage
//
//     c2go myfile.c
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"errors"
	"reflect"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/transpiler"
	"honnef.co/go/tools/lint/lintutil"
	"honnef.co/go/tools/unused"
)

// Version can be requested through the command line with:
//
//     c2go -v
//
// See https://github.com/elliotchance/c2go/wiki/Release-Process
const Version = "v0.16.2 Radium 2017-09-18"

var stderr io.Writer = os.Stderr

// ProgramArgs defines the options available when processing the program. There
// is no constructor since the zeroed out values are the appropriate defaults -
// you need only set the options you need.
//
// TODO: Better separation on CLI modes
// https://github.com/elliotchance/c2go/issues/134
type ProgramArgs struct {
	verbose     bool
	ast         bool
	inputFile   string
	outputFile  string
	packageName string

	// A private option to output the Go as a *_test.go file.
	outputAsTest bool
}

func readAST(data []byte) []string {
	uncolored := regexp.MustCompile(`\x1b\[[\d;]+m`).ReplaceAll(data, []byte{})
	return strings.Split(string(uncolored), "\n")
}

type treeNode struct {
	indent int
	node   ast.Node
}

func convertLinesToNodes(lines []string) []treeNode {
	nodes := []treeNode{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// It is tempting to discard null AST nodes, but these may
		// have semantic importance: for example, they represent omitted
		// for-loop conditions, as in for(;;).
		line = strings.Replace(line, "<<<NULL>>>", "NullStmt", 1)
		trimmed := strings.TrimLeft(line, "|\\- `")
		node := ast.Parse(trimmed)
		indentLevel := (len(line) - len(trimmed)) / 2
		nodes = append(nodes, treeNode{indentLevel, node})
	}

	return nodes
}

// buildTree converts an array of nodes, each prefixed with a depth into a tree.
func buildTree(nodes []treeNode, depth int) []ast.Node {
	if len(nodes) == 0 {
		return []ast.Node{}
	}

	// Split the list into sections, treat each section as a tree with its own
	// root.
	sections := [][]treeNode{}
	for _, node := range nodes {
		if node.indent == depth {
			sections = append(sections, []treeNode{node})
		} else {
			sections[len(sections)-1] = append(sections[len(sections)-1], node)
		}
	}

	results := []ast.Node{}
	for _, section := range sections {
		slice := []treeNode{}
		for _, n := range section {
			if n.indent > depth {
				slice = append(slice, n)
			}
		}

		children := buildTree(slice, depth+1)
		for _, child := range children {
			section[0].node.AddChild(child)
		}
		results = append(results, section[0].node)
	}

	return results
}

func toJSON(tree []interface{}) []map[string]interface{} {
	r := make([]map[string]interface{}, len(tree))

	for j, n := range tree {
		rn := reflect.ValueOf(n).Elem()
		r[j] = make(map[string]interface{})
		r[j]["node"] = rn.Type().Name()

		for i := 0; i < rn.NumField(); i++ {
			name := strings.ToLower(rn.Type().Field(i).Name)
			value := rn.Field(i).Interface()

			if name == "children" {
				v := value.([]interface{})

				if len(v) == 0 {
					continue
				}

				value = toJSON(v)
			}

			r[j][name] = value
		}
	}

	return r
}

func check(prefix string, e error) {
	if e != nil {
		panic(prefix + e.Error())
	}
}

// Start begins transpiling an input file.
func Start(args ProgramArgs) error {
	if os.Getenv("GOPATH") == "" {
		return fmt.Errorf("The $GOPATH must be set")
	}

	// 1. Compile it first (checking for errors)
	_, err := os.Stat(args.inputFile)
	if err != nil {
		return fmt.Errorf("Input file is not found")
	}

	// 2. Preprocess
	var pp []byte
	{
		// See : https://clang.llvm.org/docs/CommandGuide/clang.html
		// clang -E <file>    Run the preprocessor stage.
		cmd := exec.Command("clang", "-E", args.inputFile)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
		}
		pp = []byte(out.String())
	}

	ppFilePath := path.Join("/tmp", "pp.c")
	err = ioutil.WriteFile(ppFilePath, pp, 0644)
	if err != nil {
		return fmt.Errorf("writing to %s failed: %v", ppFilePath, err)
	}

	// 3. Generate JSON from AST
	astPP, err := exec.Command("clang", "-Xclang", "-ast-dump", "-fsyntax-only", ppFilePath).Output()
	if err != nil {
		// If clang fails it still prints out the AST, so we have to run it
		// again to get the real error.
		errBody, _ := exec.Command("clang", ppFilePath).CombinedOutput()

		panic("clang failed: " + err.Error() + ":\n\n" + string(errBody))
	}

	lines := readAST(astPP)
	if args.ast {
		for _, l := range lines {
			fmt.Println(l)
		}
		fmt.Println()

		return nil
	}

	p := program.NewProgram()
	p.Verbose = args.verbose
	p.OutputAsTest = true // args.outputAsTest

	nodes := convertLinesToNodes(lines)
	tree := buildTree(nodes, 0)
	ast.FixPositions(tree)

	// Repair the floating literals. See RepairFloatingLiteralsFromSource for
	// more information.
	floatingErrors := ast.RepairFloatingLiteralsFromSource(tree[0], ppFilePath)

	for _, fErr := range floatingErrors {
		message := fmt.Sprintf("could not read exact floating literal: %s",
			fErr.Err.Error())
		p.AddMessage(p.GenerateWarningMessage(errors.New(message), fErr.Node))
	}

	err = transpiler.TranspileAST(args.inputFile, args.packageName, p, tree[0].(ast.Node))
	if err != nil {
		panic(err)
	}

	outputFilePath := args.outputFile

	if outputFilePath == "" {
		cleanFileName := filepath.Clean(filepath.Base(args.inputFile))
		extension := filepath.Ext(args.inputFile)
		outputFilePath = cleanFileName[0:len(cleanFileName)-len(extension)] + ".go"
	}

	err = ioutil.WriteFile(outputFilePath, []byte(p.String()), 0755)
	if err != nil {
		return fmt.Errorf("writing C output file failed: %v", err)
	}

	return nil
}

// newTempFile - returns temp file
func newTempFile(dir, prefix, suffix string) (*os.File, error) {
	for index := 1; index < 10000; index++ {
		path := filepath.Join(dir, fmt.Sprintf("%s%03d%s", prefix, index, suffix))
		if _, err := os.Stat(path); err != nil {
			return os.Create(path)
		}
	}
	return nil, fmt.Errorf("could not create file: %s%03d%s", prefix, 1, suffix)
}

var (
	versionFlag       = flag.Bool("v", false, "print the version and exit")
	transpileCommand  = flag.NewFlagSet("transpile", flag.ContinueOnError)
	verboseFlag       = transpileCommand.Bool("V", false, "print progress as comments")
	outputFlag        = transpileCommand.String("o", "", "output Go generated code to the specified file")
	packageFlag       = transpileCommand.String("p", "main", "set the name of the generated package")
	transpileHelpFlag = transpileCommand.Bool("h", false, "print help information")
	astCommand        = flag.NewFlagSet("ast", flag.ContinueOnError)
	astHelpFlag       = astCommand.Bool("h", false, "print help information")
)

func main() {
	code := runCommand()
	if code != 0 {
		os.Exit(code)
	}
}

func runCommand() int {
	flag.Usage = func() {
		usage := "Usage: %s [-v] [<command>] [<flags>] file.c\n\n"
		usage += "Commands:\n"
		usage += "  transpile\ttranspile an input C source file to Go\n"
		usage += "  ast\t\tprint AST before translated Go code\n\n"

		usage += "Flags:\n"
		fmt.Fprintf(stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	transpileCommand.SetOutput(stderr)
	astCommand.SetOutput(stderr)

	flag.Parse()

	if *versionFlag {
		// Simply print out the version and exit.
		fmt.Println(Version)
		return 0
	}

	if flag.NArg() < 1 {
		flag.Usage()
		return 1
	}

	args := ProgramArgs{verbose: *verboseFlag, ast: false}

	switch os.Args[1] {
	case "ast":
		err := astCommand.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("ast command cannot parse: %v", err)
			return 1
		}

		if *astHelpFlag || astCommand.NArg() == 0 {
			fmt.Fprintf(stderr, "Usage: %s ast file.c\n", os.Args[0])
			astCommand.PrintDefaults()
			return 1
		}

		args.ast = true
		args.inputFile = astCommand.Arg(0)
	case "transpile":
		err := transpileCommand.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("transpile command cannot parse: %v", err)
			return 1
		}

		if *transpileHelpFlag || transpileCommand.NArg() == 0 {
			fmt.Fprintf(stderr, "Usage: %s transpile [-V] [-o file.go] [-p package] file.c\n", os.Args[0])
			transpileCommand.PrintDefaults()
			return 1
		}

		args.inputFile = transpileCommand.Arg(0)
		args.outputFile = *outputFlag
		args.packageName = *packageFlag
	default:
		flag.Usage()
		return 1
	}

	if err := Start(args); err != nil {
		fmt.Printf("Error: %v", err)
		return 1
	}

	var mode unused.CheckMode
	mode |= unused.CheckConstants
	//mode |= unused.CheckFields
	mode |= unused.CheckFunctions
	mode |= unused.CheckTypes
	mode |= unused.CheckVariables
	checker := unused.NewChecker(mode)
	l := unused.NewLintChecker(checker)

	fs := lintutil.FlagSet("unused")

	//TODO wrong
	_ = fs.Parse([]string{"hello.go"})

	//TODO more elegant
	ps, _ /* lprog*/, err := lintutil.Lint(l, fs.Args(), &lintutil.Options{})
	if err != nil {
		return 0
	}

	// create ast tree
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, "hello.go", nil, 0)
	//fmt.Println("tree = ", tree)
	//fmt.Println("err  = ", err)
	_ = tree
	_ = err

Back:
	for _, p := range ps {
		/*
			pos := lprog.Fset.Position(p.Position)
				if strings.Contains(p.Text, unusedConstantans) {
					continue
				}
				if strings.Contains(p.Text, unusedFunction) {
					continue
				}
				if strings.Contains(p.Text, unusedType) {
					continue
				}
				if strings.Contains(p.Text, unusedVariable) {
					continue
				}
				fmt.Printf("%v|||| %s\n", pos, p.Text)
		*/
		if strings.Contains(p.Text, unusedType) {
			name := strings.TrimSpace(p.Text[len(unusedType) : len(p.Text)-len(postfix)])
			// fmt.Printf("%s\n", name)

			// Print the AST.

			// var b bytes.Buffer
			// _ = goast.Fprint(&b, fset, tree, nil)
			// fmt.Printf("%s\n", b.String())
			// find const
			for i := range tree.Decls {
				decl := tree.Decls[i]
				if gen, ok := decl.(*goast.GenDecl); ok && gen.Tok != token.CONST {
					/*
						for ff := range gen.Specs {
							if s, ok := gen.Specs[ff].(*goast.ValueSpec); ok {
								for j := range s.Names {
									fmt.Println("Value Spec", j)
								}
								fmt.Printf("Value type = %#v\n", s.Type)
								for j := range s.Values {
									fmt.Println("Valu val ", s.Values[j])
								}
								fmt.Println("")
							}
						}
					*/
					if s, ok := gen.Specs[0].(*goast.TypeSpec); ok {
						if strings.Contains(s.Name.String(), name) {
							var rr []goast.Decl
							if i != 0 {
								rr = append(rr, tree.Decls[0:i]...)
							}
							if i != len(tree.Decls)-1 {
								rr = append(rr, tree.Decls[i+1:len(tree.Decls)-1]...)
							}
							fmt.Println("Len = ", len(rr), "|", len(tree.Decls))
							tree.Decls = rr
							fmt.Println("name = ", name)
							goto Back

							//					tree.Decls = append(tree.Decls[i], tree.Decls[i+1]...)
							// copy(tree.Decls[i:], tree.Decls[i+1:])
							// tree.Decls = tree.Decls[:len(tree.Decls)-1]
						}
					}
					/*
						// if dd, ok := gen.Specs[0].(*goast.ValueSpec); ok {
						fmt.Println("FIND =>", tree.Decls[i].Pos())
						fmt.Println("1: ", gen.Tok)
						fmt.Printf("2: %#v\n", gen.Specs)
						for h := range gen.Specs {
							fmt.Printf("2.%v = %v\n", h, gen.Specs[h])
							fmt.Printf("2.%v = %+v\n", h, gen.Specs[h])
							fmt.Printf("2.%v = %#v\n", h, gen.Specs[h])
						}
						fmt.Println("3: ", gen.Lparen)
						fmt.Println("4: ", gen.Rparen)
					*/
					// fmt.Println(" DD= ", dd)
					// }
					// copy(f.Decls[i:], f.Decls[i+1:])
					// f.Decls = f.Decls[:len(f.Decls)-1]
					// for j := i; j < len(tree.Decls); j++ {
					// 	decl2 := tree.Decls[j]
					// 	if gen, ok := decl2.(*goast.ValueSpec); ok {
					// 		fmt.Println(">==", gen.Names)
					// 		break
					// 	}
					// }
				}
			}
			// remove from ast tree
			//

			//break
			//continue
		}
	}

	// var b bytes.Buffer
	// _ = goast.Fprint(&b, fset, tree, nil)
	// _, _ = fmt.Printf("%s\n", b.String())

	//_ = goast.Fprint(fset, tree)

	var buf bytes.Buffer
	_ = printer.Fprint(&buf, fset, tree) // funcAST.Body)
	_, _ = fmt.Printf("%s\n", buf.String())

	_, _ = fmt.Printf("%+v\n", tree)

	return 0
}

const (
	unusedConstantans = "const"
	unusedFunction    = "func"
	unusedType        = "type"
	unusedVariable    = "var"
	postfix           = " is unused (U1000)"
)

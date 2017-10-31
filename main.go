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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"errors"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/transpiler"
)

// Version can be requested through the command line with:
//
//     c2go -v
//
// See https://github.com/elliotchance/c2go/wiki/Release-Process
const Version = "v0.16.10 Radium 2017-10-31"

var stderr io.Writer = os.Stderr

// ProgramArgs defines the options available when processing the program. There
// is no constructor since the zeroed out values are the appropriate defaults -
// you need only set the options you need.
//
// TODO: Better separation on CLI modes
// https://github.com/elliotchance/c2go/issues/134
//
// Do not instantiate this directly. Instead use DefaultProgramArgs(); then
// modify any specific attributes.
type ProgramArgs struct {
	verbose     bool
	ast         bool
	inputFile   string
	outputFile  string
	packageName string

	// A private option to output the Go as a *_test.go file.
	outputAsTest bool
}

// DefaultProgramArgs default value of ProgramArgs
func DefaultProgramArgs() ProgramArgs {
	return ProgramArgs{
		verbose:      false,
		ast:          false,
		packageName:  "main",
		outputAsTest: false,
	}
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
	nodes := make([]treeNode, len(lines))
	var counter int
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
		nodes[counter] = treeNode{indentLevel, node}
		counter++
	}
	nodes = nodes[0:counter]

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

// Start begins transpiling an input file.
func Start(args ProgramArgs) (err error) {
	if args.verbose {
		fmt.Println("Start tanspiling ...")
	}

	if os.Getenv("GOPATH") == "" {
		return fmt.Errorf("The $GOPATH must be set")
	}

	// 1. Compile it first (checking for errors)
	_, err = os.Stat(args.inputFile)
	if err != nil {
		return fmt.Errorf("Input file is not found")
	}

	// 2. Preprocess
	if args.verbose {
		fmt.Println("Running clang preprocessor...")
	}
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
		pp = out.Bytes()
	}

	if args.verbose {
		fmt.Println("Writing preprocessor ...")
	}
	dir, err := ioutil.TempDir("", "c2go")
	if err != nil {
		return fmt.Errorf("Cannot create temp folder: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	ppFilePath := path.Join(dir, "pp.c")
	err = ioutil.WriteFile(ppFilePath, pp, 0644)
	if err != nil {
		return fmt.Errorf("writing to %s failed: %v", ppFilePath, err)
	}

	// 3. Generate JSON from AST
	if args.verbose {
		fmt.Println("Running clang for AST tree...")
	}
	astPP, err := exec.Command("clang", "-Xclang", "-ast-dump", "-fsyntax-only", ppFilePath).Output()
	if err != nil {
		// If clang fails it still prints out the AST, so we have to run it
		// again to get the real error.
		errBody, _ := exec.Command("clang", ppFilePath).CombinedOutput()

		panic("clang failed: " + err.Error() + ":\n\n" + string(errBody))
	}

	if args.verbose {
		fmt.Println("Reading clang AST tree...")
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
	p.OutputAsTest = args.outputAsTest

	// Converting to nodes
	if args.verbose {
		fmt.Println("Converting to nodes...")
	}
	nodes := convertLinesToNodes(lines)

	// build tree
	if args.verbose {
		fmt.Println("Building tree...")
	}
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

	// transpile ast tree
	if args.verbose {
		fmt.Println("Transpiling tree...")
	}
	err = transpiler.TranspileAST(args.inputFile, args.packageName, p, tree[0].(ast.Node))
	if err != nil {
		return fmt.Errorf("cannot transpile AST : %v", err)
	}

	outputFilePath := args.outputFile

	if outputFilePath == "" {
		cleanFileName := filepath.Clean(filepath.Base(args.inputFile))
		extension := filepath.Ext(args.inputFile)
		outputFilePath = cleanFileName[0:len(cleanFileName)-len(extension)] + ".go"
	}

	// write the output Go code
	if args.verbose {
		fmt.Println("Writing the output Go code...")
	}
	err = ioutil.WriteFile(outputFilePath, []byte(p.String()), 0644)
	if err != nil {
		return fmt.Errorf("writing Go output file failed: %v", err)
	}

	return nil
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

	args := DefaultProgramArgs()

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
		args.verbose = *verboseFlag
	default:
		flag.Usage()
		return 1
	}

	if err := Start(args); err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}

	return 0
}

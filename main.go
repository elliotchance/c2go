// Package c2go contains the main function for running the executable.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/transpiler"
)

// Version can be requested through the command line with:
//
//     c2go -v
//
// See https://github.com/elliotchance/c2go/wiki/Release-Process
const Version = "0.13.4"

// ProgramArgs - arguments of program
type ProgramArgs struct {
	verbose     bool
	ast         bool
	inputFile   string
	outputFile  string
	packageName string
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

		indentAndType := regexp.MustCompile("^([|\\- `]*)(\\w+)").FindStringSubmatch(line)
		if len(indentAndType) == 0 {
			panic(fmt.Sprintf("Cannot understand line '%s'", line))
		}

		offset := len(indentAndType[1])
		node := ast.Parse(line[offset:])

		indentLevel := len(indentAndType[1]) / 2
		nodes = append(nodes, treeNode{indentLevel, node})
	}

	return nodes
}

// buildTree convert an array of nodes, each prefixed with a depth into a tree.
func buildTree(nodes []treeNode, depth int) []ast.Node {
	if len(nodes) == 0 {
		return []ast.Node{}
	}

	// Split the list into sections, treat each section as a a tree with its own root.
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

// ToJSON - tree convert to JSON
func ToJSON(tree []interface{}) []map[string]interface{} {
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

				value = ToJSON(v)
			}

			r[j][name] = value
		}
	}

	return r
}

// Start - base function
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

	ppFilePath := path.Join(os.TempDir(), "pp.c")
	err = ioutil.WriteFile(ppFilePath, pp, 0644)
	if err != nil {
		return fmt.Errorf("writing to /tmp/pp.c failed: %v", err)
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
	}

	nodes := convertLinesToNodes(lines)
	tree := buildTree(nodes, 0)

	p := program.NewProgram()
	p.Verbose = args.verbose

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

func main() {
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

	flag.Usage = func() {
		usage := "Usage: %s [-v] [<command>] [<flags>] file.c\n\n"
		usage += "Commands:\n"
		usage += "  transpile\ttranspile an input C source file to Go\n"
		usage += "  ast\t\tprint AST before translated Go code\n\n"

		usage += "Flags:\n"
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionFlag {
		// Simply print out the version and exit.
		fmt.Println(Version)
		return
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	args := ProgramArgs{verbose: *verboseFlag, ast: false}

	switch os.Args[1] {
	case "ast":
		err := astCommand.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("ast command cannot parse: %v", err)
			os.Exit(1)
		}

		if *astHelpFlag || astCommand.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "Usage: %s ast file.c\n", os.Args[0])
			astCommand.PrintDefaults()
			os.Exit(1)
		}

		args.ast = true
		args.inputFile = astCommand.Arg(0)
	case "transpile":
		err := transpileCommand.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("transpile command cannot parse: %v", err)
			os.Exit(1)
		}

		if *transpileHelpFlag || transpileCommand.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "Usage: %s transpile [-V] [-o file.go] [-p package] file.c\n", os.Args[0])
			transpileCommand.PrintDefaults()
			os.Exit(1)
		}

		args.inputFile = transpileCommand.Arg(0)
		args.outputFile = *outputFlag
		args.packageName = *packageFlag
	default:
		flag.Usage()
		os.Exit(1)
	}

	if err := Start(args); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

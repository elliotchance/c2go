// Package c2go contains the main function for running the executable.
package main

import (
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
const Version = "0.12.4"

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

func Check(prefix string, e error) {
	if e != nil {
		panic(prefix + e.Error())
	}
}

func Start(args ProgramArgs) {
	if os.Getenv("GOPATH") == "" {
		panic("The $GOPATH must be set.")
	}

	// 1. Compile it first (checking for errors)
	_, err := os.Stat(args.inputFile)
	Check("", err)

	// 2. Preprocess
	pp, err := exec.Command("clang", "-E", args.inputFile).Output()
	Check("preprocess failed: ", err)

	pp_file_path := path.Join(os.TempDir(), "pp.c")
	err = ioutil.WriteFile(pp_file_path, pp, 0644)
	Check("writing to /tmp/pp.c failed: ", err)

	// 3. Generate JSON from AST
	ast_pp, err := exec.Command("clang", "-Xclang", "-ast-dump", "-fsyntax-only", pp_file_path).Output()
	if err != nil {
		// If clang fails it still prints out the AST, so we have to run it
		// again to get the real error.
		errBody, _ := exec.Command("clang", pp_file_path).CombinedOutput()

		panic("clang failed: " + err.Error() + ":\n\n" + string(errBody))
	}

	lines := readAST(ast_pp)
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
	Check("writing C output file failed: ", err)
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
		astCommand.Parse(os.Args[2:])

		if *astHelpFlag || astCommand.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "Usage: %s ast file.c\n", os.Args[0])
			astCommand.PrintDefaults()
			os.Exit(1)
		}

		args.ast = true
		args.inputFile = astCommand.Arg(0)

		Start(args)
	case "transpile":
		transpileCommand.Parse(os.Args[2:])

		if *transpileHelpFlag || transpileCommand.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "Usage: %s transpile [-V] [-o file.go] [-p package] file.c\n", os.Args[0])
			transpileCommand.PrintDefaults()
			os.Exit(1)
		}

		args.inputFile = transpileCommand.Arg(0)
		args.outputFile = *outputFlag
		args.packageName = *packageFlag

		Start(args)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

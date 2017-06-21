// Package c2go contains the main function for running the executable.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Konstantin8105/c2go/analyze"
)

// Version can be requested through the command line with:
//
//     c2go -v
//
// See https://github.com/elliotchance/c2go/wiki/Release-Process
const Version = "0.13.3"

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

	args := analyze.ProgramArgs{Verbose: *verboseFlag, Ast: false}

	switch os.Args[1] {
	case "ast":
		err := astCommand.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("Ast command cannot parse: %v", err)
			os.Exit(1)
		}

		if *astHelpFlag || astCommand.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "Usage: %s ast file.c\n", os.Args[0])
			astCommand.PrintDefaults()
			os.Exit(1)
		}

		args.Ast = true
		args.InputFile = astCommand.Arg(0)

		if err = analyze.Start(args); err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	case "transpile":
		err := transpileCommand.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("Transpile command cannot parse: %v", err)
			os.Exit(1)
		}

		if *transpileHelpFlag || transpileCommand.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "Usage: %s transpile [-V] [-o file.go] [-p package] file.c\n", os.Args[0])
			transpileCommand.PrintDefaults()
			os.Exit(1)
		}

		args.InputFile = transpileCommand.Arg(0)
		args.OutputFile = *outputFlag
		args.PackageName = *packageFlag

		if err = analyze.Start(args); err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

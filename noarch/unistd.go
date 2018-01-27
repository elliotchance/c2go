package noarch

import (
	"strings"
)

// See documentation:
// https://www.gnu.org/software/libc/manual/html_node/Using-Getopt.html#Using-Getopt

// Optarg - is implementation of global variable "optarg" from "unistd.h"
// This variable is set by getopt to point at the value of the option argument,
// for those options that accept arguments.
var Optarg []byte /* argument associated with option */

// Opterr - is implementation of global variable "opterr" from "unistd.h"
// If the value of this variable is nonzero, then getopt prints an error
// message to the standard error stream if it encounters an unknown option
// character or an option with a missing required argument. This is the default
// behavior. If you set this variable to zero, getopt does not print any
// messages, but it still returns the character ? to indicate an error.
var Opterr int = 1 /* if error message should be printed */

// Optind - is implementation of global variable "optind" from "unistd.h"
// This variable is set by getopt to the index of the next element of the argv
// array to be processed. Once getopt has found all of the option arguments,
// you can use this variable to determine where the remaining non-option
// arguments begin. The initial value of this variable is 1.
var Optind int = 1 /* index into parent argv vector */

// Optopt - is implementation of global variable "optopt" from "unistd.h"
// When getopt encounters an unknown option character or an option with
// a missing required argument, it stores that option character in this
// variable. You can use this for providing your own diagnostic messages.
var Optopt int /* character checked for validity */

var sp = 1

const EOF = -1

// Getopt - is implementation of function "getopt" from "unistd.h"
// The getopt function gets the next option argument from the argument list
// specified by the argv and argc arguments. Normally these values come
// directly from the arguments received by main.
func Getopt(argc int, argv [][]byte, options []byte) int {
	// Source prototype see :
	// http://web.mit.edu/ghudson/dev/third/krb5/src/util/windows/getopt.c
	var c int
	var cp int

	if sp == 1 {
		if Optind >= argc ||
			(len(argv[0]) > 0 && argv[Optind][0] != '-') ||
			len(argv[0]) == 1 {
			return EOF
		} else if string(argv[Optind]) == "--" {
			Optind++
			return EOF
		}
	}
	c = int(argv[Optind][sp])
	Optopt = c
	cp = strings.Index(string(options), string(c))
	if c == ':' || cp == -1 {
		if Opterr != 0 {
			println(": illegal option --", string(c))
		}
		sp++
		if len(argv[Optind]) == sp {
			Optind++
			sp = 1
		}
		return '?'
	}
	cp++
	if cp < len(options) && options[cp] == ':' {
		if len(argv[Optind]) > sp+1 {
			Optarg = argv[Optind][sp+1 : sp+2]
			Optind++
		} else {
			Optind++
			if Optind >= argc {
				if Opterr != 0 {
					println(": option requires an argument --", string(c))
				}
				sp = 1
				return '?'
			} else {
				Optarg = argv[Optind]
				Optind++
			}
		}
		sp = 1
	} else {
		sp++
		if len(argv[Optind]) == sp {
			sp = 1
			Optind++
		}
		Optarg = []byte("")
	}
	return int(c)
}

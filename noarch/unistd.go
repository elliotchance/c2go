package noarch

import (
	"fmt"
	"os"
	"unicode"
	"unsafe"
)

// See documentation:
// https://www.gnu.org/software/libc/manual/html_node/Using-Getopt.html#Using-Getopt

// Optarg - is implementation of global variable "optarg" from "unistd.h"
// This variable is set by getopt to point at the value of the option argument,
// for those options that accept arguments.
var Optarg []byte

// Opterr - is implementation of global variable "opterr" from "unistd.h"
// If the value of this variable is nonzero, then getopt prints an error
// message to the standard error stream if it encounters an unknown option
// character or an option with a missing required argument. This is the default
// behavior. If you set this variable to zero, getopt does not print any
// messages, but it still returns the character ? to indicate an error.
var Opterr int

// Optind - is implementation of global variable "optind" from "unistd.h"
// This variable is set by getopt to the index of the next element of the argv
// array to be processed. Once getopt has found all of the option arguments,
// you can use this variable to determine where the remaining non-option
// arguments begin. The initial value of this variable is 1.
var Optind int

// Optopt - is implementation of global variable "optopt" from "unistd.h"
// When getopt encounters an unknown option character or an option with
// a missing required argument, it stores that option character in this
// variable. You can use this for providing your own diagnostic messages.
var Optopt int

// Getopt - is implementation of function "getopt" from "unistd.h"
// The getopt function gets the next option argument from the argument list
// specified by the argv and argc arguments. Normally these values come
// directly from the arguments received by main.
// func Getopt(argc int, argv [][]byte, options []byte) int {
// 	defer func() {
// 		Optind++
// 	}()
// 	if argc <= Optind {
// 		return -1
// 	}
// 	// TODO: create correct algorithm
// 	if argv != nil {
// 		if len(argv) > 0 {
// 			if len(argv[Optind]) > 0 {
// 				if argv[Optind][0] == '-' {
// 					return int(argv[Optind][1])
// 				}
// 			}
// 		}
// 	}
//
// 	return -1
// }

func Getopt(argc int, argv [][]byte, optstring []byte) int {
	var Optpos int = 1
	var arg []byte
	_ = argc
	if Optind == 0 {
		Optind = 1
		Optpos = 1
	}
	if len(argv) <= Optind {
		return -1
	}
	arg = argv[Optind]
	if arg != nil && Strcmp(arg, []byte("--\x00")) == 0 {
		Optind += 1
		return -1
	} else {
		if arg == nil || int(arg[0]) != int('-') || !unicode.IsNumber(rune(arg[1])) /*isalnum(arg[1]) */ {
			return -1
		} else {
			var opt []byte = Strchr(optstring, int(arg[Optpos]))
			Optopt = int(arg[Optpos])
			if opt == nil {
				if Opterr != 0 && int(optstring[0]) != int(':') {
					fmt.Fprintf(os.Stderr, "%s: illegal option: %c\n", argv[0], Optopt)
				}
				return int('?')
			} else {
				if int(opt[1]) == int(':') {
					if arg[Optpos+1] != 0 {
						Optarg = (*(*[1000000]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&(*(*[1]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&arg[0])) + (uintptr)(Optpos)*unsafe.Sizeof(arg[0]))))[:][0])) + (uintptr)(1)*unsafe.Sizeof((*(*[100000]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&arg[0])) + (uintptr)(Optpos)*unsafe.Sizeof(arg[0]))))[:][0]))))[:]
						Optind += 1
						Optpos = 1
						return Optopt
					} else {
						if !CStringIsNull(argv[Optind+1]) {
							Optarg = argv[Optind+1]
							Optind += 2
							Optpos = 1
							return Optopt
						} else {
							if Opterr != 0 && int(optstring[0]) != int(':') {
								fmt.Fprintf(os.Stderr, "%s: option requires an argument: %c\n", argv[0], Optopt)
							}
							return func() int {
								if int(optstring[0]) == int(':') {
									return int(':')
								} else {
									return int('?')
								}
							}()
						}
					}
				} else {
					Optpos += 1
					if arg[Optpos] != '\x00' {
						Optind += 1
						Optpos = 1
					}
					return Optopt
				}
			}
		}
	}
	return -1
}

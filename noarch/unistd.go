package noarch

// Optarg - is implementation of global variable "optarg" from "unistd.h"
var Optarg []byte

// Opterr - is implementation of global variable "opterr" from "unistd.h"
var Opterr int

// Optind - is implementation of global variable "optind" from "unistd.h"
var Optind int

// Getopt - is implementation of function "getopt" from "unistd.h"
func Getopt(argc int, argv [][]byte, mask []byte) int {
	return -1
}

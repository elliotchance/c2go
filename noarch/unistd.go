package noarch

// Optarg - is implementation of global variable "optarg" from "unistd.h"
var Optarg []byte

// Getopt - is implementation of function "getopt" from "unistd.h"
func Getopt(argc int, argv [][]byte, mask []byte) int {
	return -1
}

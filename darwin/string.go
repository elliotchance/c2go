package darwin

import "github.com/elliotchance/c2go/noarch"

func BuiltinStrcpy(dest, src []byte, size int) []byte {
	return noarch.Strcpy(dest, src)
}

func BuiltinObjectSize(ptr []byte, theType int) int {
	return 5
}

package darwin

import "github.com/elliotchance/c2go/noarch"

func BuiltinSprintfChk(buffer []byte, _ int, _ int, format []byte, args ...interface{}) int {
	return noarch.Sprintf(buffer, format, args)
}

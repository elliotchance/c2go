/*
	Package main - transpiled by c2go version: v0.21.10 Zinc 2018-02-14

	If you have found any issues, please raise an issue at:
	https://github.com/elliotchance/c2go/
*/

package code_quality

// switch_function - transpiled function from  /home/konstantin/go/src/github.com/elliotchance/c2go/tests/code_quality/switch.c:1
func switch_function() {
	var i int = 34
	switch i {
	case (0):
		fallthrough
	case (2):
		{
			_ = (i)
			return
		}
	case 3:
		{
			var c int
			return
		}
	case 4:
	case 5:
	case 6:
		fallthrough
	case 7:
		var d int
	case (1):
		{
			return
		}
	}
}
func init() {
}

type _Bool int

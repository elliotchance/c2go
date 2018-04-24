/*
	Package main - transpiled by c2go version: v0.22.4 Aluminium 2018-04-24

	If you have found any issues, please raise an issue at:
	https://github.com/elliotchance/c2go/
*/

package code_quality

// switch_function - transpiled function from  tests/code_quality/switch.c:1
func switch_function() {
	var i int32 = 34
	switch i {
	case (0):
		fallthrough
	case (1):
		{
			return
		}
	case (2):
		{
			_ = (i)
			return
		}
	case 3:
		{
			var c int32
			return
		}
	case 4:
	case 5:
	case 6:
		fallthrough
	case 7:
		{
			var d int32
			break
		}
	}
}
func init() {
}

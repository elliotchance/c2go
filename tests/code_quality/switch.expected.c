/*
	Package main - transpiled by c2go version: v0.23.0 Berkelium 2018-04-27

	If you have found any issues, please raise an issue at:
	https://github.com/elliotchance/c2go/
*/

package code_quality

// switch_function - transpiled function from  tests/code_quality/switch.c:1
func switch_function() {
	var i int32 = int32(34)
	switch i {
	case (int32(0)):
		fallthrough
	case (int32(1)):
		{
			return
		}
	case (int32(2)):
		{
			_ = (i)
			return
		}
	case int32(3):
		{
			var c int32
			return
		}
	case int32(4):
		{
		}
	case int32(5):
		{
		}
	case int32(6):
		fallthrough
	case int32(7):
		{
			var d int32
			break
		}
	}
}
func init() {
}

/*
	Package main - transpiled by c2go version: v0.22.4 Aluminium 2018-04-24

	If you have found any issues, please raise an issue at:
	https://github.com/elliotchance/c2go/
*/

package code_quality

// if_1 - transpiled function from  tests/code_quality/if.c:1
func if_1() {
	var a int32 = int32(5)
	var b int32 = int32(2)
	var c int32 = int32(4)
	if a > b {
		return
	} else if c <= a {
		a = int32(0)
	}
	_ = (a)
	_ = (b)
	_ = (c)
}
func init() {
}

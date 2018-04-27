/*
	Package main - transpiled by c2go version: v0.22.4 Aluminium 2018-04-24

	If you have found any issues, please raise an issue at:
	https://github.com/elliotchance/c2go/
*/

package code_quality

// f1 - transpiled function from  tests/code_quality/for.c:1
func f1() {
	var i int32
	for i = int32(0); i < int32(10); i++ {
	}
}

// f2 - transpiled function from  tests/code_quality/for.c:7
func f2() {
	var i int32
	for i = int32(10); i > int32(0); i-- {
	}
}

// f3 - transpiled function from  tests/code_quality/for.c:13
func f3() {
	{
		var i int32 = int32(0)
		for ; i < int32(10); i++ {
		}
	}
}
func init() {
}

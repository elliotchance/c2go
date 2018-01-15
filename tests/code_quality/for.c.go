/* Package main - transpiled by c2go

If you have found any issues, please raise an issue at:
https://github.com/elliotchance/c2go/
*/

package main

// f1 - transpiled function from file : tests/code_quality/for.c , line : 1
func f1() {
	var i int
	for i = 0; i < 10; i++ {
	}
}

// f2 - transpiled function from file : tests/code_quality/for.c , line : 6
func f2() {
	var i int
	for i = 10; i > 0; func() int {
		i -= 1
		return i
	}() {
	}
}

// f3 - transpiled function from file : tests/code_quality/for.c , line : 11
func f3() {
	{
		var i int = 0
		for ; i < 10; i++ {
		}
	}
}

// main - transpiled function from file : tests/code_quality/for.c , line : 16
func main() {
	f1()
	f2()
	f3()
	return
}
func init() {
}

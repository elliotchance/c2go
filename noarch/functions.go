package noarch

func BoolToInt(x bool) int {
	if x {
		return 1
	}

	return 0
}

func NotInt(x int) int {
	if x == 0 {
		return 1
	}

	return 0
}

func NotUint16(x uint16) uint16 {
	if x == 0 {
		return 1
	}

	return 0
}

func Ternary(a bool, b, c func() interface{}) interface{} {
	if a {
		return b()
	}

	return c()
}

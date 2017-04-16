package noarch

// FIXME
type __builtin_va_list int64

func BoolToInt(x bool) int {
	if x {
		return 1
	}

	return 0
}

func __bool_to_uint32(x bool) int {
	if x {
		return 1
	}

	return 0
}

func __not_uint32(x uint32) uint32 {
	if x == 0 {
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

func Ternary(a bool, b, c func() interface{}) interface{} {
	if a {
		return b()
	}

	return c()
}

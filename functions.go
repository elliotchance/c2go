package main

func __bool_to_int(x bool) int {
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

func __not_int(x int) int {
    if x == 0 {
        return 1
    }

    return 0
}

func __ternary(a bool, b interface{}, c interface{}) interface{} {
    if a {
        return b
    }

    return c
}

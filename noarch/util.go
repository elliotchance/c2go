package noarch

import (
	"reflect"
)

// NullTerminatedByteSlice returns a string that contains all the bytes in the
// provided C string up until the first NULL character.
func NullTerminatedByteSlice(s []byte) string {
	if s == nil {
		return ""
	}

	end := -1
	for i, b := range s {
		if b == 0 {
			end = i
			break
		}
	}

	if end == -1 {
		end = len(s)
	}

	newSlice := make([]byte, end)
	copy(newSlice, s)

	return string(newSlice)
}

// CStringIsNull will test if a C string is NULL. This is equivalent to:
//
//    s == NULL
func CStringIsNull(s []byte) bool {
	if s == nil || len(s) < 1 {
		return true
	}

	return s[0] == 0
}

func CPointerToGoPointer(a interface{}) interface{} {
	t := reflect.TypeOf(a).Elem()

	return reflect.New(t).Elem().Addr().Interface()
}

func GoPointerToCPointer(destination interface{}, value interface{}) {
	v := reflect.ValueOf(destination).Elem()
	reflect.ValueOf(value).Index(0).Set(v)
}

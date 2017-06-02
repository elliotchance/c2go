package transpiler

import (
    "fmt"
    "reflect"
    "unsafe"
)

// Common type for all union structures
type _UnionType [4]byte

// Get casted value
func (self *_UnionType) cast(t reflect.Type) reflect.Value {
    return reflect.NewAt(t, unsafe.Pointer(&self[0])).Elem()
}

// Assign
func (self *_UnionType) assign(v interface{}) {
    value := reflect.ValueOf(v).Elem()

    value.Set(self.cast(value.Type()))
}

// Setter
func (self *_UnionType) Set(v interface{}) {
    value := reflect.ValueOf(v)

    self.cast(value.Type()).Set(value)
}

// F1 Getter
func (self *_UnionType) F1() int32 {
    var res int32

    self.assign(&res)

    return res
}

// F2 Getter
func (self *_UnionType) F2() uint32 {
    var res uint32

    self.assign(&res)

    return res
}

// F3 Getter
func (self *_UnionType) F3() byte {
    var res byte

    self.assign(&res)

    return res
}

// F4 Getter
func (self *_UnionType) F4() int16 {
    var res int16

    self.assign(&res)

    return res
}

func fake_main() {
    // Create the union
    var u _UnionType

    // Set a value
    u.Set(0x12345678)

    // Get values
    f1 := u.F1()
    f2 := u.F2()
    f3 := u.F3()
    f4 := u.F4()

    // Print the results
    fmt.Printf("%x %x %x %x\n", f1, f2, f3, f4)
}

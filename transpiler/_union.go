package main

import(
    "fmt"
    "unsafe"
    "reflect"
)

// Common type for all union structures
type UnionType []byte

// Get casted value
func (self UnionType) cast(t reflect.Type) reflect.Value {
    return reflect.NewAt(t, unsafe.Pointer(&self[0])).Elem()
}

// Getter
func (self UnionType) Get(v interface{}) {
    value := reflect.ValueOf(v).Elem()

    value.Set(self.cast(value.Type()))
}

// Setter
func (self UnionType) Set(v interface{}) {
    value := reflect.ValueOf(v)

    self.cast(value.Type()).Set(value)
}

func main(){
    // Create the union
    u := make(UnionType, 4)

    // Set a value
    u.Set(0x12345678)

    var f1 int32
    var f2 uint32
    var f3 byte
    var f4 int16

    // Get values
    u.Get(&f1)
    u.Get(&f2)
    u.Get(&f3)
    u.Get(&f4)

    // Print the results
    fmt.Printf("%x %x %x %x\n", f1, f2, f3, f4)
}

package main

import (
    "fmt"
    "reflect"
    "unsafe"
)

// Union type
type SampleType [4]byte

// Get casted pointer
func (self *SampleType) cast(t reflect.Type) reflect.Value {
    return reflect.NewAt(t, unsafe.Pointer(&self[0]))
}

// Assign value from an union field (used by getters)
func (self *SampleType) assign(v interface{}) {
    value := reflect.ValueOf(v).Elem()

    value.Set(self.cast(value.Type()).Elem())
}

// Get typed pointer
func (self *SampleType) pointer(v interface{}) {
    value := reflect.ValueOf(v).Elem()

    value.Set(self.cast(value.Type().Elem()))
}

// Generic setter
func (self *SampleType) UntypedSet(v interface{}) {
    value := reflect.ValueOf(v)

    self.cast(value.Type()).Elem().Set(value)
}

// Setters
func (self *SampleType) SetF1(v int32)  { self.UntypedSet(v) }
func (self *SampleType) SetF2(v uint32) { self.UntypedSet(v) }
func (self *SampleType) SetF3(v byte)   { self.UntypedSet(v) }
func (self *SampleType) SetF4(v int16)  { self.UntypedSet(v) }

// Getters
func (self *SampleType) GetF1() (res int32)  { self.assign(&res); return }
func (self *SampleType) GetF2() (res uint32) { self.assign(&res); return }
func (self *SampleType) GetF3() (res byte)   { self.assign(&res); return }
func (self *SampleType) GetF4() (res int16)  { self.assign(&res); return }

// Pointers
func (self *SampleType) PtrF1() (res *int32)  { self.pointer(&res); return }
func (self *SampleType) PtrF2() (res *uint32) { self.pointer(&res); return }
func (self *SampleType) PtrF3() (res *byte)   { self.pointer(&res); return }
func (self *SampleType) PtrF4() (res *int16)  { self.pointer(&res); return }

func main() {
    // Create the union
    var u SampleType

    // Set a value
    u.UntypedSet(0x12345678)

    // Get values
    f1 := u.GetF1()
    f2 := u.GetF2()
    f3 := u.GetF3()
    f4 := u.GetF4()

    // Print the results
    fmt.Printf("Values:              %x %x %x %x\n", f1, f2, f3, f4)

    // Get pointers
    p1 := u.PtrF1()
    p2 := u.PtrF2()
    p3 := u.PtrF3()
    p4 := u.PtrF4()

    // Print values before modification
    fmt.Printf("Before modification: %x %x %x %x\n", *p1, *p2, *p3, *p4)

    // modification
    *p2 = 0x12344321

    // Print values after modification
    fmt.Printf("After modification:  %x %x %x %x\n", *p1, *p2, *p3, *p4)
}

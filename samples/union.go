package main

import (
    "fmt"
    "reflect"
    "unsafe"
)

/*
SampleType is a Go union type example for the C union type:
    union SampleType
    {
        long f1;
        unsigned long f2;
        unsigned char f3;
        short f4;
    };
*/
type SampleType [4]byte

// Get casted pointer
func (st *SampleType) cast(t reflect.Type) reflect.Value {
    return reflect.NewAt(t, unsafe.Pointer(&st[0]))
}

// Assign value from an union field (used by getters)
func (st *SampleType) assign(v interface{}) {
    value := reflect.ValueOf(v).Elem()

    value.Set(st.cast(value.Type()).Elem())
}

// Get typed pointer
func (st *SampleType) pointer(v interface{}) {
    value := reflect.ValueOf(v).Elem()

    value.Set(st.cast(value.Type().Elem()))
}

// UntypedSet is the generic setter
func (st *SampleType) UntypedSet(v interface{}) {
    value := reflect.ValueOf(v)

    st.cast(value.Type()).Elem().Set(value)
}

/* -- Pointers -- */

// PtrF1 gets pointer on F1 field
func (st *SampleType) PtrF1() (res *int32) { st.pointer(&res); return }

// PtrF2 gets pointer on F2 field
func (st *SampleType) PtrF2() (res *uint32) { st.pointer(&res); return }

// PtrF3 gets pointer on F3 field
func (st *SampleType) PtrF3() (res *byte) { st.pointer(&res); return }

// PtrF4 gets pointer on F4 field
func (st *SampleType) PtrF4() (res *int16) { st.pointer(&res); return }

/* -- Setters -- */

// SetF1 sets F1 field
func (st *SampleType) SetF1(v int32) { st.UntypedSet(v) }

// SetF2 sets F2 field
func (st *SampleType) SetF2(v uint32) { st.UntypedSet(v) }

// SetF3 sets F3 field
func (st *SampleType) SetF3(v byte) { st.UntypedSet(v) }

// SetF4 sets F4 field
func (st *SampleType) SetF4(v int16) { st.UntypedSet(v) }

/* -- Getters -- */

// GetF1 gets F1 field
func (st *SampleType) GetF1() (res int32) { st.assign(&res); return }

// GetF2 gets F2 field
func (st *SampleType) GetF2() (res uint32) { st.assign(&res); return }

// GetF3 gets F3 field
func (st *SampleType) GetF3() (res byte) { st.assign(&res); return }

// GetF4 gets F4 field
func (st *SampleType) GetF4() (res int16) { st.assign(&res); return }

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

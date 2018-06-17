package noarch

import (
	"reflect"
	"sync"
	"unsafe"
)

// CStringToString returns a string that contains all the bytes in the
// provided C string up until the first NULL character.
func CStringToString(s *byte) string {
	if s == nil {
		return ""
	}

	end := -1
	for i := 0; ; i++ {
		if *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + uintptr(i))) == 0 {
			end = i
			break
		}
	}

	if end == -1 {
		return ""
	}

	return string(toByteSlice(s, int32(end)))
}

// StringToCString returns the C string (also known as a null terminated string)
// to be as used as a string in C.
func StringToCString(s string) *byte {
	cString := make([]byte, len(s)+1)
	copy(cString, []byte(s))
	cString[len(s)] = 0

	return &cString[0]
}

// CStringIsNull will test if a C string is NULL. This is equivalent to:
//
//    s == NULL
func CStringIsNull(s *byte) bool {
	if s == nil {
		return true
	}

	return *s == 0
}

// CPointerToGoPointer converts a C-style pointer into a Go-style pointer.
//
// C pointers are represented as slices that have one element pointing to where
// the original C pointer would be referencing. This isn't useful if the pointed
// value needs to be passed to another Go function in these libraries.
//
// See also GoPointerToCPointer.
func CPointerToGoPointer(a interface{}) interface{} {
	t := reflect.TypeOf(a).Elem()

	return reflect.New(t).Elem().Addr().Interface()
}

// toByteSlice returns a byte slice to a with the given length.
func toByteSlice(a *byte, length int32) []byte {
	header := reflect.SliceHeader{
		uintptr(unsafe.Pointer(a)),
		int(length),
		int(length),
	}
	return (*(*[]byte)(unsafe.Pointer(&header)))[:]
}

// GoPointerToCPointer does the opposite of CPointerToGoPointer.
//
// A Go pointer (simply a pointer) is converted back into the original slice
// structure (of the original slice reference) so that the calling functions
// will be able to see the new data of that pointer.
func GoPointerToCPointer(destination interface{}, value interface{}) {
	v := reflect.ValueOf(destination).Elem()
	reflect.ValueOf(value).Index(0).Set(v)
}

// UnsafeSliceToSlice takes a slice and transforms it into a slice of a different type.
// For this we need to adjust the length and capacity in accordance with the sizes
// of the underlying types.
func UnsafeSliceToSlice(a interface{}, fromSize int32, toSize int32) *reflect.SliceHeader {
	v := reflect.ValueOf(a)

	// v might not be addressable, use this trick to get v2 = v,
	// with v2 being addressable
	v2 := reflect.New(v.Type()).Elem()
	v2.Set(v)

	// get a pointer to the SliceHeader
	// Calling Pointer() on the slice directly only gets a pointer to the 1st element, not the slice header,
	// which is why we first call Addr()
	ptr := unsafe.Pointer(v2.Addr().Pointer())

	// adjust header to adjust sizes for the new type
	header := *(*reflect.SliceHeader)(ptr)
	header.Len = (header.Len * int(fromSize)) / int(toSize)
	header.Cap = (header.Cap * int(fromSize)) / int(toSize)
	return &header
}

// UnsafeSliceToSliceUnlimited takes a slice and transforms it into a slice of a different type.
// The length and capacity will be set to unlimited.
func UnsafeSliceToSliceUnlimited(a interface{}) *reflect.SliceHeader {
	v := reflect.ValueOf(a)

	// v might not be addressable, use this trick to get v2 = v,
	// with v2 being addressable
	v2 := reflect.New(v.Type()).Elem()
	v2.Set(v)

	// get a pointer to the SliceHeader
	// Calling Pointer() on the slice directly only gets a pointer to the 1st element, not the slice header,
	// which is why we first call Addr()
	ptr := unsafe.Pointer(v2.Addr().Pointer())

	// adjust header to adjust sizes for the new type
	header := *(*reflect.SliceHeader)(ptr)
	header.Len = 1000000000
	header.Cap = 1000000000
	return &header
}

// Safe contains a thread-safe value
type Safe struct {
	value interface{}
	lock  sync.RWMutex
}

// NewSafe create a new Safe instance given a value
func NewSafe(value interface{}) *Safe {
	return &Safe{value: value, lock: sync.RWMutex{}}
}

// Get returns the value
func (s *Safe) Get() interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.value
}

// Set sets a new value
func (s *Safe) Set(value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.value = value
}

var (
	interfaceMgmt map[reflect.Value]*InterfaceWrapper
	interfaceSync sync.Mutex
)

func init() {
	interfaceMgmt = make(map[reflect.Value]*InterfaceWrapper)
}

// InterfaceWrapper is used for those case where we cannot use unsafe.Pointer
// the default catch all data type for c2go.
// For the moment this is the case for function pointers.
type InterfaceWrapper struct {
	X interface{}
}

// CastInterfaceToPointer will take an interface and store it in a map by its reflect.Value
// the unsafe.Pointer to its containing InterfaceWrapper is returned.
// Since no element is ever removed from the map this should only be used for a limited amount
// of elements, like function pointers.
func CastInterfaceToPointer(in interface{}) unsafe.Pointer {
	interfaceSync.Lock()
	defer interfaceSync.Unlock()
	val := reflect.ValueOf(in)
	if iw, ok := interfaceMgmt[val]; ok {
		return unsafe.Pointer(iw)
	} else {
		iw := &InterfaceWrapper{
			X: in,
		}
		interfaceMgmt[val] = iw
		return unsafe.Pointer(iw)
	}
}

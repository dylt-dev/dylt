package eco

import (
	"fmt"
	"reflect"

	"github.com/dylt-dev/dylt/common"
)

type Flavor reflect.Kind
const (
	InvalidFlavor Flavor = iota
	Map
	Pointer
	Scalar
	Slice
	Struct
)

func NewFlavor (knd reflect.Kind) Flavor {
	switch knd {
	// simple case for simple types
	case reflect.Bool,
	     reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	     reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	     reflect.Float32,
	     reflect.Float64,
	     reflect.String:
		return Scalar

	case reflect.Pointer:
		return Pointer

	case reflect.Slice:
		return Slice

	case reflect.Map:
		return Map

	case reflect.Struct:
		return Struct

	default:
		return InvalidFlavor
	}
}



// Check if the argument refers to a pointer to an allocated variable, or
// if it is a pointer to a nil pointer, or a pointer to a Slice or Map
// If the argument is none of these, return an error
func IsPointerAllocated (ctx *common.EcoContext, a any) (bool, error) {
	ctx.Logger.Signature("IsPointerAllocated")
	ctx.Inc()
	defer ctx.Dec()

	rv := reflect.ValueOf(a)
	if !rv.IsValid() { return false, fmt.Errorf("Invalid value") }
	if rv.Kind() != reflect.Pointer { return false, fmt.Errorf("expected pointer, got %s", rv.Kind().String()) }
	if rv.IsNil() {
		if RvPointer(rv).IsReference(ctx) {
			return false, nil
		} else {
			return false, fmt.Errorf("Expecting non-nil value")
		}
	}

	rvElem := rv.Elem()
	if rvElem.Kind() == reflect.Pointer {
		if !rvElem.IsNil() {
			return false, fmt.Errorf("pointer to a non-nil pointer")
		} else {
			return false, nil
		}
	}

	return true, nil
}


func MakeSlice (typ reflect.Type, n int) any {
	rvSlice := reflect.MakeSlice(typ, n, n)
	// Slices are not addressable. We need to allocate a pointer
	// to a slice, assign our new slice to the pointer's Elem,
	// and then return the Interface() of the new pointer
	rvSlicePtr := reflect.New(typ)
	rvSlicePtr.Elem().Set(rvSlice)
	return rvSlicePtr.Interface()
}
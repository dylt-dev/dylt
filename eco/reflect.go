package eco

import (
	"fmt"
	"reflect"

	"github.com/dylt-dev/dylt/common"
)

// Check if the argument refers to a pointer to an allocated variable, or
// if it is a pointer to a nil pointer, or a pointer to a Slice or Map
// If the argument is none of these, return an error
func IsPointerAllocated(ctx *common.EcoContext, a any) (bool, error) {
	ctx.Logger.Signature("IsPointerAllocated")
	ctx.Inc()
	defer ctx.Dec()

	rv := reflect.ValueOf(a)
	if !rv.IsValid() {
		return false, fmt.Errorf("Invalid value")
	}
	if rv.Kind() != reflect.Pointer {
		return false, fmt.Errorf("expected pointer, got %s", rv.Kind().String())
	}
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

func NewMap(typ reflect.Type, n int) any {
	rv := reflect.MakeMapWithSize(typ, n)
	// Slices are not addressable. We need to allocate a pointer
	// to a slice, assign our new slice to the pointer's Elem,
	// and then return the Interface() of the new pointer
	rvPtr := reflect.New(typ)
	rvPtr.Elem().Set(rv)
	return rvPtr.Interface()
}

func NewSlice(typ reflect.Type, n int) any {
	ev := reflect.MakeSlice(typ, n, n)
	// Slices are not addressable. We need to allocate a pointer
	// to a slice, assign our new slice to the pointer's Elem,
	// and then return the Interface() of the new pointer
	rvPtr := reflect.New(typ)
	rvPtr.Elem().Set(ev)
	return rvPtr.Interface()
}

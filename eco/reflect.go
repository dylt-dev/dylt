package eco

import (
	"fmt"
	"reflect"
)

type RvPointer reflect.Value


func (this RvPointer) CreateOrGet (any, error) {
	
}

func (this RvPointer) Walk() (any, error) {
	rv := reflect.Value(this)
	if !rv.IsValid() { return nil, fmt.Errorf("Invalid value") }
	if rv.IsNil() { return nil, fmt.Errorf("Expecting non-nilt value") }
	if rv.Kind() != reflect.Pointer { return nil, fmt.Errorf("expected pointer, got %s", rv.Kind().String()) }

	// We definitely have a non-nil pointer
	// Now, we care about 3 success conditions and one error condition
	// - rv is a pointer to a Reference type (slice or map) - return
	// - rv is a non-nil pointer to a Value type (scalar struct) - return
	// - rv is a pointer to a non-pointer - return Error
	// - rv is a pointer to a nil value pointer - return
	
	for {
		rvElem := rv.Elem()
		
		// check rv is a pointer to a Reference type (slice or map)
		if rvElem.Kind() == reflect.Map || rvElem.Kind() == reflect.Slice {
			return rv.Interface(), nil
		}

		// check rv is a non-nil pointer to a Value type (scalar struct) - return
		// - rv is a pointer to a nil value pointer - return
		// check rv is a pointer to a non-pointer - return Error
		switch rvElem.Kind() {
			case reflect.Bool,
				 reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				 reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
				 reflect.Float32,
				 reflect.Float64,
				 reflect.String,
				 reflect.Struct: {
					return rv.Interface(), nil
			}
			case reflect.Pointer: {
				if rvElem.IsNil() {
					return rv.Interface(), nil
				}			
			}

			default: {
				return nil, fmt.Errorf("pointer to unsupported kind (%s)", rvElem.Kind().String())
			}
		}

		// If we've gotten this far, we have a pointer to a non-nil pointer, so we keep looking
		rv = rvElem
	}
}
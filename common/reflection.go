package common

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type TestMap map[string]string
type TestStruct struct {
	Name        string  `eco:"name"`
	LuckyNumber float64 `eco:"lucky_number"`
	NoTag       string
}

// Create a new map or return a reference to an allocated map
// Maps are nillable. So the pointer might be nil, or the pointer
// might point to a map that's nil. In either case the map needs
// to be allocated.
//
// rv    A pointer to an allocated map, or a pointer to a map pointer
//
//	If pointer-to-map, return the pointer
//	If pointer-to-pointer, allocate a new map and return the address
//
// returns	A pointer to the new map or the existing map
//
// @note battle harden and test vs other invalid data, and possibly return error
//
//	instead of bool
func CreateOrGetMap(ctx *EcoContext, rv reflect.Value) (any, bool) {
	// rv is either a pointer to a map, or a pointer-to-pointer to a map
	// if it's a non-nil pointer to a map, we don't have to allocate it
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return nil, false
	}
	var i any
	// If rv is a pointer-to-a-pointer, allocate a new map
	if rv.Elem().Type().Kind() == reflect.Pointer {
		ctx.Logger.Comment("Allocating a new map")
		typMap, err := GetUnderlyingMapType(ctx, rv)
		if err != nil {
			return nil, false
		}
		// This creates a map, but we need to figure out a way to
		rvMap := reflect.MakeMap(typMap)
		rvMapP := reflect.New(typMap)
		rvMapP.Elem().Set(rvMap)
		rv.Elem().Set(rvMapP)
		i = rvMapP.Interface()
	} else if rv.Elem().IsNil() {
		ctx.Logger.Comment("rv points to a nil map - allocating a new map")
		typMap, err := GetUnderlyingMapType(ctx, rv)
		if err != nil {
			return nil, false
		}
		rvMap := reflect.MakeMap(typMap)
		rv.Elem().Set(rvMap)
		i = rv.Interface()
	} else {
		i = rv.Elem().Addr().Interface()
	}

	return i, true
}

// Create a new struct or return a reference to an allocated struct
//
// rv    A pointer to an allocated struct, or a pointer to a struct pointer
//
//	If pointer-to-struct, return the pointer
//	If pointer-to-pointer, allocate a new struct and return the address
//
// returns	A pointer to the new struct or the existing struct
// @note battle harden and test vs other invalid data, and possibly return error
//
//	instead of bool
func CreateOrGetStruct(ctx *EcoContext, rv reflect.Value) (any, bool) {
	// rv is either a pointer to a struct, or a pointer-to-pointer to a struct
	// if it's a non-nil pointer to a struct, we don't have to allocate it
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return nil, false
	}
	var i any
	// If rv is a pointer-to-a-pointer, allocate a new struct
	if rv.Elem().Type().Kind() == reflect.Pointer {
		typStruct, err := GetUnderlyingStructType(ctx, rv)
		if err != nil {
			return nil, false
		}
		rvStruct := reflect.New(typStruct)
		rv.Elem().Set(rvStruct)
		i = rvStruct.Interface()
	} else {
		i = rv.Elem().Addr().Interface()
	}

	return i, true
}

func GetUnderlyingMapType(ctx *EcoContext, i any) (reflect.Type, error) {
	// Check if i is a reflect.Value and if so use i's Interface() as i
	i = ToInterface(i)

	// Get i's type and kind
	var typ reflect.Type = reflect.TypeOf(i)
	var knd reflect.Kind = typ.Kind()

	// Confirm i is a pointer
	ctx.Logger.Comment("Confirming i is a pointer ...")
	ctx.Logger.Infof("knd=%s", knd.String())
	if knd != reflect.Pointer {
		return nil, fmt.Errorf("Unsupported type (kind=%s)", knd.String())
	}

	// If *i is a map, success
	ctx.Logger.Comment("Checking if *i is a map ...")
	knd = typ.Elem().Kind()
	ctx.Logger.Infof("knd=%s", knd.String())
	if knd == reflect.Map {
		return typ.Elem(), nil
	}

	// If i is not a **, error
	ctx.Logger.Comment("Checking if i is not a ** ...")
	if knd != reflect.Pointer {
		return nil, fmt.Errorf("Unsupported type (kind=%s)", knd.String())
	}

	// Check if **i is a struct
	ctx.Logger.Comment("Checking if **i is a map")
	typ = typ.Elem().Elem()
	knd = typ.Kind()
	ctx.Logger.Infof("knd=%s", knd.String())
	if knd != reflect.Map {
		return nil, fmt.Errorf("Unsupported type (kind=%s)", knd.String())
	}

	return typ, nil

	// var knd reflect.Kind

	// // Confirm ppMapType is a pointer
	// knd = ppMapType.Kind()
	// if knd != reflect.Pointer {
	// 	return nil, fmt.Errorf("Expecting a pointer-to-pointer-to map (kind=%s)", knd.String())
	// }

	// // Confirm *ppMapType is a pointer
	// knd = ppMapType.Elem().Kind()
	// if knd != reflect.Pointer {
	// 	return nil, fmt.Errorf("Expecting a pointer-to-pointer-to map (kind=%s)", knd.String())
	// }

	// // Confirm **ppMapType is a map
	// knd = ppMapType.Elem().Kind()
	// if knd != reflect.Pointer {
	// 	return nil, fmt.Errorf("Expecting a pointer-to-pointer-to map (kind=%s)", knd.String())
	// }

	// // We're good :)
	// return ppMapType.Elem().Elem(), nil
}

// @note maybe these functions make more sense as either loops or recursive
// funcs, that loop until they find a struct or a non-pointer, which also
// doing cycle detection. I have an exmaple of that somewhere.
func GetUnderlyingStructType(ctx *EcoContext, i any) (reflect.Type, error) {
	// Check if i is a reflect.Value and if so use i's Interface() as i
	i = ToInterface(i)

	// Get i's type and kind
	var typ reflect.Type = reflect.TypeOf(i)
	var knd reflect.Kind = typ.Kind()

	// Confirm i is a pointer
	ctx.Logger.Comment("Confirming i is a pointer ...")
	ctx.Logger.Infof("knd=%s", knd.String())
	if knd != reflect.Pointer {
		return nil, fmt.Errorf("Unsupported type (kind=%s)", knd.String())
	}

	// If *i is a struct, success
	ctx.Logger.Comment("Checking if *i is a struct ...")
	knd = typ.Elem().Kind()
	ctx.Logger.Infof("knd=%s", knd.String())
	if knd == reflect.Struct {
		return typ.Elem(), nil
	}

	// If i is not a **, error
	ctx.Logger.Comment("Checking if i is not a ** ...")
	if knd != reflect.Pointer {
		return nil, fmt.Errorf("Unsupported type (kind=%s)", knd.String())
	}

	// Check if **i is a struct
	ctx.Logger.Comment("Checking if **i is a struct")
	typ = typ.Elem().Elem()
	knd = typ.Kind()
	ctx.Logger.Infof("knd=%s", knd.String())
	if knd != reflect.Struct {
		return nil, fmt.Errorf("Unsupported type (kind=%s)", knd.String())
	}

	return typ, nil
}

func IsZero(rv reflect.Value) bool {
	return rv == reflect.Zero(reflect.TypeFor[reflect.Value]())
}

func SetStructField(pStruct any, fieldName string, val any) error {
	rvStruct := ToRv(pStruct)
	rvVal := ToRv(val)
	field := rvStruct.Elem().FieldByName(fieldName)
	if IsZero(field) {
		return fmt.Errorf("Unable to lookup field (%s)", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("Unable to Set field (%s)", fieldName)
	}
	rvStruct.Elem().FieldByName(fieldName).Set(rvVal)

	return nil
}

func ToInterface(val any) any {
	rv, is := val.(reflect.Value)
	var i any
	if is {
		i = rv.Interface()
	} else {
		i = val
	}

	return i
}

func ToRv(val any) reflect.Value {
	rvVal, is := val.(reflect.Value)
	if !is {
		rvVal = reflect.ValueOf(val)
	}

	return rvVal
}

func UnmarshalMapKey(s string, p any) error {
	ps, is := p.(*string)
	if is {
		*ps = s
	} else {
		err := json.Unmarshal([]byte(s), p)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnmarshalStructField(pStruct any, fieldName string, buf []byte) error {
	rvStruct := ToRv(pStruct)
	field := rvStruct.Elem().FieldByName(fieldName)
	if IsZero(field) {
		return fmt.Errorf("Unable to lookup field (%s)", fieldName)
	}
	if !field.CanAddr() {
		return fmt.Errorf("Unable to Address field (%s)", fieldName)
	}
	addr := field.Addr().Interface()
	err := json.Unmarshal(buf, addr)
	if err != nil {
		return err
	}

	return nil
}

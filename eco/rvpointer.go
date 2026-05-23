package eco

import (
	"fmt"
	"reflect"

	"github.com/dylt-dev/dylt/common"
)

type RvPointer reflect.Value

func NewRvPointer(a any) (*RvPointer, error) {
	rv := common.Reflect(a)
	if rv.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("expected pointer, not %s", rv.Kind().String())
	}
	if !rv.IsValid() {
		return nil, fmt.Errorf("invalid pointer")
	}

	rvp := RvPointer(rv)
	return &rvp, nil
}

func (rvp RvPointer) CreateOrGet(ctx *common.EcoContext, n int) (*NormPtr, error) {
	ctx.Logger.Signature("RvPointer.CreateOrGet")
	ctx.Inc()
	defer ctx.Dec()

	normPtr, err := rvp.Walk(ctx)
	if err != nil {
		return nil, err
	}
	ctx.Logger.Comment("Getting elemType ...")
	elemType := normPtr.ElemType(ctx)
	ctx.Logger.Infof("elemType.Kind()=%v", elemType.Kind())

	flavor := common.NewFlavor(elemType.Kind())
	ctx.Logger.Infof("flavor=%v", flavor)
	switch flavor {
	case common.Map: return CreateOrGetMap(ctx, normPtr, n)
	case common.Scalar: return CreateOrGetScalar(ctx, normPtr)
	case common.Slice: return CreateOrGetSlice(ctx, normPtr, n)
	case common.Struct: return CreateOrGetStruct(ctx, normPtr)
	default: return nil, fmt.Errorf("unsupported type (%s)", elemType.Kind().String())
	}
}


func (rvp RvPointer) Elem() RvPointer {
	rv := reflect.Value(rvp)
	rvElem := rv.Elem()
	rvpElem := RvPointer(rvElem)
	return rvpElem
}


func (rvp RvPointer) ElemType(ctx *common.EcoContext) reflect.Type {
	ctx.Logger.Signature("RvPointer.ElemType")
	ctx.Inc()
	defer ctx.Dec()

	// Walk the pointer
	normPtr, err := rvp.Walk(ctx)
	if err != nil {
		return nil
	}
	ptr := normPtr.Value

	// If non-pointer, return element type
	// Otherwise return pointer element type
	rv := reflect.ValueOf(ptr)
	ctx.Logger.Infof("rv.Kind()=%v", rv.Kind())
	rvElemType := rv.Type().Elem()
	if rvElemType.Kind() != reflect.Pointer {
		return rvElemType
	}
	return rvElemType.Elem()
}


func (rvp RvPointer) Flavor () common.Flavor {
	rv := reflect.Value(rvp)
	return common.NewFlavor(rv.Type().Elem().Kind())
}


func (rvp RvPointer) IsNil(ctx *common.EcoContext) bool {
	ctx.Logger.Signature("IsNil")
	ctx.Inc()
	defer ctx.Dec()

	rv := reflect.Value(rvp)
	return rv.IsNil()
}


func (rvp RvPointer) IsPointer(ctx *common.EcoContext) bool {
	ctx.Logger.Signature("RvPointer.IsPointer")
	ctx.Inc()
	defer ctx.Dec()

	rv := reflect.Value(rvp)
	if !rv.IsValid() {
		return false
	}
	elemKind := rv.Type().Elem().Kind()
	ctx.Logger.Infof("elemKind=%v\n", elemKind)
	if elemKind == reflect.Pointer {
		ctx.Logger.Comment("pointer to pointer - returning true")
		return true
	}

	ctx.Logger.Comment("fall through -- returning false")
	return false
}

func (rvp RvPointer) IsReference(ctx *common.EcoContext) bool {
	ctx.Logger.Signature("RvPointer.IsReference")
	ctx.Inc()
	defer ctx.Dec()

	rv := reflect.Value(rvp)
	if !rv.IsValid() {
		return false
	}
	elemKind := rv.Type().Elem().Kind()
	ctx.Logger.Infof("elemKind=%v\n", elemKind)
	if elemKind == reflect.Slice || elemKind == reflect.Map {
		ctx.Logger.Comment("slice or map pointer - returning true")
		return true
	}

	ctx.Logger.Comment("fall through -- returning false")
	return false
}

func (rvp RvPointer) IsSlice(ctx *common.EcoContext) bool {
	ctx.Logger.Signature("RvPointer.IsSlice")
	ctx.Inc()
	defer ctx.Dec()

	rv := reflect.Value(rvp)
	if !rv.IsValid() {
		return false
	}
	elemKind := rv.Type().Elem().Kind()
	ctx.Logger.Infof("elemKind=%v\n", elemKind)
	if elemKind == reflect.Slice {
		ctx.Logger.Comment("pointer to slice - returning true")
		return true
	}

	ctx.Logger.Comment("fall through -- returning false")
	return false
}

func (rvp RvPointer) IsValue(ctx *common.EcoContext) bool {
	ctx.Logger.Signature("RvPointer.IsValue")
	ctx.Inc()
	defer ctx.Dec()

	rv := reflect.Value(rvp)
	if !rv.IsValid() {
		return false
	}
	elemKind := rv.Type().Elem().Kind()
	ctx.Logger.Infof("elemKind=%v\n", elemKind)
	flavor := common.NewFlavor(elemKind)
	ctx.Logger.Infof("flavor=%v\n", flavor)
	if flavor == common.Scalar || flavor == common.Struct {
		ctx.Logger.Comment("scalar or struct pointer - returning true")
		return true
	}

	ctx.Logger.Comment("fall through -- returning false")
	return false
}


func (rvp RvPointer) Set (a any) {
	reflect.Value(rvp).Elem().Set(common.Reflect(a))
}


func (rvp RvPointer) Walk(ctx *common.EcoContext) (*NormPtr, error) {
	ctx.Logger.Signature("RvPointer.Walk")
	ctx.Inc()
	defer ctx.Dec()

	if rvp.IsNil(ctx) {
		return nil, fmt.Errorf("Expecting non-nil value")
	}

	// We definitely have a non-nil pointer
	// Now, we loop until a success or error condition
	// - rv is a non-nil pointer to a Value type (scalar struct) - success
	// - rv is a non-nil pointer to a non-nil Reference type (slice or map) - success
	// - rv is a pointer to a non-pointer - error
	// - rv is a pointer to a nil Value pointer - success
	// - rv is a pointer to a nil Reference pointer - success
	// - rv is a pointer to a refernce pointer to a nil reference  - success
	// - else error
	ctx.Logger.Comment("we have a non-nil pointer")
	for {
		flavor := rvp.Flavor()
		switch flavor {
		case common.Map, common.Scalar, common.Slice, common.Struct:
			return NewNormPtr(ctx, reflect.Value(rvp).Interface())
		case common.Pointer:
			if rvp.Elem().IsNil(ctx) {
				return NewNormPtr(ctx, reflect.Value(rvp).Interface())
			}
		default:
			err := fmt.Errorf("Unsupported type (%s)", reflect.Value(rvp).Type().Elem().Kind().String())
			ctx.Logger.Error(err.Error())
			return nil, err	
		}
		rvp = rvp.Elem()
	}
	/*
		rv := reflect.Value(rvp)
		ctx.Logger.Infof("rv.Type()=%v", rv.Type())
		// - rv is a non-nil pointer to a Value type (scalar struct) - success
		ctx.Logger.Comment("Check if pointer is non-nil and is a reference or value")
		if !rvp.IsNil(ctx) && (rvp.IsReference(ctx) || rvp.IsValue(ctx)) {
			normPtr, err := NewNormPtr(ctx, rv.Interface())
			return normPtr, err
		}

		// rv is a pointer to a non-nil pointer
		ctx.Logger.Comment("Checking if pointer-to-pointer")
		if rvp.IsPointer(ctx) {
			ctx.Logger.Comment("Creating new pointer with pointer dereference")
			rvpElem, err := NewRvPointer(rv.Elem())
			if err != nil {
				return nil, err
			}
			if rvpElem.IsNil(ctx) {
				ctx.Logger.Comment("pointer to nil pointer; success")
				normPtr, err := NewNormPtr(ctx, rv.Interface())
				return normPtr, err
			}
		} else {
			ctx.Logger.Infof("unsupported type (%s) - error", rv.Kind().String())
			return nil, fmt.Errorf("unsupported type (%s)", rv.Kind().String())
		}

		// If we've gotten this far, we have a pointer to a non-nil pointer, so we keep looking
		ctx.Logger.Comment("non-nil pointer found; keep looking")
		prvp, err := NewRvPointer(rv.Elem().Interface())
		if err != nil {
			return nil, err
		}
		rvp = *prvp
	*/
}

func CreateOrGetMap(ctx *common.EcoContext, normPtr *NormPtr, n int) (*NormPtr, error) {
	ctx.Logger.Signature("RvPointer.CreateOrGetMap")
	ctx.Inc()
	defer ctx.Dec()

	ctx.Logger.Infof("normPtr.Value.Type()=%v", reflect.ValueOf(normPtr.Value).Type())
	if normPtr.IsAllocated(ctx) {
		return normPtr, nil
	}
	
	elemType := normPtr.ElemType(ctx)
	ctx.Logger.Infof("elemType=%s", elemType)
	ctx.Logger.Comment("Allocating new map")
	mapPtr := NewMap(elemType, n)
	if normPtr.IsPointer(ctx) {
		normPtr.Set(mapPtr)
	} else {
		rvMap := reflect.ValueOf(mapPtr).Elem()
		normPtr.Set(rvMap)
	}
	
	return &NormPtr{mapPtr}, nil
}

func CreateOrGetScalar(ctx *common.EcoContext, normPtr *NormPtr) (*NormPtr, error) {
	if !normPtr.IsPointer(ctx) {
		return normPtr, nil
	}
	
	elemType := normPtr.ElemType(ctx)

	rvNew := reflect.New(elemType)
	normPtr.Set(rvNew)

	newNormPtr, err := NewNormPtr(ctx, rvNew.Interface())
	if err != nil {
		return nil, err
	}

	return newNormPtr, nil
}


func CreateOrGetSlice(ctx *common.EcoContext, normPtr *NormPtr, n int) (*NormPtr, error) {
	ctx.Logger.Signature("CreateOrGetSlice", reflect.TypeOf(normPtr.Value).Kind().String(), n)
	ctx.Inc()
	defer ctx.Dec()

	ctx.Logger.Infof("ptr=%#v", normPtr)
	isAllocated := normPtr.IsAllocated(ctx)
	ctx.Logger.Infof("isAllocated=%#v", isAllocated)
	
	isBigEnough := normPtr.IsBigEnough(ctx, n)
	ctx.Logger.Infof("isBigEnough=%#v", isBigEnough)

	if isAllocated && isBigEnough {
		return normPtr, nil
	}

	elemType := normPtr.ElemType(ctx)
	ctx.Logger.Infof("elemType=%s", elemType)

	ctx.Logger.Comment("Allocating new slice")
	slicePtr := NewSlice(elemType, n)
	if normPtr.IsPointer(ctx) {
		reflect.ValueOf(normPtr.Value).Elem().Set(reflect.ValueOf(slicePtr))
	} else {
		reflect.ValueOf(normPtr.Value).Elem().Set(reflect.ValueOf(slicePtr).Elem())
	}
	
	return &NormPtr{slicePtr}, nil
/*
	// Since we need to allocate. we need a pointer-to-a-pointer
	if !normPtr.IsPointer(ctx) {
		err := fmt.Errorf("Allocation requires a pointer to a pointer (%s)", reflect.ValueOf(normPtr.Value).Elem().Type())
		ctx.Logger.Error(err.Error())
		return nil, err
	}

	ctx.Logger.Comment("Allocating new slice")
	elemType := normPtr.ElemType(ctx)
	ctx.Logger.Infof("elemType=%s", elemType)
	slicePtr := MakeSlice(elemType, n)
	reflect.ValueOf(normPtr.Value).Elem().Set(reflect.ValueOf(slicePtr))
	normPtr = &NormPtr{slicePtr}

	ctx.Logger.Infof("ptr=%#v", normPtr)
	return normPtr, nil
*/
}

func CreateOrGetStruct(ctx *common.EcoContext, normPtr *NormPtr) (*NormPtr, error) {
	if !normPtr.IsPointer(ctx) {
		return normPtr, nil
	}
	
	elemType := normPtr.ElemType(ctx)
	rv := reflect.New(elemType)
	normPtr.Set(rv)
	newNormPtr, err := NewNormPtr(ctx, rv.Interface())
	if err != nil {
		return nil, err
	}

	return newNormPtr, nil
}

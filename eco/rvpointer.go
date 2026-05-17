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

func (rvp RvPointer) CreateOrGet(ctx *common.EcoContext, kvSlice KvSeries) (*NormPtr, error) {
	ctx.Logger.Signature("RvPointer.CreateOrGet")
	ctx.Inc()
	defer ctx.Dec()

	normPtr, err := rvp.Walk(ctx)
	if err != nil {
		return nil, err
	}
	ctx.Logger.Comment("Getting elemType ...")
	elemType, err := normPtr.ElemType(ctx)
	ctx.Logger.Infof("elemType.Kind()=%v", elemType.Kind())
	if err != nil {
		return nil, err
	}

	flavor := NewFlavor(elemType.Kind())
	ctx.Logger.Infof("flavor=%v", flavor)
	switch flavor {
	case Map: return CreateOrGetMap(ctx, normPtr, kvSlice.Len())
	case Scalar: return CreateOrGetScalar(ctx, normPtr)
	case Slice: return rvp.CreateOrGetSlice(ctx, kvSlice)
	case Struct: return CreateOrGetStruct(ctx, normPtr)
	default: return nil, fmt.Errorf("unsupported type (%s)", elemType.Kind().String())
	}
}

func CreateOrGetMap(ctx *common.EcoContext, normPtr *NormPtr, size int) (*NormPtr, error) {
	ctx.Logger.Signature("RvPointer.CreateOrGetMap")
	ctx.Inc()
	defer ctx.Dec()

	ctx.Logger.Infof("normPtr.Value.Type()=%v", reflect.ValueOf(normPtr.Value).Type())
	if !normPtr.IsPointer(ctx) {
		return normPtr, nil
	}
	
	elemType, err := normPtr.ElemType(ctx)
	if err != nil {
		return nil, err
	}

	rvMap := reflect.MakeMapWithSize(elemType, size)
	rv := reflect.New(elemType)
	rv.Elem().Set(rvMap)
	newNormPtr, err := NewNormPtr(ctx, rv.Interface())
	if err != nil {
		return nil, err
	}

	return newNormPtr, nil
}

func CreateOrGetScalar(ctx *common.EcoContext, normPtr *NormPtr) (*NormPtr, error) {
	if !normPtr.IsPointer(ctx) {
		return normPtr, nil
	}
	
	elemType, err := normPtr.ElemType(ctx)
	if err != nil {
		return nil, err
	}

	rv := reflect.New(elemType)
	newNormPtr, err := NewNormPtr(ctx, rv.Interface())
	if err != nil {
		return nil, err
	}

	return newNormPtr, nil

	
}

func (rvp RvPointer) CreateOrGetSlice(ctx *common.EcoContext, kvSlice KvSeries) (*NormPtr, error) {
	ctx.Logger.Signature("RvPointer.CreateOrGetSlice")
	ctx.Inc()
	defer ctx.Dec()
	normPtr, err := rvp.Walk(ctx)
	if err != nil {
		return nil, err
	}

	ctx.Logger.Infof("ptr=%#v", normPtr)
	isAllocated, err := normPtr.IsAllocated(ctx)
	if err != nil {
		return nil, err
	}
	ctx.Logger.Infof("isAllocated=%#v", isAllocated)
	n := kvSlice.MaxIndex() + 1
	ctx.Logger.Infof("n=%#v", n)
	isBigEnough, err := normPtr.IsBigEnough(ctx, n)
	if err != nil {
		return nil, err
	}
	ctx.Logger.Infof("isBigEnough=%#v", isBigEnough)
	if !isAllocated || !isBigEnough {
		ctx.Logger.Comment("Allocating new slice")
		elemType, err := normPtr.ElemType(ctx)
		if err != nil {
			return nil, err
		}
		ctx.Logger.Infof("elemType=%s", elemType)
		slice := MakeSlice(elemType, n)
		normPtr = &NormPtr{slice}
	}

	ctx.Logger.Infof("ptr=%#v", normPtr)
	return normPtr, nil
}

func CreateOrGetStruct(ctx *common.EcoContext, normPtr *NormPtr) (*NormPtr, error) {
	if !normPtr.IsPointer(ctx) {
		return normPtr, nil
	}
	
	elemType, err := normPtr.ElemType(ctx)
	if err != nil {
		return nil, err
	}

	rv := reflect.New(elemType)
	newNormPtr, err := NewNormPtr(ctx, rv.Interface())
	if err != nil {
		return nil, err
	}

	return newNormPtr, nil
}

func (rvp RvPointer) ElemType(ctx *common.EcoContext) (reflect.Type, error) {
	ctx.Logger.Signature("RvPointer.ElemType")
	ctx.Inc()
	defer ctx.Dec()

	// Walk the pointer
	normPtr, err := rvp.Walk(ctx)
	if err != nil {
		return nil, err
	}
	ptr := normPtr.Value

	// If non-pointer, return element type
	// Otherwise return pointer element type
	rv := reflect.ValueOf(ptr)
	ctx.Logger.Infof("rv.Kind()=%v", rv.Kind())
	rvElemType := rv.Type().Elem()
	if rvElemType.Kind() != reflect.Pointer {
		return rvElemType, nil
	}
	return rvElemType.Elem(), nil
}

func (rvp RvPointer) IsNil(ctx *common.EcoContext) bool {
	ctx.Logger.Signature("IsNil")
	ctx.Inc()
	defer ctx.Dec()

	rv := reflect.Value(rvp)
	if rv.IsNil() {
		ctx.Logger.Comment("rv.IsNil() is true")
		return true
	}
	if rvp.IsReference(ctx) && rv.Elem().IsNil() {
		ctx.Logger.Comment("rvp is a Reference and rv.IsNil() is true")
		return true
	}

	ctx.Logger.Comment("pointer is not nil")
	return false
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
	flavor := NewFlavor(elemKind)
	ctx.Logger.Infof("flavor=%v\n", flavor)
	if flavor == Scalar || flavor == Struct {
		ctx.Logger.Comment("scalar or struct pointer - returning true")
		return true
	}

	ctx.Logger.Comment("fall through -- returning false")
	return false
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
	}
}

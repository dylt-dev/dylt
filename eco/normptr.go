package eco

import (
	"fmt"
	"reflect"

	"github.com/dylt-dev/dylt/common"
)

type NormPtr struct {
	Value any
}

// Rules for NormPtr
//	-- Must be a pointer
//  --     to a Scalar or a struct
//  --     or a non-nil Map or Slice
//  --     or a pointer
//            that IsNil
//            or is Type().Elem() == map or slice and Elem().IsNil()
func NewNormPtr (ctx *common.EcoContext, a any) (*NormPtr, error) {
	ctx.Logger.Signature("NewNormPtr")
	ctx.Inc()
	defer ctx.Dec()
	
	if !IsNormPointer(ctx, a) {
		return nil, fmt.Errorf("not a valid norm pointer")
	}

	return &NormPtr{a}, nil	
}


func (this NormPtr) Elem () RvPointer {
	rv := reflect.ValueOf(this.Value)
	rvElem := rv.Elem()
	rvpElem := RvPointer(rvElem)

	return rvpElem
}


func (this NormPtr) ElemType(ctx *common.EcoContext) reflect.Type {
	ctx.Logger.Signature("RvPointer.ElemType")
	ctx.Inc()
	defer ctx.Dec()

	return RvPointer(reflect.ValueOf(this.Value)).ElemType(ctx)
}


func (this NormPtr) IsAllocated (ctx *common.EcoContext) bool {
	rv := reflect.ValueOf(this.Value)
	rvElemType := rv.Type().Elem()

	if rvElemType.Kind() != reflect.Pointer {
		if this.IsReference(ctx) && this.Elem().IsNil(ctx) {
			return false
		}
		return true
	}

	return false
}


func (this NormPtr) IsBigEnough (ctx *common.EcoContext, n int) bool {
	ctx.Logger.Signature("NormPtr.IsBigEnough")
	ctx.Inc()
	defer ctx.Dec()

	ctx.Logger.Infof("this.Value.Type()=%v", reflect.ValueOf(this.Value).Type())
	
	// Confirm that this is ultimately a slice ptr
	rvp := RvPointer(reflect.ValueOf(this.Value))
	elemType := rvp.ElemType(ctx)
	ctx.Logger.Infof("elemType=%v", elemType)
	if elemType.Kind() != reflect.Slice {
		ctx.Logger.Infof("type is not slice (%s) - returning true", elemType.Kind())
		return true
	}

	// nil pointers always return false
	if rvp.IsNil(ctx) {
		ctx.Logger.Info("pointer is nil - returning false")
		return false
	}

	// ptr to slice
	if rvp.IsSlice(ctx) {
		rv := reflect.ValueOf(this.Value)
		ctx.Logger.Comment("pointer to slice")
		cap := rv.Elem().Cap()
		ctx.Logger.Infof("cap=%v", cap)
		return cap >= n
	}

	// ptr to ptr - by NormPtr rules these are always nil so return false
	if rvp.IsPointer(ctx) {
		ctx.Logger.Comment("pointer to pointer")
		return false
	}

	// Non-slice and non-pointer, so return true
	return true
}


func (this NormPtr) IsPointer(ctx *common.EcoContext) bool {
	return RvPointer(reflect.ValueOf(this.Value)).IsPointer(ctx)
}


func (this NormPtr) IsReference(ctx *common.EcoContext) bool {
	return RvPointer(reflect.ValueOf(this.Value)).IsReference(ctx)
}


func (this NormPtr) Set(a any) {
	reflect.ValueOf(this.Value).Elem().Set(common.Reflect(a))
}


func IsNormPointer (ctx *common.EcoContext, a any) bool {
	ctx.Logger.Signature("IsNormPointer", reflect.ValueOf(a).Type())
	ctx.Inc()
	defer ctx.Dec()

	rvp := RvPointer(reflect.ValueOf(a))
	flavor := rvp.Flavor()
	switch flavor {
	case common.Map, common.Scalar, common.Slice, common.Struct:
		if rvp.IsNil(ctx) {
			return false
		}
		return true
	case common.Pointer:
		if rvp.Elem().IsNil(ctx) {
			return true
		} else {
			return false
		}
	default:
		return false
	}
	
/*
	rv := common.Reflect(a)
	if rv.Kind() != reflect.Pointer {
		ctx.Logger.Infof("expected pointer, got %s instead", rv.Kind())
		return false
	}
	if rv.IsNil() {
		ctx.Logger.Infof("nil pointer")
		return false
	}

	typeElem := rv.Type().Elem()
	flavor := NewFlavor(typeElem.Kind())
	switch flavor{
	case Scalar, Struct: return true
	case Map, Slice: {
		if rv.Elem().IsNil() {
			ctx.Logger.Infof("nil pointer to %s", typeElem.Kind().String())
			return false
		}
		return true
	}
	case Pointer:{
		if rv.Elem().IsNil() {
			return true
		}
		rvp, err := NewRvPointer(rv.Elem())
		if err != nil {
			return false
		}
		if rvp.IsReference(ctx) && rv.Elem().Elem().IsNil() {
			return true
		}
		ctx.Logger.Infof("pointer to non-nil pointer")
		return false
	}

	default: {
		ctx.Logger.Infof("unsupported type (%v)", rv.Kind().String())
		return false
	}
	}
*/
}
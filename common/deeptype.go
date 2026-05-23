package common

import (
	"fmt"
	"reflect"
)

/*
	Type declaration	type deepType struct{Data []map[string]struct{Slice []map[string]struct{ Val []struct{N int}}}}


	Object creation     structTree1 := NewValueTree(ctx, "N", 13)
	                    sliceTree1  := NewValueTree(ctx, 3, structTree1)
	                    valTree     := NewValueTree(ctx, "Val", sliceTree1)
	                    mapTree     := NewValueTree(ctx, "foo", valTree)
	                    sliceTree2  := NewValueTree(ctx, 0, mapTree)
	                    structTree2 := NewValueTree(ctx, "Slice", sliceTree2)
	                    mapTree2    := NewValueTree(ctx, "bar", structTree2)
	                    sliceTree3  := NewValueTree(ctx, 2, mapTree2)
	                    structTree3 := NewValueTree(ctx, "Data", sliceTree3)

	Field access        x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
	                    x.Data[0][""].Slice[0][""].Val[0].N
*/

type DeepType struct{ typ reflect.Type }
type ScalarType DeepType
type MapType DeepType
type SliceType DeepType
type StructType DeepType

func NewDeepType(typ reflect.Type) DeepType {
	switch NewFlavor(typ.Kind()) {
	case Map,
		Scalar,
		Slice,
		Struct:
		return DeepType{typ}
	default:
		panic(fmt.Sprintf("Unsupported type (%s)", typ))
	}
}

func NewMapType(dt DeepType) MapType {
	if dt.Flavor() != Map {
		panic(fmt.Errorf("expecting map (%s)", dt.Flavor().String()))
	}

	return MapType(dt)
}

func NewScalarType(dt DeepType) ScalarType {
	if dt.Flavor() != Scalar {
		panic(fmt.Errorf("expecting map (%s)", dt.Flavor().String()))
	}

	return ScalarType(dt)
}

func NewSliceType(dt DeepType) SliceType {
	if dt.Flavor() != Slice {
		panic(fmt.Errorf("expecting map (%s)", dt.Flavor().String()))
	}

	return SliceType(dt)
}

func NewStructType(dt DeepType) StructType {
	if dt.Flavor() != Struct {
		panic(fmt.Errorf("expecting map (%s)", dt.Flavor().String()))
	}

	return StructType(dt)
}


func (dt DeepType) EmitFieldRef() {
	switch dt.Flavor() {
	case Map:
		NewMapType(dt).EmitFieldRef()
	case Scalar:
		NewScalarType(dt).EmitFieldRef()
	case Slice:
		NewSliceType(dt).EmitFieldRef()
	case Struct:
		NewStructType(dt).EmitFieldRef()
	default:
		panic(fmt.Errorf("unsupported type (%s)", dt.typ))
	}
}

func (dt DeepType) EmitTree(plevel *int) {
	switch dt.Flavor() {
	case Map:
		NewMapType(dt).EmitTree(plevel)
	case Scalar:
		NewScalarType(dt).EmitTree(plevel)
	case Slice:
		NewSliceType(dt).EmitTree(plevel)
	case Struct:
		NewStructType(dt).EmitTree(plevel)
	default:
		panic(fmt.Errorf("unsupported type (%s)", dt.typ))
	}
}

func (dt DeepType) Flavor() Flavor {
	return NewFlavor(dt.typ.Kind())
}

func (t MapType) EmitFieldRef() {
	// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
	if t.typ.Key().Kind() == reflect.String {
		fmt.Printf("[%q]", t.KeyName())
	} else {
		fmt.Printf("[%s]", t.KeyName())
	}
	if !t.IsScalar() {
		NewDeepType(t.typ.Elem()).EmitFieldRef()
	}
}

func (t MapType) EmitTree(plevel *int) {
	if t.IsScalar() {
		if t.typ.Elem().Kind() == reflect.String {
			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %q)\n", *plevel, t.KeyName(), t.ZeroValue())
		} else {
			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %s)\n", *plevel, t.KeyName(), t.ZeroValue())
		}
	} else {
		NewDeepType(t.typ.Elem()).EmitTree(plevel)
		*plevel++
		fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", tree%d)\n", *plevel, t.KeyName(), *plevel-1)
	}
}

func (t MapType) IsScalar() bool {
	return NewDeepType(t.typ.Elem()).Flavor() == Scalar
}

func (t MapType) KeyName() string {
	return fmt.Sprintf("%v", reflect.Zero(t.typ.Key()).Interface())
}

func (t MapType) ZeroValue() string {
	return fmt.Sprintf("%v", reflect.Zero(t.typ.Elem()))
}

func (t SliceType) EmitFieldRef() {
	// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
	if t.typ.Elem().Kind() == reflect.String {
		fmt.Printf("[%q]", t.KeyName())
	} else {
		fmt.Printf("[%s]", t.KeyName())
	}
	if !t.IsScalar() {
		NewDeepType(t.typ.Elem()).EmitFieldRef()
	}
}

func (t ScalarType) EmitTree(plevel *int) {
}

func (t ScalarType) KeyName() {

}

func (t ScalarType) ZeroValue() reflect.Value {
	return reflect.Zero(t.typ)
}

func (t ScalarType) EmitFieldRef() {
}

func (t SliceType) EmitTree(plevel *int) {
	if t.IsScalar() {
		if t.typ.Elem().Kind() == reflect.String {
			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %q)\n", *plevel, t.KeyName(), t.ZeroValue())
		} else {
			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %s)\n", *plevel, t.KeyName(), t.ZeroValue())
		}
	} else {
		NewDeepType(t.typ.Elem()).EmitTree(plevel)
		*plevel++
		fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", tree%d)\n", *plevel, t.KeyName(), *plevel-1)
	}
}

func (t SliceType) IsScalar() bool {
	return NewDeepType(t.typ.Elem()).Flavor() == Scalar
}

func (t SliceType) KeyName() string {
	return "0"
}

func (t SliceType) ZeroValue() string {
	return fmt.Sprintf("%v", reflect.Zero(t.typ.Elem()).Interface())
}

func (t StructType) EmitFieldRef() {
	// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
	fmt.Printf(".%s", t.KeyName())
	if !t.IsScalar() {
		NewDeepType(t.typ.Field(0).Type).EmitFieldRef()
	}
}

func (t StructType) EmitTree(plevel *int) {
	if t.IsScalar() {
		if t.typ.Field(0).Type.Kind() == reflect.String {
			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %q)\n", *plevel, t.KeyName(), t.ZeroValue())
		} else {
			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %s)\n", *plevel, t.KeyName(), t.ZeroValue())
		}
	} else {
		NewDeepType(t.typ.Field(0).Type).EmitTree(plevel)
		*plevel++
		fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", tree%d)\n", *plevel, t.KeyName(), *plevel-1)
	}
}

func (t StructType) IsScalar() bool {
	return NewDeepType(t.typ.Field(0).Type).Flavor() == Scalar
}

func (t StructType) KeyName() string {
	return t.typ.Field(0).Name
}

func (t StructType) ZeroValue() string {
	return fmt.Sprintf("%v", reflect.Zero(t.typ.Field(0).Type).Interface())
}

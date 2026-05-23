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

// type ScalarType DeepType
type MapType DeepType
type SliceType DeepType
type StructType DeepType

type DeepSubType interface {
	emitKey()
	keyName() string
	keyType() DeepType
	nextType() DeepType
	zeroValue() string
}

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

// func NewScalarType(dt DeepType) ScalarType {
// 	if dt.Flavor() != Scalar {
// 		panic(fmt.Errorf("expecting map (%s)", dt.Flavor().String()))
// 	}

// 	return ScalarType(dt)
// }

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
	// switch dt.Flavor() {
	// case Map:
	// 	NewMapType(dt).EmitFieldRef()
	// case Slice:
	// 	NewSliceType(dt).EmitFieldRef()
	// case Struct:
	// 	NewStructType(dt).EmitFieldRef()
	// default:
	// 	panic(fmt.Errorf("unsupported type (%s)", dt.typ))
	// }
	dt.subType().emitKey()
	if !dt.isScalar() {
		dt.subType().nextType().EmitFieldRef()
	}
}

func (dt DeepType) EmitTree(plevel *int) {
	// switch dt.Flavor() {
	// case Map:
	// 	NewMapType(dt).EmitTree(plevel)
	// // case Scalar:
	// // 	NewScalarType(dt).EmitTree(plevel)
	// case Slice:
	// 	NewSliceType(dt).EmitTree(plevel)
	// case Struct:
	// 	NewStructType(dt).EmitTree(plevel)
	// default:
	// 	panic(fmt.Errorf("unsupported type (%s)", dt.typ))
	// }
	if dt.isScalar() {
		// if dt.subType().nextType().typ.Kind() == reflect.String {
		if dt.subType().keyType().typ.Kind() == reflect.String {
			fmt.Printf("tree%d := NewValueTree(ctx, %q, %s)\n", *plevel, dt.subType().keyName(), dt.subType().zeroValue())
		} else {
			fmt.Printf("tree%d := NewValueTree(ctx, %s, %s)\n", *plevel, dt.subType().keyName(), dt.subType().zeroValue())
		}
	} else {
		dt.subType().nextType().EmitTree(plevel)
		*plevel++
		var keyName string
		if dt.subType().keyType().typ.Kind() == reflect.String {
			keyName = fmt.Sprintf("%q", dt.subType().keyName())
		} else {
			keyName = dt.subType().keyName()
		}

		fmt.Printf("tree%d := NewValueTree(ctx, %s, tree%d)\n", *plevel, keyName, *plevel-1)
	}
}

func (dt DeepType) Flavor() Flavor {
	return NewFlavor(dt.typ.Kind())
}

func (dt DeepType) isScalar() bool {
	return dt.subType().nextType().Flavor() == Scalar
}

func (dt DeepType) subType() DeepSubType {
	switch dt.Flavor() {
	case Map:
		return MapType(dt)
	case Slice:
		return SliceType(dt)
	case Struct:
		return StructType(dt)
	default:
		panic(fmt.Errorf("unsupported type (%s)", dt.typ))
	}
}

// func (t MapType) EmitFieldRef() {
// 	// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
// 	t.emitKey()
// 	if !t.IsScalar() {
// 		t.nextType().EmitFieldRef()
// 	}
// }

// func (t MapType) EmitTree(plevel *int) {
// 	if t.IsScalar() {
// 		if t.typ.Elem().Kind() == reflect.String {
// 			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %q)\n", *plevel, t.keyName(), t.ZeroValue())
// 		} else {
// 			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %s)\n", *plevel, t.keyName(), t.ZeroValue())
// 		}
// 	} else {
// 		NewDeepType(t.typ.Elem()).EmitTree(plevel)
// 		*plevel++
// 		fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", tree%d)\n", *plevel, t.keyName(), *plevel-1)
// 	}
// }

// func (t MapType) IsScalar() bool {
// 	return t.nextType().Flavor() == Scalar
// }

func (t MapType) emitKey() {
	// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
	if t.typ.Key().Kind() == reflect.String {
		fmt.Printf("[%q]", t.keyName())
	} else {
		fmt.Printf("[%s]", t.keyName())
	}
}

func (t MapType) keyType() DeepType {
	return NewDeepType(t.typ.Key())
}

func (t MapType) keyName() string {
	return fmt.Sprintf("%s", reflect.Zero(t.typ.Key()))
	// buf, err := json.Marshal(reflect.Zero(t.typ.Key()).Interface())
	// if err != nil {
	// 	panic(err)
	// }
	// return fmt.Sprintf("%s", string(buf))
}

func (t MapType) nextType() DeepType {
	return NewDeepType(t.typ.Elem())
}

func (t MapType) zeroValue() string {
	return fmt.Sprintf("%v", reflect.Zero(t.nextType().typ))
}

// func (t ScalarType) EmitTree(plevel *int) {
// }

// func (t ScalarType) KeyName() {

// }

// func (t ScalarType) ZeroValue() reflect.Value {
// 	return reflect.Zero(t.typ)
// }

// func (t ScalarType) EmitFieldRef() {
// }

// func (t SliceType) EmitFieldRef() {
// 	// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
// 	t.emitKey()
// 	if !t.IsScalar() {
// 		t.nextType().EmitFieldRef()
// 	}
// }

// func (t SliceType) EmitTree(plevel *int) {
// 	if t.IsScalar() {
// 		if t.typ.Elem().Kind() == reflect.String {
// 			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %q)\n", *plevel, t.KeyName(), t.ZeroValue())
// 		} else {
// 			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %s)\n", *plevel, t.KeyName(), t.ZeroValue())
// 		}
// 	} else {
// 		NewDeepType(t.typ.Elem()).EmitTree(plevel)
// 		*plevel++
// 		fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", tree%d)\n", *plevel, t.KeyName(), *plevel-1)
// 	}
// }

// func (t SliceType) IsScalar() bool {
// 	return t.nextType().Flavor() == Scalar
// }

func (t SliceType) emitKey() {
	if t.typ.Elem().Kind() == reflect.String {
		fmt.Printf("[%q]", t.keyName())
	} else {
		fmt.Printf("[%s]", t.keyName())
	}
}

func (t SliceType) keyType() DeepType {
	return NewDeepType(reflect.TypeFor[int]())
}

func (t SliceType) keyName() string {
	return "0"
}

func (t SliceType) nextType() DeepType {
	return NewDeepType(t.typ.Elem())
}

func (t SliceType) zeroValue() string {
	return fmt.Sprintf("%v", reflect.Zero(t.nextType().typ))
}

// func (t StructType) EmitFieldRef() {
// 	// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
// 	t.emitKey()
// 	if !t.IsScalar() {
// 		t.nextType().EmitFieldRef()
// 	}
// }

// func (t StructType) EmitTree(plevel *int) {
// 	if t.IsScalar() {
// 		if t.typ.Field(0).Type.Kind() == reflect.String {
// 			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %q)\n", *plevel, t.KeyName(), t.ZeroValue())
// 		} else {
// 			fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", %s)\n", *plevel, t.KeyName(), t.ZeroValue())
// 		}
// 	} else {
// 		NewDeepType(t.typ.Field(0).Type).EmitTree(plevel)
// 		*plevel++
// 		fmt.Printf("tree%d := NewValueTree(ctx, \"%s\", tree%d)\n", *plevel, t.KeyName(), *plevel-1)
// 	}
// }

// func (t StructType) IsScalar() bool {
// 	return t.nextType().Flavor() == Scalar
// }

func (t StructType) emitKey() {
	fmt.Printf(".%s", t.keyName())
}

func (t StructType) keyName() string {
	return fmt.Sprintf("%s", t.typ.Field(0).Name)
}

func (t StructType) keyType() DeepType {
	return NewDeepType(reflect.TypeFor[string]())
}

func (t StructType) nextType() DeepType {
	return NewDeepType(t.typ.Field(0).Type)
}

func (t StructType) zeroValue() string {
	return fmt.Sprintf("%v", reflect.Zero(t.nextType().typ).Interface())
}

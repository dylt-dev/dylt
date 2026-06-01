package common

import (
	"fmt"
	"io"
	"reflect"
	"strings"
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

type MapType DeepType
type SliceType DeepType
type StructType DeepType

type DeepSubType interface {
	keyName() string
	keyType() DeepType
	nextType() DeepType
	writeKeyRef([]any, io.Writer)
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


func (dt DeepType) Flavor() Flavor {
	return NewFlavor(dt.typ.Kind())
}


func (dt DeepType) WriteTreeDecl(plevel *int, values []any, w io.Writer) {
	// If scalar, emit the current tree
	// Else, recurse on the element type, skipping the first value,
	//       then bump the level and emit the current tree
	if dt.isScalar() {
		key := dt.createValueTreeKey(values[0])
		var val any = values[1]
		if reflect.TypeOf(val).Kind() == reflect.String {
			fmt.Fprintf(w, "tree%d := NewValueTree(ctx, %v, %q)\n", *plevel, key, val)
		} else {
			fmt.Fprintf(w, "tree%d := NewValueTree(ctx, %v, %v)\n", *plevel, key, val)
		}
	} else {
		key := dt.createValueTreeKey(values[0])
		dt.nextType().WriteTreeDecl(plevel, values[1:], w)
		*plevel++
		fmt.Fprintf(w, "tree%d := NewValueTree(ctx, %v, tree%d)\n", *plevel, key, *plevel-1)
	}
}


func (dt DeepType) WriteValueRef(values []any, w io.Writer) {
	dt.writeKeyRef(values, w)
	// if not scalar, recurse
	if !dt.isScalar() {
		dt.nextType().WriteValueRef(values[1:], w)
	}
}


func (dt DeepType) writeKeyRef(values []any, w io.Writer) {
	dt.subType().writeKeyRef(values, w)
}


func (dt DeepType) isScalar() bool {
	return dt.nextType().Flavor() == Scalar
}


func (dt DeepType) createValueTreeKey(a any) any {
	var key any
	if dt.keyType().typ.Kind() == reflect.String {
		key = fmt.Sprintf("%q", a)
	} else {
		key = a
	}

	return key
}


func (dt DeepType) keyName() string {
	return dt.subType().keyName()
}


func (dt DeepType) keyType() DeepType {
	return dt.subType().keyType()
}


func (dt DeepType) nextType() DeepType {
	return dt.subType().nextType()
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


func (dt DeepType) zeroValue() string {
	return dt.subType().zeroValue()
}


func (t MapType) writeKeyRef(values []any, w io.Writer) {
	// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
	if t.typ.Key().Kind() == reflect.String {
		fmt.Fprintf(w, "[%q]", values[0])
	} else {
		fmt.Fprintf(w, "[%v]", values[0])
	}
}


func (t MapType) keyType() DeepType {
	return NewDeepType(t.typ.Key())
}


func (t MapType) keyName() string {
	return fmt.Sprintf("%s", reflect.Zero(t.typ.Key()))
}


func (t MapType) nextType() DeepType {
	return NewDeepType(t.typ.Elem())
}


func (t MapType) zeroValue() string {
	return fmt.Sprintf("%v", reflect.Zero(t.nextType().typ))
}


func (t SliceType) writeKeyRef(values []any, w io.Writer) {
	if t.typ.Elem().Kind() == reflect.String {
		fmt.Fprintf(w, "[%q]", values[0])
	} else {
		fmt.Fprintf(w, "[%v]", values[0])
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


func (t StructType) writeKeyRef(values []any, w io.Writer) {
	fmt.Fprintf(w, ".%s", values[0])
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




func isElemScalar(typ reflect.Type) bool {
	elemType := typ.Elem()
	elemDeep := DeepType{elemType}
	elemFlavor := elemDeep.Flavor()
	return elemFlavor == Scalar
}


func GetDeclFromType(rt reflect.Type) string {
	sb := strings.Builder{}

	switch NewFlavor(rt.Kind()) {
	case Map:
		{
			sb.WriteString("map[")
			sb.WriteString(GetDeclFromType(rt.Key()))
			sb.WriteString("]")
			sb.WriteString(GetDeclFromType(rt.Elem()))
		}
	case Scalar:
		{
			sb.WriteString(rt.String())
		}

	case Slice:
		{
			sb.WriteString("[]")
			sb.WriteString(GetDeclFromType(rt.Elem()))
		}
	case Struct:
		{
			sb.WriteString("struct{")
			for sf := range rt.Fields() {
				sb.WriteString(sf.Name)
				sb.WriteString(" ")
				sb.WriteString(GetDeclFromType(sf.Type))
				sb.WriteString(";")
			}
			sb.WriteString("}")
		}
	}
	return sb.String()
}

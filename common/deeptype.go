package common

import (
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strings"

	"github.com/jaswdr/faker"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	emitKeyRef([]any)
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


func (dt DeepType) EmitTreeDecl(plevel *int, values []any) {
	if dt.isScalar() {
		key := values[0]
		if dt.keyType().typ.Kind() == reflect.String {
			key = fmt.Sprintf("%q", key)
		}
		val := values[1]
		// if dt.nextType().typ.Kind() == reflect.String {
		// 	val = fmt.Sprintf("%q", key)
		// }
		fmt.Printf("tree%d := NewValueTree(ctx, %v, %v)\n", *plevel, key, val)
	} else {
		// Use first value as the key name. strings get quoted, other values used as is
		var keyName any
		if dt.keyType().typ.Kind() == reflect.String {
			keyName = fmt.Sprintf("%q", values[0])
		} else {
			keyName = values[0]
		}

		// advance past first value & recurse
		dt.nextType().EmitTreeDecl(plevel, values[1:])

		// bump current level up to reflect a sub tree has been emitted
		*plevel++

		fmt.Printf("tree%d := NewValueTree(ctx, %v, tree%d)\n", *plevel, keyName, *plevel-1)
	}
}


func (dt DeepType) EmitValueRef(values []any) {
	dt.emitKeyRef(values)
	// if not scalar, recurse
	if !dt.isScalar() {
		dt.nextType().EmitValueRef(values[1:])
	}
}


func (dt DeepType) Flavor() Flavor {
	return NewFlavor(dt.typ.Kind())
}


func (dt DeepType) emitKeyRef(values []any) {
	dt.subType().emitKeyRef(values)
}


func (dt DeepType) isScalar() bool {
	return dt.nextType().Flavor() == Scalar
}


func (dt DeepType) keyName () string {
	return dt.subType().keyName()
}


func (dt DeepType) keyType () DeepType {
	return dt.subType().keyType()
}


func (dt DeepType) nextType () DeepType {
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


func (t MapType) emitKeyRef(values []any) {
	// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
	if t.typ.Key().Kind() == reflect.String {
		fmt.Printf("[%q]", values[0])
	} else {
		fmt.Printf("[%v]", values[0])
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


func (t SliceType) emitKeyRef(values []any) {
	if t.typ.Elem().Kind() == reflect.String {
		fmt.Printf("[%q]", values[0])
	} else {
		fmt.Printf("[%v]", values[0])
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


// .keyName
func (t StructType) emitKeyRef(values []any) {
	fmt.Printf(".%s", values[0])
}


// name of first field
func (t StructType) keyName() string {
	return fmt.Sprintf("%s", t.typ.Field(0).Name)
}


// always string
func (t StructType) keyType() DeepType {
	return NewDeepType(reflect.TypeFor[string]())
}


func (t StructType) nextType() DeepType {
	return NewDeepType(t.typ.Field(0).Type)
}


func (t StructType) zeroValue() string {
	return fmt.Sprintf("%v", reflect.Zero(t.nextType().typ).Interface())
}


func genDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("genDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()
	
	flavor := getRandFlavor(ctx)
	switch flavor {
	case Map:
		genMapDeclaration(ctx, n, r, w)
	case Slice:
		genSliceDeclaration(ctx, n, r, w)
	case Struct:
		genStructDeclaration(ctx, n, r, w)
	default:
		panic(fmt.Errorf("How'd I get a flavor of %s???", flavor))
	}
}


func genMapDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("genMapDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()
	
	if n < 1 {
		return
	}

	w.Write([]byte("map["))
	w.Write([]byte(getRandScalar(ctx).String()))
	w.Write([]byte("]"))
	writeScalarOrRecurse(ctx, n, r, w)
}


func genMapValues(ctx *EcoContext, typ reflect.Type, r rand.Source, values *[]any, n int) {
	ctx.Signature("genMapValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	for range n {
		// emit random key value
		ctx.Commentf("generating %d values ...", n)
		keyType := typ.Key()
		keyDeep := DeepType{keyType}
		keyFlavor := keyDeep.Flavor()
		if keyFlavor != Scalar {
			panic("inconthievalble!")
		}
		a := genRandScalarValue(ctx, keyType, r)
		if keyType.Kind() == reflect.String {
			a = strings.ToLower(a.(string))
		}
		ctx.Infof("value=%v", a)
		*values = append(*values, a)

		// emit value if scalar, else recurse
		elemType := typ.Elem()
		elemDeep := DeepType{elemType}
		elemFlavor := elemDeep.Flavor()
		if elemFlavor == Scalar {
			a := genRandScalarValue(ctx, elemType, r)
			*values = append(*values, a)
		} else {
			genScalarValues(ctx, elemType, r, values)
		}
	}
}


func genRandScalarValue(ctx *EcoContext, typ reflect.Type, r rand.Source) any {
	ctx.Signature("genRandScalarValue", typ)
	ctx.Inc()
	defer ctx.Dec()
	
	switch typ.Kind() {
	case reflect.Bool:
		return faker.NewWithSeed(r).Bool()
	case reflect.Int:
		return int(faker.NewWithSeed(r).Int16Between(0, 999))
	case reflect.String:
		return faker.NewWithSeed(r).Lorem().Word()
	default:
		panic("inconthievalble!")
	}
}


func genScalarValues(ctx *EcoContext, typ reflect.Type, r rand.Source, values *[]any) {
	ctx.Signature("genScalarValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	deep := DeepType{typ}
	flavor := deep.Flavor()
	switch flavor {
	case Map:
		genMapValues(ctx, typ, r, values, 1)
	case Slice:
		genSliceValues(ctx, typ, r, values, 1)
	case Struct:
		genStructValues(ctx, typ, r, values)
	default:
		panic("inconthievalble!")
	}		
}


func genSliceValues(ctx *EcoContext, typ reflect.Type, r rand.Source, values *[]any, n int) {
	ctx.Signature("genSliceValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	// emit random slice index from [0, n)
	*values = append(*values, rand.Intn(n))

	// if slice type is scalar, emit scalar, else recurse
	tyElem := typ.Elem()
	flavorElem := NewFlavor(tyElem.Kind())
	for range n {
		if flavorElem == Scalar {
			a := genRandScalarValue(ctx, tyElem, r)
			*values = append(*values, a)
		} else {
			genScalarValues(ctx, tyElem, r, values)
		}
	}
}


func genStructValues(ctx *EcoContext, typ reflect.Type, r rand.Source, values *[]any) {
	ctx.Signature("genStructValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	field := typ.Field(0)
	fieldType := field.Type
	fieldDeep := DeepType{fieldType}
	fieldFlavor := fieldDeep.Flavor()
	
	// emit first struct field name
	*values = append(*values, field.Name)

	// emit field value if scalar, else recurse
	if fieldFlavor == Scalar {
		a := genRandScalarValue(ctx, fieldType, r)
		*values = append(*values, a)
	} else {
		genScalarValues(ctx, fieldType, r, values)
	}
}


func genSliceDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("genSliceDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()
	
	if n < 1 {
		return
	}

	sb := strings.Builder{}
	w.Write([]byte("[]"))
	writeScalarOrRecurse(ctx, n, r, w)
	fmt.Println(sb.String())
}


func genStructDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("genStructDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()
	
	if n < 1 {
		return
	}

	fieldName := genRandScalarValue(ctx, reflect.TypeFor[string](), r).(string)	
	bufSrc := []byte(fieldName)
	bufDst := make([]byte, len(bufSrc))
	caser := cases.Title(language.English)
	nDst, nSrc, err := caser.Transform(bufDst, bufSrc, true)
	if err != nil {
		panic(err)
	}
	if nDst < len(bufDst) {
		panic("nDst too small")
	}
	if nSrc < len(bufSrc) {
		panic("nSrc too small")
	}
	fieldName = string(bufDst)
	
	sb := strings.Builder{}
	w.Write([]byte(fmt.Sprintf("struct{%s ", fieldName)))
	writeScalarOrRecurse(ctx, n, r, w)
	w.Write([]byte("}"))
	fmt.Println(sb.String())
}


func getRandFlavor(ctx *EcoContext) Flavor {
	ctx.Signature("getRandFlavor")
	ctx.Inc()
	defer ctx.Dec()


	nMax := 3
	switch rand.Intn(nMax) {
	case 0:
		return Map
	case 1:
		return Slice
	case 2:
		return Struct
	default:
		panic("inconthievable!")
	}
}


func getRandScalar(ctx *EcoContext) reflect.Kind {
	ctx.Signature("getRandScalar")
	ctx.Inc()
	defer ctx.Dec()
	
	nMax := int(reflect.UnsafePointer)
	for {
		n := rand.Intn(nMax)
		knd := reflect.Kind(n)
		switch knd {
		case reflect.Bool,
			reflect.Int,
			reflect.String:
			return knd
		default:
			continue
		}
	}
}


func writeScalarOrRecurse(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("writeScalarOrRecurse", n)
	ctx.Inc()
	defer ctx.Dec()
	
	if n == 1 {
		w.Write([]byte(getRandScalar(ctx, ).String()))
	} else {
		genDeclaration(ctx, n-1, r, w)
	}
}

package common

import (
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strings"
	"time"

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

// .keyName
func (t StructType) writeKeyRef(values []any, w io.Writer) {
	fmt.Fprintf(w, ".%s", values[0])
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

func GenDeclaration(ctx *EcoContext, n int, r rand.Source) string {
	ctx.Signature("genDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return ""
	}

	if r == nil {
		r = rand.NewSource(time.Now().UTC().UnixNano())
	}

	flavor := getRandFlavor(ctx)
	var decl string
	switch flavor {
	case Map:
		decl = genMapDeclaration(ctx, n, r)
	case Slice:
		decl = genSliceDeclaration(ctx, n, r)
	case Struct:
		decl = genStructDeclaration(ctx, n, r)
	default:
		panic(fmt.Errorf("How'd I get a flavor of %s???", flavor))
	}

	return decl
}

func GenScalarValues(ctx *EcoContext, typ reflect.Type, r rand.Source, values *[]any) {
	ctx.Signature("genScalarValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	flavor := NewDeepType(typ).Flavor()
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

func WriteDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("genDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return
	}

	if r == nil {
		r = rand.NewSource(time.Now().UTC().UnixNano())
	}

	flavor := getRandFlavor(ctx)
	switch flavor {
	case Map:
		writeMapDeclaration(ctx, n, r, w)
	case Slice:
		writeSliceDeclaration(ctx, n, r, w)
	case Struct:
		writeStructDeclaration(ctx, n, r, w)
	default:
		panic(fmt.Errorf("How'd I get a flavor of %s???", flavor))
	}
}

func genMapKeyString(ctx *EcoContext, r rand.Source) string {
	mapKey := genRandScalarValue(ctx, reflect.TypeFor[string](), r).(string)
	return strings.ToLower(mapKey)
}

func genMapKeyValue(ctx *EcoContext, typ reflect.Type, r rand.Source) any {
	var a any
	keyType := typ.Key()
	keyDeep := DeepType{keyType}
	keyFlavor := keyDeep.Flavor()
	if keyFlavor != Scalar {
		panic("inconthievalble!")
	}
	if keyType.Kind() == reflect.String {
		a = genMapKeyString(ctx, r)
	} else {
		a = genRandScalarValue(ctx, keyType, r)
	}

	return a
}

func genMapValues(ctx *EcoContext, typ reflect.Type, r rand.Source, values *[]any, n int) {
	ctx.Signature("genMapValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	for range n {
		// emit random key value
		ctx.Commentf("generating %d value(s) ...", n)
		a := genMapKeyValue(ctx, typ, r)
		ctx.Infof("value=%v", a)
		*values = append(*values, a)

		// emit value if scalar, else recurse
		if isElemScalar(typ) {
			a := genRandScalarValue(ctx, typ.Elem(), r)
			*values = append(*values, a)
		} else {
			GenScalarValues(ctx, typ.Elem(), r, values)
		}
	}
}

func genMapDeclaration(ctx *EcoContext, n int, r rand.Source) string {
	ctx.Signature("genMapDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("map[")
	sb.WriteString(getRandScalar(ctx).String())
	sb.WriteString("]")
	sb.WriteString(genScalarOrRecurse(ctx, n, r))
	
	return sb.String()
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

func genScalarOrRecurse(ctx *EcoContext, n int, r rand.Source) string {
	ctx.Signature("genScalarOrRecurse", n)
	ctx.Inc()
	defer ctx.Dec()

	if n == 1 {
		return getRandScalar(ctx).String()
	}
		
	return GenDeclaration(ctx, n-1, r)
}

func genSliceDeclaration(ctx *EcoContext, n int, r rand.Source) string {
	ctx.Signature("genSliceDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("[]")
	return genScalarOrRecurse(ctx, n, r)
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
			GenScalarValues(ctx, tyElem, r, values)
		}
	}
}

func genStructFieldName(ctx *EcoContext, r rand.Source) string {
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

	return fieldName
}

func genStructDeclaration(ctx *EcoContext, n int, r rand.Source) string {
	ctx.Signature("genStructDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return ""
	}

	fieldName := genStructFieldName(ctx, r)
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "struct{%s ", fieldName)
	sb.WriteString(genScalarOrRecurse(ctx, n, r))
	sb.WriteString("}")

	return sb.String()
}

func genStructValues(ctx *EcoContext, typ reflect.Type, r rand.Source, values *[]any) {
	ctx.Signature("genStructValues", typ, len(*values))
	ctx.Inc()
	defer ctx.Dec()

	if typ.NumField() == 0 {
		return
	}

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
		GenScalarValues(ctx, fieldType, r, values)
	}
}

// @testme
func isElemScalar(typ reflect.Type) bool {
	elemType := typ.Elem()
	elemDeep := DeepType{elemType}
	elemFlavor := elemDeep.Flavor()
	return elemFlavor == Scalar
}

func writeMapDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("writeMapDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return
	}

	s := genMapDeclaration(ctx, n, r)
	w.Write([]byte(s))
}

func writeScalarOrRecurse(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("writeScalarOrRecurse", n)
	ctx.Inc()
	defer ctx.Dec()

	s := genScalarOrRecurse(ctx, n, r)
	w.Write([]byte(s))
}

func writeSliceDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("writeSliceDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return
	}

	s := genSliceDeclaration(ctx, n, r)
	w.Write([]byte(s))
}

func writeStructDeclaration(ctx *EcoContext, n int, r rand.Source, w io.Writer) {
	ctx.Signature("writeStructDeclaration", n)
	ctx.Inc()
	defer ctx.Dec()

	if n < 1 {
		return
	}

	s := genStructDeclaration(ctx, n, r)
	w.Write([]byte(s))
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

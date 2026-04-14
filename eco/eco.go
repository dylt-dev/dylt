package eco

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"path"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/dylt-dev/dylt/color"
	"github.com/dylt-dev/dylt/common"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
)

type Decoder interface {
	Decode(*ecoContext, []*mvccpb.KeyValue, string, reflect.Value) error
}
type MapDecoder struct{}
type MainDecoder struct{}
type ScalarDecoder[U any] struct{}
type SliceDecoder struct{}
type StructDecoder struct{}
type DecoderMap map[reflect.Kind]Decoder

type MapData map[string][]byte
type SliceData map[int][]byte


var decoderMap DecoderMap = DecoderMap{
	reflect.Bool:    &ScalarDecoder[bool]{},
	reflect.Int:     &ScalarDecoder[int]{},
	reflect.Int8:    &ScalarDecoder[int8]{},
	reflect.Int16:   &ScalarDecoder[int16]{},
	reflect.Int32:   &ScalarDecoder[int32]{},
	reflect.Int64:   &ScalarDecoder[int64]{},
	reflect.Uint:    &ScalarDecoder[uint]{},
	reflect.Uint8:   &ScalarDecoder[uint8]{},
	reflect.Uint16:  &ScalarDecoder[uint16]{},
	reflect.Uint32:  &ScalarDecoder[uint32]{},
	reflect.Uint64:  &ScalarDecoder[uint64]{},
	reflect.Float32: &ScalarDecoder[float32]{},
	reflect.Float64: &ScalarDecoder[float64]{},
	reflect.String:  &ScalarDecoder[string]{},
	reflect.Array:   &SliceDecoder{},
	reflect.Slice:   &SliceDecoder{},
	reflect.Map:     &MapDecoder{},
	reflect.Struct:  &StructDecoder{},
}

func (d *MainDecoder) Decode(ctx *ecoContext, kvs []*mvccpb.KeyValue, key string, rv reflect.Value) error {
	// Get the decoder from the decoder map, if it exists
	pKind, err := getUnderlyingPointerKind(rv)
	if err != nil {
		return err
	}
	decoder, is := decoderMap[pKind]
	if !is {
		return fmt.Errorf("Unsupported pointer type (kind=%s)", pKind.String())
	}

	return decoder.Decode(ctx, kvs, key, rv)
}

// Decode the kvs at the key into a map
// The map is specified by a pointer-to-a-pointer-to-map (ppm). The pointer to
// the map (pm) is assumed to be nil, and it is this function's job to allocate
// the map. and then assign the address of the allocated map to the ppm. This is
// how a function can allocation a value and 'return' it via an incoming
// parameter.
//
// @note it might not make a lot of sense to deal with double indirection just
// to support allocating a new value to an incoming parameter. It's what
// json.Unmarshal() does but that's because json.Unmarshal() also supports
// unmarshalling into an existing data structure, as well as allocation. If a
// a function is always allocating and unmarshalling into a new object, it might
// make sense to just return it.
//
// ctx	Context for logging+etcd client
// kvs  key-value pairs which comprise the data
// key  key that prefixes all map keys
// rv   reflection pointer-to-pointer-to-map
func (d *MapDecoder) Decode(ctx *ecoContext, kvs []*mvccpb.KeyValue, key string, ppMap reflect.Value) error {
	// get the reflect.Type for the map to allocate
	typ, err := getUnderlyingMapType(ppMap.Type())
	if err != nil {
		return err
	}

	// allocate the new map + save the value type
	rMap := reflect.MakeMap(typ)
	typValue := rMap.Type().Elem()

	// get the map data from the kvs+key
	mapData := getMapData(kvs, key)

	// populate the new map with the data
	for k, v := range mapData {
		ctx.logger.Infof("Decoding %s ...", k)	
		// create a new map item
		pnew := reflect.New(typValue)

		// get the address of the new element and unmarshal the mapData value
		addr := pnew.Elem().Addr()
		i := addr.Interface()
		err := json.Unmarshal(v, i)
		if err != nil {
			return err
		}

		// Create reflect.Value for mapData key and add key+val to new map 
		rk := reflect.ValueOf(k)
		rMap.SetMapIndex(rk, pnew.Elem())
	}

	// Create a new map pointer and assign the new map to it
	pMap := reflect.New(typ)
	pMap.Elem().Set(rMap)

	// assign the new map to the rv
	ppMap.Elem().Set(pMap)

	// done :)
	return nil
}


func (d *ScalarDecoder[U]) Decode(ctx *ecoContext, kvs []*mvccpb.KeyValue, key string, rv reflect.Value) error {
	data := kvs[0].Value
	ctx.logger.Infof("data=%#v", data)
	i := rv.Interface()
	err := json.Unmarshal(data, i)
	return err
}

func (d *SliceDecoder) Decode(ctx *ecoContext, kvs []*mvccpb.KeyValue, key string, rv reflect.Value) error {
	sliceData := getSliceData(kvs, key)
	maxIndex := sliceData.MaxIndex()
	typSlice := rv.Type().Elem().Elem()
	len := maxIndex + 1
	cap := maxIndex + 1
	rvSlice := reflect.MakeSlice(typSlice, len, cap)

	// Unmarshal all the elements
	for i, data := range sliceData {
		// Get a pointer to the slice element at the specified index
		el := rvSlice.Index(i)
		addr := el.Addr()
		pEl := addr.Interface()

		// Unmarshal the specified data into the element pointer
		err := json.Unmarshal(data, pEl)
		if err != nil {
			return err
		}
	}

	// Make a slice pointer + assign the new slice to the pointer's Elem()
	rvNew := reflect.New(typSlice)
	rvNew.Elem().Set(rvSlice)

	// Assign the new slice pointer to the incoming rv
	rv.Elem().Set(rvNew)

	return nil
}

func (d *StructDecoder) Decode(ctx *ecoContext, kvs []*mvccpb.KeyValue, key string, rv reflect.Value) error {
	return nil
}

func Encode(ctx *ecoContext, key string, i any) ([]etcd.Op, error) {
	ctx.logger.signature("Encode", key, reflect.TypeOf(i))
	if _, ok := i.(reflect.Value); ok {
		ctx.logger.info("arg i is of type reflect.Value; did you mean to call i.Interface()?")
	}
	ctx.inc()
	defer ctx.dec()

	var ty reflect.Type = reflect.TypeOf(i)
	// var _ reflect.Value = reflect.ValueOf(i)
	var ops = []etcd.Op{}
	var kind reflect.Kind = ty.Kind()
	var val reflect.Value = reflect.ValueOf(i)
	var err error

	// ctx.println(color.Styledstring("Check object type to confirm it can be encoded").Fg(color.X11.CornflowerBlue))
	ctx.logger.comment("Check object type to confirm it can be encoded")
	ctx.logger.Infof("Switching on kind=%s ...", kind.String())
	switch kind {

	// simple case for simple types
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		ops, err = encodeScalar(ctx, key, val)

	case reflect.Array,
		reflect.Slice:
		ops, err = encodeSlice(ctx, key, val)

	case reflect.Map:
		ops, err = encodeMap(ctx, key, val)

	case reflect.Struct:
		ops, err = encodeStruct(ctx, key, val)

	default:
		err = fmt.Errorf("unsupported reflection kind (%s)", kind.String())
	}

	if err != nil {
		return nil, err
	}
	return ops, nil
}

type ecoContext struct {
	context.Context
	depth  int
	logger *ecoLogger
}

func newEcoContext(w io.Writer) *ecoContext {
	var ctx = &ecoContext{
		Context: context.Background(),
		depth:   0,
	}
	ctx.logger = newEcoLogger(w, ctx)

	return ctx
}

func (ctx *ecoContext) dec() *ecoContext {
	ctx.depth--
	return ctx
}

func (ctx *ecoContext) Depth() int {
	return ctx.depth
}

func (ctx *ecoContext) inc() *ecoContext {
	ctx.depth++
	return ctx
}

// func (ctx *ecoContext) indent() string {
// 	const tab = "  "
// 	return strings.Repeat(tab, ctx.level)
// }

// func (ctx *ecoContext) printf(format string, a ...any) (int, error) {
// 	format = fmt.Sprintf("%s%s", ctx.indent(), format)
// 	return fmt.Printf(format, a...)
// }

// func (ctx *ecoContext) println(a ...any) (int, error) {
// 	args := fmt.Sprintln(a...)
// 	return fmt.Printf("%s%s", ctx.indent(), args)
// }

type kind uint

const (
	Invalid kind = iota
	InvalidSlice
	Bool
	Number
	String
	SimpleArray
	SimpleInterface
	SimpleMap
	SimplePointer
	SimpleSlice
	SimpleStruct
)

func (k kind) IsScalar() bool {
	switch k {
	case Bool,
		Number,
		String:
		return true
	default:
		return false
	}
}

func (k kind) IsSimple() bool {
	switch k {
	case SimpleArray,
		SimpleInterface,
		SimpleMap,
		SimplePointer,
		SimpleSlice,
		SimpleStruct:
		return true
	default:
		return false
	}
}

func (k kind) String() string {
	switch k {
	case Invalid:
		return "Invalid"
	case InvalidSlice:
		return "InvalidSlice"
	case Bool:
		return "Bool"
	case Number:
		return "Number"
	case String:
		return "String"
	case SimpleArray:
		return "SimpleArray"
	case SimpleInterface:
		return "SimpleInterface"
	case SimpleMap:
		return "SimpleMap"
	case SimplePointer:
		return "SimplePointer"
	case SimpleSlice:
		return "SimpleSlice"
	case SimpleStruct:
		return "SimpleStruct"
	default:
		return fmt.Sprintf("Unknown kind :%d", k)
	}
}

func arrayKind(ctx *ecoContext, ty reflect.Type) kind {
	ctx.logger.signature("arrayKind", ty)
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Array {
		ctx.logger.info("type is not a array; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	ctx.logger.Infof("Checking element type (%s) ... ", common.FullTypeName(tyElem))
	if isTypeScalar(ty.Elem()) {
		ctx.logger.Infof("element type (%s) is scalar; returning SimpleArray", common.FullTypeName(tyElem))
		return SimpleArray
	}
	ctx.logger.info("conditions were not met; returning Invalid")
	return Invalid
}

func decodeScalar(ctx *ecoContext, key string) (etcd.Op, error) {
	ctx.logger.signature("decodeScalar", key)
	ctx.inc()
	defer ctx.dec()

	op := etcd.OpGet(key)
	return op, nil
}

func encodeMap(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.logger.signature("encodeMap", key, val.Type())
	ctx.inc()
	defer ctx.dec()

	ty := val.Type()
	ctx.logger.Infof("Confirming type (%s) is SimpleMap ...", common.FullTypeName(ty))
	if getTypeKind(ctx, ty) != SimpleMap {
		ctx.logger.comment("incorrect.")
		return nil, fmt.Errorf("expecting SimpleMap; got %s", common.FullTypeName(ty))
	}

	ctx.logger.comment("confirmed.")
	ctx.logger.info("Encoding keys and values ...")
	var ops = []etcd.Op{}
	mapIter := val.MapRange()
	for mapIter.Next() {
		miKey := fmt.Sprintf("%v", mapIter.Key().Interface())
		elKey := filepath.Join(key, string(miKey))
		elVal := mapIter.Value()
		elOps, err := Encode(ctx, elKey, elVal.Interface())
		if err != nil {
			return nil, err
		}
		ops = append(ops, elOps...)
	}

	return ops, nil
}

func encodeScalar(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.logger.signature("encodeDefault", key, val.Type())
	ctx.inc()
	defer ctx.dec()

	i := val.Interface()
	j, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	opPut := etcd.OpPut(key, string(j))

	return []etcd.Op{opPut}, nil
}

func encodeSlice(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.logger.signature("encodeSlice", key, val.Type())
	ctx.inc()
	defer ctx.dec()

	ty := val.Type()
	if getTypeKind(ctx, ty) != SimpleSlice {
		return nil, fmt.Errorf("expecting SimpleSlice; got %s", common.FullTypeName(ty))
	}

	n := val.Len()
	ops := []etcd.Op{}
	for i := range n {
		el := val.Index(i)
		elKey := path.Join(key, strconv.Itoa(i))
		op, err := Encode(ctx, elKey, el.Interface())
		if err != nil {
			return nil, err
		}
		ops = slices.Concat(ops, op)
	}
	// j, err := json.Marshal(val.Interface())
	// if err != nil {
	// 	return nil, err
	// }
	// op := etcd.OpPut(key, string(j))

	return ops, nil
}

func encodeString(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.logger.signature("encodeString", key, val.Type())
	ctx.inc()
	defer ctx.dec()

	s := val.String()
	opPut := etcd.OpPut(key, string(s))

	return []etcd.Op{opPut}, nil
}

func encodeStruct(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.logger.signature("encodeStruct", key, val.Type())
	ctx.inc()
	defer ctx.dec()

	ty := val.Type()
	ctx.logger.commentf("Confirming type (%s) is SimpleStruct ...", common.FullTypeName(ty))
	if getTypeKind(ctx, ty) != SimpleStruct {
		ctx.logger.comment("incorrect.")
		return nil, fmt.Errorf("expecting SimpleStruct; got %s", common.FullTypeName(ty))
	}

	ctx.logger.info("confirmed.")
	ctx.logger.info("Encoding fields ...")
	var ops = []etcd.Op{}
	for i := range ty.NumField() {
		sf := ty.Field(i)
		sfName := getFieldKey(sf)
		sfVal := val.Field(i)
		sfKey := filepath.Join(key, sfName)
		sfOps, err := Encode(ctx, sfKey, sfVal.Interface())
		if err != nil {
			return nil, err
		}
		ops = append(ops, sfOps...)
	}

	return ops, nil
}

func fieldNameMap(i any) (map[string]reflect.Value, error) {
	var tyElem reflect.Type
	var valElem reflect.Value
	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("expecting pointer; got %s", common.FullTypeName(ty))
	}
	tyElem = ty.Elem()
	if tyElem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expecting pointer to struct; got %s", tyElem.Kind().String())
	}
	val := reflect.ValueOf(i)
	valElem = val.Elem()

	fieldNameMap := map[string]reflect.Value{}
	for i := range tyElem.NumField() {
		tyField := tyElem.Field(i)
		fieldKey := getFieldKey(tyField)
		fieldName := tyField.Name
		valField := valElem.FieldByName(fieldName)
		fieldNameMap[fieldKey] = valField
	}

	return fieldNameMap, nil
}

func getFieldKey(sf reflect.StructField) string {
	tagValue, ok := sf.Tag.Lookup("eco")
	var fieldName string
	if ok {
		fieldName = tagValue
	} else {
		fieldName = sf.Name
	}

	return fieldName
}

func getFieldValue(val reflect.Value) (string, error) {
	var s string
	kind := val.Type().Kind()
	switch kind {
	case reflect.Bool:
		s = strconv.FormatBool(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s = strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32:
		s = strconv.FormatFloat(val.Float(), 'f', -1, 32)
	case reflect.Float64:
		s = strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case reflect.String:
		s = string(val.String())
	default:
		return "", fmt.Errorf("unsupported kind: %s", kind.String())
	}

	return s, nil
}


// All the data that comprises the map specified by the parentKey. This is the
// subset of kvs with keys that match {parentKey}/mapKey, where mapKey is a
// single key segment. 
func getMapData (kvs []*mvccpb.KeyValue, parentKey string) MapData {
	mapData := MapData{}

	for _, kv := range kvs {
		key := string(kv.Key)
		itemKey, is := getMapItemKey(parentKey, key)
		if is {
			mapData[itemKey] = kv.Value
		}
	}

	return mapData
}


// Get the index portion of this slice key, which is the trailing portion
// of the key after the final slash, other than an optional trailing slash.
// The trailing portion must be an integer.
//
// This function also validates that the key is actually a slice key index.
// If it isn't, the bool return value will be false.
func getMapItemKey(key string, subkey string) (string, bool) {
	if !strings.HasPrefix(subkey, key) {
		return "", false
	}

	// Trim the last character if its a slash
	lastChar := subkey[len(subkey)-1:]
	if lastChar == "/" {
		subkey = subkey[0 : len(subkey)-1]
	}

	// Confirm key contains at least one slash
	iLastSlash := strings.LastIndex(subkey, "/")
	if iLastSlash == -1 {
		return "", false
	}

	// Confirm the part of the subkey preceding the last slash marches the map key
	prefix := subkey[:iLastSlash]
	if prefix != key {
		return "", false
	}

	// Use the piece following the trailing slash as
	itemKey := subkey[iLastSlash+1:]

	return itemKey, true
}

// All the data that comprises the slice. Uses all keys of the
// form /key/N, where N is any integer. The data is returned as a map form,
// where the map keys are the N index values for each element
//
// @note I think Im allowing negative indexes
func getSliceData(kvs []*mvccpb.KeyValue, sliceKey string) SliceData {
	sliceData := SliceData{}
	
	for _, kv := range kvs {
		key := string(kv.Key)
		itemKey, is := getSliceItemKey(sliceKey, key)
		if is {
			sliceData[itemKey] = kv.Value
		}
	}

	return sliceData
}

func getKind(ctx *ecoContext, i any) kind {
	ctx.logger.signature("getKind", reflect.TypeOf(i))
	ctx.inc()
	defer ctx.dec()

	ty := reflect.TypeOf(i)
	if common.FullTypeName(ty) == "reflect.Type" {
		ctx.logger.Warn("Warning - GetKind() called with reflect.Type(). Did you mean GetTypeKind()?")
	}

	return getTypeKind(ctx, reflect.TypeOf(i))
}

// Get the index portion of this slice key, which is the trailing portion
// of the key after the final slash, other than an optional trailing slash.
// The trailing portion must be an integer.
//
// This function also validates that the key is actually a slice key index.
// If it isn't, the bool return value will be false.
func getSliceItemKey(key string, subkey string) (int, bool) {
	if !strings.HasPrefix(subkey, key) {
		return -1, false
	}
	// Trim the last character if its a slash
	lastChar := subkey[len(subkey)-1:]
	if lastChar == "/" {
		subkey = subkey[0 : len(subkey)-1]
	}

	// Confirm key contains at least one slash
	iLastSlash := strings.LastIndex(subkey, "/")
	if iLastSlash == -1 {
		return -1, false
	}
	sIndex := subkey[iLastSlash+1:]
	index, err := strconv.Atoi(sIndex)
	if err != nil || index < 0 {
		return -1, false
	}
	return index, true
}

func (m SliceData) MaxIndex() int {
	maxIndex := 0
	for key := range m {
		if key > maxIndex {
			maxIndex = key
		}

	}
	return maxIndex
}

func getTypeKind(ctx *ecoContext, ty reflect.Type) kind {
	reflectKind := ty.Kind()
	// fmt.Printf("ty=%s reflectKind=%s\n", fullTypeName(ty), reflectKind.String())
	switch reflectKind {
	case reflect.Bool:
		return Bool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32,
		reflect.Float64:
		return Number
	case reflect.String:
		return String
	case reflect.Array:
		return arrayKind(ctx, ty)
	case reflect.Interface:
		return interfaceKind(ctx, ty)
	case reflect.Map:
		return mapKind(ctx, ty)
	case reflect.Pointer:
		return pointerKind(ctx, ty)
	case reflect.Slice:
		return sliceKind(ctx, ty)
	case reflect.Struct:
		return structKind(ctx, ty)
	case reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.Func, reflect.UnsafePointer:
		return Invalid
	default:
		return Invalid
	}
}

func getUnderlyingPointerKind(p any) (reflect.Kind, error) {
	if p == nil {
		return reflect.Invalid, fmt.Errorf("expecting a pointer or pointer-to-pointer")
	}

	if reflect.TypeOf(p) == reflect.TypeFor[reflect.Value]() {
		v := p.(reflect.Value)
		p = v.Interface()
	}

	var typ reflect.Type
	var knd reflect.Kind

	// Confirm p is a pointer
	typ = reflect.TypeOf(p)
	knd = typ.Kind()
	if knd != reflect.Pointer {
		return reflect.Invalid, fmt.Errorf("expecting a pointer or pointer-to-pointer")
	}

	// If *p is not pointer we're done
	typ = typ.Elem()
	knd = typ.Kind()
	if knd != reflect.Pointer {
		return knd, nil
	}

	// Confirm **p is _not_ a pointer
	typ = typ.Elem()
	knd = typ.Kind()
	if knd == reflect.Pointer {
		return reflect.Invalid, fmt.Errorf("expecting a pointer or pointer-to-pointer")
	}

	return knd, nil
}

// Unmarshalling functions like json.Unmarshaller() sometimes take an argument
// of type any, that can be either a pointer, or a pointer to a pointer. The
// idea is that if the argument is a pointer, then the pointer refers to an
// initialized data structure that the Unmarshalling function is expected to
// populate. If on the other hand the argument is a pointer to a pointer, then
// it is assumed the caller has not initialized a data structure to receive
// the unmarshalled data, and the Unmarshaller is expected to allocate a new
// structure, then dereference the pointer-to-pointer and assign the new
// structure's address. If this is confusing it's because pointer-to-pointer
// scenarios are always confusing to mere mortals, eg your author.
//
// Note in these Unmarshalling scenarios, it is often an error if the argument
// is nil. When the argument is nil there is no way to create a new structure
// and assign its address to the dereferenced argument, because if the argument
// is nil it cannot be dereferenced. This function is unopinionated regarding
// the value of the argument. It only cares about the argument type.
func getUnderlyingSliceType(p any) (reflect.Type, error) {
	// Check if the type is a pointer
	pType := reflect.TypeOf(p)
	pKind := pType.Kind()
	if pKind != reflect.Pointer {
		return pType, fmt.Errorf("Not a pointer (kind=%s)", pKind.String())
	}

	// If not a pointer to a pointer, return the element type
	elType := pType.Elem()
	elKind := elType.Kind()
	if elKind != reflect.Pointer {
		return elType, nil
	}

	// Pointer-to-pointer, so get the element type again
	elType = elType.Elem()
	elKind = elType.Kind()
	if elKind == reflect.Pointer {
		return elType, fmt.Errorf("**p must not be a pointer - that's too deep")
	}

	// We have a valid pointer-to-pointer, so we're done
	return elType, nil
}

func getUnderlyingMapType(ppMapType reflect.Type) (reflect.Type, error) {
	var knd reflect.Kind

	// Confirm ppMapType is a pointer
	knd = ppMapType.Kind()
	if knd != reflect.Pointer {
		return nil, fmt.Errorf("Expecting a pointer-to-pointer-to map (kind=%s)", knd.String())
	}

	// Confirm *ppMapType is a pointer
	knd = ppMapType.Elem().Kind()
	if knd != reflect.Pointer {
		return nil, fmt.Errorf("Expecting a pointer-to-pointer-to map (kind=%s)", knd.String())
	}

	// Confirm **ppMapType is a map
	knd = ppMapType.Elem().Kind()
	if knd != reflect.Pointer {
		return nil, fmt.Errorf("Expecting a pointer-to-pointer-to map (kind=%s)", knd.String())
	}

	// We're good :)
	return ppMapType.Elem().Elem(), nil
}

func interfaceKind(ctx *ecoContext, ty reflect.Type) kind {
	ctx.logger.signature("interfaceKind", ty)
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Interface {
		ctx.logger.info("type is not an interface; returning Invalid")
		return Invalid
	}

	return Invalid
}

func isNormalPointer(p any) bool {
	if p == nil {
		return false
	}

	var typ reflect.Type
	var knd reflect.Kind

	// Confirm p is a pointer
	typ = reflect.TypeOf(p)
	knd = typ.Kind()
	if knd != reflect.Pointer {
		return false
	}

	// Confirm *p is _not_ a pointer
	typ = typ.Elem()
	knd = typ.Kind()
	if knd == reflect.Pointer {
		return false
	}

	return true
}

func isPointerToPointer(p any) bool {
	if p == nil {
		return false
	}

	var typ reflect.Type
	var knd reflect.Kind

	// Confirm p is a pointer
	typ = reflect.TypeOf(p)
	knd = typ.Kind()
	if knd != reflect.Pointer {
		return false
	}

	// Confirm *p is a pointer
	typ = typ.Elem()
	knd = typ.Kind()
	if knd != reflect.Pointer {
		return false
	}

	// Confirm **p is _not_ a pointer
	typ = typ.Elem()
	knd = typ.Kind()
	if knd == reflect.Pointer {
		return false
	}

	return true
}

func isScalar(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

func isTypeScalar(ty reflect.Type) bool {
	return isScalar(ty.Kind())
}

func isValidPointer(p any) bool {
	if p == nil {
		return false
	}

	var (
		typP  reflect.Type
		// typPP reflect.Type
		kndP  reflect.Kind
		// kndPP reflect.Kind
		// valP  reflect.Value
	)
	// Confirm p is a pointer
	typP = reflect.TypeOf(p)
	kndP = typP.Kind()
	if kndP != reflect.Pointer {
		return false
	}

	// // Confirm *p is a pointer
	// typPP = typP.Elem()
	// kndPP = typP.Kind()
	// if kndPP != reflect.Pointer {
	// 	// *p is a pointer, so p is a pointer, not a pointer to pointer.
	// 	// We need to confirm its not nil.
	// 	valP := valueOf(p)

	// }

	// // Confirm **p is _not_ a pointer
	// typ = typ.Elem()
	// knd = typ.Kind()
	// if knd == reflect.Pointer {
	// 	return false
	// }

	return true
}

func mapKind(ctx *ecoContext, ty reflect.Type) kind {
	sig := common.FullTypeName(ty)
	ctx.logger.Infof("%s(%s)", common.Highlight("mapKind"), common.Lowlight(sig))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Map {
		ctx.logger.info("type is not a map; returning Invalid")
		return Invalid
	}

	ctx.logger.Infof("%-70s", common.Lowlight(fmt.Sprintf("checking key (%s) ...", common.FullTypeName(ty.Key()))))
	if !isTypeScalar(ty.Key()) {
		ctx.logger.Infof("%-32s: %-32s; %s", "key", "non-scalar", common.Highlight("returning Invalid"))
		return Invalid
	}
	ctx.logger.Infof("%-16s %-16s; continuing", "key", "scalar")

	tyElem := ty.Elem()
	ctx.logger.info(common.Lowlight(fmt.Sprintf("checking element type (%s) ...", common.FullTypeName(tyElem))))
	ctx.inc()
	kindElem := getTypeKind(ctx, tyElem)
	ctx.dec()
	if isTypeScalar(tyElem) {
		ctx.logger.Infof("%-16s %-16s; %s", "type", "scalar", common.Highlight("returning SimpleMap"))
		return SimpleMap
	}
	ctx.logger.Infof("%-16s %-16s; continuing", common.FullTypeName(tyElem), "not scalar")

	ctx.logger.Infof("%-70s", common.Lowlight(fmt.Sprintf("Checking element kind (%s) ...", kindElem.String())))
	if kindElem == SimpleMap ||
		kindElem == SimpleStruct ||
		kindElem == SimpleSlice {
		ctx.logger.Infof("%s: simple; returning SimpleMap", kindElem.String())
		return SimpleMap
	}
	ctx.logger.info("type: not simple; continuing")

	ctx.logger.info(common.Highlight("conditions were not met; returning Invalid"))
	return Invalid
}

func pointerKind(ctx *ecoContext, ty reflect.Type) kind {
	ctx.logger.signature("pointerKind", ty)
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Pointer {
		ctx.logger.info("type is not a pointer; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	ctx.logger.Infof("Checking pointer type (%s) ... ", common.FullTypeName(tyElem))
	if isTypeScalar(tyElem) {
		ctx.logger.info("pointer type is scalar; returning SimplePointer")
		return SimplePointer
	}

	ctx.logger.info("conditions were not met; returning Invalid")
	return Invalid
}

func sliceKind(ctx *ecoContext, ty reflect.Type) kind {
	ctx.logger.signature("sliceKind", ty)
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Slice {
		ctx.logger.info("type is not a slice; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	kind := getTypeKind(ctx, tyElem)
	ctx.logger.commentf("Checking element type (%s) ... ", common.FullTypeName(tyElem))
	ctx.logger.Appendf("IsScalar(%s) ...", common.FullTypeName(tyElem))
	if kind.IsScalar() {
		ctx.logger.AppendAndFlush(slog.LevelDebug, "true; returning SimpleSlice")
		return SimpleSlice
	} else {
		ctx.logger.Appendf("IsSimple(%s) ...", common.FullTypeName(tyElem))
		if kind.IsSimple() {
			ctx.logger.AppendAndFlush(slog.LevelDebug, "true; returning SimpleSlice")
			return SimpleSlice
		} else {
			ctx.logger.AppendAndFlush(slog.LevelDebug, "false")
		}
	}

	ctx.logger.info("conditions were not met; returning InvalidSlice")
	return InvalidSlice
}

func structKind(ctx *ecoContext, ty reflect.Type) kind {
	// ctx.printf("%s(%s)\n", highlight("structKind"), lowlight(fullTypeName(ty)))
	ctx.logger.signature("structKind", common.FullTypeName(ty))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Struct {
		ctx.logger.info("type is not a struct; returning Invalid")

		return Invalid
	}

	ctx.logger.Infof("%d field(s)", ty.NumField())
	for i := range ty.NumField() {
		sf := ty.Field(i)
		sfType := sf.Type
		ctx.logger.Appendf("%-70s", common.Lowlight(fmt.Sprintf("checking field '%s' (%s) ...", sf.Name, common.FullTypeName(sfType))))
		sfReflectKind := sfType.Kind()

		if isTypeScalar(sfType) {
			ctx.logger.AppendfAndFlush(slog.LevelInfo, "%-16s %-16s; %s", sfType, "scalar", "continuing")
			continue
		}

		if sfReflectKind == reflect.Map && mapKind(ctx, sfType) == SimpleMap {
			ctx.logger.AppendAndFlush(slog.LevelInfo, "field type is SimpleMap; continuing")
			continue
		}

		if sfReflectKind == reflect.Slice && sliceKind(ctx, sfType) == SimpleSlice {
			ctx.logger.AppendAndFlush(slog.LevelInfo, "field type is SimpleSlice; continuing")
			continue
		}

		if sfReflectKind == reflect.Struct && structKind(ctx, sf.Type) == SimpleStruct {
			ctx.logger.AppendAndFlush(slog.LevelInfo, "field type is SimpleStruct; continuing")
			continue
		}

		return Invalid
	}

	ctx.logger.Infof("%s; %s", "All fields passed", common.Highlight("returning SimpleStruct"))
	return SimpleStruct
}

/* log styles */

type Depther interface {
	Depth() int
}

type ecoLogger struct {
	buf []byte
	*slog.Logger
	depther Depther
}

func newEcoLogger(w io.Writer, depther Depther) *ecoLogger {
	options := color.ColorOptions{Level: slog.LevelDebug}
	handler := color.NewColorHandler(w, options)
	return &ecoLogger{
		Logger:  slog.New(handler),
		depther: depther,
		buf:     make([]byte, 200),
	}
}

func (l *ecoLogger) Append(s string) *ecoLogger {
	l.buf = slices.Concat(l.buf, []byte(s))

	return l
}

func (l *ecoLogger) Appendf(sfmt string, args ...any) *ecoLogger {
	s := fmt.Sprintf(sfmt, args...)
	return l.Append(s)
}

func (l *ecoLogger) AppendAndFlush(level slog.Level, s string) {
	l.Append(s)
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.Flush(level)
}

func (l *ecoLogger) AppendfAndFlush(level slog.Level, sfmt string, args ...any) {
	msg := fmt.Sprintf(sfmt, args...)
	l.Append(msg)
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.Flush(level)
}

func (l *ecoLogger) Debugf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Debug(l.indent() + s)
}

func (l *ecoLogger) DebugContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.DebugContext(ctx, l.indent()+s)
}

func (l *ecoLogger) Errorf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Error(l.indent() + s)
}

func (l *ecoLogger) ErrorContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.ErrorContext(ctx, l.indent()+s)
}

func (l *ecoLogger) Infof(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Info(l.indent() + s)
}

func (l *ecoLogger) InfoContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.InfoContext(ctx, l.indent()+s)
}

func (l *ecoLogger) Flush(level slog.Level) {
	l.Logger.Log(context.Background(), level, string(l.buf))
	l.buf = make([]byte, 200)
}

func (l *ecoLogger) Warnf(sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.Warn(l.indent() + s)
}

func (l *ecoLogger) WarnContextf(ctx context.Context, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.WarnContext(ctx, l.indent()+s)
}

func (l *ecoLogger) comment(msg string) {
	l.Logger.Info(l.indent() + string(color.Styledstring(msg).Fg(color.X11.CornflowerBlue)))
}

func (l *ecoLogger) commentf(sFmt string, args ...any) {
	msg := fmt.Sprintf(sFmt, args...)
	l.comment(msg)
}

func (l *ecoLogger) indent() string {
	const tab = "  "
	return strings.Repeat(tab, l.depther.Depth())
}

func (l *ecoLogger) info(s string) {
	l.Logger.Info(l.indent() + s)
}

func (l *ecoLogger) signature(name string, args ...any) {
	sig := createSignature(name, args...)
	l.Logger.Info(l.indent() + sig)
}

func allocateSlice[U any](pslice **[]U, len int, cap int) {
	typ := reflect.TypeFor[[]U]()
	rSlice := reflect.MakeSlice(typ, len, cap)
	slice := rSlice.Interface().([]U)
	*pslice = &slice
}

func createSignature(name string, args ...any) string {
	// highlight, concat, all that good stuff
	sFmt := fmt.Sprintf("%%s(%s)", strings.Repeat("%v, ", len(args)-1)+"%v")
	args2 := make([]any, len(args)+1)
	args2[0] = common.Highlight(name)
	for i, arg := range args {
		typ, is := arg.(reflect.Type)
		var sArg string
		if is {
			sArg = fmt.Sprintf("-%s-", common.FullTypeName(typ))
		} else {
			_, is := arg.(string)
			if is {
				sArg = fmt.Sprintf("\"%s\"", arg)
			} else {
				sArg = fmt.Sprintf("%v", arg)
			}
		}
		args2[i+1] = common.Lowlight(sArg)
	}
	s := fmt.Sprintf(sFmt, args2...)

	return s
}

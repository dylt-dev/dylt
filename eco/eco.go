package eco

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"path"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/dylt-dev/dylt/common"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
)

func Encode_Legacy(ctx *common.EcoContext, key string, i any) ([]etcd.Op, error) {
	ctx.Logger.Signature("Encode", key, reflect.TypeOf(i))
	ctx.Inc()
	defer ctx.Dec()

	if _, ok := i.(reflect.Value); ok {
		ctx.Logger.Info("arg i is of type reflect.Value; did you mean to call i.Interface()?")
	}
	ctx.Inc()
	defer ctx.Dec()

	var ty reflect.Type = reflect.TypeOf(i)
	// var _ reflect.Value = reflect.ValueOf(i)
	var ops = []etcd.Op{}
	var kind reflect.Kind = ty.Kind()
	var val reflect.Value = reflect.ValueOf(i)
	var err error

	// ctx.println(color.Styledstring("Check object type to confirm it can be encoded").Fg(color.X11.CornflowerBlue))
	ctx.Logger.Comment("Check object type to confirm it can be encoded")
	ctx.Logger.Infof("Switching on kind=%s ...", kind.String())
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

func arrayKind(ctx *common.EcoContext, ty reflect.Type) kind {
	ctx.Logger.Signature("arrayKind", ty)
	ctx.Inc()
	defer ctx.Dec()

	if ty.Kind() != reflect.Array {
		ctx.Logger.Info("type is not a array; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	ctx.Logger.Infof("Checking element type (%s) ... ", common.FullTypeName(tyElem))
	if common.IsTypeScalar(ty.Elem()) {
		ctx.Logger.Infof("element type (%s) is scalar; returning SimpleArray", common.FullTypeName(tyElem))
		return SimpleArray
	}
	ctx.Logger.Info("conditions were not met; returning Invalid")
	return Invalid
}


func encodeMap(ctx *common.EcoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.Logger.Signature("encodeMap", key, val.Type())
	ctx.Inc()
	defer ctx.Dec()

	ty := val.Type()
	ctx.Logger.Infof("Confirming type (%s) is SimpleMap ...", common.FullTypeName(ty))
	if getTypeKind(ctx, ty) != SimpleMap {
		ctx.Logger.Comment("incorrect.")
		return nil, fmt.Errorf("expecting SimpleMap; got %s", common.FullTypeName(ty))
	}

	ctx.Logger.Comment("confirmed.")
	ctx.Logger.Info("Encoding keys and values ...")
	var ops = []etcd.Op{}
	mapIter := val.MapRange()
	for mapIter.Next() {
		miKey := fmt.Sprintf("%v", mapIter.Key().Interface())
		elKey := filepath.Join(key, string(miKey))
		elVal := mapIter.Value()
		elOps, err := Encode_Legacy(ctx, elKey, elVal.Interface())
		if err != nil {
			return nil, err
		}
		ops = append(ops, elOps...)
	}

	return ops, nil
}

func encodeScalar(ctx *common.EcoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.Logger.Signature("encodeDefault", key, val.Type())
	ctx.Inc()
	defer ctx.Dec()

	i := val.Interface()
	j, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	opPut := etcd.OpPut(key, string(j))

	return []etcd.Op{opPut}, nil
}

func encodeSlice(ctx *common.EcoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.Logger.Signature("encodeSlice", key, val.Type())
	ctx.Inc()
	defer ctx.Dec()

	ty := val.Type()
	if getTypeKind(ctx, ty) != SimpleSlice {
		return nil, fmt.Errorf("expecting SimpleSlice; got %s", common.FullTypeName(ty))
	}

	n := val.Len()
	ops := []etcd.Op{}
	for i := range n {
		el := val.Index(i)
		elKey := path.Join(key, strconv.Itoa(i))
		op, err := Encode_Legacy(ctx, elKey, el.Interface())
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

// func encodeString(ctx *common.EcoContext, key string, val reflect.Value) ([]etcd.Op, error) {
// 	ctx.Logger.Signature("encodeString", key, val.Type())
// 	ctx.Inc()
// 	defer ctx.Dec()

// 	s := val.String()
// 	opPut := etcd.OpPut(key, string(s))

// 	return []etcd.Op{opPut}, nil
// }

func encodeStruct(ctx *common.EcoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.Logger.Signature("encodeStruct", key, val.Type())
	ctx.Inc()
	defer ctx.Dec()

	ty := val.Type()
	ctx.Logger.Commentf("Confirming type (%s) is SimpleStruct ...", common.FullTypeName(ty))
	if getTypeKind(ctx, ty) != SimpleStruct {
		ctx.Logger.Comment("incorrect.")
		return nil, fmt.Errorf("expecting SimpleStruct; got %s", common.FullTypeName(ty))
	}

	ctx.Logger.Info("confirmed.")
	ctx.Logger.Info("Encoding fields ...")
	var ops = []etcd.Op{}
	for i := range ty.NumField() {
		sf := ty.Field(i)
		sfName := getFieldKey(sf)
		sfVal := val.Field(i)
		sfKey := filepath.Join(key, sfName)
		sfOps, err := Encode_Legacy(ctx, sfKey, sfVal.Interface())
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
func getMapData(kvs []*mvccpb.KeyValue, parentKey string) DecoderMapData {
	mapData := DecoderMapData{}

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
func getSliceData(kvs []*mvccpb.KeyValue, sliceKey string) DecoderSliceData {
	sliceData := DecoderSliceData{}

	for _, kv := range kvs {
		key := string(kv.Key)
		itemKey, is := getSliceItemKey(sliceKey, key)
		if is {
			sliceData[itemKey] = kv.Value
		}
	}

	return sliceData
}

func getKind(ctx *common.EcoContext, i any) kind {
	ctx.Logger.Signature("getKind", reflect.TypeOf(i))
	ctx.Inc()
	defer ctx.Dec()

	ty := reflect.TypeOf(i)
	if common.FullTypeName(ty) == "reflect.Type" {
		ctx.Logger.Warn("Warning - GetKind() called with reflect.Type(). Did you mean GetTypeKind()?")
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

func (m DecoderSliceData) MaxIndex() int {
	maxIndex := 0
	for key := range m {
		if key > maxIndex {
			maxIndex = key
		}

	}

	return maxIndex
}

func getEtcdKvs(ctx *common.EcoContext, cli *EtcdClient, rootKey string) ([]*mvccpb.KeyValue, error) {
	// Create an Op to get all keys by prefix
	op := etcd.OpGet(rootKey, etcd.WithPrefix())

	// Get all keys in a single Txn Commit
	txn := cli.Txn(ctx)
	resp, err := txn.Then(op).Commit()
	if err != nil {
		return nil, err
	}

	// Treat response as a ResponseRange and get the kvs
	if !resp.Succeeded {
		return nil, fmt.Errorf("Bad. Bad etcd.")
	}
	if len(resp.Responses) != 1 {
		return nil, fmt.Errorf("Expected 1 response, not %d", len(resp.Responses))
	}
	rangeResp := resp.Responses[0].GetResponseRange()
	if rangeResp == nil {
		buf, ints := resp.Responses[0].Descriptor()
		return nil, fmt.Errorf("Response was not a range response (%s, %#v)", string(buf), ints)
	}

	return rangeResp.Kvs, nil
}

func getTypeKind(ctx *common.EcoContext, ty reflect.Type) kind {
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
	// Check if the type is a reflect.Value - if so, use the Interface()
	rv, is := p.(reflect.Value)
	if is {
		p = rv.Interface()
		common.Logger.Info("Incoming pointer is reflect.Value -- using Interface()")
	}

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

func interfaceKind(ctx *common.EcoContext, ty reflect.Type) kind {
	ctx.Logger.Signature("interfaceKind", ty)
	ctx.Inc()
	defer ctx.Dec()

	if ty.Kind() != reflect.Interface {
		ctx.Logger.Info("type is not an interface; returning Invalid")
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

// A valid pointer is any pointer that is not nil, but points
// to an actual variable. This target variable might be a pointer,
// which is perfectly valid, even if the targeted pointer is nil.
func isValidPointer(i any) bool {
	if i == nil {
		return false
	}

	// Check for reflect.Value
	v, is := i.(reflect.Value)
	if is {
		// Check for nil
		if v.IsNil() {
			return false
		}
		i = v.Interface()
	}

	// Check that i is a (non nil) pointer
	typP := reflect.TypeOf(i)
	kndP := typP.Kind()
	if kndP != reflect.Pointer {
		return false
	}

	return true
}


// @note testme
func keyFromValues (values []any) KeyString {
	if len(values) == 0 {
		return ""
	}

	sb := strings.Builder{}
	// Skip last element - the value is not part of the key
	for i := 0; i < len(values)-1; i++ {
		sb.WriteString("/")
		sb.WriteString(fmt.Sprint(values[i]))
	}
	key := KeyString(sb.String())

	return key
}


func mapKind(ctx *common.EcoContext, ty reflect.Type) kind {
	sig := common.FullTypeName(ty)
	ctx.Logger.Infof("%s(%s)", common.Highlight("mapKind"), common.Lowlight(sig))
	ctx.Inc()
	defer ctx.Dec()

	if ty.Kind() != reflect.Map {
		ctx.Logger.Info("type is not a map; returning Invalid")
		return Invalid
	}

	ctx.Logger.Infof("%-70s", common.Lowlight(fmt.Sprintf("checking key (%s) ...", common.FullTypeName(ty.Key()))))
	if !common.IsTypeScalar(ty.Key()) {
		ctx.Logger.Infof("%-32s: %-32s; %s", "key", "non-scalar", common.Highlight("returning Invalid"))
		return Invalid
	}
	ctx.Logger.Infof("%-16s %-16s; continuing", "key", "scalar")

	tyElem := ty.Elem()
	ctx.Logger.Info(common.Lowlight(fmt.Sprintf("checking element type (%s) ...", common.FullTypeName(tyElem))))
	ctx.Inc()
	kindElem := getTypeKind(ctx, tyElem)
	ctx.Dec()
	if common.IsTypeScalar(tyElem) {
		ctx.Logger.Infof("%-16s %-16s; %s", "type", "scalar", common.Highlight("returning SimpleMap"))
		return SimpleMap
	}
	ctx.Logger.Infof("%-16s %-16s; continuing", common.FullTypeName(tyElem), "not scalar")

	ctx.Logger.Infof("%-70s", common.Lowlight(fmt.Sprintf("Checking element kind (%s) ...", kindElem.String())))
	if kindElem == SimpleMap ||
		kindElem == SimpleStruct ||
		kindElem == SimpleSlice {
		ctx.Logger.Infof("%s: simple; returning SimpleMap", kindElem.String())
		return SimpleMap
	}
	ctx.Logger.Info("type: not simple; continuing")

	ctx.Logger.Info(common.Highlight("conditions were not met; returning Invalid"))
	return Invalid
}

func pointerKind(ctx *common.EcoContext, ty reflect.Type) kind {
	ctx.Logger.Signature("pointerKind", ty)
	ctx.Inc()
	defer ctx.Dec()

	if ty.Kind() != reflect.Pointer {
		ctx.Logger.Info("type is not a pointer; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	ctx.Logger.Infof("Checking pointer type (%s) ... ", common.FullTypeName(tyElem))
	if common.IsTypeScalar(tyElem) {
		ctx.Logger.Info("pointer type is scalar; returning SimplePointer")
		return SimplePointer
	}

	ctx.Logger.Info("conditions were not met; returning Invalid")
	return Invalid
}

func sliceKind(ctx *common.EcoContext, ty reflect.Type) kind {
	ctx.Logger.Signature("sliceKind", ty)
	ctx.Inc()
	defer ctx.Dec()

	if ty.Kind() != reflect.Slice {
		ctx.Logger.Info("type is not a slice; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	kind := getTypeKind(ctx, tyElem)
	ctx.Logger.Commentf("Checking element type (%s) ... ", common.FullTypeName(tyElem))
	ctx.Logger.Appendf("IsScalar(%s) ...", common.FullTypeName(tyElem))
	if kind.IsScalar() {
		ctx.Logger.AppendAndFlush(slog.LevelDebug, "true; returning SimpleSlice")
		return SimpleSlice
	} else {
		ctx.Logger.Appendf("IsSimple(%s) ...", common.FullTypeName(tyElem))
		if kind.IsSimple() {
			ctx.Logger.AppendAndFlush(slog.LevelDebug, "true; returning SimpleSlice")
			return SimpleSlice
		} else {
			ctx.Logger.AppendAndFlush(slog.LevelDebug, "false")
		}
	}

	ctx.Logger.Info("conditions were not met; returning InvalidSlice")
	return InvalidSlice
}

func structKind(ctx *common.EcoContext, ty reflect.Type) kind {
	// ctx.printf("%s(%s)\n", highlight("structKind"), lowlight(fullTypeName(ty)))
	ctx.Logger.Signature("structKind", common.FullTypeName(ty))
	ctx.Inc()
	defer ctx.Dec()

	if ty.Kind() != reflect.Struct {
		ctx.Logger.Info("type is not a struct; returning Invalid")

		return Invalid
	}

	ctx.Logger.Infof("%d field(s)", ty.NumField())
	for i := range ty.NumField() {
		sf := ty.Field(i)
		sfType := sf.Type
		ctx.Logger.Appendf("%-70s", common.Lowlight(fmt.Sprintf("checking field '%s' (%s) ...", sf.Name, common.FullTypeName(sfType))))
		sfReflectKind := sfType.Kind()

		if common.IsTypeScalar(sfType) {
			ctx.Logger.AppendfAndFlush(slog.LevelInfo, "%-16s %-16s; %s", sfType, "scalar", "continuing")
			continue
		}

		if sfReflectKind == reflect.Map && mapKind(ctx, sfType) == SimpleMap {
			ctx.Logger.AppendAndFlush(slog.LevelInfo, "field type is SimpleMap; continuing")
			continue
		}

		if sfReflectKind == reflect.Slice && sliceKind(ctx, sfType) == SimpleSlice {
			ctx.Logger.AppendAndFlush(slog.LevelInfo, "field type is SimpleSlice; continuing")
			continue
		}

		if sfReflectKind == reflect.Struct && structKind(ctx, sf.Type) == SimpleStruct {
			ctx.Logger.AppendAndFlush(slog.LevelInfo, "field type is SimpleStruct; continuing")
			continue
		}

		return Invalid
	}

	ctx.Logger.Infof("%s; %s", "All fields passed", common.Highlight("returning SimpleStruct"))
	return SimpleStruct
}

/* log styles */

func allocateSlice[U any](pslice **[]U, len int, cap int) {
	typ := reflect.TypeFor[[]U]()
	rSlice := reflect.MakeSlice(typ, len, cap)
	slice := rSlice.Interface().([]U)
	*pslice = &slice
}

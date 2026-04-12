package eco

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
)

type KeyValue mvccpb.KeyValue
type Decode func(ctx *ecoContext, kvs []*mvccpb.KeyValue, key string, p any) error

func decode(ctx *ecoContext, cli *EtcdClient, key string, pp any) error {
	ctx.logger.signature("decode", key, reflect.TypeOf(pp).Elem())
	ctx.inc()
	defer ctx.dec()

	// Confirm p is a 'normal pointer', ie a pointer that is not a pointer-to-a-pointer
	if !isPointerToPointer(pp) {
		return fmt.Errorf("p must be a pointer-to-a-pointer, with an element type that is not a pointer (kind=%s)",
			reflect.TypeOf(pp).Kind().String())
	}

	// Get object from etcd + make sure there's only 1
	op := etcd.OpGet(key, etcd.WithPrefix())
	txn := cli.Txn(ctx)
	resp, err := txn.Then(op).Commit()
	if err != nil {
		return nil
	}

	rangeResponse := resp.Responses[0].GetResponseRange()
	kvs := rangeResponse.Kvs
	decoder := MainDecoder{}
	rv := reflect.ValueOf(pp)
	err = decoder.Decode(ctx, kvs, key, rv)
	return err

	// // Simple objects are easy to deal with. Just use json.Unmarhsal()
	// if isScalar(ty.Elem().Kind()) {
	// 	// Get object from etcd + make sure there's only 1
	// 	resp, err := cli.Client.Get(ctx, key)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(resp.Kvs) != 1 {
	// 		return fmt.Errorf("expected one key; got %d", len(resp.Kvs))
	// 	}

	// 	// Unmarshal the result
	// 	getVal := resp.Kvs[0].Value
	// 	ctx.logger.Infof("getVal()=%v (%s)", getVal, getVal)
	// 	err = json.Unmarshal(getVal, i)
	// 	if err != nil {
	// 		ctx.logger.Errorf("Unmarshalling error: %s (%#v)", err.Error(), getVal)
	// 		return err
	// 	}
	// 	// @note - should we return here?
	// 	return nil
	// }

	// Some non-simple type are supported. The rest of the function checks for them.
	// Note - we want the type of the underlying element, not the type of the pointer
	// kindElem := getTypeKind(ctx, ty.Elem())

	// switch kindElem {
	// case SimpleMap: return decodeMap(ctx, cli, key, i)
	// case SimpleSlice: return decodeSlice(ctx, cli, key, i)
	// case SimpleStruct: return decodeStruct(ctx, cli, key, i)

	// default:
	// 	return errors.New("unsupported type")
	// }
}

// eco stores maps as a number of sub-KVs with a common prefix. Go requires all
// KV values in a map are the same type, but etcd has no way of enforcing this. To
// etcd they're all just KVs.
func decodeMap(ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.logger.signature("decodeMap", "-etcdClient-", key, reflect.TypeOf(i))
	ctx.inc()
	defer ctx.dec()

	ctx.logger.Infof("i=%v ValueOf(i)=%v Elem()=%v ValueOf(Elem())=%v", i, reflect.ValueOf(i), reflect.ValueOf(i).Elem(), reflect.ValueOf(reflect.ValueOf(i).Elem()))

	ty := reflect.TypeOf(i)
	// ctx.println(subtle(fmt.Sprintf("ty=%s", fullTypeName(ty))))
	// Only pointers are supported
	if ty.Kind() != reflect.Pointer {
		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty))
	}

	// Only simple maps are supported
	kind := getTypeKind(ctx, ty.Elem())
	if kind != SimpleMap {
		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty.Elem()))
	}

	// add trailing slash. a/key => a/key/
	if !strings.HasSuffix(key, string(filepath.Separator)) {
		key += string(filepath.Separator)
	}

	// Get entire object tree
	// @note this might be quite large. ideally pagination would avoid issues with huge maps
	resp, err := etcdClient.Client.Get(ctx, key, etcd.WithPrefix())
	ctx.logger.info(common.Highlight("Keys"))
	var valMap reflect.Value
	// The caller may have specified a nil map, or an existing map
	// If nil, create a new map. If not, use the existing map
	if reflect.ValueOf(i).Elem().IsNil() {
		ctx.logger.info("map is nil; initializing new map")
		valMap = reflect.MakeMap(ty.Elem())
		reflect.ValueOf(i).Elem().Set(valMap)
	} else {
		ctx.logger.info("pointer is not nil; using existing map")
		valMap = reflect.Indirect(reflect.ValueOf(i))
	}
	for _, kv := range resp.Kvs {
		// Print a nice log statement
		// @note this is a lot of clutter for logging, esp when the real code
		// is a simple json.Unmarshal()
		skey := strings.TrimPrefix(string(kv.Key), key)
		skeyQuoted := fmt.Sprintf("\"%s\"", skey)
		ctx.logger.Infof("%-16s %-16s", skeyQuoted, kv.Value)
		// (*i)[skey] = kv.Value
		// simple json.Unmarshal() of value
		// @note this only supports maps of scalars. it needs to support nested maps since those are allowed. I think.
		var sval string
		err = json.Unmarshal(kv.Value, &sval)
		if err != nil {
			return err
		}
		// set a map value, reflection-style
		valMap.SetMapIndex(reflect.ValueOf(skey), reflect.ValueOf(sval))
	}

	if err != nil {
		return err
	}

	ctx.logger.info(common.Highlight("returning nil"))
	return nil
}

func decodeResponse(ctx *ecoContext, key string, kvs []*mvccpb.KeyValue, i any) error {
	// Confirm that the incoming variable is a pointer
	iType := reflect.TypeOf(i)
	if iType.Kind() != reflect.Pointer {
		return fmt.Errorf("unsupported type (%s) -  must be pointer", common.FullTypeName(iType))
	}

	// Get the kind of incoming pointer, to determine how to unmarshal
	elemKind := iType.Elem().Kind()
	var decoder Decoder = decoderMap[elemKind]
	rv := reflect.ValueOf(i)
	err := decoder.Decode(ctx, kvs, key, rv)
	return err
}

func decodeResponseMap(ctx *ecoContext, key string, kvs []*mvccpb.KeyValue, i any) error { return nil }

func decodeResponseScalar(ctx *ecoContext, key string, kvs []*mvccpb.KeyValue, i any) error {
	kv := kvs[0]
	err := json.Unmarshal(kv.Value, i)
	return err
}

func decodeResponseSlice(ctx *ecoContext, kvs []*mvccpb.KeyValue, key string, p any) error {
	var iMax int = -1
	for _, kv := range kvs {
		index, is := getSliceKeyIndex(string(kv.Key))
		if !is {
			return fmt.Errorf("key is not a valid slice element key (key='%s')", key)
		}
		if index > iMax {
			iMax = index
		}
	}
	ctx.logger.Infof("iMax=%d", iMax)

	// typeP := reflect.TypeOf(p)
	// typeEl := typeP.Elem()
	// size := iMax+1
	// slice := reflect.MakeSlice(typeEl, size, size)
	// for _, kv := range kvs {
	// index, err := getSliceKeyIndex(string(kv.Key))
	// if err != nil {
	// 	return err
	// }
	// err = json.Unmarshal(kv.Value, &slice[index])
	// if err == nil {
	// 	return err
	// }
	// }

	// *p = slice
	return nil

	// slice := make([]bool, respRange.Count)
	// for i := range slice {
	// 	err = json.Unmarshal(respRange.Kvs[i].Value, &slice[i])
	// 	ctx.logger.Infof("respRange.Kvs[i].Value=%#v slice[%d]=%#v", string(respRange.Kvs[i].Value), i, slice[i])
	// 	require.NoError(t, err)
	// }
}

func decodeResponseStruct(ctx *ecoContext, key string, kvs []*mvccpb.KeyValue, i any) error {
	return nil
}

func decodeSlice(ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.logger.signature("decodeSlice", key, reflect.TypeOf(i).Elem())
	ctx.inc()
	defer ctx.dec()

	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer {
		return fmt.Errorf("unsupported type (%s) -  must be pointer", common.FullTypeName(ty))
	}
	kind := getTypeKind(ctx, ty.Elem())
	if kind != SimpleSlice {
		return fmt.Errorf("unsupported kind (%s) - must be SimpleSlice", common.FullTypeName(ty.Elem()))
	}

	kindElem := getTypeKind(ctx, ty.Elem().Elem())
	ctx.logger.Infof("kindElem=%s", kindElem.String())
	// Get slice keys
	sliceKeys, err := getSliceKeys(ctx, etcdClient, key)
	if err != nil {
		return err
	}
	ctx.logger.Infof("sliceKeys=%v", sliceKeys)
	// slice := reflect.MakeSlice(ty.Elem(), len(sliceKeys), len(sliceKeys))

	// for _, sliceKey := range sliceKeys {
	// 	elKey := path.Join(key, elKey)
	// 	// I have an element type.
	// 	// How do I create a variable to hold that type, and then decode a byte string into it?
	// }

	// Dynamically allocate array
	// For each slice key
	//	get index thingee
	//	somehow do a decode to the reflect.Value, even though I don't know how to do that

	// ctx.logger.Appendf("Getting key %s ...", key)
	// resp, err := etcdClient.Client.Get(ctx, key)
	// if err != nil {
	// 	ctx.logger.Flush(slog.LevelError, err.Error())
	// 	return err
	// }
	// if len(resp.Kvs) != 1 {
	// 	ctx.logger.Flush(slog.LevelError, "error")
	// 	return fmt.Errorf("expected one key; got %d", len(resp.Kvs))
	// }
	// ctx.logger.Flush(slog.LevelInfo, "ok")

	// getVal := resp.Kvs[0].Value
	// ctx.logger.Infof("getVal()=%v (%s)", getVal, getVal)
	// err = json.Unmarshal(getVal, i)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func decodeStruct(ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.logger.signature("decodeStruct", "-etcdClient", key, reflect.TypeOf(i))
	ctx.inc()
	defer ctx.dec()

	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer {
		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty))
	}
	tyElem := ty.Elem()
	kind := getTypeKind(ctx, tyElem)
	if kind != SimpleStruct {
		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty.Elem()))
	}
	nFields := tyElem.NumField()
	ctx.logger.Infof("%-16s %-16d", "nFields", nFields)
	for iField := range nFields {
		field := tyElem.Field(iField)
		ctx.logger.info(string(common.Lowlight(fmt.Sprintf("%-16d %-16s %-16s", iField, field.Name, field.Tag.Get("eco")))))
	}

	if !strings.HasSuffix(key, string(filepath.Separator)) {
		key += string(filepath.Separator)
	}
	resp, err := etcdClient.Client.Get(ctx, key, etcd.WithPrefix())
	if err != nil {
		return err
	}

	fieldNameMap, err := fieldNameMap(i)
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		skey := strings.TrimPrefix(string(kv.Key), key)
		skeyQuoted := fmt.Sprintf("\"%s\"", skey)
		var sval any
		err = json.Unmarshal(kv.Value, &sval)
		if err != nil {
			return err
		}
		field := fieldNameMap[skey]
		field.Set(reflect.ValueOf(sval))
		ctx.logger.Infof("%-16s %-16v", skeyQuoted, sval)
	}

	return nil
}

func TestBoolSlice(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := "/test/boolslice"
	expectedVal := []bool{true, true, false}
	putAndTest(t, ctx, cli, key, expectedVal)

	var p *[]bool
	pp := &p
	err := decode(ctx, cli, key, pp)
	require.NoError(t, err)
	require.Equal(t, expectedVal, p)
	t.Log(p)
}

func TestDecodeBool(t *testing.T) {
	testDecodeScalar(t, "/test/bool", true)
}

func TestDecodeFloat(t *testing.T) {
	testDecodeScalar(t, "/test/float", float32(42.0))
}

func TestDecodeInt(t *testing.T) {
	testDecodeScalar(t, "/test/int", int(-13.0))
}

func TestDecodeString(t *testing.T) {
	testDecodeScalar(t, "/test/string", `This\nis\a\<difficult>\nstring\n\to\n\e"s'c"a'p"e\n`)
}

func TestDecodeUint(t *testing.T) {
	testDecodeScalar(t, "/test/uint", uint(13.0))
}

func TestGetBool(t *testing.T) {
	testGetScalar(t, "test/bool", true)
}

func TestGetBoolSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := "/test/boolSlice"
	expectedVals := []bool{true, false, true, true}
	for i, expectedVal := range expectedVals {
		subkey := fmt.Sprintf("%s/%d", key, i)
		putAndTest(t, ctx, cli, subkey, expectedVal)
	}

	// Get kvs for seeded data
	// // Create etcd OpGet
	op := etcd.OpGet(key, etcd.WithPrefix())
	txn := createTxn(t, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)
	require.NotNil(t, resp)
	respRange := resp.Responses[0].GetResponseRange()
	ctx.logger.Infof("respRange.Count=%d", respRange.Count)
	kvs := respRange.Kvs

	// Decode the slice
	var p *[]bool
	pp := &p
	rv := reflect.ValueOf(pp)
	decoder := SliceDecoder{}
	err = decoder.Decode(ctx, kvs, key, rv)
	i := rv.Interface()
	pp = i.(**[]bool)
	require.NoError(t, err)
	require.Equal(t, expectedVals, **pp)

	// // Create etcd OpGet
	// op := etcd.OpGet(key, etcd.WithPrefix())
	// require.True(t, op.IsGet())
	// require.Equal(t, key, string(op.KeyBytes()))

	// txn := createTxn(t, cli)
	// resp, err := txn.Then(op).Commit()
	// require.NoError(t, err)
	// require.NotNil(t, resp)

	// respRange := resp.Responses[0].GetResponseRange()
	// ctx.logger.Infof("respRange.Count=%d", respRange.Count)
	// kvs := respRange.Kvs
	// for i, kv := range kvs {
	// 	ctx.logger.Infof("i=%d key=%s val=%s", i, string(kv.Key), string(kv.Value))
	// }
	// // ctx.logger.Infof("%#v", resp.Responses)
	// // ctx.logger.Infof("%#v", resp.Responses[0])

	// slice := make([]bool, respRange.Count)
	// err = decodeResponseSlice(ctx, kvs, key, &slice)
	// require.NoError(t, err)
	// // for i := range slice {
	// // 	err = json.Unmarshal(respRange.Kvs[i].Value, &slice[i])
	// // 	ctx.logger.Infof("respRange.Kvs[i].Value=%#v slice[%d]=%#v", string(respRange.Kvs[i].Value), i, slice[i])
	// // 	require.NoError(t, err)
	// // }

	// ctx.logger.Infof("*pslice=%#v", *pp)
}

func TestGetFloat(t *testing.T) {
	testGetScalar(t, "/test/float", float32(42.0))
}

func TestGetInt(t *testing.T) {
	testGetScalar(t, "/test/int", int(-13))
}

func TestGetString(t *testing.T) {
	testGetScalar(t, "/test/string", "hello world")
}

func TestGetUint(t *testing.T) {
	testGetScalar(t, "/test/uint", uint(13))
}

func TestFieldNameMap(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var ecoTest = EcoTest{}
	var p *EcoTest = &ecoTest

	fieldNameMap, err := fieldNameMap(p)
	require.NoError(t, err)
	require.NotNil(t, fieldNameMap)
	t.Logf("%#v", fieldNameMap)

	fieldNameMap["Anon"].Set(reflect.ValueOf("(...)"))
	fieldNameMap["name"].Set(reflect.ValueOf("Me"))
	fieldNameMap["lucky_number"].Set(reflect.ValueOf(13.0))
	t.Logf("%#v", p)
}

func TestFloatSlice(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	ctx, cli := initAndTest(t)

	key := "/test/float32slice"
	val := []float32{42.0, 1764.0, 6.54321}
	putAndTest(t, ctx, cli, key, val)

	type float32slice []float32
	var decodedVal float32slice
	var i = &decodedVal
	err := decode(ctx, cli, key, any(i).(**any))
	require.NoError(t, err)
	require.Equal(t, float32slice(val), decodedVal)
	t.Log(decodedVal)
}

func TestDecode_IntSlice(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	ctx, etcdClient := initAndTest(t)

	key := "/test/intSlice"
	val := []int{5, 8, 13}

	type intslice []int
	var decodedVal intslice
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, any(i).(**any))
	require.NoError(t, err)
	require.Equal(t, intslice(val), decodedVal)
	t.Log(decodedVal)
}

func TestMapStringString(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	ctx, cli := initAndTest(t)

	key := "/test/map/stringstring"
	val1 := "meat"
	val2 := "Meat"
	val3 := "MEEEEAT"
	type mapstringstring map[string]string
	val := mapstringstring{
		"foo": val1,
		"bar": val2,
		"bum": val3,
	}
	for k, v := range val {
		putAndTest(t, ctx, cli, filepath.Join(key, k), v)
	}

	var decodedVal mapstringstring = nil
	type pmapstringstring *mapstringstring
	var i pmapstringstring = &decodedVal
	err := decode(ctx, cli, key, any(&i).(**any))
	require.NoError(t, err)
	require.Equal(t, (val), decodedVal)
	t.Log(decodedVal)
}

func TestMisc(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	etcdClient, err := NewEtcdClientFromConfig()
	ctx := newEcoContext(os.Stdout)

	key1 := "/test/f"
	key2 := "/test/f"
	opGet1 := etcd.OpGet(key1)
	opGet2 := etcd.OpGet(key2, etcd.WithPrefix())
	require.NoError(t, err)
	txn := etcdClient.Txn(ctx)
	resp, err := txn.Then(opGet1, opGet2).Commit()
	require.NoError(t, err)
	for _, resp2 := range resp.Responses {
		t.Logf("%d", resp2.GetResponseRange().Count)
	}
}

func TestNilMap(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var m map[string]string = nil
	var pm *map[string]string = &m
	t.Logf("reflect.ValueOf(pm).IsNil()=%v", reflect.ValueOf(pm).IsNil())
	t.Logf("reflect.ValueOf(pm).Elem().IsNil()=%v", reflect.ValueOf(pm).Elem().IsNil())

	// Get type of underlying object from pointer
	tyElem := reflect.TypeOf(pm).Elem()
	t.Logf("tyElem=%s", common.FullTypeName(tyElem))

	// Create map and assign a value
	valMap := reflect.MakeMap(tyElem)
	reflect.ValueOf(pm).Elem().Set(valMap)
	valMap.SetMapIndex(reflect.ValueOf("foo"), reflect.ValueOf("13"))

	//
	t.Logf("%#v", *pm)
}

func TestNilMapPointer(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var pm *map[string]string = nil
	pmValue := reflect.ValueOf(pm)
	t.Logf("pmValue.IsNil()=%v", pmValue.IsNil())
}

func TestNilPointer(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var m map[string]string = nil
	var pm *map[string]string = &m
	t.Logf("reflect.ValueOf(pm).IsNil()=%v", reflect.ValueOf(pm).IsNil())
	t.Logf("reflect.ValueOf(pm).Elem().IsNil()=%v", reflect.ValueOf(pm).Elem().IsNil())
}


func TestStringSlice(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	ctx := newEcoContext(os.Stdout)
	cli, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, cli)

	key := "/test/stringslice"
	val := []string{"foo", "bar", "bum"}
	putAndTest(t, ctx, cli, key, val)

	type stringslice []string
	var decodedVal stringslice
	var i = &decodedVal
	err = decode(ctx, cli, key, any(i).(**any))
	require.NoError(t, err)
	require.Equal(t, stringslice(val), decodedVal)
	t.Log(decodedVal)
}

func TestStructEcoTest(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	ctx, cli := initAndTest(t)

	key := "/test/struct/ecotest"
	name := "Me"
	luckyNumber := 13
	val := NewEcoTest(name, float64(luckyNumber))
	putAndTest(t, ctx, cli, filepath.Join(key, "name"), val.Name)
	putAndTest(t, ctx, cli, filepath.Join(key, "lucky_number"), val.LuckyNumber)
	putAndTest(t, ctx, cli, filepath.Join(key, "Anon"), val.Anon)

	var decodedVal EcoTest
	type pEcoTest *EcoTest
	var i pEcoTest = &decodedVal
	err := decode(ctx, cli, key, any(i).(**any))
	require.NoError(t, err)
	require.Equal(t, *val, *i)
	t.Log(decodedVal)
}

func TestStructSetField0(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var st = EcoTest{}
	type pEcoTest *EcoTest
	var pst pEcoTest = &st

	val := reflect.ValueOf(pst).Elem()
	val.Field(0).Set(reflect.ValueOf("me"))
	val.Field(1).Set(reflect.ValueOf(13.0))
	t.Logf("%#v", *pst)

}

func TestStructSetField1(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var st = EcoTest{}
	p := reflect.ValueOf(&st)
	val := p.Elem()

	t.Logf("fullTypeName(val.Type())=%s", common.FullTypeName(val.Type()))
	t.Logf("val.Type().Kind()=%s", val.Type().Kind().String())

	require.True(t, val.CanSet())
	require.True(t, val.Field(0).CanSet())

	val.Field(0).SetString("me")
	val.Field(1).Set(reflect.ValueOf(13.0))
	t.Logf("%#v", st)
}

func TestUintSlice(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	ctx := newEcoContext(os.Stdout)
	cli, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, cli)

	key := "/test/uintslice"
	val := []uint{5, 12, 13}
	putAndTest(t, ctx, cli, key, val)

	type uintslice []uint
	var decodedVal uintslice
	var i = &decodedVal
	err = decode(ctx, cli, key, any(i).(**any))
	require.NoError(t, err)
	require.Equal(t, uintslice(val), decodedVal)
	t.Log(decodedVal)
}

func initAndTest(t *testing.T) (*ecoContext, *EtcdClient) {
	ctx := newEcoContext(os.Stdout)
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	return ctx, etcdClient
}

// With the EtcdClient, Put a value to etcd, then Get it back to confirm the
// Put succeeded
func putAndTest(t *testing.T, ctx *ecoContext, etcdClient *EtcdClient, key string, i any) {
	// resp, err := etcdClient.Put(context.Background(), key, val)
	ctx.inc()
	defer ctx.dec()

	ctx.logger.Infof("Writing to %s... ", key)
	j, err := json.Marshal(i)
	require.NoError(t, err)
	resp, err := etcdClient.Put(ctx, key, string(j))
	require.NoError(t, err)
	require.NotNil(t, resp)

	ctx.logger.Infof("Reading %s... ", key)
	buf, err := etcdClient.Get(key)
	require.NoError(t, err)
	require.Equal(t, j, buf)
	require.Equal(t, string(j), string(buf))
	ctx.logger.Infof("%#v", resp)
}

func testDecodeScalar[U any](t *testing.T, key string, expectedVal U) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	ctx, cli := initAndTest(t)

	// Seed test data
	putAndTest(t, ctx, cli, key, expectedVal)

	var decodedVal *U
	var i = &decodedVal
	err := decode(ctx, cli, key, any(i).(**any))
	require.NoError(t, err)
	require.Equal(t, expectedVal, *decodedVal)
	t.Log(decodedVal)
}

func testGetScalar[U any](t *testing.T, key string, expectedVal U) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	putAndTest(t, ctx, cli, key, expectedVal)

	// Create the GET Op
	op := etcd.OpGet(key)

	// Get the response from etcd
	txn := createTxn(t, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)

	// Get the KVs from the response
	p := new(U)
	pp := &p
	rv := reflect.ValueOf(pp)
	rangeResp := resp.Responses[0].GetResponseRange()
	kvs := rangeResp.Kvs

	// Get the decoder from the DecoderMap and decode
	decoder := &ScalarDecoder[U]{}
	err = decoder.Decode(ctx, kvs, key, rv)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, expectedVal, *p)
}

/*
ap[string]reflect.Value
{
	"Anon":reflect.Value{
		typ_:(*abi.Type)(0x1014f6460),
		ptr:(unsafe.Pointer)(0x140002f0b28),
		flag:0x98
	},
	"lucky_number":reflect.Value{
		typ_:(*abi.Type)(0x1014f6820),
		ptr:(unsafe.Pointer)(0x140002f0b20),
		flag:0x8e
	},
	"name":reflect.Value{
		typ_:(*abi.Type)(0x1014f6460),
		ptr:(unsafe.Pointer)(0x140002f0b10),
		flag:0x98}
	}

	reflect.Value{
		"Anon":reflect.Value{
			typ_:(*abi.Type)(0x104afa460),
			ptr:(unsafe.Pointer)(0x140002f0b28),
			flag:0x198
		},
		"lucky_number":reflect.Value{
			typ_:(*abi.Type)(0x104afa820),
			ptr:(unsafe.Pointer)(0x140002f0b20),
			flag:0x18e
		},
		"name":reflect.Value{
			typ_:(*abi.Type)(0x104afa460),
			ptr:(unsafe.Pointer)(0x140002f0b10),
			flag:0x198}
		}
*/

func TestIsNormalPointer_Int(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var expectedVal bool = false
	var p int
	flag := isNormalPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsNormalPointer_Pointer(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var expectedVal bool = true
	var p *int
	flag := isNormalPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsNormalPointer_PointerPointer(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var expectedVal bool = false
	var p **int
	flag := isNormalPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsPointerToPointer_Int(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var expectedVal bool = false
	var p int
	flag := isPointerToPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsPointerToPointer_Pointer(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var expectedVal bool = false
	var p *int
	flag := isPointerToPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsPointerToPointer_PointerPointer(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var expectedVal bool = true
	var p **int
	flag := isPointerToPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsPointerToPointer_PointerPointerPointer(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	var expectedVal bool = false
	var p ***int
	flag := isPointerToPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestAllocOrDont(t *testing.T) {
	defer func() { pa := recover(); if pa != nil { t.Error(pa) } }()
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	var err error
	expectedValue := 13
	buf := fmt.Append(nil, expectedValue)

	var val int
	err = allocOrDont[int](t, buf, &val)
	require.NoError(t, err)
	require.Equal(t, expectedValue, val)

	var pval *int
	err = allocOrDont[int](t, buf, &pval)
	require.NoError(t, err)
	require.Equal(t, expectedValue, *pval)

	var pval2 *int = new(int)
	err = allocOrDont[int](t, buf, pval2)
	require.NoError(t, err)
	require.Equal(t, expectedValue, *pval2)

	var val2 int
	var pval3 *int = &val2
	err = allocOrDont[int](t, buf, &pval3)
	require.NoError(t, err)
	require.Equal(t, expectedValue, val)

}

func allocOrDont[U any, V *U | **U](t *testing.T, buf []byte, v V) error {
	var pp **U

	t.Log("checking if variable is pointer-to-pointer")
	pp, is := any(v).(**U)
	if !is {
		t.Log("False - it's not a pointer to a pointer. So we assume it's a pointer and get its address")
		pp = (any(&v)).(**U)
	} else {
		t.Log("True: it's a pointer to a pointer")
	}
	err := json.Unmarshal(buf, pp)
	return err
}

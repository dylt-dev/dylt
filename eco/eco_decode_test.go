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
		index, is := getSliceItemKey(key, string(kv.Key))
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

func TestDecodeBool(t *testing.T) {
	decodeAndTestScalar(t, "/test/bool", true)
}

func TestDecodeBoolSlice(t *testing.T) {
	decodeAndTestSlice(t,
		"/test/boolslice",
		[]bool{true, true, false})
}

func TestDecodeFloat(t *testing.T) {
	decodeAndTestScalar(t, "/test/float", float32(42.0))
}

func TestDecodeFloatSlice(t *testing.T) {
	decodeAndTestSlice(t,
		"/test/float32slice",
		[]float32{42.0, 1764.0, 6.54321})
}

func TestDecodeInt(t *testing.T) {
	decodeAndTestScalar(t, "/test/int", int(-13.0))
}

func TestDecodeIntSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	key := "/test/intSlice"
	expectedData := []int{5, 8, 13}
	putAndTestSlice(t, ctx, cli, key, expectedData)

	var pSlice *[]int
	err := decode(ctx, cli, key, &pSlice)
	require.NoError(t, err)
	require.Equal(t, expectedData, *pSlice)
	t.Log(*pSlice)
}

func TestDecodeString(t *testing.T) {
	decodeAndTestScalar(t, "/test/string", `This\nis\a\<difficult>\nstring\n\to\n\e"s'c"a'p"e\n`)
}

func TestDecodeStringSlice(t *testing.T) {
	decodeAndTestSlice(t,
		"/test/stringslice",
		[]string{"foo", "bar", "bum"})
}

func TestDecodeUint(t *testing.T) {
	decodeAndTestScalar(t, "/test/uint", uint(13.0))
}

func TestDecodeUintSlice(t *testing.T) {
	decodeAndTestSlice(t,
		"/test/uintslice",
		[]uint{5, 12, 13})
}

func TestGetBool(t *testing.T) {
	testGetScalar(t, "test/bool", true)
}

func TestGetBoolSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := "/test/boolSlice"
	expectedData := []bool{true, false, true, true}
	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}

func TestGetFloat(t *testing.T) {
	testGetScalar(t, "/test/float", float32(42.0))
}

func TestGetFloatSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := "/test/floatSlice"
	expectedData := []float32{42.0, 1764.0, 6.54321}

	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}

func TestGetInt(t *testing.T) {
	testGetScalar(t, "/test/int", int(-13))
}

func TestGetIntSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := "/test/intSlice"
	expectedData := []int{5, 8, 13}

	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}

func TestGetMap(t *testing.T) {
	ctx, cli := initAndTest(t)

	key := "/test/team/astros/Players/altuve"
	expectedBorn := "Venezuela"
	expectedId := "1"
	expectedIsActive := "true"
	expectedName := "Jose Altuve"
	expectedWeight := "160"
	expectedData := map[string]string{
		"Born": expectedBorn,
		"Id": expectedId,
		"IsActive": expectedIsActive,
		"Name": expectedName,
		"Weight": expectedWeight,
	}
	putAndTestMap(t, ctx, cli, key, expectedData)

	// Create the GET Op
	op := etcd.OpGet(key, etcd.WithPrefix())

	// Get the response from etcd
	txn := createTxn(t, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)

	// Get the KVs from the response
	rangeResp := resp.Responses[0].GetResponseRange()
	kvs := rangeResp.Kvs

	// Decode the map
	var pmap *map[string]string
	decoder := MapDecoder{}
	err = decoder.Decode(ctx, kvs, key, reflect.ValueOf(&pmap))
	require.NoError(t, err)

	// Check the contents
	require.Equal(t, expectedData, *pmap)

	t.Log(*pmap)
}

func TestGetString(t *testing.T) {
	testGetScalar(t, "/test/string", "hello world")
}

func TestGetStringSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := "/test/stringSlice"
	expectedData := []string{"foo", "bar", "bum"}
	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}

func TestGetUint(t *testing.T) {
	testGetScalar(t, "/test/uint", uint(13))
}

func TestGetUintSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := "/test/uintSlice"
	expectedData := []uint{5, 12, 13}

	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}

func TestFieldNameMap(t *testing.T) {
	t.Skip("I honestly don't know what this test is for")
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

func TestMisc(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
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
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
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
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var pm *map[string]string = nil
	pmValue := reflect.ValueOf(pm)
	t.Logf("pmValue.IsNil()=%v", pmValue.IsNil())
}

func TestNilPointer(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var m map[string]string = nil
	var pm *map[string]string = &m
	t.Logf("reflect.ValueOf(pm).IsNil()=%v", reflect.ValueOf(pm).IsNil())
	t.Logf("reflect.ValueOf(pm).Elem().IsNil()=%v", reflect.ValueOf(pm).Elem().IsNil())
}

func TestStructEcoTest(t *testing.T) {
	ctx, cli := initAndTest(t)

	key := "/test/struct/ecotest"
	expectedName := "Me"
	expectedLuckyNumber := 13
	expectedNoTag := "no-tag-value"
	expectedData := EcoTest{
		Name:        expectedName,
		LuckyNumber: float64(expectedLuckyNumber),
		NoTag:       expectedNoTag,
	}
	putAndTestStruct(t, ctx, cli, key, expectedData)

	var data *EcoTest
	err := decode(ctx, cli, key, &data)
	require.NoError(t, err)
	// require.Equal(t, expectedData, *data)
	require.Nil(t, data)
	// t.Log(*decodedVal)
}

func TestStructSetField0(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var st = EcoTest{}
	type pEcoTest *EcoTest
	var pst pEcoTest = &st

	val := reflect.ValueOf(pst).Elem()
	val.Field(0).Set(reflect.ValueOf("me"))
	val.Field(1).Set(reflect.ValueOf(13.0))
	t.Logf("%#v", *pst)

}

func TestStructSetField1(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
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

func decodeAndTestSlice[U any](t *testing.T, key string, expectedData []U) {
	ctx, cli := initAndTest(t)

	putAndTestSlice(t, ctx, cli, key, expectedData)

	var pSlice *[]U
	err := decode(ctx, cli, key, &pSlice)
	require.NoError(t, err)
	require.Equal(t, expectedData, *pSlice)
	t.Log(*pSlice)

}

func getAndTestMap[U any](t *testing.T, ctx *ecoContext, expectedData map[string]U, kvs []*mvccpb.KeyValue, key string) {
	var pData *map[string]U
	rv := reflect.ValueOf(&pData)
	decoder := SliceDecoder{}
	err := decoder.Decode(ctx, kvs, key, rv)
	require.NoError(t, err)
	require.Equal(t, expectedData, *pData)
}

func getAndTestSlice[U any](t *testing.T, ctx *ecoContext, expectedData []U, kvs []*mvccpb.KeyValue, key string) {
	var pSlice *[]U
	rv := reflect.ValueOf(&pSlice)
	decoder := SliceDecoder{}
	err := decoder.Decode(ctx, kvs, key, rv)
	require.NoError(t, err)
	require.Equal(t, expectedData, *pSlice)
}

func getAndTestSliceKVs(t *testing.T, ctx *ecoContext, cli *EtcdClient, key string) []*mvccpb.KeyValue {
	op := etcd.OpGet(key, etcd.WithPrefix())
	txn := createTxn(t, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)
	require.NotNil(t, resp)
	respRange := resp.Responses[0].GetResponseRange()
	ctx.logger.Infof("respRange.Count=%d", respRange.Count)
	kvs := respRange.Kvs
	return kvs
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
func putAndTestMap[U any](t *testing.T, ctx *ecoContext, cli *EtcdClient, key string, data map[string]U) {
	ctx.logger.Infof("Writing slice at %s ...", key)
	ops := []etcd.Op{}
	for itemKey, val := range data {
		subkey := fmt.Sprintf("%s/%s", key, itemKey)
		bufVal, err := json.Marshal(val)
		require.NoError(t, err)
		op := etcd.OpPut(subkey, string(bufVal))
		ops = append(ops, op)
	}
	txn := createTxn(t, cli)
	require.NotNil(t, txn)
	txn.Then(ops...).Commit()
}

func putAndTestScalar(t *testing.T, ctx *ecoContext, etcdClient *EtcdClient, key string, i any) {
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

func putAndTestSlice[U any](t *testing.T, ctx *ecoContext, cli *EtcdClient, key string, data []U) {
	ctx.logger.Infof("Writing map at %s ...", key)
	ops := []etcd.Op{}
	for i, val := range data {
		subkey := fmt.Sprintf("%s/%d", key, i)
		bufVal, err := json.Marshal(val)
		require.NoError(t, err)
		op := etcd.OpPut(subkey, string(bufVal))
		ops = append(ops, op)
	}
	txn := createTxn(t, cli)
	require.NotNil(t, txn)
	txn.Then(ops...).Commit()
}

func putAndTestStruct[U any](t *testing.T, ctx *ecoContext, cli *EtcdClient, key string, data U) {

}

func decodeAndTestScalar[U any](t *testing.T, key string, expectedVal U) {
	ctx, cli := initAndTest(t)

	// Seed test data
	putAndTestScalar(t, ctx, cli, key, expectedVal)

	var decodedVal *U
	var i = &decodedVal
	err := decode(ctx, cli, key, i)
	require.NoError(t, err)
	require.Equal(t, expectedVal, *decodedVal)
	t.Log(decodedVal)
}

func testGetScalar[U any](t *testing.T, key string, expectedVal U) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	putAndTestScalar(t, ctx, cli, key, expectedVal)

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
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var expectedVal bool = false
	var p int
	flag := isNormalPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsNormalPointer_Pointer(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var expectedVal bool = true
	var p *int
	flag := isNormalPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsNormalPointer_PointerPointer(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var expectedVal bool = false
	var p **int
	flag := isNormalPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsPointerToPointer_Int(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var expectedVal bool = false
	var p int
	flag := isPointerToPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsPointerToPointer_Pointer(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var expectedVal bool = false
	var p *int
	flag := isPointerToPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsPointerToPointer_PointerPointer(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var expectedVal bool = true
	var p **int
	flag := isPointerToPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestIsPointerToPointer_PointerPointerPointer(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var expectedVal bool = false
	var p ***int
	flag := isPointerToPointer(p)
	require.Equal(t, expectedVal, flag)
}

func TestAllocOrDont(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
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

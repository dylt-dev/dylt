package eco

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
)

// eco stores maps as a number of sub-KVs with a common prefix. Go requires all
// KV values in a map are the same type, but etcd has no way of enforcing this. To
// etcd they're all just KVs.
// func decodeMap(ctx *common.EcoContext, etcdClient *EtcdClient, key string, i any) error {
// 	ctx.Logger.signature("decodeMap", "-etcdClient-", key, reflect.TypeOf(i))
// 	ctx.inc()
// 	defer ctx.dec()

// 	ctx.Logger.Infof("i=%v ValueOf(i)=%v Elem()=%v ValueOf(Elem())=%v", i, reflect.ValueOf(i), reflect.ValueOf(i).Elem(), reflect.ValueOf(reflect.ValueOf(i).Elem()))

// 	ty := reflect.TypeOf(i)
// 	// ctx.println(subtle(fmt.Sprintf("ty=%s", fullTypeName(ty))))
// 	// Only pointers are supported
// 	if ty.Kind() != reflect.Pointer {
// 		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty))
// 	}

// 	// Only simple maps are supported
// 	kind := getTypeKind(ctx, ty.Elem())
// 	if kind != SimpleMap {
// 		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty.Elem()))
// 	}

// 	// add trailing slash. a/key => a/key/
// 	if !strings.HasSuffix(key, string(filepath.Separator)) {
// 		key += string(filepath.Separator)
// 	}

// 	// Get entire object tree
// 	// @note this might be quite large. ideally pagination would avoid issues with huge maps
// 	resp, err := etcdClient.Client.Get(ctx, key, etcd.WithPrefix())
// 	ctx.Logger.info(common.Highlight("Keys"))
// 	var valMap reflect.Value
// 	// The caller may have specified a nil map, or an existing map
// 	// If nil, create a new map. If not, use the existing map
// 	if reflect.ValueOf(i).Elem().IsNil() {
// 		ctx.Logger.info("map is nil; initializing new map")
// 		valMap = reflect.MakeMap(ty.Elem())
// 		reflect.ValueOf(i).Elem().Set(valMap)
// 	} else {
// 		ctx.Logger.info("pointer is not nil; using existing map")
// 		valMap = reflect.Indirect(reflect.ValueOf(i))
// 	}
// 	for _, kv := range resp.Kvs {
// 		// Print a nice log statement
// 		// @note this is a lot of clutter for logging, esp when the real code
// 		// is a simple json.Unmarshal()
// 		skey := strings.TrimPrefix(string(kv.Key), key)
// 		skeyQuoted := fmt.Sprintf("\"%s\"", skey)
// 		ctx.Logger.Infof("%-16s %-16s", skeyQuoted, kv.Value)
// 		// (*i)[skey] = kv.Value
// 		// simple json.Unmarshal() of value
// 		// @note this only supports maps of scalars. it needs to support nested maps since those are allowed. I think.
// 		var sval string
// 		err = json.Unmarshal(kv.Value, &sval)
// 		if err != nil {
// 			return err
// 		}
// 		// set a map value, reflection-style
// 		valMap.SetMapIndex(reflect.ValueOf(skey), reflect.ValueOf(sval))
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	ctx.Logger.info(common.Highlight("returning nil"))
// 	return nil
// }

// func decodeResponse(ctx *common.EcoContext, key string, kvs []*mvccpb.KeyValue, i any) error {
// 	// Confirm that the incoming variable is a pointer
// 	iType := reflect.TypeOf(i)
// 	if iType.Kind() != reflect.Pointer {
// 		return fmt.Errorf("unsupported type (%s) -  must be pointer", common.FullTypeName(iType))
// 	}

// 	// Get the kind of incoming pointer, to determine how to unmarshal
// 	elemKind := iType.Elem().Kind()
// 	var decoder Decoder = decoderMap[elemKind]
// 	rv := reflect.ValueOf(i)
// 	err := decoder.Decode(ctx, kvs, key, rv)
// 	return err
// }

// func decodeResponseMap(ctx *common.EcoContext, key string, kvs []*mvccpb.KeyValue, i any) error { return nil }

// func decodeResponseScalar(ctx *common.EcoContext, key string, kvs []*mvccpb.KeyValue, i any) error {
// 	kv := kvs[0]
// 	err := json.Unmarshal(kv.Value, i)
// 	return err
// }

// func decodeResponseSlice(ctx *common.EcoContext, kvs []*mvccpb.KeyValue, key string, p any) error {
// 	var iMax int = -1
// 	for _, kv := range kvs {
// 		index, is := getSliceItemKey(key, string(kv.Key))
// 		if !is {
// 			return fmt.Errorf("key is not a valid slice element key (key='%s')", key)
// 		}
// 		if index > iMax {
// 			iMax = index
// 		}
// 	}
// 	ctx.Logger.Infof("iMax=%d", iMax)

// 	// typeP := reflect.TypeOf(p)
// 	// typeEl := typeP.Elem()
// 	// size := iMax+1
// 	// slice := reflect.MakeSlice(typeEl, size, size)
// 	// for _, kv := range kvs {
// 	// index, err := getSliceKeyIndex(string(kv.Key))
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// err = json.Unmarshal(kv.Value, &slice[index])
// 	// if err == nil {
// 	// 	return err
// 	// }
// 	// }

// 	// *p = slice
// 	return nil

// 	// slice := make([]bool, respRange.Count)
// 	// for i := range slice {
// 	// 	err = json.Unmarshal(respRange.Kvs[i].Value, &slice[i])
// 	// 	ctx.Logger.Infof("respRange.Kvs[i].Value=%#v slice[%d]=%#v", string(respRange.Kvs[i].Value), i, slice[i])
// 	// 	require.NoError(t, err)
// 	// }
// }

// func decodeResponseStruct(ctx *common.EcoContext, key string, kvs []*mvccpb.KeyValue, i any) error {
// 	return nil
// }

// func decodeSlice(ctx *common.EcoContext, etcdClient *EtcdClient, key string, i any) error {
// 	ctx.Logger.signature("decodeSlice", key, reflect.TypeOf(i).Elem())
// 	ctx.inc()
// 	defer ctx.dec()

// 	ty := reflect.TypeOf(i)
// 	if ty.Kind() != reflect.Pointer {
// 		return fmt.Errorf("unsupported type (%s) -  must be pointer", common.FullTypeName(ty))
// 	}
// 	kind := getTypeKind(ctx, ty.Elem())
// 	if kind != SimpleSlice {
// 		return fmt.Errorf("unsupported kind (%s) - must be SimpleSlice", common.FullTypeName(ty.Elem()))
// 	}

// 	kindElem := getTypeKind(ctx, ty.Elem().Elem())
// 	ctx.Logger.Infof("kindElem=%s", kindElem.String())
// 	// Get slice keys
// 	sliceKeys, err := getSliceKeys(ctx, etcdClient, key)
// 	if err != nil {
// 		return err
// 	}
// 	ctx.Logger.Infof("sliceKeys=%v", sliceKeys)
// 	// slice := reflect.MakeSlice(ty.Elem(), len(sliceKeys), len(sliceKeys))

// 	// for _, sliceKey := range sliceKeys {
// 	// 	elKey := path.Join(key, elKey)
// 	// 	// I have an element type.
// 	// 	// How do I create a variable to hold that type, and then decode a byte string into it?
// 	// }

// 	// Dynamically allocate array
// 	// For each slice key
// 	//	get index thingee
// 	//	somehow do a decode to the reflect.Value, even though I don't know how to do that

// 	// ctx.Logger.Appendf("Getting key %s ...", key)
// 	// resp, err := etcdClient.Client.Get(ctx, key)
// 	// if err != nil {
// 	// 	ctx.Logger.Flush(slog.LevelError, err.Error())
// 	// 	return err
// 	// }
// 	// if len(resp.Kvs) != 1 {
// 	// 	ctx.Logger.Flush(slog.LevelError, "error")
// 	// 	return fmt.Errorf("expected one key; got %d", len(resp.Kvs))
// 	// }
// 	// ctx.Logger.Flush(slog.LevelInfo, "ok")

// 	// getVal := resp.Kvs[0].Value
// 	// ctx.Logger.Infof("getVal()=%v (%s)", getVal, getVal)
// 	// err = json.Unmarshal(getVal, i)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

// func decodeStruct(ctx *common.EcoContext, etcdClient *EtcdClient, key string, i any) error {
// 	ctx.Logger.signature("decodeStruct", "-etcdClient", key, reflect.TypeOf(i))
// 	ctx.inc()
// 	defer ctx.dec()

// 	ty := reflect.TypeOf(i)
// 	if ty.Kind() != reflect.Pointer {
// 		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty))
// 	}
// 	tyElem := ty.Elem()
// 	kind := getTypeKind(ctx, tyElem)
// 	if kind != SimpleStruct {
// 		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty.Elem()))
// 	}
// 	nFields := tyElem.NumField()
// 	ctx.Logger.Infof("%-16s %-16d", "nFields", nFields)
// 	for iField := range nFields {
// 		field := tyElem.Field(iField)
// 		ctx.Logger.info(string(common.Lowlight(fmt.Sprintf("%-16d %-16s %-16s", iField, field.Name, field.Tag.Get("eco")))))
// 	}

// 	if !strings.HasSuffix(key, string(filepath.Separator)) {
// 		key += string(filepath.Separator)
// 	}
// 	resp, err := etcdClient.Client.Get(ctx, key, etcd.WithPrefix())
// 	if err != nil {
// 		return err
// 	}

// 	fieldNameMap, err := fieldNameMap(i)
// 	if err != nil {
// 		return err
// 	}
// 	for _, kv := range resp.Kvs {
// 		skey := strings.TrimPrefix(string(kv.Key), key)
// 		skeyQuoted := fmt.Sprintf("\"%s\"", skey)
// 		var sval any
// 		err = json.Unmarshal(kv.Value, &sval)
// 		if err != nil {
// 			return err
// 		}
// 		field := fieldNameMap[skey]
// 		field.Set(reflect.ValueOf(sval))
// 		ctx.Logger.Infof("%-16s %-16v", skeyQuoted, sval)
// 	}

// 	return nil
// }


func TestCreateOrGetStruct(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	st := common.TestStruct{
		Name:        "foo",
		LuckyNumber: 13,
		NoTag:       "bum",
	}

	pst, is := common.CreateOrGetStruct(ctx, reflect.ValueOf(&st))
	require.True(t, is)
	require.NotNil(t, pst)
	pEco, is := pst.(*common.TestStruct)
	require.True(t, is)
	require.Equal(t, pEco, &st)
}

// Pass CreateOtGetStruct a **struct, and confirm it allocates a struct,
// and that the pointer to the struct can be obtained by deferencing **struct
func TestCreateOrGetStructAlloc(t *testing.T) {
	expectedName := "foo"
	expectedLuckyNumber := float64(13.0)
	expectedNoTag := "bar"
	expectedStruct := common.TestStruct{
		Name:        expectedName,
		LuckyNumber: expectedLuckyNumber,
		NoTag:       expectedNoTag,
	}

	ctx := common.NewEcoContext(os.Stdout)
	var pst *common.TestStruct = nil
	ppst := &pst

	pNewSt, is := common.CreateOrGetStruct(ctx, reflect.ValueOf(ppst))
	require.True(t, is)
	require.NotNil(t, pNewSt)
	pEco, is := pNewSt.(*common.TestStruct)
	require.True(t, is)
	require.NotNil(t, pEco)
	require.Equal(t, pst, pEco)
	*pEco = expectedStruct
	require.Equal(t, expectedName, (**ppst).Name)
	require.Equal(t, expectedLuckyNumber, (**ppst).LuckyNumber)
	require.Equal(t, expectedNoTag, (**ppst).NoTag)
}

func TestDecodeBool(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := true
	decoder := MainDecoder{}
	
	buf := strconv.AppendBool([]byte{}, expected)
	tree := &ValueTree{Value: buf}
	var x bool
	p := &x

	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, expected, x)
}

func TestDecodeBool2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := true
	decoder := MainDecoder{}
	
	tree := &ValueTree{}
	tree.setBool(expected)
	var p *bool = nil
	pp := &p

	err := decoder.Decode(ctx, tree, pp)
	require.NoError(t, err)
	require.Equal(t, expected, *p)
}

func TestDecodeFloat1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := float64(169.0)
	decoder := MainDecoder{}
	
	tree := &ValueTree{}
	tree.setFloat(expected)
	var x float64
	p := &x

	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, expected, x)
}

func TestDecodeFloat2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := float64(169.0)
	decoder := MainDecoder{}
	
	tree := &ValueTree{}
	tree.setFloat(expected)
	var p *float64 = nil
	pp := &p

	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	require.Equal(t, expected, x)
}

func TestDecodeInt(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := int64(13)
	decoder := MainDecoder{}
	
	tree := &ValueTree{}
	tree.setInt(expected)
	var x int64
	p := &x

	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, expected, x)
}

func TestDecodeInt2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := int64(13)
	decoder := MainDecoder{}
	
	tree := &ValueTree{}
	tree.setInt(expected)
	var p *int64 = nil
	pp := &p

	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	require.Equal(t, expected, x)
}

func TestDecodeString(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := "meat"
	decoder := MainDecoder{}
	
	tree := &ValueTree{}
	tree.setString(expected)
	var x string
	p := &x

	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, expected, x)
}

func TestDecodeString2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := "meat"
	decoder := MainDecoder{}
	
	tree := &ValueTree{}
	tree.setString(expected)
	var p *string = nil
	pp := &p

	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	require.Equal(t, expected, x)
}

func TestDecodeUint1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := uint64(169.0)
	decoder := MainDecoder{}
	
	tree := &ValueTree{}
	tree.setUint(expected)
	var x uint64
	p := &x

	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, expected, x)
}

func TestDecodeUint2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := uint64(169.0)
	decoder := MainDecoder{}
	
	tree := &ValueTree{}
	tree.setUint(expected)
	var p *uint64 = nil
	pp := &p

	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	require.Equal(t, expected, x)
}

func TestGetBool1(t *testing.T) {
	testGetScalar(t, "/test/scalar/bool", true)
}

func TestGetBool2(t *testing.T) {
	testGetScalar2(t, "test/scalar/bool2", true)
}

func TestGetFloat(t *testing.T) {
	testGetScalar(t, "/test/float", float32(42.0))
}

func TestGetFloat2(t *testing.T) {
	testGetScalar2(t, "/test/float2", float32(42.0))
}

func TestGetInt(t *testing.T) {
	testGetScalar(t, "/test/int", int(-13))
}

func TestGetInt2(t *testing.T) {
	testGetScalar2(t, "/test/int2", int(-13))
}

func TestGetString(t *testing.T) {
	testGetScalar(t, "/test/string", "hello world")
}

func TestGetString2(t *testing.T) {
	testGetScalar2(t, "/test/string2", "hello world")
}

func TestGetStructAndUnmarshalField(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedName := "foo"
	expectedNameBuf, err := json.Marshal(expectedName)
	require.NoError(t, err)
	var st common.TestStruct = common.TestStruct{}
	var pst *common.TestStruct = &st

	pEco, is := common.CreateOrGetStruct(ctx, reflect.ValueOf(pst))
	require.True(t, is)
	require.NotNil(t, pEco)
	_, is = pEco.(*common.TestStruct)
	require.True(t, is)

	err = common.UnmarshalStructField(pst, "Name", expectedNameBuf)
	require.NoError(t, err)
	require.Equal(t, expectedName, (*pst).Name)
}

func TestGetStructFieldKey1(t *testing.T) {
	expectedData := "name"
	typ := reflect.TypeFor[common.TestStruct]()
	fld, is := typ.FieldByName("Name")
	require.True(t, is)
	key := GetStructFieldKey(fld)
	require.Equal(t, expectedData, key)
}

func TestGetStructFieldKey2(t *testing.T) {
	expectedData := "lucky_number"
	typ := reflect.TypeFor[common.TestStruct]()
	fld, is := typ.FieldByName("LuckyNumber")
	require.True(t, is)
	key := GetStructFieldKey(fld)
	require.Equal(t, expectedData, key)
}

func TestGetStructFieldKey3(t *testing.T) {
	expectedData := "NoTag"
	typ := reflect.TypeFor[common.TestStruct]()
	fld, is := typ.FieldByName("NoTag")
	require.True(t, is)
	key := GetStructFieldKey(fld)
	require.Equal(t, expectedData, key)
}

func TestGetUint(t *testing.T) {
	testGetScalar(t, "/test/uint", uint(13))
}

func TestGetUint2(t *testing.T) {
	testGetScalar2(t, "/test/uint2", uint(13))
}

func TestFieldNameMap(t *testing.T) {
	t.Skip("I honestly don't know what this test is for")
	var ecoTest = common.TestStruct{}
	var p *common.TestStruct = &ecoTest

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
	ctx := common.NewEcoContext(os.Stdout)

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

func TestStructSetField0(t *testing.T) {
	defer func() {
		pa := recover()
		if pa != nil {
			t.Error(pa)
		}
	}()
	var st = common.TestStruct{}
	type pEcoTest *common.TestStruct
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
	var st = common.TestStruct{}
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

func decodeAndTestScalar[U any](t *testing.T, key KeyString, expectedVal U) {
	ctx, cli := initAndTest(t)

	// Seed test data
	putAndTestScalar(t, ctx, cli, key, expectedVal)

	var p *U = nil
	var pp = &p
	err := Decode(ctx, cli, string(key), pp)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, expectedVal, *p)
	t.Log(p)
}

func decodeAndTestScalar2[U any](t *testing.T, key KeyString, expectedVal U) {
	ctx, cli := initAndTest(t)

	// Seed test data
	putAndTestScalar(t, ctx, cli, key, expectedVal)

	var v U
	var p *U = &v
	err := Decode(ctx, cli, string(key), p)
	require.NoError(t, err)
	require.Equal(t, expectedVal, *p)
	t.Log(p)
}


// func getAndTestMap[U any](t *testing.T, ctx *common.EcoContext, expectedData map[string]U, kvs []*mvccpb.KeyValue, key string) {
// 	var pData *map[string]U
// 	rv := reflect.ValueOf(&pData)
// 	decoder := SliceDecoder{}
// 	err := decoder.Decode(ctx, kvs, key, rv)
// 	require.NoError(t, err)
// 	require.Equal(t, expectedData, *pData)
// }

func deleteObjectFromCluster(t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key KeyString, testPrefix KeyString) {
	ctx.Logger.Commentf("Deleting object keys (%s) ...", key)
	var txn etcd.Txn

	// Confirm that key matches a test prefix - this is too dangerous otherwise
	if key.HasPrefix(testPrefix) {
		ctx.Logger.Infof("Key (%s) does not begin with /test/map - possibly unsafe to delete all subkeys.", key)
		return
	}

	// Delete all keys recursively whose prefix matches the object key
	opDelete := etcd.OpDelete(string(key), etcd.WithPrefix())
	txn = createTxn(t, ctx, cli)
	require.NotNil(t, txn)
	txn.Then(opDelete).Commit()
}

func putAndTestScalar(t *testing.T, ctx *common.EcoContext, etcdClient *EtcdClient, key KeyString, i any) {
	// resp, err := etcdClient.Put(context.Background(), key, val)
	ctx.Inc()
	defer ctx.Dec()

	ctx.Logger.Infof("Writing to %s... ", key)
	j, err := json.Marshal(i)
	require.NoError(t, err)
	resp, err := etcdClient.Put(ctx, string(key), string(j))
	require.NoError(t, err)
	require.NotNil(t, resp)

	ctx.Logger.Infof("Reading %s... ", key)
	buf, err := etcdClient.Get(string(key))
	require.NoError(t, err)
	require.Equal(t, j, buf)
	require.Equal(t, string(j), string(buf))
	ctx.Logger.Infof("%#v", resp)
}

func testGetScalar[U any](t *testing.T, key KeyString, expectedVal U) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	ctx.Logger.Comment("Writing scalar seed data to cluster ...")
	putAndTestScalar(t, ctx, cli, key, expectedVal)

	// Create the GET Op
	op := etcd.OpGet(string(key))

	// Get the response from etcd
	ctx.Logger.Comment("Getting scalar value from the cluster ...")
	txn := createTxn(t, ctx, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)

	// Get the KVs from the response
	p := new(U)
	pp := &p
	rangeResp := resp.Responses[0].GetResponseRange()
	etcdKvs := rangeResp.Kvs
	require.Equal(t, 1, len(etcdKvs))
	ctx.Logger.Comment("Done with etcd")
	ctx.Logger.Comment()
	ctx.Logger.Comment("Decoding data ...")
	kvSeries, err := NewKvSeries(key, etcdKvs)
	require.NoError(t, err)
	require.Equal(t, kvSeries.RootKey, key)
	require.Equal(t, 1, len(kvSeries.Kvs))
	tree, err := NewValueTree(ctx, kvSeries)
	require.NoError(t, err)
	require.Equal(t, 0, len(tree.ChildMap))
	buf, err := json.Marshal(expectedVal)
	require.NoError(t, err)
	require.Equal(t, buf, tree.Value)
	// Get the decoder from the DecoderMap and decode
	decoder := &ScalarDecoder[U]{}
	err = decoder.Decode(ctx, tree, pp)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, expectedVal, *p)
}

func testGetScalar2[U any](t *testing.T, key KeyString, expectedVal U) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	putAndTestScalar(t, ctx, cli, key, expectedVal)

	// Create the GET Op
	op := etcd.OpGet(string(key))

	// Get the response from etcd
	txn := createTxn(t, ctx, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)

	// Get the KVs from the response
	var v U
	p := &v
	rangeResp := resp.Responses[0].GetResponseRange()
	etcdKvs := rangeResp.Kvs

	kvSeries, err := NewKvSeries(key, etcdKvs)
	require.NoError(t, err)
	tree, err := NewValueTree(ctx, kvSeries)
	require.NoError(t, err)

	// Get the decoder from the DecoderMap and decode
	decoder := &ScalarDecoder[U]{}
	err = decoder.Decode(ctx, tree, p)
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

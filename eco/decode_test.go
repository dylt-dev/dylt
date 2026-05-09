package eco

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/mvccpb"
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

func TestCreateMapAndUnmarshalValue(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedKey := "foo"
	expectedVal := "bar"
	expectedValBuf, err := json.Marshal(expectedVal)
	require.NoError(t, err)
	var m common.TestMap = common.TestMap{}
	var pm *common.TestMap = &m

	p, is := common.CreateOrGetMap(ctx, reflect.ValueOf(pm))
	require.True(t, is)
	require.NotNil(t, p)
	pTestMap, is := p.(*common.TestMap)
	require.True(t, is)
	require.Equal(t, pm, pTestMap)

	// Get types for map, map key, and map value
	typMap := reflect.TypeFor[common.TestMap]()
	// typKey := typMap.Key()
	typVal := typMap.Elem()

	// Create a new map value
	rvVal := reflect.New(typVal)
	pVal := rvVal.Interface()
	err = json.Unmarshal(expectedValBuf, pVal)
	rvTestMap := reflect.ValueOf(*pTestMap)
	rvTestMap.SetMapIndex(reflect.ValueOf(expectedKey), rvVal.Elem())

	// Check the unmarshal worked
	val, is := m[expectedKey]
	require.True(t, is)
	require.Equal(t, expectedVal, val)

	// err = common.UnmarshalStructField(pst, "Name", expectedNameBuf)
	// require.NoError(t, err)
	// require.Equal(t, expectedName, (*pst).Name)
}

func TestCreateOrGetMapExists(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedKey := "foo"
	expectedVal := "bar"
	expectedData := common.TestMap{
		expectedKey: expectedVal,
	}
	var pExpectedData *common.TestMap = &expectedData

	p, is := common.CreateOrGetMap(ctx, reflect.ValueOf(pExpectedData))
	require.True(t, is)
	require.NotNil(t, p)
	pTestMap, is := p.(*common.TestMap)
	require.True(t, is)
	// Check that the key exists, and its value matches the expected value
	testMap := *pTestMap
	val, is := testMap[expectedKey]
	require.True(t, is)
	require.Equal(t, expectedVal, val)
}

// Pass CreateOtGetStruct a **struct, and confirm it allocates a struct,
// and that the pointer to the struct can be obtained by deferencing **struct
func TestCreateOrGetMapAlloc(t *testing.T) {
	expectedKey := "foo"
	expectedVal := "bar"
	expectedData := common.TestMap{expectedKey: expectedVal}
	ctx := common.NewEcoContext(os.Stdout)
	var pmap *common.TestMap = nil
	ppmap := &pmap

	pNewMap, is := common.CreateOrGetMap(ctx, reflect.ValueOf(ppmap))
	require.True(t, is)
	require.NotNil(t, pNewMap)
	pTestMap, is := pNewMap.(*common.TestMap)
	require.True(t, is)
	require.NotNil(t, pTestMap)
	require.Equal(t, pmap, pTestMap)
	(*pTestMap)[expectedKey] = expectedVal
	require.Equal(t, expectedData, *pTestMap)
	value, is := (**ppmap)[expectedKey]
	require.True(t, is)
	require.Equal(t, expectedVal, value)
}

func TestCreateOrGetMapNil(t *testing.T) {
	expectedKey := "foo"
	expectedVal := "bar"
	expectedData := common.TestMap{expectedKey: expectedVal}
	ctx := common.NewEcoContext(os.Stdout)
	var m common.TestMap = nil
	var pm *common.TestMap = &m

	p, is := common.CreateOrGetMap(ctx, reflect.ValueOf(pm))
	require.True(t, is)
	require.NotNil(t, p)
	pTestMap, is := p.(*common.TestMap)
	require.True(t, is)
	require.NotNil(t, pTestMap)
	require.NotNil(t, *pTestMap)
	require.Equal(t, pm, pTestMap)
	(*pTestMap)[expectedKey] = expectedVal
	require.Equal(t, expectedData, *pTestMap)
	value, is := (*pTestMap)[expectedKey]
	require.True(t, is)
	require.Equal(t, expectedVal, value)
}


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
	decodeAndTestScalar(t, "/test/bool", true)
}

func TestDecodeBool2(t *testing.T) {
	decodeAndTestScalar2(t, "/test/bool2", true)
}

func TestDecodeBoolSlice(t *testing.T) {
	decodeAndTestSlice(t,
		"/test/boolslice",
		[]bool{true, true, false})
}

func TestDecodeFloat(t *testing.T) {
	decodeAndTestScalar(t, "/test/float", float32(42.0))
}

func TestDecodeFloat2(t *testing.T) {
	decodeAndTestScalar2(t, "/test/float2", float32(42.0))
}

func TestDecodeFloatSlice(t *testing.T) {
	decodeAndTestSlice(t,
		"/test/float32slice",
		[]float32{42.0, 1764.0, 6.54321})
}

func TestDecodeInt(t *testing.T) {
	decodeAndTestScalar(t, "/test/int", int(-13.0))
}

func TestDecodeInt2(t *testing.T) {
	decodeAndTestScalar2(t, "/test/int2", int(-13.0))
}

func TestDecodeIntSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	key := "/test/intSlice"
	expectedData := []int{5, 8, 13}
	putAndTestSlice(t, ctx, cli, key, expectedData)

	var pSlice *[]int
	err := Decode(ctx, cli, key, &pSlice)
	require.NoError(t, err)
	require.Equal(t, expectedData, *pSlice)
	t.Log(*pSlice)
}

func TestDecodeString(t *testing.T) {
	decodeAndTestScalar(t, "/test/string", `This\nis\a\<difficult>\nstring\n\to\n\e"s'c"a'p"e\n`)
}

func TestDecodeString2(t *testing.T) {
	decodeAndTestScalar2(t, "/test/string2", `This\nis\a\<difficult>\nstring\n\to\n\e"s'c"a'p"e\n`)
}

func TestDecodeStringSlice(t *testing.T) {
	decodeAndTestSlice(t,
		"/test/stringslice",
		[]string{"foo", "bar", "bum"})
}

func TestDecodeUint(t *testing.T) {
	decodeAndTestScalar(t, "/test/uint", uint(13.0))
}

func TestDecodeUint2(t *testing.T) {
	decodeAndTestScalar2(t, "/test/uint2", uint(13.0))
}

func TestDecodeUintSlice(t *testing.T) {
	decodeAndTestSlice(t,
		"/test/uintslice",
		[]uint{5, 12, 13})
}

func TestGetBool(t *testing.T) {
	testGetScalar(t, "test/bool", true)
}

func TestGetBool2(t *testing.T) {
	testGetScalar2(t, "test/bool2", true)
}

func TestGetBoolSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	ctx.Logger.Comment("Write test data to cluster ...")
	key := "/test/boolSlice"
	expectedData := []bool{true, false, true, true}
	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	ctx.Logger.Comment("Read test data from cluster ...")
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	ctx.Logger.Comment("Done with etcd")
	ctx.Logger.Comment()
	ctx.Logger.Comment("Decoding data ...")
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}

func TestGetFloat(t *testing.T) {
	testGetScalar(t, "/test/float", float32(42.0))
}

func TestGetFloat2(t *testing.T) {
	testGetScalar2(t, "/test/float2", float32(42.0))
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

func TestGetInt2(t *testing.T) {
	testGetScalar2(t, "/test/int2", int(-13))
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

	key := "/test/map"
	expectedBorn := "Venezuela"
	expectedId := "1"
	expectedIsActive := "true"
	expectedName := "Jose Altuve"
	expectedWeight := "160"
	expectedData := map[string]string{
		"Born":     expectedBorn,
		"Id":       expectedId,
		"IsActive": expectedIsActive,
		"Name":     expectedName,
		"Weight":   expectedWeight,
	}
	putAndTestMap(t, ctx, cli, key, expectedData)

	// Create the GET Op
	op := etcd.OpGet(key, etcd.WithPrefix())

	// Get the response from etcd
	txn := createTxn(t, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)

	// Get the etcd KVs from the response
	rangeResp := resp.Responses[0].GetResponseRange()
	etcdKvs := rangeResp.Kvs

	// Create a slice of KeyValue objects
	kvs := createKvSlice(etcdKvs)

	// Create a KeyValueTree from the KVs
	kvTree := createKvTree(ctx, key, kvs)

	// Decode the map
	var pmap *map[string]string
	decoder := MapDecoder{}
	err = decoder.Decode(ctx, kvTree, key, reflect.ValueOf(&pmap))
	require.NoError(t, err)

	// Check the contents
	require.Equal(t, expectedData, *pmap)

	t.Log(*pmap)
}

func TestGetMapOfMaps(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := "/test/stros"
	map0 := map[string]string{"Name": "Altuve", "Position": "2B"}
	map1 := map[string]string{"Name": "Pena", "Position": "SS"}
	map2 := map[string]string{"Name": "Javier", "Position": "P"}
	mapStros := map[int]map[string]string{27: map0, 3: map1, 53: map2}
	ops, err := Encode(common.NewEcoContext(os.Stdout), key, mapStros)
	require.NoError(t, err)

	ctx.Logger.Comment("writing keys ...")
	txn := createTxn(t, cli)
	resp, err := txn.Then(ops...).Commit()
	require.NoError(t, err)
	require.NotNil(t, resp)
	ctx.Logger.Info("done - writing keys")
	ctx.Logger.Info()

	// Create the GET Op
	op := etcd.OpGet(key, etcd.WithPrefix())

	// Get the response from etcd
	txn = createTxn(t, cli)
	resp, err = txn.Then(op).Commit()
	require.NoError(t, err)

	// Get the etcd KVs from the response
	rangeResp := resp.Responses[0].GetResponseRange()
	etcdKvs := rangeResp.Kvs

	// Create a slice of KeyValue objects
	kvs := createKvSlice(etcdKvs)

	// Create a KeyValueTree from the KVs
	kvTree := createKvTree(ctx, key, kvs)

	// Decode the map
	ctx.Logger.Comment("decoding data ...")
	var pmap *map[int]map[string]string
	decoder := MapDecoder{}
	err = decoder.Decode(ctx, kvTree, key, reflect.ValueOf(&pmap))
	ctx.Logger.Info("done - decoding data")
	require.NoError(t, err)
	t.Log(*pmap)
	mTeam := *pmap
	mapAltuve, is := mTeam[27]
	require.True(t, is)
	require.Equal(t, "Altuve", mapAltuve["Name"])
	require.Equal(t, "2B", mapAltuve["Position"])
	mapPena, is := mTeam[3]
	require.True(t, is)
	require.Equal(t, "Pena", mapPena["Name"])
	require.Equal(t, "SS", mapPena["Position"])
	mapJavier, is := mTeam[53]
	require.True(t, is)
	require.Equal(t, "Javier", mapJavier["Name"])
	require.Equal(t, "P", mapJavier["Position"])
}

func TestGetString(t *testing.T) {
	testGetScalar(t, "/test/string", "hello world")
}

func TestGetString2(t *testing.T) {
	testGetScalar2(t, "/test/string2", "hello world")
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
	ctx, _ := initAndTest(t)

	// Setup keys and values
	key := "/test/struct/ecotest"
	expectedData := common.TestStruct{
		Name:        "Me",
		LuckyNumber: 169.0,
		NoTag:       "tagless",
	}
	// expectedName := ex
	// expectedLuckyNumber := 13
	// expectedNoTag := "no-tag-value"
	keyName := fmt.Sprintf("%s/%s", key, "name")
	keyLuckyNumber := fmt.Sprintf("%s/%s", key, "lucky_number")
	keyNoTag := fmt.Sprintf("%s/%s", key, "NoTag")

	// Encode []byte values for struct fields
	bufName, err := json.Marshal(expectedData.Name)
	require.NoError(t, err)
	bufLuckyNumber, err := json.Marshal(expectedData.LuckyNumber)
	require.NoError(t, err)
	bufNoTag, err := json.Marshal(expectedData.NoTag)
	require.NoError(t, err)
	etcdKvs := []*mvccpb.KeyValue{
		{Key: []byte(keyName), Value: bufName},
		{Key: []byte(keyLuckyNumber), Value: bufLuckyNumber},
		{Key: []byte(keyNoTag), Value: bufNoTag},
	}
	kvs := createKvSlice(etcdKvs)
	kvTree := createKvTree(ctx, key, kvs)
	// putAndTestStruct(t, ctx, cli, key, kvs)

	data := common.TestStruct{}
	rv := reflect.ValueOf(&data)
	decoder := StructDecoder{}
	decoder.Decode(ctx, kvTree, key, rv)
	t.Logf("%#v", data)

	// expectedData := EcoTest{
	// 	Name:        expectedName,
	// 	LuckyNumber: float64(expectedLuckyNumber),
	// 	NoTag:       expectedNoTag,
	// }
	// var data *EcoTest
	// err = Decode(ctx, cli, key, &expectedData)
	// require.NoError(t, err)
	// // require.Equal(t, expectedData, *data)
	// require.Nil(t, data)
	// t.Log(*decodedVal)
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

func decodeAndTestScalar[U any](t *testing.T, key string, expectedVal U) {
	ctx, cli := initAndTest(t)

	// Seed test data
	putAndTestScalar(t, ctx, cli, key, expectedVal)

	var p *U
	var i = &p
	err := Decode(ctx, cli, key, i)
	require.NoError(t, err)
	require.Equal(t, expectedVal, *p)
	t.Log(p)
}

func decodeAndTestScalar2[U any](t *testing.T, key string, expectedVal U) {
	ctx, cli := initAndTest(t)

	// Seed test data
	putAndTestScalar(t, ctx, cli, key, expectedVal)

	var v U
	var p *U = &v
	err := Decode(ctx, cli, key, p)
	require.NoError(t, err)
	require.Equal(t, expectedVal, *p)
	t.Log(p)
}

func decodeAndTestSlice[U any](t *testing.T, key string, expectedData []U) {
	ctx, cli := initAndTest(t)

	putAndTestSlice(t, ctx, cli, key, expectedData)

	var pSlice *[]U
	err := Decode(ctx, cli, key, &pSlice)
	require.NoError(t, err)
	require.Equal(t, expectedData, *pSlice)
	t.Log(*pSlice)

}

// func getAndTestMap[U any](t *testing.T, ctx *common.EcoContext, expectedData map[string]U, kvs []*mvccpb.KeyValue, key string) {
// 	var pData *map[string]U
// 	rv := reflect.ValueOf(&pData)
// 	decoder := SliceDecoder{}
// 	err := decoder.Decode(ctx, kvs, key, rv)
// 	require.NoError(t, err)
// 	require.Equal(t, expectedData, *pData)
// }

func deleteObjectFromCluster(t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key string, testPrefix string) {
	ctx.Logger.Commentf("Deleting object keys (%s) ...", key)
	var txn etcd.Txn

	// Confirm that key matches a test prefix - this is too dangerous otherwise
	if !strings.HasPrefix(key, testPrefix) {
		ctx.Logger.Infof("Key (%s) does not begin with /test/map - possibly unsafe to delete all subkeys.", key)
		return
	}

	// Delete all keys recursively whose prefix matches the object key
	opDelete := etcd.OpDelete(key, etcd.WithPrefix())
	txn = createTxn(t, cli)
	require.NotNil(t, txn)
	txn.Then(opDelete).Commit()
}

func getAndTestSlice[U any](t *testing.T, ctx *common.EcoContext, expectedData []U, etcdKvs []*mvccpb.KeyValue, key string) {
	var pSlice *[]U
	rv := reflect.ValueOf(&pSlice)
	decoder := SliceDecoder{}

	kvs := createKvSlice(etcdKvs)
	kvTree := createKvTree(ctx, key, kvs)
	ctx.Logger.Debugf("kvTree.Children=%#v", kvTree.Children)

	err := decoder.Decode(ctx, kvTree, key, rv)
	require.NoError(t, err)
	require.Equal(t, expectedData, *pSlice)
}

func getAndTestSliceKVs(t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key string) []*mvccpb.KeyValue {
	op := etcd.OpGet(key, etcd.WithPrefix())
	txn := createTxn(t, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)
	require.NotNil(t, resp)
	respRange := resp.Responses[0].GetResponseRange()
	ctx.Logger.Infof("respRange.Count=%d", respRange.Count)
	kvs := respRange.Kvs
	return kvs
}

// With the EtcdClient, Put a value to etcd, then Get it back to confirm the
// Put succeeded
func putAndTestMap[U any](t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key string, data map[string]U) {
	deleteObjectFromCluster(t, ctx, cli, key, "/test/map")

	ctx.Logger.Infof("Writing map at %s ...", key)
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

func putAndTestScalar(t *testing.T, ctx *common.EcoContext, etcdClient *EtcdClient, key string, i any) {
	// resp, err := etcdClient.Put(context.Background(), key, val)
	ctx.Inc()
	defer ctx.Dec()

	ctx.Logger.Infof("Writing to %s... ", key)
	j, err := json.Marshal(i)
	require.NoError(t, err)
	resp, err := etcdClient.Put(ctx, key, string(j))
	require.NoError(t, err)
	require.NotNil(t, resp)

	ctx.Logger.Infof("Reading %s... ", key)
	buf, err := etcdClient.Get(key)
	require.NoError(t, err)
	require.Equal(t, j, buf)
	require.Equal(t, string(j), string(buf))
	ctx.Logger.Infof("%#v", resp)
}

func putAndTestSlice[U any](t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key string, data []U) {
	ctx.Logger.Infof("Writing slice at %s ...", key)
	ops := []etcd.Op{}
	for i, val := range data {
		ctx.Inc()
		subkey := fmt.Sprintf("%s/%d", key, i)
		bufVal, err := json.Marshal(val)
		require.NoError(t, err)
		ctx.Logger.Infof("%s => %s", subkey, string(bufVal))
		op := etcd.OpPut(subkey, string(bufVal))
		ops = append(ops, op)
		ctx.Dec()
	}
	txn := createTxn(t, cli)
	require.NotNil(t, txn)
	txn.Then(ops...).Commit()
}

func putAndTestStruct(t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key string, kvs []*mvccpb.KeyValue) {
	deleteObjectFromCluster(t, ctx, cli, key, "/test/struct")

	ctx.Logger.Infof("Writing struct at %s ...", key)
	ops := []etcd.Op{}
	for _, kv := range kvs {
		subkey := fmt.Sprintf("%s/%s", key, kv.Key)
		op := etcd.OpPut(subkey, string(kv.Value))
		ops = append(ops, op)
	}
	txn := createTxn(t, cli)
	require.NotNil(t, txn)
	txn.Then(ops...).Commit()
}

func testGetScalar[U any](t *testing.T, key string, expectedVal U) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	ctx.Logger.Comment("Writing scalar seed data to cluster ...")
	putAndTestScalar(t, ctx, cli, key, expectedVal)

	// Create the GET Op
	op := etcd.OpGet(key)

	// Get the response from etcd
	ctx.Logger.Comment("Getting scalar value from the cluster ...")
	txn := createTxn(t, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)

	// Get the KVs from the response
	p := new(U)
	pp := &p
	rv := reflect.ValueOf(pp)
	rangeResp := resp.Responses[0].GetResponseRange()
	etcdKvs := rangeResp.Kvs

	ctx.Logger.Comment("Done with etcd")
	ctx.Logger.Comment()
	ctx.Logger.Comment("Decoding data ...")
	kvs := createKvSlice(etcdKvs)
	kvTree := createKvTree(ctx, key, kvs)

	// Get the decoder from the DecoderMap and decode
	decoder := &ScalarDecoder[U]{}
	err = decoder.Decode(ctx, kvTree, key, rv)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, expectedVal, *p)
}

func testGetScalar2[U any](t *testing.T, key string, expectedVal U) {
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
	var v U
	p := &v
	rv := reflect.ValueOf(p)
	rangeResp := resp.Responses[0].GetResponseRange()
	etcdKvs := rangeResp.Kvs

	kvs := createKvSlice(etcdKvs)
	kvTree := createKvTree(ctx, key, kvs)

	// Get the decoder from the DecoderMap and decode
	decoder := &ScalarDecoder[U]{}
	err = decoder.Decode(ctx, kvTree, key, rv)
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

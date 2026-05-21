package eco

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
)

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


func TestDecodeMap1 (t *testing.T) {
	expected := map[string]int64 { "foo": 13, "bar": 169 }

	decodeAndTestMap(t, expected)
}


func TestDecodeMap2 (t *testing.T) {
	expected := map[string]int64 { "foo": 13, "bar": 169 }

	decodeAndTestNilMap(t, expected)
}


func TestGetMap(t *testing.T) {
	ctx, cli := initAndTest(t)

	key := KeyString("/test/map")
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
	op := etcd.OpGet(string(key), etcd.WithPrefix())

	// Get the response from etcd
	txn := createTxn(t, ctx, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)

	// Get the etcd KVs from the response
	rangeResp := resp.Responses[0].GetResponseRange()
	etcdKvs := rangeResp.Kvs

	kvSeries, err := NewKvSeries(key, etcdKvs)
	require.NoError(t, err)
	tree, err := NewValueTree(ctx, kvSeries)
	require.NoError(t, err)

	// Decode the map
	var p *map[string]string = nil
	pp := &p
	decoder := MapDecoder{}
	err = decoder.Decode(ctx, tree, pp)
	require.NoError(t, err)
	require.NotNil(t, p)

	// Check the contents
	require.Equal(t, expectedData, *p)

	t.Log(*p)
}

func TestGetMapOfMaps(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := KeyString("/test/stros")
	map0 := map[string]string{"Name": "Altuve", "Position": "2B"}
	map1 := map[string]string{"Name": "Pena", "Position": "SS"}
	map2 := map[string]string{"Name": "Javier", "Position": "P"}
	mapStros := map[int]map[string]string{27: map0, 3: map1, 53: map2}
	ops, err := Encode(common.NewEcoContext(os.Stdout), string(key), mapStros)
	require.NoError(t, err)

	ctx.Logger.Comment("writing keys ...")
	txn := createTxn(t, ctx, cli)
	resp, err := txn.Then(ops...).Commit()
	require.NoError(t, err)
	require.NotNil(t, resp)
	ctx.Logger.Info("done - writing keys")
	ctx.Logger.Info()

	// Create the GET Op
	op := etcd.OpGet(string(key), etcd.WithPrefix())

	// Get the response from etcd
	txn = createTxn(t, ctx, cli)
	resp, err = txn.Then(op).Commit()
	require.NoError(t, err)

	// Get the etcd KVs from the response
	rangeResp := resp.Responses[0].GetResponseRange()
	etcdKvs := rangeResp.Kvs

	// Create a ValueTree
	kvSeries, err := NewKvSeries(key, etcdKvs)
	require.NoError(t, err)
	tree, err := NewValueTree(ctx, kvSeries)
	require.NoError(t, err)

	// Decode the map
	ctx.Logger.Comment("decoding data ...")
	var p *map[int]map[string]string = nil
	pp := &p 
	decoder := MapDecoder{}
	err = decoder.Decode(ctx, tree, pp)
	ctx.Logger.Info("done - decoding data")
	require.NoError(t, err)
	require.NotNil(t, p)
	t.Log(*p)
	mTeam := *p
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


func decodeAndTestMap (t *testing.T, expected map[string]int64) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := MapDecoder{}
	
	tree := &ValueTree{}
	for k, v := range expected {
		key := fmt.Sprintf("/%s", k)
		tree.addInt(ctx, key, v)
	}
	require.Equal(t, len(expected), len(tree.ChildMap))

	var x map[string]int64 = nil
	p := &x
	err := decoder.Decode(ctx, tree, p)

	require.NoError(t, err)
	require.Equal(t, len(expected), len(x))
	require.Equal(t, expected, x)
}


func decodeAndTestNilMap (t *testing.T, expected map[string]int64) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := MapDecoder{}
	
	tree := &ValueTree{}
	for k, v := range expected {
		key := fmt.Sprintf("/%s", k)
		tree.addInt(ctx, key, v)
	}
	require.Equal(t, len(expected), len(tree.ChildMap))

	var x map[string]int64 = nil
	p := &x
	err := decoder.Decode(ctx, tree, p)

	require.NoError(t, err)
	require.Equal(t, len(expected), len(x))
	require.Equal(t, expected, x)

}


// With the EtcdClient, Put a value to etcd, then Get it back to confirm the
// Put succeeded
func putAndTestMap[U any](t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key KeyString, data map[string]U) {
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
	txn := createTxn(t, ctx, cli)
	require.NotNil(t, txn)
	txn.Then(ops...).Commit()
}

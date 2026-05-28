package eco

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
)

func TestGetStruct1(t *testing.T) {
	// test data
	rootKey := KeyString("/test/struct")
	expected := common.TestStruct{
		Name:        fmt.Sprintf("%s", "meat"),
		LuckyNumber: 13,
		NoTag:       fmt.Sprintf("%s", "tagless"),
	}
	getAndTestStruct(t, expected, rootKey)

	/*
	   // create OpPuts for struct
	   rootKey := KeyString("/test/struct")
	   rvs := rvStruct(reflect.ValueOf(expected))
	   ops := rvs.Ops(t, ctx, rootKey)

	   // Write OpPuts to etcd
	   putOps(t, ctx, cli, ops)

	   // Read back OpPuts as etcd.KVs
	   etcdKvs := getAndTestKVs(t, ctx, cli, rootKey)

	   // Convert etcd.KVs into KvSeries
	   kvSeries, err := NewKvSeries(rootKey, etcdKvs)
	   require.NoError(t, err)

	   // Convert KvSeries into ValueTree
	   tree, err := NewValueTree(ctx, kvSeries)
	   require.NoError(t, err)

	   // Decode
	   decoder := StructDecoder{}
	   var x common.TestStruct
	   p := &x
	   err = decoder.Decode(ctx, tree, p)
	   require.NoError(t, err)
	   require.Equal(t, expected, x)
	*/
}

func TestStructDecoder1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)

	decoder := StructDecoder{}
	tree := &ValueTree{}
	tree.add(ctx, "/name", "Smitty")
	tree.add(ctx, "/lucky_number", 13)
	tree.add(ctx, "/NoTag", "tagless")

	var x common.TestStruct
	p := &x
	err := decoder.Decode(ctx, tree, p)

	require.NoError(t, err)
	require.Equal(t, "Smitty", x.Name)
	require.Equal(t, float64(13), x.LuckyNumber)
	require.Equal(t, "tagless", x.NoTag)
}

func TestStructDecoder2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)

	decoder := StructDecoder{}
	tree := &ValueTree{}
	tree.add(ctx, "/name", "Smitty")
	tree.add(ctx, "/lucky_number", 13)
	tree.add(ctx, "/NoTag", "tagless")

	var p *common.TestStruct = nil
	pp := &p
	err := decoder.Decode(ctx, tree, pp)
	x := *p

	require.NoError(t, err)
	require.Equal(t, "Smitty", x.Name)
	require.Equal(t, float64(13), x.LuckyNumber)
	require.Equal(t, "tagless", x.NoTag)
}

func TestStructEcoTest(t *testing.T) {
	ctx, _ := initAndTest(t)

	// Setup keys and values
	key := KeyString("/test/struct/ecotest")
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
	kvSeries, err := NewKvSeries(key, etcdKvs)
	require.NoError(t, err)
	tree, err := NewValueTreeFromKvSeries(ctx, kvSeries)
	require.NoError(t, err)

	x := common.TestStruct{}
	p := &x
	decoder := StructDecoder{}
	decoder.Decode(ctx, tree, p)
	t.Logf("%#v", x)
}

func getAndTestStruct(t *testing.T, expected any, rootKey KeyString) {
	ctx, cli := initAndTest(t)

	// create OpPuts for struct
	rvs, err := NewRvStruct(expected)
	ops := rvs.Ops(t, ctx, rootKey)

	// Write OpPuts to etcd
	putOps(t, ctx, cli, ops)

	// Read back OpPuts as etcd.KVs
	etcdKvs := getAndTestKVs(t, ctx, cli, rootKey)

	// Convert etcd.KVs into KvSeries
	kvSeries, err := NewKvSeries(rootKey, etcdKvs)
	require.NoError(t, err)

	// Convert KvSeries into ValueTree
	tree, err := NewValueTreeFromKvSeries(ctx, kvSeries)
	require.NoError(t, err)

	// Decode
	decoder := StructDecoder{}

	// allocate a pointer to a struct of the same type as expected
	typStruct := reflect.ValueOf(expected).Type()
	rvPtr := reflect.New(typStruct)

	// decode into the newly allocated pointer
	p := rvPtr.Interface()
	err = decoder.Decode(ctx, tree, p)
	require.NoError(t, err)

	// Get the new struct via the allocated pointer + test
	x := rvPtr.Elem().Interface()
	require.Equal(t, expected, x)
}

// func putAndTestStruct(t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key KeyString, kvs []*mvccpb.KeyValue) {
// 	deleteObjectFromCluster(t, ctx, cli, key, "/test/struct")

// 	ctx.Logger.Infof("Writing struct at %s ...", key)
// 	ops := []etcd.Op{}
// 	for _, kv := range kvs {
// 		subkey := fmt.Sprintf("%s/%s", key, kv.Key)
// 		op := etcd.OpPut(subkey, string(kv.Value))
// 		ops = append(ops, op)
// 	}
// 	txn := createTxn(t, ctx, cli)
// 	require.NotNil(t, txn)
// 	txn.Then(ops...).Commit()
// }

type RvStruct reflect.Value

func NewRvStruct(a any) (*RvStruct, error) {
	rv := common.Reflect(a)
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expecting struct, got %s", rv.Kind().String())
	}
	rvs := RvStruct(rv)
	return &rvs, nil
}

func (rvs RvStruct) Ops(t *testing.T, ctx *common.EcoContext, rootKey KeyString) []etcd.Op {
	ops := []etcd.Op{}
	rv := reflect.Value(rvs)
	ctx.Logger.Infof("rv.Type()=%s", rv.Type())
	for structField := range rv.Type().Fields() {
		ctx.Logger.Infof("structField.Name=%s", structField.Name)
		key := fmt.Sprintf("%s/%s", rootKey, GetStructFieldKey(structField))
		rvFieldVal := rv.FieldByName(structField.Name)
		buf, err := json.Marshal(rvFieldVal.Interface())
		require.NoError(t, err)
		val := string(buf)
		ctx.Logger.Infof("key=%s val=%s", key, val)
		ops = append(ops, etcd.OpPut(key, val))
	}

	return ops
}

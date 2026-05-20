package eco

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
)


func TestDecodeBoolSlice(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected0 := true
	expected1 := true
	expected9 := true
	decoder := MainDecoder{}

	tree := &ValueTree{}
	tree.addBool(ctx, "/0", expected0)
	tree.addBool(ctx, "/1", expected1)
	tree.addBool(ctx, "/9", expected9)
	var p *[]bool = nil
	pp := &p
	
	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	require.Equal(t, expected0, x[0])
	require.Equal(t, expected1, x[1])
	require.Equal(t, expected9, x[9])
	require.False(t, x[2])
}


func TestDecodeFloatSlice(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected0 := 42.0
	expected1 := 1764.0
	expected9 := 6.54321
	decoder := MainDecoder{}

	tree := &ValueTree{}
	tree.addFloat(ctx, "/0", expected0)
	tree.addFloat(ctx, "/1", expected1)
	tree.addFloat(ctx, "/9", expected9)
	var p *[]float64 = nil
	pp := &p
	
	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	require.Equal(t, expected0, x[0])
	require.Equal(t, expected1, x[1])
	require.Equal(t, expected9, x[9])
	require.Zero(t, x[2])

}


func TestDecodeIntSlice(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var expected0 int64 = 13
	var expected1 int64 = 169
	var expected9 int64 = -1997
	decoder := MainDecoder{}

	tree := &ValueTree{}
	tree.addInt(ctx, "/0", expected0)
	tree.addInt(ctx, "/1", expected1)
	tree.addInt(ctx, "/9", expected9)
	var p *[]int64 = nil
	pp := &p
	
	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	require.Equal(t, expected0, x[0])
	require.Equal(t, expected1, x[1])
	require.Equal(t, expected9, x[9])
	require.Zero(t, x[2])
}


func TestDecodeIntSlice1000(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := 1000
	expected := make([]int64, n)
	for i := range expected {
		expected[i] = int64(i*10)
	}
	decoder := MainDecoder{}

	tree := &ValueTree{}
	for i, val := range expected {
		key := fmt.Sprintf("/%d", i)
		tree.addInt(ctx, key, val)
	}
	var p *[]int64 = nil
	pp := &p

	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	for i, val := range x {
		require.Equal(t, expected[i], val)
	}
}


func TestDecodeStringSlice(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected0 := "foo"
	expected1 := "bar"
	expected9 := "bum"
	decoder := MainDecoder{}

	tree := &ValueTree{}
	tree.addString(ctx, "/0", expected0)
	tree.addString(ctx, "/1", expected1)
	tree.addString(ctx, "/9", expected9)
	var p *[]string = nil
	pp := &p
	
	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	require.Equal(t, expected0, x[0])
	require.Equal(t, expected1, x[1])
	require.Equal(t, expected9, x[9])
	require.Zero(t, x[2])
}


func TestDecodeUintSlice(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected0 := int64(13)
	expected1 := int64(169)
	expected9 := int64(1997)
	decoder := MainDecoder{}

	tree := &ValueTree{}
	tree.addInt(ctx, "/0", expected0)
	tree.addInt(ctx, "/1", expected1)
	tree.addInt(ctx, "/9", expected9)
	var p *[]int64 = nil
	pp := &p
	
	err := decoder.Decode(ctx, tree, pp)
	x := *p
	require.NoError(t, err)
	require.Equal(t, expected0, x[0])
	require.Equal(t, expected1, x[1])
	require.Equal(t, expected9, x[9])
	require.Zero(t, x[2])
}


func TestGetBoolSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	ctx.Logger.Comment("Write test data to cluster ...")
	key := KeyString("/test/boolSlice")
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


func TestGetFloatSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := KeyString("/test/floatSlice")
	expectedData := []float32{42.0, 1764.0, 6.54321}

	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}


func TestGetIntSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := KeyString("/test/intSlice")
	expectedData := []int{5, 8, 13}

	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}


func TestGetStringSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := KeyString("/test/stringSlice")
	expectedData := []string{"foo", "bar", "bum"}
	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}


func TestGetUintSlice(t *testing.T) {
	ctx, cli := initAndTest(t)

	// Seed etcd with test data
	key := KeyString("/test/uintSlice")
	expectedData := []uint{5, 12, 13}

	putAndTestSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	kvs := getAndTestSliceKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}


func TestSliceDecoder1 (t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := SliceDecoder{}
	tree := &ValueTree{}
	tree.Add(ctx, KeyString("/0"), []byte("\"foo\""))
	tree.Add(ctx, KeyString("/1"), []byte("\"bar\""))
	tree.Add(ctx, KeyString("/9"), []byte("\"bum\""))
	require.Equal(t, 3, len(tree.ChildMap))

	var x []string = nil
	p := &x

	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, 10, len(x))
	require.Equal(t, "foo", x[0])
	require.Equal(t, "bar", x[1])
	require.Equal(t, "", x[2])
	require.Equal(t, "bum", x[9])
}


func decodeAndTestSlice[U any](t *testing.T, key KeyString, expectedData []U) {
	ctx, cli := initAndTest(t)

	putAndTestSlice(t, ctx, cli, key, expectedData)

	var p*[]U = nil
	pp := &p
	err := Decode(ctx, cli, string(key), pp)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, expectedData, pp)
	t.Log(p)

}


func getAndTestSlice[U any](t *testing.T, ctx *common.EcoContext, expectedData []U, etcdKvs []*mvccpb.KeyValue, key KeyString) {
	decoder := SliceDecoder{}

	kvSeries, err := NewKvSeries(key, etcdKvs)
	require.NoError(t, err)
	tree, err := NewValueTree(ctx, kvSeries)
	require.NoError(t, err)
	ctx.Logger.Debugf("tree.ChildMap=%#v", tree.ChildMap)

	var p *[]U = nil
	pp := &p
	err = decoder.Decode(ctx, tree, pp)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.NoError(t, err)
	require.Equal(t, expectedData, *p)
}


func getAndTestSliceKVs(t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key KeyString) []*mvccpb.KeyValue {
	op := etcd.OpGet(string(key), etcd.WithPrefix())
	txn := createTxn(t, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)
	require.NotNil(t, resp)
	respRange := resp.Responses[0].GetResponseRange()
	ctx.Logger.Infof("respRange.Count=%d", respRange.Count)
	kvs := respRange.Kvs
	return kvs
}


func putAndTestSlice[U any](t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key KeyString, data []U) {
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

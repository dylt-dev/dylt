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
		expected[i] = int64(i * 10)
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
	putSlice(t, ctx, cli, key, expectedData)

	// Get kvs for seeded data
	ctx.Logger.Comment("Read test data from cluster ...")
	kvs := getAndTestKVs(t, ctx, cli, key)

	// Decode the slice and test expected values
	ctx.Logger.Comment("Done with etcd")
	ctx.Logger.Comment()
	ctx.Logger.Comment("Decoding data ...")
	getAndTestSlice(t, ctx, expectedData, kvs, key)
}


func TestGetFloatSlice(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := KeyString("/test/floatSlice")
	expected := []float32{42.0, 1764.0, 6.54321}
	testSlice(t, ctx, cli, key, expected)
}


func TestGetIntSlice(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := KeyString("/test/intSlice")
	expected := []int{5, 8, 13}
	testSlice(t, ctx, cli, key, expected)
}


func TestGetStringSlice(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := KeyString("/test/stringSlice")
	expected := []string{"foo", "bar", "bum"}
	testSlice(t, ctx, cli, key, expected)
}


func TestGetUintSlice(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := KeyString("/test/uintSlice")
	expectedData := []uint{5, 12, 13}
	testSlice(t, ctx, cli, key, expectedData)
}


func TestSliceDecoder1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := SliceDecoder{}
	tree := &ValueTree{}
	tree.Add(ctx, KeyString("/0"), []byte("\"foo\""))
	tree.Add(ctx, KeyString("/1"), []byte("\"bar\""))
	tree.Add(ctx, KeyString("/9"), []byte("\"bum\""))

	var x []string = nil
	p := &x

	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, 10, len(x))
	require.Equal(t, "foo", x[0])
	require.Equal(t, "bar", x[1])
	require.Zero(t, x[2])
	require.Equal(t, "bum", x[9])
}


func decodeAndTestSlice[U any](t *testing.T, key KeyString, expectedData []U) {
	ctx, cli := initAndTest(t)

	putSlice(t, ctx, cli, key, expectedData)

	var p *[]U = nil
	pp := &p
	err := Decode(ctx, cli, string(key), pp)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, expectedData, pp)
	t.Log(p)

}


// Convert KVs into a tree, decode, and compare against expected data
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


// Confirm that kvs written to etcd (eg by `putSlice()`) were successfully written
func getAndTestKVs(t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key KeyString) []*mvccpb.KeyValue {
	op := etcd.OpGet(string(key), etcd.WithPrefix())
	txn := createTxn(t, ctx, cli)
	resp, err := txn.Then(op).Commit()
	require.NoError(t, err)
	require.NotNil(t, resp)
	respRange := resp.Responses[0].GetResponseRange()
	ctx.Logger.Infof("respRange.Count=%d", respRange.Count)
	kvs := respRange.Kvs
	return kvs
}


// Convert slices to etcd Ops and write them
func putSlice[U any](t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key KeyString, data []U) {
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
	putOps(t, ctx, cli, ops)
}

func testSlice[U any](t *testing.T, ctx *common.EcoContext, cli *EtcdClient, key KeyString, expected []U) {
	// Write slice data to etcd
	ctx.Logger.Commentf("Writing test data to cluster (%s) ...", key)
	putSlice(t, ctx, cli, key, expected)

	// Retrieve slice data
	ctx.Logger.Commentf("Reading test data from cluster (%s) ...", key)
	kvs := getAndTestKVs(t, ctx, cli, key)

	// Decode slice data and confirm it matches original
	ctx.Logger.Comment("Done with etcd")
	ctx.Logger.Comment()
	ctx.Logger.Comment("Decoding data ...")
	getAndTestSlice(t, ctx, expected, kvs, key)
}

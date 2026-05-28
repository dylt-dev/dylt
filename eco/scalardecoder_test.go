package eco

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
)

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
	tree.set(expected)
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
	tree.set(expected)
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
	tree.set(expected)
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
	tree.set(expected)
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
	tree.set(expected)
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
	tree.set(expected)
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
	tree.set(expected)
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
	tree.set(expected)
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
	tree.set(expected)
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

func TestGetUint(t *testing.T) {
	testGetScalar(t, "/test/uint", uint(13))
}

func TestGetUint2(t *testing.T) {
	testGetScalar2(t, "/test/uint2", uint(13))
}

func TestScalarDecoderBool(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := ScalarDecoder[bool]{}
	tree := &ValueTree{}
	tree.Value = []byte(strconv.AppendBool([]byte{}, true))
	var x bool
	p := &x
	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, true, x)
}

func TestScalarDecoderBoolNil(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := ScalarDecoder[bool]{}
	tree := &ValueTree{}
	tree.Value = []byte(strconv.AppendBool([]byte{}, true))
	var p *bool
	pp := &p
	err := decoder.Decode(ctx, tree, pp)
	require.NoError(t, err)
	require.Equal(t, true, **pp)
}

func TestUnmarshalPp(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *bool = nil
	var pp = &p
	rvp, err := NewRvPointer(pp)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, 0)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	require.NotNil(t, normPtr.Value)
	require.NotNil(t, *pp)
	buf := strconv.AppendBool([]byte{}, true)
	err = json.Unmarshal(buf, p)
	require.NoError(t, err)
	require.Equal(t, true, *p)
}

// func decodeAndTestScalar[U any](t *testing.T, key KeyString, expectedVal U) {
// 	ctx, cli := initAndTest(t)

// 	// Seed test data
// 	putAndTestScalar(t, ctx, cli, key, expectedVal)

// 	var p *U = nil
// 	var pp = &p
// 	err := Decode(ctx, cli, string(key), pp)
// 	require.NoError(t, err)
// 	require.NotNil(t, p)
// 	require.Equal(t, expectedVal, *p)
// 	t.Log(p)
// }

// func decodeAndTestScalar2[U any](t *testing.T, key KeyString, expectedVal U) {
// 	ctx, cli := initAndTest(t)

// 	// Seed test data
// 	putAndTestScalar(t, ctx, cli, key, expectedVal)

// 	var v U
// 	var p *U = &v
// 	err := Decode(ctx, cli, string(key), p)
// 	require.NoError(t, err)
// 	require.Equal(t, expectedVal, *p)
// 	t.Log(p)
// }

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
	tree, err := NewValueTreeFromKvSeries(ctx, kvSeries)
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
	tree, err := NewValueTreeFromKvSeries(ctx, kvSeries)
	require.NoError(t, err)

	// Get the decoder from the DecoderMap and decode
	decoder := &ScalarDecoder[U]{}
	err = decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, expectedVal, *p)
}

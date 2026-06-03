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

func TestEncode1(t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ string
	x := typ("foo")

	key := "/test/scalars/scalar1"
	opDelete := etcd.OpDelete(string(key), etcd.WithPrefix())
	txn := cli.Txn(ctx)
	_, err := txn.Then(opDelete).Commit()
	require.NoError(t, err)
	err = Encode(ctx, cli, x, key)
	require.NoError(t, err)

	// decode. cmon we can do this
	var y typ
	p := &y
	err = Decode(ctx, cli, key, p)
	require.NoError(t, err)
	require.Equal(t, x, y)
}

func TestEncode2(t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ int
	x := typ(13)

	key := "/test/scalars/scalar2"
	opDelete := etcd.OpDelete(string(key), etcd.WithPrefix())
	txn := cli.Txn(ctx)
	_, err := txn.Then(opDelete).Commit()
	require.NoError(t, err)
	err = Encode(ctx, cli, x, key)
	require.NoError(t, err)

	// decode. cmon we can do this
	var y typ
	p := &y
	err = Decode(ctx, cli, key, p)
	require.NoError(t, err)
	require.Equal(t, x, y)
}

func TestEncode3(t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ bool
	x := typ(true)

	key := "/test/scalars/scalar3"
	opDelete := etcd.OpDelete(string(key), etcd.WithPrefix())
	txn := cli.Txn(ctx)
	_, err := txn.Then(opDelete).Commit()
	require.NoError(t, err)
	err = Encode(ctx, cli, x, key)
	require.NoError(t, err)

	// decode. cmon we can do this
	var y typ
	p := &y
	err = Decode(ctx, cli, key, p)
	require.NoError(t, err)
	require.Equal(t, x, y)
}

func TestEncode4(t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ map[string]int
	x := typ{
		"foo": 13,
		"bar": 169,
		"bum": 1997,
	}

	key := "/test/maps/map1"
	opDelete := etcd.OpDelete(string(key), etcd.WithPrefix())
	txn := cli.Txn(ctx)
	_, err := txn.Then(opDelete).Commit()
	require.NoError(t, err)
	err = Encode(ctx, cli, x, key)
	require.NoError(t, err)

	// decode. cmon we can do this
	var y typ
	p := &y
	err = Decode(ctx, cli, key, p)
	require.NoError(t, err)
	require.Equal(t, x, y)
}

func TestEncode5(t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ []string
	x := typ{"foo", "bar", "bum"}

	key := "/test/slices/slice1"
	opDelete := etcd.OpDelete(string(key), etcd.WithPrefix())
	txn := cli.Txn(ctx)
	_, err := txn.Then(opDelete).Commit()
	require.NoError(t, err)
	err = Encode(ctx, cli, x, key)
	require.NoError(t, err)

	// decode. cmon we can do this
	var y typ
	p := &y
	err = Decode(ctx, cli, string(key), p)
	require.NoError(t, err)
	require.Equal(t, x, y)
}

func TestEncode6(t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ struct{ Val int }
	x := typ{Val: 13}

	key := "/test/structs/struct1"
	opDelete := etcd.OpDelete(string(key), etcd.WithPrefix())
	txn := cli.Txn(ctx)
	_, err := txn.Then(opDelete).Commit()
	require.NoError(t, err)
	err = Encode(ctx, cli, x, key)
	require.NoError(t, err)

	// decode. cmon we can do this
	var y typ
	p := &y
	err = Decode(ctx, cli, string(key), p)
	require.NoError(t, err)
	require.Equal(t, x, y)
}

func TestEncode7(t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ struct{ Vals []int }
	x := typ{Vals: []int{13, 169, 1997}}

	key := "/test/structs/struct2"
	opDelete := etcd.OpDelete(string(key), etcd.WithPrefix())
	txn := cli.Txn(ctx)
	_, err := txn.Then(opDelete).Commit()
	require.NoError(t, err)
	err = Encode(ctx, cli, x, key)
	require.NoError(t, err)

	// decode. cmon we can do this
	var y typ
	p := &y
	err = Decode(ctx, cli, string(key), p)
	require.NoError(t, err)
	require.Equal(t, x, y)
}

func TestEncode8(t *testing.T) {
	/*
	   ctx, cli := initAndTest(t)
	   type typ struct{ Vals []string }
	   type typ0 string
	   type typ1 []typ0
	   type typ2 struct{ Vals typ1 }
	   x0 := "foo"
	   x1 := typ1{"foo", "bar", "bum"}
	   x2 := typ2{Vals: x1}
	   x := x2
	   key := "/test/structs/struct3"
	   var y typ
	*/
}

func TestEncode10(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)

	type typ [][]struct {
		Tempora struct {
			Eum map[string]struct{ Dolorem map[int][]map[bool][]string }
		}
	}

	x0 := "meat"
	x1 := []string{x0}
	x2 := map[bool][]string{true: x1}
	x3 := []map[bool][]string{x2}
	x4 := map[int][]map[bool][]string{13: x3}
	x5 := struct{ Dolorem map[int][]map[bool][]string }{Dolorem: x4}
	x6 := map[string]struct{ Dolorem map[int][]map[bool][]string }{"foo": x5}
	x7 := struct {
		Eum map[string]struct{ Dolorem map[int][]map[bool][]string }
	}{Eum: x6}
	x8 := struct {
		Tempora struct {
			Eum map[string]struct{ Dolorem map[int][]map[bool][]string }
		}
	}{Tempora: x7}
	x9 := []struct {
		Tempora struct {
			Eum map[string]struct{ Dolorem map[int][]map[bool][]string }
		}
	}{x8}
	x10 := [][]struct {
		Tempora struct {
			Eum map[string]struct{ Dolorem map[int][]map[bool][]string }
		}
	}{x9}
	var x typ = x10

	expected, err := json.Marshal(x0)
	require.NoError(t, err)
	kvs := encode(ctx, x)
	require.NotNil(t, kvs)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, KeyString("/0/0/Tempora/Eum/foo/Dolorem/13/0/true/0"), kvs[0].Key)
	require.Equal(t, expected, kvs[0].Value)
	fmt.Fprint(t.Output(), kvs)
}

func testEncode(t *testing.T, rt reflect.Type, expected any, expectedKey KeyString, expectedVal any) {
	ctx, cli := initAndTest(t)
	var err error
	// p := &y

	// test encode() - just verify KVs
	expectedBuf := common.MarshalAndTest(t, expectedVal)
	kvs := encode(ctx, expected)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, expectedKey, kvs[0].Key)
	require.Equal(t, expectedBuf, kvs[0].Value)

	// delete existing cluster entries for test key
	rootKey := KeyString("/test/testEncode")
	key := rootKey.Add(KeyString(expectedKey))
	ctx.Commentf("Deleting all subkeys of %s ...", key)
	opDelete := etcd.OpDelete(string(key), etcd.WithPrefix())
	txn := cli.Txn(ctx)
	_, err = txn.Then(opDelete).Commit()
	require.NoError(t, err)

	// test Encode - actually to cluster
	err = Encode(ctx, cli, expected, string(key))
	require.NoError(t, err)
// 
	// Create a pointer to the reflected type
	rtNew := reflect.New(rt)
	p := rtNew.Interface()

	// decode & see if we get the object we started with
	err = Decode(ctx, cli, string(key), p)
	require.NoError(t, err)
	require.Equal(t, expected, rtNew.Elem().Interface())
	

}

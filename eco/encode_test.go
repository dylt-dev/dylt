package eco

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
)


func TestEncode1 (t *testing.T) {
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


func TestEncode2 (t *testing.T) {
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


func TestEncode3 (t *testing.T) {
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


func TestEncode4 (t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ map[string]int
	x := typ {
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


func TestEncode5 (t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ []string
	x := typ { "foo", "bar", "bum" }

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


func TestEncode6 (t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ struct{Val int}
	x := typ { Val: 13 }

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


func TestEncode7 (t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ struct{Vals []int}
	x := typ { Vals: []int{13, 169, 1997}}

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


func TestEncode8 (t *testing.T) {
	ctx, cli := initAndTest(t)

	type typ struct{Vals []string}
	x := typ { Vals: []string{"foo", "bar", "bum"}}

	key := "/test/structs/struct3"
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


func marshal(t *testing.T, a any) []byte {
	var buf []byte
	var err error

	_, is := a.(string)
	if is {
		// buf = []byte(s)
		buf, err = json.Marshal(a)
		require.NoError(t, err)
	} else {
		buf, err = json.Marshal(a)
		require.NoError(t, err)
	}

	return buf
}

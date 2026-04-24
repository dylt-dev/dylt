package eco

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

func TestRemoveKeyFromSlice (t *testing.T) {
	ctx, _ := initAndTest(t)
	keyToRemove := "/test/slice/1"
	etcdKvs := []*mvccpb.KeyValue{
		{Key: []byte("/test/slice/0"), Value: []byte{}},
		{Key: []byte(keyToRemove), Value: []byte{}},
		{Key: []byte("/test/slice/2"), Value: []byte{}},
	}

	kvs := createKvSlice(etcdKvs)
	kvs = deleteKeyFromSlice(ctx, kvs, keyToRemove)
	require.Equal(t, 2, len(kvs))
}

func TestRemoveKeyFromSliceEmpty (t *testing.T) {
	ctx, _ := initAndTest(t)
	keyToRemove := "/test/slice/666"
	etcdKvs := []*mvccpb.KeyValue{
	}

	kvs := createKvSlice(etcdKvs)
	kvs = deleteKeyFromSlice(ctx, kvs, keyToRemove)
	require.Equal(t, 0, len(kvs))
}
func TestRemoveKeyFromSliceMissing (t *testing.T) {
	ctx, _ := initAndTest(t)
	keyToRemove := "/test/slice/666"
	etcdKvs := []*mvccpb.KeyValue{
		{Key: []byte("/test/slice/0"), Value: []byte{}},
		{Key: []byte("/test/slice/1"), Value: []byte{}},
		{Key: []byte("/test/slice/2"), Value: []byte{}},
	}

	kvs := createKvSlice(etcdKvs)
	kvs = deleteKeyFromSlice(ctx, kvs, keyToRemove)
	require.Equal(t, 3, len(kvs))
}
package eco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKvSeriesAdd(t *testing.T) {
	kvSeries := KvSeries{"/foo", []KeyValue{}}
	var is bool
	is = kvSeries.Add(KeyValue{"/foo/bar1", []byte{}})
	require.True(t, is)
	is = kvSeries.Add(KeyValue{"/foo/bar2", []byte{}})
	require.True(t, is)
	is = kvSeries.Add(KeyValue{"/foo/bum", []byte{}})
	require.True(t, is)
	is = kvSeries.Add(KeyValue{"XXX", []byte{}})
	require.False(t, is)
	require.Equal(t, 3, kvSeries.Len())
}

func TestKvSeriesIsOwner1(t *testing.T) {
	rootKey := KeyString("/foo")
	kvSeries := KvSeries{rootKey, []KeyValue{}}
	require.True(t, kvSeries.IsOwner("/foo/bar"))
}


func TestKvSeriesIsOwner2(t *testing.T) {
	rootKey := KeyString("/foo")
	kvSeries := KvSeries{rootKey, []KeyValue{}}
	require.True(t, kvSeries.IsOwner("/foo/bar/bum"))
}


func TestKvSeriesIsOwner3(t *testing.T) {
	rootKey := KeyString("/foo")
	kvSeries := KvSeries{rootKey, []KeyValue{}}
	require.False(t, kvSeries.IsOwner("/bar"))
}


func TestKvSeriesMaxIndex(t *testing.T) {
	expectedData := int(99)
	kvs := []KeyValue{
		{"/slice/0", nil},
		{"/slice/1", nil},
		{"/slice/99", nil},
		{"/slice/foo", nil},
		{"/slice/foo/bar", nil},
		{"/slice/foo/bar/13", nil},
		{"/fakeslice/169", nil},
		{"/1313", nil},
	}
	kvSlice := KvSeries{"/slice", kvs}
	n := kvSlice.MaxIndex()
	require.Equal(t, expectedData, n)
}

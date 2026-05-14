package eco

import (
	"testing"

	"github.com/stretchr/testify/require"
)


func TestKvSliceMaxIndex (t *testing.T) {
	expectedData := uint64(99)
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
	kvSlice := KvSlice{kvs, "/slice"}
	n := kvSlice.MaxIndex()
	require.Equal(t, expectedData, n)
}

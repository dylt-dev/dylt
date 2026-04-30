package eco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSliceKeyIndex(t *testing.T) {
	key := "/test/boolSlice/13"
	expectedVal := 13
	index, is := getSliceItemKey("/test/boolSlice", key)
	require.True(t, is)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyIndexEmpty(t *testing.T) {
	key := ""
	expectedVal := -1
	index, is := getSliceItemKey("/slice", key)
	require.False(t, is)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyIndexNoSlash(t *testing.T) {
	key := "/slice/barbum0"
	expectedVal := -1
	index, is := getSliceItemKey("/foo", key)
	require.False(t, is)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyIndexNoTrailingInt(t *testing.T) {
	key := "/foo/bar/bum"
	expectedVal := -1
	index, is := getSliceItemKey("/foo/bar", key)
	require.False(t, is)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyTrailingSlash(t *testing.T) {
	key := "/test/boolSlice/13/"
	expectedVal := 13
	index, is := getSliceItemKey("/test/boolSlice", key)
	require.True(t, is)
	require.Equal(t, expectedVal, index)
}

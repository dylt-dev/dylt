package eco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSliceKeyIndex(t *testing.T) {
	key := "/test/boolSlice/13"	
	expectedVal := 13
	index, is := getSliceKeyIndex(key)
	require.True(t, is)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyIndexEmpty(t *testing.T) {
	key := ""	
	expectedVal := -1
	index, is := getSliceKeyIndex(key)
	require.False(t, is)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyIndexNoSlash(t *testing.T) {
	key := "foobarbum0"	
	expectedVal := -1
	index, is := getSliceKeyIndex(key)
	require.False(t, is)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyIndexNoTrailingInt(t *testing.T) {
	key := "/foo/bar/bum"	
	expectedVal := -1
	index, is := getSliceKeyIndex(key)
	require.False(t, is)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyTrailingSlash(t *testing.T) {
	key := "/test/boolSlice/13/"	
	expectedVal := 13
	index, is := getSliceKeyIndex(key)
	require.True(t, is)
	require.Equal(t, expectedVal, index)
}

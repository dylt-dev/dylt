package eco

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSliceKeyIndex(t *testing.T) {
	key := "/test/boolSlice/13"	
	expectedVal := 13
	index, err := getSliceKeyIndex(key)
	require.NoError(t, err)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyIndexEmpty(t *testing.T) {
	key := ""	
	expectedVal := 0
	index, err := getSliceKeyIndex(key)
	require.Error(t, err)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyIndexNoSlash(t *testing.T) {
	key := "foobarbum0"	
	expectedVal := 0
	index, err := getSliceKeyIndex(key)
	require.Error(t, err)
	require.Equal(t, expectedVal, index)
}

func TestGetSliceKeyIndexNoTrailingInt(t *testing.T) {
	key := "/foo/bar/bum"	
	expectedVal := 0
	index, err := getSliceKeyIndex(key)
	require.Error(t, err)
	require.Equal(t, expectedVal, index)
}

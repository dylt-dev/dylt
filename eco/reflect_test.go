package eco

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// RvPointer.Walk() - non-nil pointer
func TestRvPointerWalk1 (t *testing.T) {
	var n int = 13
	var pn *int = &n

	rv := reflect.ValueOf(pn)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, *ptr, n)
}

// RvPointer.Walk() - pointer to nil pointer
func TestRvPointerWalk2 (t *testing.T) {
	var pn *int = nil
	var ppn **int = &pn

	rv := reflect.ValueOf(ppn)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	pptr, is := a.(**int)
	require.NotNil(t, pptr)
	require.True(t, is)
	require.Nil(t, *pptr)
}

// RvPointer.Walk() - pointer chain to non-nil pointer
func TestRvPointerWalk3 (t *testing.T) {
	var n int = 13
	var pn *int = &n
	var ppn **int = &pn
	var pppn ***int = &ppn
	var ppppn ****int = &pppn
	var pppppn *****int = &ppppn

	rv := reflect.ValueOf(pppppn)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, *ptr, n)
}

// RvPointer.Walk() - pointer chain to nil pointer
func TestRvPointerWalk4 (t *testing.T) {
	var pn *int = nil
	var ppn **int = &pn
	var pppn ***int = &ppn
	var ppppn ****int = &pppn
	var pppppn *****int = &ppppn

	rv := reflect.ValueOf(pppppn)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	pptr, is := a.(**int)
	require.NotNil(t, pptr)
	require.True(t, is)
	require.Nil(t, *pptr)
	require.Equal(t, pptr, ppn)
}

// pointer to nil slice
func TestRvPointerWalk5 (t *testing.T) {
	var slice []int = nil
	var pslice = &slice

	rv := reflect.ValueOf(pslice)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*[]int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Nil(t, *ptr)
	
}

// pointer to non-nil slice
func TestRvPointerWalk6 (t *testing.T) {
	var slice []int = []int{1, 2, 3, 4, 5}
	var pslice = &slice

	rv := reflect.ValueOf(pslice)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*[]int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, 3, (*ptr)[2])
	
}
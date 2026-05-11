package eco

import (
	"reflect"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)


func TestRvPointerElemTypeInt1 (t *testing.T) {
	expectedData := reflect.TypeFor[int]()
	var pn *int = nil
	var ppn **int = &pn
	rv := reflect.ValueOf(ppn)

	typ, err := RvPointer(rv).ElemType()
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}


func TestRvPointerElemTypeInt2 (t *testing.T) {
	expectedData := reflect.TypeFor[int]()
	var n int = 13
	var pn *int = &n
	rv := reflect.ValueOf(pn)

	typ, err := RvPointer(rv).ElemType()
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}


// RvPointer.ElemType() - pointer chain to nil pointer
func TestRvPointerElemTypeInt3(t *testing.T) {
	expectedData := reflect.TypeFor[int]()
	var pn *int = nil
	var ppn **int = &pn
	var pppn ***int = &ppn
	var ppppn ****int = &pppn
	var pppppn *****int = &ppppn
	rv := reflect.ValueOf(pppppn)

	typ, err := RvPointer(rv).ElemType()
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// RvPointer.ElemType() - pointer chain to non-nil pointer
func TestRvPointerElemTypeInt4(t *testing.T) {
	expectedData := reflect.TypeFor[int]()
	var n int = 13
	var pn *int = &n
	var ppn **int = &pn
	var pppn ***int = &ppn
	var ppppn ****int = &pppn
	var pppppn *****int = &ppppn
	rv := reflect.ValueOf(pppppn)

	typ, err := RvPointer(rv).ElemType()
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// pointer to nil map
func TestRvPointerElemTypeMap1(t *testing.T) {
	expectedData := reflect.TypeFor[map[int]string]()
	var m map[int]string = nil
	var pm = &m
	rv := reflect.ValueOf(pm)

	typ, err := RvPointer(rv).ElemType()
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// pointer to non-nil map
func TestRvPointerElemTypeMap2(t *testing.T) {
	expectedData := reflect.TypeFor[map[int]string]()
	var m map[int]string = map[int]string{13: "meat"}
	var pm = &m
	rv := reflect.ValueOf(pm)

	typ, err := RvPointer(rv).ElemType()
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// pointer chain to nil map
func TestRvPointerElemTypeMap3(t *testing.T) {
	expectedData := reflect.TypeFor[map[int]string]()
	var m map[int]string = nil
	var pm = &m
	var ppm = &pm
	var pppm = &ppm
	var ppppm = &pppm
	rv := reflect.ValueOf(ppppm)

	typ, err := RvPointer(rv).ElemType()
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// pointer chain to non-nil map
func TestRvPointerElemTypeMap4(t *testing.T) {
	expectedData := reflect.TypeFor[map[int]string]()
	var m map[int]string = map[int]string{13: "meat"}
	var pm = &m
	var ppm = &pm
	var pppm = &ppm
	var ppppm = &pppm
	rv := reflect.ValueOf(ppppm)

	typ, err := RvPointer(rv).ElemType()
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// RvPointer.Walk() - pointer to nil pointer
func TestRvPointerWalkInt1(t *testing.T) {
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

// RvPointer.Walk() - non-nil pointer
func TestRvPointerWalkInt2(t *testing.T) {
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

// RvPointer.Walk() - pointer chain to nil pointer
func TestRvPointerWalkInt3(t *testing.T) {
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

// RvPointer.Walk() - pointer chain to non-nil pointer
func TestRvPointerWalkInt4(t *testing.T) {
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

// pointer to nil map
func TestRvPointerWalkMap1(t *testing.T) {
	var m map[int]string = nil
	var pm = &m

	rv := reflect.ValueOf(pm)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Nil(t, *ptr)
}

// pointer to non-nil map
func TestRvPointerWalkMap2(t *testing.T) {
	var m map[int]string = map[int]string{13: "meat"}
	var pm = &m

	rv := reflect.ValueOf(pm)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, "meat", (*ptr)[13])
}

// pointer chain to nil map
func TestRvPointerWalkMap3(t *testing.T) {
	var m map[int]string = nil
	var pm = &m
	var ppm = &pm
	var pppm = &ppm
	var ppppm = &pppm

	rv := reflect.ValueOf(ppppm)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Nil(t, *ptr)
}

// pointer chain to non-nil map
func TestRvPointerWalkMap4(t *testing.T) {
	var m map[int]string = map[int]string{13: "meat"}
	var pm = &m
	var ppm = &pm
	var pppm = &ppm
	var ppppm = &pppm

	rv := reflect.ValueOf(ppppm)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, "meat", (*ptr)[13])
}

// pointer to nil slice
func TestRvPointerWalkSlice1(t *testing.T) {
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
func TestRvPointerWalkSlice2(t *testing.T) {
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

// pointer chain to nil slice
func TestRvPointerWalkSlice3(t *testing.T) {
	var slice []int = nil
	var pslice = &slice
	var ppslice = &pslice
	var pppslice = &ppslice
	var ppppslice = &pppslice

	rv := reflect.ValueOf(ppppslice)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*[]int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Nil(t, *ptr)

}

// pointer chain to non-nil slice
func TestRvPointerWalkSlice4(t *testing.T) {
	var slice []int = []int{1, 2, 3, 4, 5}
	var pslice = &slice
	var ppslice = &pslice
	var pppslice = &ppslice
	var ppppslice = &pppslice

	rv := reflect.ValueOf(ppppslice)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*[]int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, 3, (*ptr)[2])
}

// RvPointer.Walk() - nil struct pointer
func TestRvPointerWalk1Struct1(t *testing.T) {
	var pst *common.TestStruct = nil
	var ppst **common.TestStruct = &pst

	rv := reflect.ValueOf(ppst)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	pptr, is := a.(**common.TestStruct)
	require.NotNil(t, pptr)
	require.True(t, is)
	require.Nil(t, *pptr)
}

// RvPointer.Walk() - non-nil pointer to struct
func TestRvPointerWalkStruct2(t *testing.T) {
	name := "meat"
	var st common.TestStruct = common.TestStruct{Name: name}
	var pst *common.TestStruct = &st

	rv := reflect.ValueOf(pst)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*common.TestStruct)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, *ptr, st)
	require.Equal(t, name, (*ptr).Name)
}

// RvPointer.Walk() - nil struct pointer chain
func TestRvPointerWalkStruct3(t *testing.T) {
	var pst *common.TestStruct = nil
	var ppst **common.TestStruct = &pst
	var pppst ***common.TestStruct = &ppst
	var ppppst ****common.TestStruct = &pppst

	rv := reflect.ValueOf(ppppst)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	pptr, is := a.(**common.TestStruct)
	require.NotNil(t, pptr)
	require.True(t, is)
	require.Nil(t, *pptr)
}


// RvPointer.Walk() - non-nil pointer chain to struct
func TestRvPointerWalkStruct4(t *testing.T) {
	name := "meat"
	var st common.TestStruct = common.TestStruct{Name: name}
	var pst *common.TestStruct = &st
	var ppst **common.TestStruct = &pst
	var pppst ***common.TestStruct = &ppst
	var ppppst ****common.TestStruct = &pppst

	rv := reflect.ValueOf(ppppst)
	a, err := RvPointer(rv).Walk()
	require.NoError(t, err)
	ptr, is := a.(*common.TestStruct)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, *ptr, st)
	require.Equal(t, name, (*ptr).Name)
}

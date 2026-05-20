package eco

import (
	"os"
	"reflect"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestIsNormPointerMap1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := map[string]int{"thirteen": 13}
	p := &n
	is := IsNormPointer(ctx, p)
	require.True(t, is)
}

func TestIsNormPointerMap2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[string]int = nil
	pp := &p
	is := IsNormPointer(ctx, pp)
	require.True(t, is)
}

func TestIsNormPointerMap3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var n map[string]int = nil
	p := &n
	is := IsNormPointer(ctx, p)
	require.True(t, is)
}

func TestIsNormPointerMap4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := map[string]int{"thirteen": 13}
	p := &n
	pp := &p
	is := IsNormPointer(ctx, pp)
	require.False(t, is)
}

func TestIsNormPointerMap5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[string]int = nil
	is := IsNormPointer(ctx, p)
	require.False(t, is)
}

func TestIsNormPointerMap6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var n map[string]int = nil
	p := &n
	is := IsNormPointer(ctx, p)
	require.True(t, is)
}

func TestIsNormPointerSlice1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := []int{13}
	p := &n
	is := IsNormPointer(ctx, p)
	require.True(t, is)
}

func TestIsNormPointerSlice2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *[]int = nil
	pp := &p
	is := IsNormPointer(ctx, pp)
	require.True(t, is)
}

func TestIsNormPointerSlice3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var n []int = nil
	p := &n
	is := IsNormPointer(ctx, p)
	require.True(t, is)
}

func TestIsNormPointerSlice4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := []int{13}
	p := &n
	pp := &p
	is := IsNormPointer(ctx, pp)
	require.False(t, is)
}

func TestIsNormPointerSlice5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *[]int = nil
	is := IsNormPointer(ctx, p)
	require.False(t, is)
}

func TestIsNormPointerSlice6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var n []int = nil
	p := &n
	is := IsNormPointer(ctx, p)
	require.True(t, is)
}

func TestIsNormPointerStruct1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := struct{}{}
	p := &n
	is := IsNormPointer(ctx, p)
	require.True(t, is)
}

func TestIsNormPointerStruct2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *struct{} = nil
	pp := &p
	is := IsNormPointer(ctx, pp)
	require.True(t, is)
}

func TestIsNormPointerStruct3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := struct{}{}
	p := &n
	pp := &p
	is := IsNormPointer(ctx, pp)
	require.False(t, is)
}

func TestIsNormPointerStruct4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *struct{} = nil
	is := IsNormPointer(ctx, p)
	require.False(t, is)
}

func TestNormPtrIsAllocatedInt1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := 13
	var pn *int = &n
	normPtr, err := NewNormPtr(ctx, pn)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.True(t, is)
}

func TestNormPtrIsAllocatedInt2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *int = nil
	var pp **int = &p
	normPtr, err := NewNormPtr(ctx, pp)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.False(t, is)
}

func TestNormPtrIsAllocatedInt3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *int = nil
	normPtr, err := NewNormPtr(ctx, p)
	require.Error(t, err)
	require.Nil(t, normPtr)
}

func TestNormPtrIsAllocatedMap1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	m := map[string]int{"thirteen": 13}
	p := &m
	normPtr, err := NewNormPtr(ctx, p)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.True(t, is)
}

func TestNormPtrIsAllocatedMap2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[string]int = nil
	pp := &p
	normPtr, err := NewNormPtr(ctx, pp)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.False(t, is)
}

func TestNormPtrIsAllocatedMap3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var m map[string]int = nil
	p := &m
	normPtr, err := NewNormPtr(ctx, p)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.False(t, is)
}

func TestNormPtrIsAllocatedSlice1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	sl := []int{13}
	var p *[]int = &sl
	normPtr, err := NewNormPtr(ctx, p)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.True(t, is)
}

func TestNormPtrIsAllocatedSlice2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *[]int = nil
	var pp **[]int = &p
	normPtr, err := NewNormPtr(ctx, pp)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.False(t, is)
}

func TestNormPtrIsAllocatedSlice3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var sl []int = nil
	var p *[]int = &sl
	normPtr, err := NewNormPtr(ctx, p)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.False(t, is)
}

func TestNormPtrIsAllocatedSlice4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x []int = nil
	p := &x
	normPtr, err := NewNormPtr(ctx, p)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.False(t, is)
}
func TestNormPtrIsAllocatedStruct1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	st := struct{}{}
	p := &st
	normPtr, err := NewNormPtr(ctx, p)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.True(t, is)
}

func TestNormPtrIsAllocatedStruct2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *struct{} = nil
	pp := &p
	normPtr, err := NewNormPtr(ctx, pp)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsAllocated(ctx)
	require.False(t, is)
}

func TestNormPtrIsBigEnough1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	x := []int{1, 2, 3, 4, 5}
	p := &x
	normPtr, err := NewNormPtr(ctx, p)
	require.NoError(t, err)
	is := normPtr.IsBigEnough(ctx, 3)
	require.True(t, is)
}

func TestNormPtrIsBigEnoughSlice1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	slice := []int{1, 2, 3, 4, 5}
	pslice := &slice
	normPtr, err := NewNormPtr(ctx, pslice)
	require.NoError(t, err)
	is := normPtr.IsBigEnough(ctx, 13)
	require.False(t, is)
}

func TestNormPtrIsBigEnoughSlice2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var slice []int = nil
	pslice := &slice
	normPtr, err := NewNormPtr(ctx, pslice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
}

func TestNormPtrIsBigEnoughSlice3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x []int = nil
	p := &x
	normPtr, err := NewNormPtr(ctx, p)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsBigEnough(ctx, 1)
	require.False(t, is)
}

func TestNormPtrIsBigEnoughSlice4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *[]int = nil
	pp := &p
	normPtr, err := NewNormPtr(ctx, pp)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	is := normPtr.IsBigEnough(ctx, 1)
	require.False(t, is)
}

func TestNormPtrIsBigEnough4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var pslice *[]int = nil
	normPtr, err := NewNormPtr(ctx, pslice)
	require.Error(t, err)
	require.Nil(t, normPtr)
}

func TestNormPtrIsBigEnough6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	st := struct{}{}
	p := &st
	pp := &p
	normPtr, err := NewNormPtr(ctx, pp)
	require.Error(t, err)
	require.Nil(t, normPtr)
}

func TestNormPtrIsBigEnough7(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[string]int = nil
	normPtr, err := NewNormPtr(ctx, &p)
	require.NoError(t, err)
	is := normPtr.IsBigEnough(ctx, 13)
	require.True(t, is)
}
func TestNormPtrIsBigEnough5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *int = nil
	normPtr, err := NewNormPtr(ctx, &p)
	require.NoError(t, err)
	is := normPtr.IsBigEnough(ctx, 13)
	require.True(t, is)
}

func TestNormPtrSet1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := 13
	var x int
	p := &x

	rvp, err := NewNormPtr(ctx, p)
	require.NoError(t, err)
	rvp.Set(expected)
	require.Equal(t, expected, x)
}

func TestNormPtrSet2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expected := 13
	var p *int = nil
	pp := &p

	normPtr, err := NewNormPtr(ctx, pp)
	require.NoError(t, err)
	rvPtr := reflect.New(reflect.TypeFor[int]())
	rvPtr.Elem().Set(reflect.ValueOf(expected))
	require.Equal(t, expected, rvPtr.Elem().Interface())
	normPtr.Set(rvPtr)
	require.Equal(t, expected, *p)
}

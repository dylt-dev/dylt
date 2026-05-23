package common

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapTypeEmitFieldRef1(t *testing.T) {
	type m map[bool]int
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	mapType.EmitFieldRef()
}

func TestMapTypeEmitFieldRef2(t *testing.T) {
	type m map[int]string
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	mapType.EmitFieldRef()
}

func TestMapTypeEmitFieldRef3(t *testing.T) {
	type m map[string]bool
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	mapType.EmitFieldRef()
}

func TestMapTypeEmitFieldRef4(t *testing.T) {
	type m map[string][]int
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	mapType.EmitFieldRef()
}
func TestSliceTypeEmitFieldRef1(t *testing.T) {
	type sl []bool
	mapType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	mapType.EmitFieldRef()
}

func TestSliceTypeEmitFieldRef2(t *testing.T) {
	type sl []int
	mapType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	mapType.EmitFieldRef()
}

func TestSliceTypeEmitFieldRef3(t *testing.T) {
	type sl []string
	mapType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	mapType.EmitFieldRef()
}

func TestStructTypeEmitFieldRef1(t *testing.T) {
	type st struct{BoolValue bool}
	mapType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	mapType.EmitFieldRef()
}

func TestStructTypeEmitFieldRef2(t *testing.T) {
	type st struct{IntValue int}
	mapType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	mapType.EmitFieldRef()
}

func TestStructTypeEmitFieldRef3(t *testing.T) {
	type st struct{StringValue string}
	mapType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	mapType.EmitFieldRef()
}

func TestMapTypeEmitTree1(t *testing.T) {
	type m map[string]int
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	n := 0
	mapType.EmitTree(&n)
}

func TestMapTypeEmitTree2(t *testing.T) {
	type m map[string]string
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	n := 0
	mapType.EmitTree(&n)
}

func TestMapTypeEmitTree3(t *testing.T) {
	type m map[string]bool
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	n := 0
	mapType.EmitTree(&n)
}

func TestMapTypeEmitTree4(t *testing.T) {
	type m map[string][]int
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	n := 0
	mapType.EmitTree(&n)
}
func TestMapTypeIsScalar1(t *testing.T) {
	type m map[string]int
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	require.True(t, mapType.IsScalar())
}

func TestMapTypeIsScalar2(t *testing.T) {
	type m map[string]struct{}
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	require.False(t, mapType.IsScalar())
}

func TestMapTypeKeyName1(t *testing.T) {
	type m map[string]int
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	require.Equal(t, "", mapType.KeyName())
}

func TestMapTypeKeyName2(t *testing.T) {
	type m map[bool]int
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	require.Equal(t, "false", mapType.KeyName())
}

func TestMapTypeKeyName3(t *testing.T) {
	type m map[int]int
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	require.Equal(t, "0", mapType.KeyName())
}

func TestMapTypeZeroValue1(t *testing.T) {
	type m map[string]string
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	require.Equal(t, "", mapType.ZeroValue())
}

func TestMapTypeZeroValue2(t *testing.T) {
	type m map[string]bool
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	require.Equal(t, "false", mapType.ZeroValue())
}

func TestMapTypeZeroValue3(t *testing.T) {
	type m map[string]int
	mapType := NewMapType(NewDeepType(reflect.TypeFor[m]()))
	require.Equal(t, "0", mapType.ZeroValue())
}

func TestSliceTypeIsScalar1(t *testing.T) {
	type sl []int
	sliceType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	require.True(t, sliceType.IsScalar())
}

func TestSliceTypeIsScalar2(t *testing.T) {
	type sl []struct{}
	sliceType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	require.False(t, sliceType.IsScalar())
}

func TestSliceTypeKeyName1(t *testing.T) {
	type sl []int
	sliceType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	require.Equal(t, "0", sliceType.KeyName())
}

func TestSliceTypeKeyName2(t *testing.T) {
	type sl []string
	sliceType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	require.Equal(t, "0", sliceType.KeyName())
}

func TestSliceTypeKeyName3(t *testing.T) {
	type sl []bool
	sliceType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	require.Equal(t, "0", sliceType.KeyName())
}

func TestSliceTypeKeyName4(t *testing.T) {
	type sl []struct{}
	sliceType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	require.Equal(t, "0", sliceType.KeyName())
}

func TestSliceTypeZeroValue1(t *testing.T) {
	type sl []string
	sliceType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	require.Equal(t, "", sliceType.ZeroValue())
}

func TestSliceTypeZeroValue2(t *testing.T) {
	type sl []bool
	sliceType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	require.Equal(t, "false", sliceType.ZeroValue())
}

func TestSliceTypeZeroValue3(t *testing.T) {
	type sl []int
	sliceType := NewSliceType(NewDeepType(reflect.TypeFor[sl]()))
	require.Equal(t, "0", sliceType.ZeroValue())
}

func TestStructTypeIsScalar1(t *testing.T) {
	type st struct{ Value int }
	structType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	require.True(t, structType.IsScalar())
}

func TestStructTypeIsScalar2(t *testing.T) {
	type st struct{ Value []int }
	structType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	require.False(t, structType.IsScalar())
}

func TestStructTypeKeyName1(t *testing.T) {
	type st struct{ Value string }
	structType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	require.Equal(t, "Value", structType.KeyName())
}

func TestStructTypeKeyName2(t *testing.T) {
	type st struct{ Value int }
	structType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	require.Equal(t, "Value", structType.KeyName())
}

func TestStructTypeKeyName3(t *testing.T) {
	type st struct{ Value struct{} }
	structType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	require.Equal(t, "Value", structType.KeyName())
}

func TestStructTypeZeroValue1(t *testing.T) {
	type st struct{ Value bool }
	structType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	require.Equal(t, "false", structType.ZeroValue())
}

func TestStructTypeZeroValue2(t *testing.T) {
	type st struct{ Value int }
	structType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	require.Equal(t, "0", structType.ZeroValue())
}

func TestStructTypeZeroValue3(t *testing.T) {
	type st struct{ Value string }
	structType := NewStructType(NewDeepType(reflect.TypeFor[st]()))
	require.Equal(t, "", structType.ZeroValue())
}

func TestEmitFieldRef(t *testing.T) {
	type deepType struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}

	typ := reflect.TypeFor[deepType]()
	fmt.Print("x")
	DeepType{typ}.EmitFieldRef()
	fmt.Println()
}
func TestEmitTree(t *testing.T) {
	type deepType struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}

	typ := reflect.TypeFor[deepType]()
	level := 0
	DeepType{typ}.EmitTree(&level)
}

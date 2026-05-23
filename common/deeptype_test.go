package common

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapTypeEmitFieldRef1(t *testing.T) {
	type m map[bool]int
	// mapType := NewDeepType(reflect.TypeFor[m]())
	mapType := NewDeepType(reflect.TypeFor[m]())
	mapType.EmitFieldRef()
}

func TestMapTypeEmitFieldRef2(t *testing.T) {
	type m map[int]string
	mapType := NewDeepType(reflect.TypeFor[m]())
	mapType.EmitFieldRef()
}

func TestMapTypeEmitFieldRef3(t *testing.T) {
	type m map[string]bool
	mapType := NewDeepType(reflect.TypeFor[m]())
	mapType.EmitFieldRef()
}

func TestMapTypeEmitFieldRef4(t *testing.T) {
	type m map[string][]int
	mapType := NewDeepType(reflect.TypeFor[m]())
	mapType.EmitFieldRef()
}
func TestSliceTypeEmitFieldRef1(t *testing.T) {
	type sl []bool
	mapType := NewDeepType(reflect.TypeFor[sl]())
	mapType.EmitFieldRef()
}

func TestSliceTypeEmitFieldRef2(t *testing.T) {
	type sl []int
	mapType := NewDeepType(reflect.TypeFor[sl]())
	mapType.EmitFieldRef()
}

func TestSliceTypeEmitFieldRef3(t *testing.T) {
	type sl []string
	mapType := NewDeepType(reflect.TypeFor[sl]())
	mapType.EmitFieldRef()
}

func TestStructTypeEmitFieldRef1(t *testing.T) {
	type st struct{ BoolValue bool }
	mapType := NewDeepType(reflect.TypeFor[st]())
	mapType.EmitFieldRef()
}

func TestStructTypeEmitFieldRef2(t *testing.T) {
	type st struct{ IntValue int }
	mapType := NewDeepType(reflect.TypeFor[st]())
	mapType.EmitFieldRef()
}

func TestStructTypeEmitFieldRef3(t *testing.T) {
	type st struct{ StringValue string }
	mapType := NewDeepType(reflect.TypeFor[st]())
	mapType.EmitFieldRef()
}

func TestMapTypeEmitTree1(t *testing.T) {
	type m map[string]int
	mapType := NewDeepType(reflect.TypeFor[m]())
	n := 0
	mapType.EmitTree(&n)
}

func TestMapTypeEmitTree2(t *testing.T) {
	type m map[string]string
	mapType := NewDeepType(reflect.TypeFor[m]())
	n := 0
	mapType.EmitTree(&n)
}

func TestMapTypeEmitTree3(t *testing.T) {
	type m map[string]bool
	mapType := NewDeepType(reflect.TypeFor[m]())
	n := 0
	mapType.EmitTree(&n)
}

func TestMapTypeEmitTree4(t *testing.T) {
	type m map[string][]int
	mapType := NewDeepType(reflect.TypeFor[m]())
	n := 0
	mapType.EmitTree(&n)
}
func TestMapTypeIsScalar1(t *testing.T) {
	type m map[string]int
	mapType := NewDeepType(reflect.TypeFor[m]())
	require.True(t, mapType.isScalar())
}

func TestMapTypeIsScalar2(t *testing.T) {
	type m map[string]struct{}
	mapType := NewDeepType(reflect.TypeFor[m]())
	require.False(t, mapType.isScalar())
}

func TestMapTypeKeyName1(t *testing.T) {
	type m map[string]int
	mapType := NewDeepType(reflect.TypeFor[m]())
	require.Equal(t, "", mapType.subType().keyName())
}

func TestMapTypeKeyName2(t *testing.T) {
	type m map[bool]int
	mapType := NewDeepType(reflect.TypeFor[m]())
	require.Equal(t, "false", mapType.subType().keyName())
}

func TestMapTypeKeyName3(t *testing.T) {
	type m map[int]int
	mapType := NewDeepType(reflect.TypeFor[m]())
	require.Equal(t, "0", mapType.subType().keyName())
}

func TestMapTypeZeroValue1(t *testing.T) {
	type m map[string]string
	mapType := NewDeepType(reflect.TypeFor[m]())
	require.Equal(t, "", mapType.subType().zeroValue())
}

func TestMapTypeZeroValue2(t *testing.T) {
	type m map[string]bool
	mapType := NewDeepType(reflect.TypeFor[m]())
	require.Equal(t, "false", mapType.subType().zeroValue())
}

func TestMapTypeZeroValue3(t *testing.T) {
	type m map[string]int
	mapType := NewDeepType(reflect.TypeFor[m]())
	require.Equal(t, "0", mapType.subType().zeroValue())
}

func TestSliceTypeIsScalar1(t *testing.T) {
	type sl []int
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.True(t, sliceType.isScalar())
}

func TestSliceTypeIsScalar2(t *testing.T) {
	type sl []struct{}
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.False(t, sliceType.isScalar())
}

func TestSliceTypeKeyName1(t *testing.T) {
	type sl []int
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.subType().keyName())
}

func TestSliceTypeKeyName2(t *testing.T) {
	type sl []string
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.subType().keyName())
}

func TestSliceTypeKeyName3(t *testing.T) {
	type sl []bool
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.subType().keyName())
}

func TestSliceTypeKeyName4(t *testing.T) {
	type sl []struct{}
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.subType().keyName())
}

func TestSliceTypeZeroValue1(t *testing.T) {
	type sl []string
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "", sliceType.subType().zeroValue())
}

func TestSliceTypeZeroValue2(t *testing.T) {
	type sl []bool
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "false", sliceType.subType().zeroValue())
}

func TestSliceTypeZeroValue3(t *testing.T) {
	type sl []int
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.subType().zeroValue())
}

func TestStructTypeIsScalar1(t *testing.T) {
	type st struct{ Value int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.True(t, structType.isScalar())
}

func TestStructTypeIsScalar2(t *testing.T) {
	type st struct{ Value []int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.False(t, structType.isScalar())
}

func TestStructTypeKeyName1(t *testing.T) {
	type st struct{ Value string }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "Value", structType.subType().keyName())
}

func TestStructTypeKeyName2(t *testing.T) {
	type st struct{ Value int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "Value", structType.subType().keyName())
}

func TestStructTypeKeyName3(t *testing.T) {
	type st struct{ Value struct{} }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "Value", structType.subType().keyName())
}

func TestStructTypeZeroValue1(t *testing.T) {
	type st struct{ Value bool }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "false", structType.subType().zeroValue())
}

func TestStructTypeZeroValue2(t *testing.T) {
	type st struct{ Value int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "0", structType.subType().zeroValue())
}

func TestStructTypeZeroValue3(t *testing.T) {
	type st struct{ Value string }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "", structType.subType().zeroValue())
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

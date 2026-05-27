package common

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
func TestEmitValueRef10(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type typ struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	values := []any{}
	genScalarValues(ctx, reflect.TypeFor[typ](), r, &values)
	t.Log(values)

	fmt.Print("x")
	DeepType{reflect.TypeFor[typ]()}.EmitValueRef(values)
	fmt.Println()

}


/*
	Object creation     structTree1 := NewValueTree(ctx, "N", 13)
						sliceTree1  := NewValueTree(ctx, 3, structTree1)
						valTree     := NewValueTree(ctx, "Val", sliceTree1)
						mapTree     := NewValueTree(ctx, "foo", valTree)
						sliceTree2  := NewValueTree(ctx, 0, mapTree)
						structTree2 := NewValueTree(ctx, "Slice", sliceTree2)
						mapTree2    := NewValueTree(ctx, "bar", structTree2)
						sliceTree3  := NewValueTree(ctx, 2, mapTree2)
						structTree3 := NewValueTree(ctx, "Data", sliceTree3)

	Field access        x.Data[2]["bar"].Slice[0]["foo"].Val[3].N
						x.Data[0][""].Slice[0][""].Val[0].N
*/
func TestGenTest10(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type typ struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	values := []any{}
	genScalarValues(ctx, reflect.TypeFor[typ](), r, &values)
	t.Log(values)

	// emit tree
	n := 0
	DeepType{reflect.TypeFor[typ]()}.EmitTreeDecl(&n, values)
	
	// emit value ref
	fmt.Print("x")
	DeepType{reflect.TypeFor[typ]()}.EmitValueRef(values)
	fmt.Println()

}


func TestGenTest100(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())
	
	type typ struct{Nemo struct{Eligendi struct{Illum map[bool]map[bool]map[bool]struct{Est struct{Tempore struct{Magnam struct{Fugiat struct{Sed struct{Et [][][]map[int]map[string]struct{Quis []map[bool]map[bool][]struct{Ullam []map[int]struct{Assumenda struct{Rerum struct{Possimus [][][]map[int]struct{Harum map[bool]map[int]struct{Aperiam struct{Incidunt struct{Beatae []struct{Nam struct{Sunt []map[string]struct{Dolore [][]map[bool][]struct{Assumenda []struct{Consequatur struct{Iste []map[bool]struct{Et [][]map[bool][][][]struct{Est map[bool]struct{Ut [][]map[string]struct{Facere struct{Velit struct{Ut map[int]map[string]struct{Quidem struct{Ipsa map[int]struct{Molestiae map[bool][]map[int]struct{Repellat [][][]struct{Reprehenderit struct{Eaque struct{Beatae map[string]struct{Unde struct{Perspiciatis [][]struct{Possimus map[bool]map[bool]struct{Eligendi struct{Modi struct{Vel []struct{Possimus []int}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}	
	values := []any{}
	genScalarValues(ctx, reflect.TypeFor[typ](), r, &values)
	t.Log(values)

	// emit tree
	n := 0
	DeepType{reflect.TypeFor[typ]()}.EmitTreeDecl(&n, values)
	
	// emit value ref
	fmt.Print("x")
	DeepType{reflect.TypeFor[typ]()}.EmitValueRef(values)
	fmt.Println()

}




func TestMapTypeEmitValueRef1a(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type typ map[bool]int
	// mapType := NewDeepType(reflect.TypeFor[m]())
	r := rand.NewSource(time.Now().UTC().UnixNano())
	values := []any{}
	genScalarValues(ctx, reflect.TypeFor[typ](), r, &values)
	t.Log(values)

	mapType := NewDeepType(reflect.TypeFor[typ]())
	mapType.EmitValueRef(values)
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
	require.Equal(t, "0", sliceType.keyName())
}

func TestSliceTypeKeyName2(t *testing.T) {
	type sl []string
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.keyName())
}

func TestSliceTypeKeyName3(t *testing.T) {
	type sl []bool
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.keyName())
}

func TestSliceTypeKeyName4(t *testing.T) {
	type sl []struct{}
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.keyName())
}

func TestSliceTypeZeroValue1(t *testing.T) {
	type sl []string
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "", sliceType.zeroValue())
}

func TestSliceTypeZeroValue2(t *testing.T) {
	type sl []bool
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "false", sliceType.zeroValue())
}

func TestSliceTypeZeroValue3(t *testing.T) {
	type sl []int
	sliceType := NewDeepType(reflect.TypeFor[sl]())
	require.Equal(t, "0", sliceType.zeroValue())
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
	require.Equal(t, "Value", structType.keyName())
}

func TestStructTypeKeyName2(t *testing.T) {
	type st struct{ Value int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "Value", structType.keyName())
}

func TestStructTypeKeyName3(t *testing.T) {
	type st struct{ Value struct{} }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "Value", structType.keyName())
}

func TestStructTypeZeroValue1(t *testing.T) {
	type st struct{ Value bool }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "false", structType.zeroValue())
}

func TestStructTypeZeroValue2(t *testing.T) {
	type st struct{ Value int }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "0", structType.zeroValue())
}

func TestStructTypeZeroValue3(t *testing.T) {
	type st struct{ Value string }
	structType := NewDeepType(reflect.TypeFor[st]())
	require.Equal(t, "", structType.zeroValue())
}

func TestEmitTree(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type deepType struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}

	typ := reflect.TypeFor[deepType]()

	values := []any{}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	genScalarValues(ctx, reflect.TypeFor[deepType](), r, &values)
	t.Log(values)
	level := 0
	DeepType{typ}.EmitTreeDecl(&level, values)
}

func TestGenDeclaration1(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())
	
	for range 10 {
		genDeclaration(ctx, 1, r, t.Output())
		t.Output().Write([]byte("\n"))
	}
}

func TestGenDeclaration2(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())
	
	for range 10 {
		genDeclaration(ctx, 2, r, t.Output())
		t.Output().Write([]byte("\n"))
	}
}

func TestGenDeclaration100(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())

	genDeclaration(ctx, 100, r, t.Output())
	t.Output().Write([]byte("\n"))
}

func TestGenScalars1(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	types := []any{
		new(map[string]bool),
		new(struct{ Field bool }),
		new(struct{ Field string }),
		new([]int),
		new(struct{ Field string }),
		new([]string),
		new([]bool),
		new(struct{ Field int }),
		new(struct{ Field bool }),
		new(map[bool]string),
	}

	r := rand.NewSource(time.Now().UTC().UnixNano())
	for _, p := range types {
		values := []any{}
		typ := reflect.TypeOf(p).Elem()
		genScalarValues(ctx, typ, r, &values)
		t.Logf("%s => %v", typ, values)
	}
}

func TestGenScalars2(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	types := []any{
		new(map[int]map[string]int),
		new(struct{ Field struct{ Field int } }),
		new([]map[string]bool),
		new(map[string][]string),
		new(struct{ Field []string }),
		new(map[string][]bool),
		new(struct{ Field map[bool]string }),
		new(struct{ Field struct{ Field bool } }),
		new(struct{ Field []bool }),
		new(struct{ Field map[int]bool }),
	}

	r := rand.NewSource(time.Now().UTC().UnixNano())
	for _, p := range types {
		values := []any{}
		typ := reflect.TypeOf(p).Elem()
		genScalarValues(ctx, typ, r, &values)
		t.Logf("%s => %v", typ, values)
	}
}


func TestGenScalars10(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	// x.Data[0][""].Slice[0][""].Val[0].N
	type deepType struct {
		Data []map[string]struct {
			Slice []map[string]struct{ Val []struct{ N int } }
		}
	}

	values := []any{}
	r := rand.NewSource(time.Now().UTC().UnixNano())
	typ := reflect.TypeFor[deepType]()
	genScalarValues(ctx, typ, r, &values)
	t.Logf("%s => %v", typ, values)
}

func TestGenScalars100(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	r := rand.NewSource(time.Now().UTC().UnixNano())

	type deepType struct{Nemo struct{Eligendi struct{Illum map[bool]map[bool]map[bool]struct{Est struct{Tempore struct{Magnam struct{Fugiat struct{Sed struct{Et [][][]map[int]map[string]struct{Quis []map[bool]map[bool][]struct{Ullam []map[int]struct{Assumenda struct{Rerum struct{Possimus [][][]map[int]struct{Harum map[bool]map[int]struct{Aperiam struct{Incidunt struct{Beatae []struct{Nam struct{Sunt []map[string]struct{Dolore [][]map[bool][]struct{Assumenda []struct{Consequatur struct{Iste []map[bool]struct{Et [][]map[bool][][][]struct{Est map[bool]struct{Ut [][]map[string]struct{Facere struct{Velit struct{Ut map[int]map[string]struct{Quidem struct{Ipsa map[int]struct{Molestiae map[bool][]map[int]struct{Repellat [][][]struct{Reprehenderit struct{Eaque struct{Beatae map[string]struct{Unde struct{Perspiciatis [][]struct{Possimus map[bool]map[bool]struct{Eligendi struct{Modi struct{Vel []struct{Possimus []int}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}}	
	values := []any{}
	typ := reflect.TypeFor[deepType]()
	genScalarValues(ctx, typ, r, &values)
	t.Logf("%s => %v", typ, values)
}

func TestGetRandomFlavor(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	flavorCount := map[Flavor]int{
		Map:    0,
		Slice:  0,
		Struct: 0,
	}
	for range 10000 {
		flavorCount[getRandFlavor(ctx)]++
	}

	t.Log(flavorCount)
}

func TestGetRandomScalarKind(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	kindCount := map[reflect.Kind]int{
		reflect.Bool:   0,
		reflect.Int:    0,
		reflect.String: 0,
	}
	for range 10000 {
		kindCount[getRandScalar(ctx)]++
	}

	t.Log(kindCount)
}

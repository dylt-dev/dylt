package eco

import (
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestRveEncode1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	x := "foo"
	expectedKey := KeyString("/test/string")
	expectedVal := marshal(t, x)

	rve := NewRvEncodable(x)
	kvs := rve.Encode(ctx, "/test/string")
	require.NotNil(t, kvs)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, expectedKey, kvs[0].Key)
	require.Equal(t, expectedVal, kvs[0].Value)

}

func TestRveEncode2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	x := 13
	expectedKey := KeyString("/test/int")
	expectedVal := marshal(t, x)

	rve := NewRvEncodable(x)
	kvs := rve.Encode(ctx, expectedKey)
	require.NotNil(t, kvs)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, expectedKey, kvs[0].Key)
	require.Equal(t, expectedVal, kvs[0].Value)

}

func TestRveEncode3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	x := false
	expectedKey := KeyString("/test/bool")
	expectedVal := marshal(t, x)

	rve := NewRvEncodable(x)
	kvs := rve.Encode(ctx, expectedKey)
	require.NotNil(t, kvs)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, expectedKey, kvs[0].Key)
	require.Equal(t, expectedVal, kvs[0].Value)

}

func TestRveEncode4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	x := []string{"foo"}
	rootKey := KeyString("/test/slice")
	expectedKey := rootKey.AddSegment("0")
	expectedVal := marshal(t, x[0])

	rve := NewRvEncodable(x)
	kvs := rve.Encode(ctx, rootKey)
	require.NotNil(t, kvs)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, expectedKey, kvs[0].Key)
	require.Equal(t, expectedVal, kvs[0].Value)
}

func TestRveEncode5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	x := map[string]int{"foo": 13}
	rootKey := KeyString("/test/map")
	expectedKey := rootKey.AddSegment("foo")
	expectedVal := marshal(t, x["foo"])

	rve := NewRvEncodable(x)
	kvs := rve.Encode(ctx, rootKey)
	require.NotNil(t, kvs)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, expectedKey, kvs[0].Key)
	require.Equal(t, expectedVal, kvs[0].Value)
}

func TestRveEncode6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	x := struct{ Flag bool }{Flag: true}
	rootKey := KeyString("/test/struct")
	expectedKey := rootKey.AddSegment("Flag")
	expectedVal := marshal(t, x.Flag)

	rve := NewRvEncodable(x)
	kvs := rve.Encode(ctx, rootKey)
	require.NotNil(t, kvs)
	require.Equal(t, 1, len(kvs))
	require.Equal(t, expectedKey, kvs[0].Key)
	require.Equal(t, expectedVal, kvs[0].Value)
}

func TestRveEncode7(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)

	type typ0 []bool
	type typ1 struct{ Flags typ0 }
	type typ2 map[string]typ1
	type typ3 []typ2
	type typ4 struct{ Data typ3 }

	x0 := typ0{true, true, false}
	x1 := typ1{Flags: x0}
	x2 := typ2{"myflags": x1}
	x3 := typ3{x2}
	x4 := typ4{Data: x3}
	x := x4

	rootKey := KeyString("/test/struct")
	// expectedKey := rootKey.AddSegment("Flag")

	rve := NewRvEncodable(x)
	kvs := rve.Encode(ctx, rootKey)
	require.NotNil(t, kvs)
	require.Equal(t, 3, len(kvs))
	t.Log(kvs)
	// require.Equal(t, expectedKey, kvs[0].Key)
	// require.Equal(t, expectedVal, kvs[0].Value)
}

func TestRveNil1(t *testing.T) {
	var x *int = nil

	rve := NewRvEncodable(x)
	require.True(t, rve.IsNil())

}

func TestRveNil2(t *testing.T) {
	var x *int = nil
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp

	rve := NewRvEncodable(pppp)
	require.True(t, rve.IsNil())

}

func TestRveNil3(t *testing.T) {
	var x int = 13

	rve := NewRvEncodable(x)
	require.False(t, rve.IsNil())

}

func TestRveNil4(t *testing.T) {
	defer func() { r := recover(); require.NotNil(t, r) }()

	var x chan<- int = nil

	rve := NewRvEncodable(x)
	require.True(t, rve.IsNil())
}

func TestRveNil5(t *testing.T) {
	var x []int = nil

	rve := NewRvEncodable(x)
	require.True(t, rve.IsNil())
}

func TestRveNil6(t *testing.T) {
	var x *map[string]int = nil

	rve := NewRvEncodable(x)
	require.True(t, rve.IsNil())
}

func TestRveNil7(t *testing.T) {
	var x map[string]int = nil
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp

	rve := NewRvEncodable(pppp)
	require.True(t, rve.IsNil())
}

func TestRveNil8(t *testing.T) {
	x := []int{13}

	rve := NewRvEncodable(x)
	require.False(t, rve.IsNil())

}

func TestRveWalkPointer1(t *testing.T) {
	var x *[]int = nil
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp

	rve := NewRvEncodable(pppp)
	rvev, err := rve.walkPointer()
	require.NoError(t, err)
	require.Equal(t, x, rvev.Interface())
}

func TestRveWalkPointer2(t *testing.T) {
	var x int = 13
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp

	rve := NewRvEncodable(pppp)
	rvev, err := rve.walkPointer()
	require.NoError(t, err)
	require.Equal(t, x, rvev.Interface())
}

func TestRveWalkPointer3(t *testing.T) {
	var x int = 13

	rve := NewRvEncodable(x)
	rvev, err := rve.walkPointer()
	require.NoError(t, err)
	require.Equal(t, x, rvev.Interface())
}

func TestRveWalkPointer4(t *testing.T) {
	var x chan<- int

	rve := NewRvEncodable(x)
	_, err := rve.walkPointer()
	require.Error(t, err)
}

func TestRveWalkPointer5(t *testing.T) {
	var x *int = nil
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp

	rve := NewRvEncodable(pppp)
	rvev, err := rve.walkPointer()
	require.NoError(t, err)
	require.Equal(t, x, rvev.Interface())
}

func TestRvMapData1(t *testing.T) {
	x := map[string]int{
		"foo": 13,
		"bar": 169,
		"bum": 1997,
	}
	expected := map[string]any{}
	for k, v := range x {
		expected[k] = v
	}

	rvs := NewRvMap(x)
	data := rvs.Data()
	require.Equal(t, expected, data)
}

func TestRvScalarData1(t *testing.T) {
	x := "foo"
	expected := map[string]any{
		"": x,
	}

	rvs := NewRvScalar(x)
	data := rvs.Data()
	require.Equal(t, expected, data)
}

func TestRvScalarData2(t *testing.T) {
	x := 13
	expected := map[string]any{
		"": x,
	}

	rvs := NewRvScalar(x)
	data := rvs.Data()
	require.Equal(t, expected, data)
}

func TestRvScalarData3(t *testing.T) {
	x := false
	expected := map[string]any{
		"": x,
	}

	rvs := NewRvScalar(x)
	data := rvs.Data()
	require.Equal(t, expected, data)
}

func TestRvSliceData1(t *testing.T) {
	x := []string{"foo", "bar", "bum"}
	expected := map[string]any{
		"0": x[0],
		"1": x[1],
		"2": x[2],
	}

	rvs := NewRvSlice(x)
	data := rvs.Data()
	require.Equal(t, expected, data)
}

func TestRvSliceData2(t *testing.T) {
	x := []any{"foo", 13, false}
	expected := map[string]any{
		"0": x[0],
		"1": x[1],
		"2": x[2],
	}

	rvs := NewRvSlice(x)
	data := rvs.Data()
	require.Equal(t, expected, data)
}

func TestRvStructData1(t *testing.T) {
	type typ struct {
		Foo int
		Bar float32
		Bum string
	}
	x := typ{
		Foo: 13,
		Bar: 169.0,
		Bum: "1997",
	}
	expected := map[string]any{
		"Foo": x.Foo,
		"Bar": x.Bar,
		"Bum": x.Bum,
	}

	rvs, err := NewRvStruct(x)
	require.NoError(t, err)
	data := rvs.Data()
	require.Equal(t, expected, data)
}


package eco

import (
	"testing"

	"github.com/stretchr/testify/require"
)


func TestRvMapData1 (t *testing.T) {
	x := map[string]int {
		"foo": 13,
		"bar": 169,
		"bum": 1997,
	}
	expected := map[string]any {}
	for k, v := range x {
		expected[k] = v
	}

	rvs := NewRvMap(x)
	data := rvs.Data()
	require.Equal(t, expected, data)
}


func TestRvSliceData1 (t *testing.T) {
	x := []string{"foo", "bar", "bum"}
	expected := map[string]any {
		"0": x[0],
		"1": x[1],
		"2": x[2],
	}

	rvs := NewRvSlice(x)
	data := rvs.Data()
	require.Equal(t, expected, data)
}


func TestRvSliceData2 (t *testing.T) {
	x := []any{"foo", 13, false}
	expected := map[string]any {
		"0": x[0],
		"1": x[1],
		"2": x[2],
	}

	rvs := NewRvSlice(x)
	data := rvs.Data()
	require.Equal(t, expected, data)
}

func TestRvStructData1 (t *testing.T) {
	type typ struct { 
		Foo int
		Bar float32
		Bum string
	}
	x := typ {
		Foo: 13,
		Bar: 169.0,
		Bum: "1997",
	}
	expected := map[string]any {
		"Foo": x.Foo,
		"Bar": x.Bar,
		"Bum": x.Bum,
	}

	rvs, err := NewRvStruct(x)
	require.NoError(t, err)
	data := rvs.Data()
	require.Equal(t, expected, data)
}


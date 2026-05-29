package eco

import (
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)



func TestDecodeDeepType2_0(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := MainDecoder{}

	// type declaration
	type deepType struct{Aut struct{Numquam bool;};}

	// value tree construction
	tree0 := NewValueTree(ctx, "Numquam", true)
tree1 := NewValueTree(ctx, "Aut", tree0)

	tree := tree1

	// decode object
	var x deepType
	p := &x
	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)

	// check leaf value
	expectedVal := true
	val := x.Aut.Numquam
	require.Equal(t, expectedVal, val)
}

func TestDecodeDeepType2_1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := MainDecoder{}

	// type declaration
	type deepType []map[string]string

	// value tree construction
	tree0 := NewValueTree(ctx, "omnis", "consectetur")
tree1 := NewValueTree(ctx, 0, tree0)

	tree := tree1

	// decode object
	var x deepType
	p := &x
	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)

	// check leaf value
	expectedVal := "consectetur"
	val := x[0]["omnis"]
	require.Equal(t, expectedVal, val)
}

func TestDecodeDeepType2_2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := MainDecoder{}

	// type declaration
	type deepType struct{Dolorum struct{Officia int;};}

	// value tree construction
	tree0 := NewValueTree(ctx, "Officia", 113)
tree1 := NewValueTree(ctx, "Dolorum", tree0)

	tree := tree1

	// decode object
	var x deepType
	p := &x
	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)

	// check leaf value
	expectedVal := 113
	val := x.Dolorum.Officia
	require.Equal(t, expectedVal, val)
}

func TestDecodeDeepType2_3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := MainDecoder{}

	// type declaration
	type deepType map[string]struct{Tempora bool;}

	// value tree construction
	tree0 := NewValueTree(ctx, "Tempora", false)
tree1 := NewValueTree(ctx, "qui", tree0)

	tree := tree1

	// decode object
	var x deepType
	p := &x
	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)

	// check leaf value
	expectedVal := false
	val := x["qui"].Tempora
	require.Equal(t, expectedVal, val)
}

func TestDecodeDeepType2_4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := MainDecoder{}

	// type declaration
	type deepType struct{Excepturi map[string]bool;}

	// value tree construction
	tree0 := NewValueTree(ctx, "quam", true)
tree1 := NewValueTree(ctx, "Excepturi", tree0)

	tree := tree1

	// decode object
	var x deepType
	p := &x
	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)

	// check leaf value
	expectedVal := true
	val := x.Excepturi["quam"]
	require.Equal(t, expectedVal, val)
}

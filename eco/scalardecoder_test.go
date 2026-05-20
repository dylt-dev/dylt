package eco

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestScalarDecoderBool(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := ScalarDecoder[bool]{}
	tree := &ValueTree{}
	tree.Value = []byte(strconv.AppendBool([]byte{}, true))
	var x bool
	p := &x
	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, true, x)
}


func TestScalarDecoderBoolNil(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := ScalarDecoder[bool]{}
	tree := &ValueTree{}
	tree.Value = []byte(strconv.AppendBool([]byte{}, true))
	var p *bool
	pp := &p
	err := decoder.Decode(ctx, tree, pp)
	require.NoError(t, err)
	require.Equal(t, true, **pp)
}


func TestUnmarshalPp(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *bool = nil
	var pp = &p
	rvp, err := NewRvPointer(pp)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, 0)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	require.NotNil(t, normPtr.Value)
	require.NotNil(t, *pp)
	buf := strconv.AppendBool([]byte{}, true)
	err = json.Unmarshal(buf, p)
	require.NoError(t, err)
	require.Equal(t, true, *p)
}
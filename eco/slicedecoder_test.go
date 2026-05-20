package eco

import (
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestSliceDecoder1 (t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := SliceDecoder{}
	tree := &ValueTree{}
	tree.Add(ctx, KeyString("/0"), []byte("\"foo\""))
	tree.Add(ctx, KeyString("/1"), []byte("\"bar\""))
	tree.Add(ctx, KeyString("/9"), []byte("\"bum\""))
	require.Equal(t, 3, len(tree.ChildMap))

	var x []string = nil
	p := &x

	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, 10, len(x))
	require.Equal(t, "foo", x[0])
	require.Equal(t, "bar", x[1])
	require.Equal(t, "", x[2])
	require.Equal(t, "bum", x[9])
}


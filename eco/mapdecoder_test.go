package eco

import (
	"os"
	"strconv"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestMapDecoder1 (t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	decoder := MapDecoder{}
	tree := &ValueTree{}
	tree.Add(ctx, KeyString("/foo"), strconv.AppendInt([]byte{}, 13, 10))
	tree.Add(ctx, KeyString("/bar"), strconv.AppendInt([]byte{}, 169, 10))
	require.Equal(t, 2, len(tree.ChildMap))

	var x map[string]int = nil
	p := &x

	err := decoder.Decode(ctx, tree, p)
	require.NoError(t, err)
	require.Equal(t, 2, len(x))
	foo, is := x["foo"]
	require.True(t, is)
	require.Equal(t, 13, foo)
	bar, is := x["bar"]
	require.True(t, is)
	require.Equal(t, 169, bar)
}

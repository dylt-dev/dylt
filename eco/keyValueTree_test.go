package eco

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	ctx, _ := initAndTest(t)

	// test data
	kvs := []*KeyValue{
		newKv("/team/1/name", "bum"),
		newKv("/team/1/color", "blue"),
		newKv("/team/2/name", "buzz"),
		newKv("/team/2/color", "green"),
		newKv("/team/3/color", "pink"),
		newKv("/team/3/stats/age", "33"),
		newKv("/team/3/stats/rating", "100"),
	}

	key := "/team"
	tree := createKvTree(ctx, key, kvs, key)
	t.Log("\n" + fmt.Sprint(tree))
	require.NotNil(t, tree)
	require.Equal(t, "", tree.Name)
	require.Nil(t, nil, tree.Value)

	// log tree
	logTree(ctx, tree)
}

func TestScalar(t *testing.T) {
	ctx, _ := initAndTest(t)

	// test data
	kvs := []*KeyValue{
		newKv("/team/1/name", "bum"),
	}

	key := "/team/1/name"
	tree := createKvTree(ctx, key, kvs, key)
	t.Log("\n" + fmt.Sprint(tree))
	require.NotNil(t, tree)
	require.Equal(t, "", tree.Name)
	require.Equal(t, "bum", string(tree.Value))
	require.Equal(t, 0, len(tree.Children))

	// log tree
	logTree(ctx, tree)
}

func logTree (ctx *ecoContext, tree *KeyValueTree) {
	ctx.logger.signature("logTree", tree.Name, string(tree.Value), len(tree.Children))
	ctx.inc()
	defer ctx.dec()
	
	for _, child := range tree.Children {
		logTree(ctx, child)
	}
}

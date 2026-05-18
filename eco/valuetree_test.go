package eco

import (
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)


func TestValueTreeAdd1 (t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	tree := &ValueTree{}
	val := []byte("13")
	tree.Add(ctx, "/foo", val)
	require.Equal(t, val, tree.Value)
	require.Nil(t, tree.ChildMap)
}

func TestValueTreeAdd2 (t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	tree := &ValueTree{}
	val := []byte("13")
	tree.Add(ctx, "/foo/bar/bum", val)

	var is bool
	// /foo - no value, one child
	require.Nil(t, tree.Value)
	require.Equal(t, 1, len(tree.ChildMap))
	tree, is = tree.ChildMap["foo"]
	require.True(t, is)
	require.NotNil(t, tree)

	// /foo/bar - no value, one child
	require.Nil(t, tree.Value)
	require.Equal(t, 1, len(tree.ChildMap))
	t.Log(tree.ChildMap)
	tree, is = tree.ChildMap["bar"]
	require.True(t, is)
	require.NotNil(t, tree)

	// /foo/bar/bum - value, no child map
	require.NotNil(t, tree.Value)
	require.Equal(t, val, tree.Value)
	require.Nil(t, tree.ChildMap)
}


func TestValueTreeAddMulti1 (t *testing.T) {
/*
	key := "/test/stros"
	map0 := map[string]string{"Name": "Altuve", "Position": "2B"}
	map1 := map[string]string{"Name": "Pena", "Position": "SS"}
	map2 := map[string]string{"Name": "Javier", "Position": "P"}
	mapStros := map[int]map[string]string{27: map0, 3: map1, 53: map2}
*/	
	ctx := common.NewEcoContext(os.Stdout)
	tree := &ValueTree{}
	tree.Add(ctx, "/27/Name", []byte("Altuve"))
	tree.Add(ctx, "/27/Position", []byte("2B"))
	tree.Add(ctx, "/3/Name", []byte("Pena"))
	tree.Add(ctx, "/3/Position", []byte("SS"))
	tree.Add(ctx, "/53/Name", []byte("Javier"))
	tree.Add(ctx, "/53/Position", []byte("SP"))
	var is bool

	require.Nil(t, tree.Value)
	require.Equal(t, 3, len(tree.ChildMap))

	tree27, is := tree.ChildMap["27"]
	require.True(t, is)
	testPlayer(t, tree27, "Altuve", "2B")

	tree3, is := tree.ChildMap["3"]
	require.True(t, is)
	testPlayer(t, tree3, "Pena", "SS")

	tree53, is := tree.ChildMap["53"]
	require.True(t, is)
	testPlayer(t, tree53, "Javier", "SP")
}


func testPlayer (t *testing.T, playerTree *ValueTree, expectedName string, expectedPosition string) {
	require.NotNil(t, playerTree)
	require.Nil(t, playerTree.Value)
	require.Equal(t, 2, len(playerTree.ChildMap))
	
	treeName, is := playerTree.ChildMap["Name"]
	require.True(t, is)
	require.NotNil(t, treeName)
	require.NotNil(t, treeName.Value)
	require.Equal(t, []byte(expectedName), treeName.Value)
	require.Nil(t, treeName.ChildMap)

	treePosition, is := playerTree.ChildMap["Position"]
	require.True(t, is)
	require.NotNil(t, treePosition)
	require.NotNil(t, treePosition.Value)
	require.Equal(t, []byte(expectedPosition), treePosition.Value)
	require.Nil(t, treePosition.ChildMap)
	
}
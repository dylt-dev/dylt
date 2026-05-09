package eco

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)


func TestKvTreeMap (t *testing.T) {
	ctx, _ := initAndTest(t)

	key := "/test/team/astros/Players/altuve"
	expectedBorn := "Venezuela"
	expectedId := "1"
	expectedIsActive := "true"
	expectedName := "Jose Altuve"
	expectedWeight := "160"
	// expectedData := map[string]string{
	// 	"Born":     expectedBorn,
	// 	"Id":       expectedId,
	// 	"IsActive": expectedIsActive,
	// 	"Name":     expectedName,
	// 	"Weight":   expectedWeight,
	// }
	keyBorn := fmt.Sprintf("%s/Born", key)
	keyId := fmt.Sprintf("%s/Id", key)
	keyIsActive := fmt.Sprintf("%s/IsActive", key)
	keyName := fmt.Sprintf("%s/Name", key)
	keyWeight := fmt.Sprintf("%s/Weight", key)

	kvs := []*KeyValue{
		newKv(keyBorn, expectedBorn),
		newKv(keyId, expectedId),
		newKv(keyIsActive, expectedIsActive),
		newKv(keyName, expectedName),
		newKv(keyWeight, expectedWeight),
	}
	tree := createKvTree(ctx, key, kvs, key)
	var child *KeyValueTree
	// Born
	child = tree.Children[KeyString(keyBorn)]
	require.NotNil(t, child)
	require.Equal(t, expectedBorn, string(child.Value))
	// Id
	child = tree.Children[KeyString(keyId)]
	require.NotNil(t, child)
	require.Equal(t, expectedId, string(child.Value))
	// IsActive
	child = tree.Children[KeyString(keyIsActive)]
	require.NotNil(t, child)
	require.Equal(t, expectedIsActive, string(child.Value))
	// Name
	child = tree.Children[KeyString(keyName)]
	require.NotNil(t, child)
	require.Equal(t, expectedName, string(child.Value))
	// Weight
	child = tree.Children[KeyString(keyWeight)]
	require.NotNil(t, child)
	require.Equal(t, expectedWeight, string(child.Value))
	t.Log("\n" + fmt.Sprint(tree))
}

func TestKvTreeSimple(t *testing.T) {
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
	require.Nil(t, nil, tree.Value)

	// log tree
	logTree(ctx, tree)
}

func TestKvTreeScalar(t *testing.T) {
	ctx, _ := initAndTest(t)

	// test data
	kvs := []*KeyValue{
		newKv("/team/1/name", "bum"),
	}

	key := "/team/1/name"
	tree := createKvTree(ctx, key, kvs, key)
	t.Log("\n" + fmt.Sprint(tree))
	require.NotNil(t, tree)
	require.Equal(t, "bum", string(tree.Value))
	require.Equal(t, 0, len(tree.Children))

	// log tree
	logTree(ctx, tree)
}

func TestKvTreeStruct(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)

	// Setup KVs
	key := "/test/struct/ecotest"
	expectedName := "Me"
	expectedLuckyNumber := 13
	expectedNoTag := "no-tag-value"
	keyName := fmt.Sprintf("%s/%s", key, "name")
	keyLuckyNumber := fmt.Sprintf("%s/%s", key, "lucky_number")
	keyNoTag := fmt.Sprintf("%s/%s", key, "NoTag")

	// Encode []byte values for struct fields
	bufName := []byte(expectedName)
	bufLuckyNumber, err := json.Marshal(expectedLuckyNumber)
	require.NoError(t, err)
	bufNoTag := []byte(expectedNoTag)
	kvs := []*KeyValue{
		{Key: keyName, Value: bufName},
		{Key: keyLuckyNumber, Value: bufLuckyNumber},
		{Key: keyNoTag, Value: bufNoTag},
	}

	kvTree := createKvTree(ctx, key, kvs, key)
	logTree(ctx, kvTree)
}

func logTree (ctx *common.EcoContext, tree *KeyValueTree) {
	ctx.Logger.Signature("logTree", string(tree.Value), len(tree.Children))
	ctx.Inc()
	defer ctx.Dec()
	
	for _, child := range tree.Children {
		logTree(ctx, child)
	}
}

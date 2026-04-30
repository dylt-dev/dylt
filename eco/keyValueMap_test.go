package eco

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd1(t *testing.T) {
	ctx, _ := initAndTest(t)

	kvm := KeyValueMap{}
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
	for _, kv := range kvs {
		kvm.add(ctx, key, kv)
	}
	t.Log(kvm)
}

func TestAdd2(t *testing.T) {
	ctx, _ := initAndTest(t)

	kvm := KeyValueMap{}
	kvs := []*KeyValue{
		newKv("/team/1/name", "bum"),
		newKv("/team/1/color", "blue"),
		newKv("/team/2/name", "buzz"),
		newKv("/team/2/color", "green"),
		newKv("/team/3/color", "pink"),
		newKv("/team/3/stats/age", "33"),
		newKv("/team/3/stats/rating", "100"),
	}
	key := "/team/1"
	for _, kv := range kvs {
		kvm.add(ctx, key, kv)
	}
	t.Log(kvm)
}

func TestAdd3(t *testing.T) {
	ctx, _ := initAndTest(t)

	kvm := KeyValueMap{}
	kvs := []*KeyValue{
		newKv("/team/1/name", "bum"),
		newKv("/team/1/color", "blue"),
		newKv("/team/2/name", "buzz"),
		newKv("/team/2/color", "green"),
		newKv("/team/3/color", "pink"),
		newKv("/team/3/stats/age", "33"),
		newKv("/team/3/stats/rating", "100"),
	}
	key := "/team/1/name"
	for _, kv := range kvs {
		kvm.add(ctx, key, kv)
	}
	t.Log(kvm)
}

func TestStruct(t *testing.T) {
	ctx, _ := initAndTest(t)
	ctx.Inc()
	defer ctx.Dec()

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

	kvMap := createKvMap(ctx, key, kvs)
	ctx.Logger.Info(kvMap)
}

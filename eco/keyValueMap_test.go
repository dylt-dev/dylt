package eco

import (
	"testing"
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
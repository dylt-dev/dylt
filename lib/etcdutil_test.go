package lib

import (
	"context"
	"fmt"
	"runtime/debug"
	"testing"

	clientV3 "go.etcd.io/etcd/client/v3"
)


func TestEtcdGet (t *testing.T) {
	cli, err := NewEtcdClient("hello.dylt.dev")
	// cli, err := NewEtcdClient("helloXXX.dylt.dev")
	// cli, err := NewEtcdClient("google.com")
	if err != nil { t.Fatal(err, debug.Stack()) }
	val, err := cli.Get("foo")
	if err != nil { t.Fatal(err, debug.Stack()) }
	fmt.Println(string(val))
}


func TestEtcdGetWithPrefix (t *testing.T){
	cli, err := NewEtcdClient("hello.dylt.dev")
	if err != nil { t.Fatal(err, debug.Stack()) }
	resp, err := cli.Client.Get(context.Background(), "/vm", clientV3.WithPrefix())
	if err != nil { t.Fatal(err, debug.Stack()) }
	t.Logf("Count    = %d", resp.Count)
	t.Logf("len(Kvs) = %d", len(resp.Kvs))
	for i, kv := range(resp.Kvs) {
		t.Logf("Kvs[%d].Version = %s\n", i, kv.Key)
	}
}
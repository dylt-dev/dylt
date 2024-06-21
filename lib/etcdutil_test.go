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


func TestEtcdGetNonExistent (t *testing.T) {
	cli, err := NewEtcdClient("hello.dylt.dev")
	key := "ABCDEF1234"
	if err != nil { t.Fatal(err, debug.Stack()) }
	val, err := cli.Get(key)
	if val != nil { t.Fatal(err, debug.Stack()) } 
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


func TestEtcdGetKeys (t *testing.T) {
	cli, err := NewEtcdClient("hello.dylt.dev")
	if err != nil { t.Fatal(err)}
	keys, err := cli.GetKeys("/vm")
	if err != nil { t.Fatal(err)}
	fmt.Printf("len(keys)=%d\n", len(keys))
	for _, key := range keys {
		fmt.Printf("key=%s\n", key)
	}
}


func TestEtcdList (t *testing.T) {
	cli, err := NewEtcdClient("hello.dylt.dev")
	if err != nil { t.Fatal(err, debug.Stack()) }
	kvs, err := cli.List()
	if err != nil { t.Fatal(err, debug.Stack()) }
	fmt.Println(len(kvs))
	for _, kv := range kvs {
		fmt.Printf("kv.Key=%s\n", kv.Key)
	}
}


func TestPointerFun (t *testing.T) {
	var vm *VmInfo = ReturnPtr()
	fmt.Printf("%v\n", vm)
	if (vm != nil) { t.Fatal("no clue")}
}

func ReturnPtr () *VmInfo { return nil}



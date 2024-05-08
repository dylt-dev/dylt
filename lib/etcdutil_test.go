package lib

import (
	"fmt"
	"testing"
)


func TestEtcdGet (t *testing.T) {
	cli, err := NewEtcdClient("hello.dylt.dev")
	// cli, err := NewEtcdClient("helloXXX.dylt.dev")
	// cli, err := NewEtcdClient("google.com")
	if err != nil {
		t.Fatal("NewEtcdClient() failed", err)
	}
	val, err := cli.Get("foo")
	if err != nil {
		t.Fatal("EtcdGet() failed", err)
	}
	fmt.Println(string(val))
}
package lib

import (
	"fmt"
	"encoding/json"
	"runtime/debug"
	"testing"
)


func TestVmAdd (t *testing.T) {
	vmClient, err := NewVmClient("hello.dylt.dev")
	if err != nil { t.Fatal(err, debug.Stack()) }
	name := "test-vm"
	address := "test-vm.dylt.dev"
	vm, err := vmClient.Add(name, address)
	if err != nil { t.Fatal(err, debug.Stack()) }
	fmt.Printf("vm=%#v\n", vm)
}


func TestVmAll (t *testing.T) {
	vmClient, err := NewVmClient("hello.dylt.dev")
	if err != nil { t.Fatal(err, debug.Stack()) }
	vmInfoMap, err := vmClient.All()
	if err != nil { t.Fatal(err, debug.Stack()) }
	for k, v := range vmInfoMap {
		fmt.Printf("VM name=%s\n", k)
		s, err := json.Marshal(v)
		if err != nil { t.Fatal(err, debug.Stack()) }
		fmt.Printf("VM info=%s\n", s)
	}
}


func TestVmGet(t *testing.T) {
	cli, err := NewVmClient("hello.dylt.dev")
	if err != nil {
		t.Fatal(err, debug.Stack())
	}
	vm, err := cli.Get("0")
	if err != nil {
		t.Fatal(err, debug.Stack())
	}
	s, _ := json.Marshal(vm)
	t.Logf("vm=%s\n", s)
}

func TestVmList(t *testing.T) {
	cli, err := NewVmClient("hello.dylt.dev")
	if err != nil {
		t.Fatal(err, debug.Stack())
	}
	names, err := cli.Names()
	if err != nil {
		t.Fatal(err, debug.Stack())
	}
	for _, name := range names {
		t.Logf("name=%s\n", name)
	}
}

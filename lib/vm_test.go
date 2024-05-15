package lib

import (
	"fmt"
	"encoding/json"
	"runtime/debug"
	"strings"
	"testing"
)


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
	t.Logf("vm.Host=%s\n", vm.Host)
	t.Logf("vm.Shr=%t\n", vm.Shr)
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
		name, _ := strings.CutPrefix(string(name), PRE_vm)
		t.Logf("name=%s\n", name)
	}
}

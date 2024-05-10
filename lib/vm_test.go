package lib

import (
	"runtime/debug"
	"strings"
	"testing"
)

func TestVmGet (t *testing.T) {
	cli, err := NewVmClient("hello.dylt.dev")
	if err != nil { t.Fatal(err, debug.Stack()) }
	vm, err := cli.Get("0")
	if err != nil { t.Fatal(err, debug.Stack()) }
	t.Logf("vm.Host=%s\n", vm.Host)
	t.Logf("vm.Shr=%t\n", vm.Shr)
}


func TestVmList (t *testing.T) {
	cli, err := NewVmClient("hello.dylt.dev")
	if err != nil { t.Fatal(err) }
	names, err := cli.Names() 
	if err != nil { t.Fatal(err) }
	for _, name := range(names) {
		name, _ := strings.CutPrefix(string(name), "/vm/")
		t.Logf("name=%s\n", name)
	}
}
package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"testing"
	"time"
)


func doit (n *int) {
	fmt.Println(*n)
}
func Test (t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			doit(&i)
		}()
	}
	time.Sleep(9)
}


func TestSimplePut (t *testing.T) {
	cli, err := NewVmClient("hello.dylt.dev")
	if err != nil { t.Fatal(err, debug.Stack()) }
	vm := VmInfo{
		Address: "hosty toasty host",
		Name: "ovh-vps0",
	}
	buf, _ := json.Marshal(vm)
	s := string(buf)
	t.Logf("s=%s\n", s)
	name := "ovh-vps0"
	key := fmt.Sprintf("/vm/%s", name)
	cli.Client.Put(context.Background(), key, s)
}

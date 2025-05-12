package eco

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init () {
//	InitConfig()
}


func TestVmAdd (t *testing.T) {
	vmClient, err := CreateVmClientFromConfig()
	if err != nil { t.Fatal(err, debug.Stack()) }
	name := "test-vm"
	address := "test-vm.dylt.dev"
	vm, err := vmClient.Add(name, address)
	if err != nil { t.Fatal(err, debug.Stack()) }
	fmt.Printf("vm=%#v\n", vm)
}


func TestVmAll (t *testing.T) {
	vmClient, err := CreateVmClientFromConfig()
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
	vmClient, err := CreateVmClientFromConfig()
	if err != nil {
		t.Fatal(err, debug.Stack())
	}
	vm, err := vmClient.Get("0")
	if err != nil {
		t.Fatal(err, debug.Stack())
	}
	s, _ := json.Marshal(vm)
	t.Logf("vm=%s\n", s)
}

func TestVmList(t *testing.T) {
	vmClient, err := CreateVmClientFromConfig()
	if err != nil {
		t.Fatal(err, debug.Stack())
	}
	names, err := vmClient.Names()
	if err != nil {
		t.Fatal(err, debug.Stack())
	}
	for _, name := range names {
		t.Logf("name=%s\n", name)
	}
}

func TestVmInfoGetAddress(t *testing.T) {
	initName := "vm-test"
	initAddress := "vm-test.dylt.dev"
	vm := VmInfo {
		Address: initAddress,
		Name: initName,
	}
	field := "Address"
	fieldValue, err := vm.Get(field)
	assert.Nil(t, err)
	assert.Equal(t, initAddress, fieldValue)
}

func TestVmInfoGetInvalid(t *testing.T) {
	initName := "vm-test"
	initAddress := "vm-test.dylt.dev"
	vm := VmInfo {
		Address: initAddress,
		Name: initName,
	}
	field := "INVALID-FIELD-NAME"
	fieldValue, err := vm.Get(field)
	assert.Nil(t, fieldValue)
	assert.NotNil(t, err)
	t.Logf("err=%s\n", err)
}


func TestVmInfoGetName(t *testing.T) {
	initName := "vm-test"
	initAddress := "vm-test.dylt.dev"
	vm := VmInfo {
		Address: initAddress,
		Name: initName,
	}
	field := "Name"
	fieldValue, err := vm.Get(field)
	assert.Nil(t, err)
	assert.Equal(t, initName, fieldValue)
}


func TestVmInfoSetAddress (t *testing.T) {
	vm := VmInfo{}
	initAddress := "vm-test.dylt.dev"
	err := vm.Set("Address", initAddress)
	assert.Nil(t, err)
	assert.Equal(t, initAddress, vm.Address)
}


func TestVmInfoSetInvalid (t *testing.T) {
	vm := VmInfo{}
	err := vm.Set("INVALID-FIELD-NAME", "")
	assert.NotNil(t, err)
	t.Logf("err=%s\n", err)
}


func TestVmInfoSetName (t *testing.T) {
	vm := VmInfo{}
	initName := "vm-test"
	err := vm.Set("Name", initName)
	assert.Nil(t, err)
	assert.Equal(t, initName, vm.Name)
}


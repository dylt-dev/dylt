package eco

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"testing"
	"unsafe"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestSimplePut(t *testing.T) {
	vmClient, err := CreateVmClientFromConfig()
	if err != nil {
		t.Fatal(err, debug.Stack())
	}
	vm := VmInfo{
		Address: "hosty toasty host",
		Name:    "ovh-vps0",
	}
	buf, _ := json.Marshal(vm)
	s := string(buf)
	t.Logf("s=%s\n", s)
	name := "ovh-vps0"
	key := fmt.Sprintf("/vm/%s", name)
	vmClient.Client.Put(context.Background(), key, s)
}

// func TestLoadConfig(t *testing.T) {
// 	cfg := Config{}
// 	err := cfg.Load()
// 	assert.Nil(t, err)
// 	domain, _ := cfg.GetEtcDomain()
// 	assert.Empty(t, domain)
// }

// func TestLoadConfig2(t *testing.T) {
// 	cfg := Config{}
// 	err := cfg.Load()
// 	assert.Nil(t, err)
// 	domain, _ := cfg.GetEtcDomain()
// 	assert.NotEmpty(t, domain)
// 	assert.Equal(t, "hello.dylt.dev", domain)
// }

// func TestSaveConfig(t *testing.T) {
// 	cfg := Config{}
// 	err := cfg.Load()
// 	assert.Nil(t, err)
// 	err = cfg.SetEtcDomain("hello.dylt.dev")
// 	assert.Nil(t, err)
// 	err = cfg.Save()
// 	assert.Nil(t, err)
// }

// func TestInitConfig(t *testing.T) {
// 	viper.SetConfigName(CFG_Filename)
// 	viper.SetConfigType(CFG_Type)
// 	viper.AddConfigPath(".")
// 	cfgFolder := GetConfigFolderPath()
// 	viper.AddConfigPath(cfgFolder)
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		panic(fmt.Errorf("fatal error config file: %s", err.Error()))
// 	}
// }

// func TestClearConfig(t *testing.T) {
// 	err := ClearConfigFile()
// 	assert.Nil(t, err)
// 	cfgFilePath := GetConfigFilePath()
// 	f, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 400)
// 	assert.Nil(t, err)
// 	defer f.Close()
// 	fi, err := f.Stat()
// 	assert.Nil(t, err)
// 	assert.NotNil(t, fi)
// 	assert.Equal(t, int64(0), fi.Size())
// }

func TestInit(t *testing.T) {
	// Init the config
	const etcdDomain = "hello.dylt.dev"
	initInfo := common.InitStruct{
		EtcdDomain: etcdDomain,
	}
	err := common.Init(&initInfo)
	assert.Nil(t, err)
	// Test the file exists
	cfg, err := common.LoadConfig()
	assert.Nil(t, err)
	t.Logf("%#v", cfg)
	assert.Nil(t, err)
	// Test the file contains the expected domain
	assert.Equal(t, etcdDomain, cfg.EtcdDomain)
}

func TestShowConfig(t *testing.T) {
	err := common.ShowConfig(os.Stdout)

	assert.Nil(t, err)
}

func TestStrings(t *testing.T) {
	str0 := "> GET /etcd-io/etcd/releases/download/v3.5.16/etcd-v3.5.16-linux_amd64.tar.gz HTTP/2"
	str1 := "> GET /etcd-io/etcd/releases/download/v3.5.16/etcd-v3.5.16-linux-amd64.tar.gz HTTP/2"
	len0 := len(str0)
	len1 := len(str1)
	assert.Equal(t, len0, len1)
	for i := range len0 {
		t.Logf("Checking index #%d ...\n", i)
		assert.Equal(t, str0[i], str1[i])
	}
	assert.Equal(t, str0, str1)
}

func TestPostViper(t *testing.T) {
	// Setup YAML string of config info
	var yml = `
etcd-domain: hello.dylt.dev
`
	cfg1 := common.ConfigStruct{
		EtcdDomain: "hello.dylt.dev",
	}

	cfg2 := common.ConfigStruct{}
	err := yaml.Unmarshal([]byte(yml), &cfg2)
	assert.Nil(t, err)
	assert.Equal(t, cfg1.EtcdDomain, cfg2.EtcdDomain)
}

func TestPostViperWrite(t *testing.T) {
	cfg := common.ConfigStruct{
		EtcdDomain: "hello.dylt.dev",
	}
	path := "/tmp/dylt.cfg"
	f, err := os.Create(path)
	assert.Nil(t, err)
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(cfg)
	assert.Nil(t, err)
	assert.Nil(t, err)
}

func TestTypeName0(t *testing.T) {
	type stringAlias string
	t.Logf("fullTypeName(stringAlias)=%s", common.FullTypeName(reflect.TypeOf(*new(stringAlias))))
	t.Logf("%#v", reflect.TypeOf(*new(stringAlias)))
	t.Logf("%#v", reflect.TypeOf(stringAlias("")))
	t.Logf("%#v", reflect.TypeOf(*new(stringAlias)).Name())
}

func TestTypeName1(t *testing.T) {
	type stringsAlias []string
	t.Logf("fullTypeName(stringsAlias)=%s", common.FullTypeName(reflect.TypeOf(*new(stringsAlias))))
	t.Logf("%#v", reflect.TypeOf(*new(stringsAlias)))
	t.Logf("%#v", reflect.TypeOf(stringsAlias{}))
	t.Logf("%#v", reflect.TypeOf(*new(stringsAlias)).Name())
	var i any = *new(stringsAlias)
	t.Logf("%#v", reflect.TypeOf(i).Name())
}

func TestReflectionIntPointer(t *testing.T) {
	var buf []byte
	var n int = 13
	var err error

	// Encode an int into a []byte using the JSON encoder
	buf, err = json.Marshal(n)
	require.NoError(t, err)

	// Create a slice 1 int long
	var valSlice reflect.Value = reflect.MakeSlice(reflect.TypeOf([]int{}), 1, 1)
	t.Logf("Before: %d", valSlice.Index(0).Int())

	// See if Addr works with the slice index
	valPtr := valSlice.Index(0).Addr()
	ptr := valPtr.Interface()
	t.Logf("valPtr Kind=%s", valPtr.Kind().String())
	t.Logf("ptr Kind=%s", reflect.TypeOf(ptr).Kind().String())
	t.Logf("ptr Elem Kind=%s", reflect.TypeOf(ptr).Elem().Kind().String())
	err = json.Unmarshal(buf, ptr)
	t.Logf("After (valPtr): %d", valSlice.Index(0).Int())

	// Get an Interface{} to the slice element. Interfaces are useful: type plus storage
	var el reflect.Value = valSlice.Index(0)
	var el2 any = el.Interface()
	t.Logf("el2 TypeOf = %s", reflect.TypeOf(el2).Name())

	// Decode the byte[] into our Interface
	err = json.Unmarshal(buf, &el2)
	t.Logf("el2 TypeOf = %s", reflect.TypeOf(el2).Name())
	require.NoError(t, err)
	t.Logf("el2=%v", el2)

	// see if updating the Interface also updated the slice
	t.Logf("After: %d", valSlice.Index(0).Int())

	// set slice element explicitlyA
	tyElem := valSlice.Type().Elem()
	t.Logf("tyElem=%s", tyElem.Name())
	valSlice.Index(0).Set(reflect.ValueOf(el2).Convert(tyElem))
	t.Logf("After-After: %d", valSlice.Index(0).Int())

	// Now try it without the slice: create a single int and get its interface
	ty := reflect.TypeFor[int]()
	valInt := reflect.New(ty)
	n2 := valInt.Interface()

	// Unmarshal into the pointer to the interface
	err = json.Unmarshal(buf, &n2)

	// Print the interface value and the Value value
	t.Logf("n2=%v", el2)
	t.Logf("valInt Kind=%s", valInt.Kind().String())
	t.Logf("valInt Elem Kind = %s", valInt.Elem().Kind().String())
	t.Logf("valInt Elem Int = %d", valInt.Elem().Int())
	t.Logf("Indirect(valInt) Int = %d", reflect.Indirect(valInt).Int())
}
func TestReflectionSlice(t *testing.T) {
	var ints = make([]int, 3)
	t.Logf("Before: %v", ints)
	var v = reflect.ValueOf(ints)
	var el = v.Index(1)
	require.True(t, el.CanSet())
	el.Set(reflect.ValueOf(13))
	t.Logf("After: %v", ints)

}

// type emptyInterface struct {
// 	typ  *abi.Type
// 	word unsafe.Pointer
// }

// func TestReflectionMisc0(t *testing.T) {
// 	var n8 int8 = 13
// 	dump(n8)
// 	var u8 uint8 = 13
// 	dump(u8)
// 	var n16 int16 = 13
// 	dump(n16)
// 	var u16 uint16 = 13
// 	dump(u16)
// 	var n32 int32 = 13
// 	dump(n32)
// 	var u32 uint32 = 13
// 	dump(u32)
// 	var n64 int64 = 13
// 	dump(n64)
// 	var u64 uint64 = 13
// 	dump(u64)
// 	var b bool = false
// 	dump(b)
// 	var s string = "ASS"
// 	dump(s)
// 	var sss string = "THAT ASSSSSSSS"
// 	dump(sss)
// }

// func dump[T any](arg T) {
// 	// ptr := unsafe.Pointer(&arg)
// 	// t.Logf("*ptr=%v", *((*T)(ptr)))

// 	var emp emptyInterface
// 	var i any = arg
// 	var iptr = unsafe.Pointer(&i)
// 	emp = *(*emptyInterface)(iptr)
// 	// t.Logf("emp.type: %#v", emp.typ)

// 	var wptr unsafe.Pointer = emp.word
// 	var valptr *T = (*T)(wptr)
// 	// t.Logf("emp.word: %v", *valptr)

// 	var typ abi.Type = *emp.typ
// 	var nkind uint8 = typ.Kind_
// 	var skind string = abi.KindNames[nkind]
// 	var size uintptr = typ.Size_
// 	// t.Logf("%d %s %d", nkind, skind, size)

// 	fmt.Printf("I have news ... your variable is of type %s, it is %d bytes long, and its value is %v\n", skind, size, *valptr)
// }

//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func TestNoEscape(t *testing.T) {
	var n int8 = 13
	uptr := unsafe.Pointer(&n)
	t.Logf("%-12s: %p", "uptr", uptr)
	nuptr := noescape(uptr)
	t.Logf("%-12s: %p", "nuptr", nuptr)
}

func TestXor(t *testing.T) {
	n := 13
	t.Logf("%d %d", n, n^0)
}

func TestPrints(t *testing.T) {
	var ns int16 = 13
	var us uint16 = 13
	t.Logf("%#v %#v", ns, us)
	var ans any = ns
	var aus any = us
	t.Logf("%#v %#v", ans, aus)
}

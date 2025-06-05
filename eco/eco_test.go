package eco

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
)

// Sample types for tests (named types look better in a log file than anonymous types)
type emptyStruct struct{}
type arrayOfInt [0]int
type arrayOfStruct [0]emptyStruct
type mapOfIntSlice map[string][]int
type pointerToInt *int
type sliceOfInt []int
type sliceOfStruct []emptyStruct
type map_emptyStruct_emptyStruct map[struct{}]struct{}
type map_emptyStruct_int map[struct{}]int
type map_int_emptyStruct map[int]struct{}
type map_string_int map[string]int
type map_string_struct map[string]EcoTest

type Stat struct {
	Name string
	Value float64
}

type MiscMap map[string]string
type StatSlice []Stat
  
type Player struct {
	Id int
	Name string
	Weight float64
	IsActive bool
	Stats StatSlice
	Misc MiscMap
}

type PlayerMap map[string]Player
  
type Team struct {
	Name string
	Players PlayerMap
}

// Sample objects for tests
var VAL_MapSimple = map_string_int{}
var VAL_MapUnsimple = map_emptyStruct_emptyStruct{}
var VAL_MapUnsimpleKey = map_emptyStruct_int{}
var VAL_MapUnsimpleValue = map_int_emptyStruct{}
var VAL_SliceSimple = []int{5, 8, 13}
var VAL_SliceUnsimple = []emptyStruct{}
var VAL_SimplePointer = new(int)
var VAL_UnsimplePointer = &(emptyStruct{})
var VAL_MapWithStructKey = map[EcoTest]string{}
var VAL_Map_String_Struct = map_string_struct{"test": *NewEcoTest("me", 13)}

var VAL_AltuveStats = StatSlice{
	{Name: "All-Stars", Value: 9},
	{Name: "Height", Value: 5.5},
}
var VAL_AltuveMisc = map[string]string {
	"Born": "Venezuela",
}
var VAL_Altuve = Player{
	Name: "Jose Altuve",
	Id: 1,
	IsActive: true,
	Weight: 1.0,
	Stats: VAL_AltuveStats,
	Misc: VAL_AltuveMisc,
}

var VAL_JavierStats = []Stat{
	{Name: "Number", Value: 53},
	{Name: "Debut", Value: 2020},
}
var VAL_JavierMisc = map[string]string{
	"Born": "Dominican Republic",
}
var VAL_Javier = Player{
	Name: "Christian Javier",
	Id: 2,
	IsActive: false,
	Weight: 0.5,
	Stats: VAL_JavierStats,
	Misc: VAL_JavierMisc,
}

var VAL_PenaStats = []Stat{
	{Name: "Gold Gloves", Value: 1},
}
var VAL_PenaMisc = map[string]string{
	"Raised": "Rhode Island",
}
var VAL_Pena = Player{
	Name: "Jeremy Pena",
	Id: 3,
	IsActive: true,
	Weight: 0.9,
	Stats: VAL_PenaStats,
	Misc: VAL_PenaMisc,
}

var VAL_Players = map[string]Player{
	"altuve": VAL_Altuve,
	"javier": VAL_Javier,
	"pena": VAL_Pena,
}

var VAL_Astros = Team{
	Name: "Astros",
	Players: VAL_Players,
}
type EcoTest struct {
	Name        string  `eco:"name"`
	LuckyNumber float64 `eco:"lucky_number"`
	Anon        string
}

type structWithMap struct {
	M map[int]string
}

type UnsimpleStruct struct {
	C complex64
	F func()
}

func NewEcoTest(name string, luckyNumber float64) *EcoTest {
	return &EcoTest{Name: name, LuckyNumber: luckyNumber}
}

// For logging
func TestCreateSignature0 (t *testing.T) {
	sig := createSignature("greatness", "foo", "bar")
	t.Log(sig)
}

func TestGetChildKeys (t *testing.T) {
	cli, err := CreateEtcdClientFromConfig()
	require.NoError(t, err)
	prefix := "/test/team/astros/Players"
	resp, err := cli.Client.Get(context.Background(), prefix, etcd.WithKeysOnly(), etcd.WithPrefix())
	require.NoError(t, err)
	var allKeys = make([][]byte, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		allKeys[i] = kv.Key
	}
	// t.Logf("%#v", allKeys)

	var srx string = fmt.Sprintf(`^%s/?\w+$`, prefix)
	var rx *regexp.Regexp = regexp.MustCompile(srx)
	for _, key := range allKeys {
		t.Logf("Matching %s ...", key)
		if rx.Match(key) {
			t.Log(string(key))
		}
	}
}


var childKeys = []string {
    "/test/team/astros/Players/altuve/Id",
    "/test/team/astros/Players/altuve/IsActive",
    "/test/team/astros/Players/altuve/Misc/Born",
    "/test/team/astros/Players/altuve/Name",
    "/test/team/astros/Players/altuve/Stats/0/Name",
    "/test/team/astros/Players/altuve/Stats/0/Value",
    "/test/team/astros/Players/altuve/Stats/1/Name",
    "/test/team/astros/Players/altuve/Stats/1/Value",
    "/test/team/astros/Players/altuve/Weight",
    "/test/team/astros/Players/javier/Id",
    "/test/team/astros/Players/javier/IsActive",
    "/test/team/astros/Players/javier/Misc/Born",
    "/test/team/astros/Players/javier/Name",
    "/test/team/astros/Players/javier/Stats/0/Name",
    "/test/team/astros/Players/javier/Stats/0/Value",
    "/test/team/astros/Players/javier/Stats/1/Name",
    "/test/team/astros/Players/javier/Stats/1/Value",
    "/test/team/astros/Players/javier/Weight",
    "/test/team/astros/Players/pena/Id",
    "/test/team/astros/Players/pena/IsActive",
    "/test/team/astros/Players/pena/Misc/Raised",
    "/test/team/astros/Players/pena/Name",
    "/test/team/astros/Players/pena/Stats/0/Name",
    "/test/team/astros/Players/pena/Stats/0/Value",
    "/test/team/astros/Players/pena/Weight",
}

// func TestGetChildKeys1 (t *testing.T) {
// 	prefix := "/test/team/astros/Players/javier/Stats"
// 	var sliceKeys []int
// 	sliceKeys, err := getSliceKeys(nil, prefix)
// 	require.NoError(t, err)
// 	t.Logf("%v", sliceKeys)
// 	var maxKey = slices.Max(sliceKeys)
// 	t.Logf("maxKey=%d", maxKey)
// 	var len = maxKey+1
// 	t.Logf("len=%d", len)
// }

func TestMatchChildKey (t *testing.T) {
	prefix := "/test/team/astros/Players"
	var srx string = fmt.Sprintf(`^%s/?\w+$`, prefix)
	var rx *regexp.Regexp = regexp.MustCompile(srx)
	key := "/test/team/astros/Players/altuve"
	require.True(t, rx.Match([]byte(key)))
	badkey := "/test/team/astros/Players/altuve/Misc"
	require.False(t, rx.Match([]byte(badkey)))
}


func TestGetObject(t *testing.T) {

}

func TestFullTypeName_StatSlice (t *testing.T) {
	s := fullTypeName(reflect.TypeOf(*new(StatSlice)))
	t.Log(s)
}

func TestKind_ArrayOfInt(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), arrayOfInt{})
	assert.Equal(t, SimpleArray, kind)
}

func TestKind_ArrayOfStruct(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), arrayOfStruct{})
	assert.Equal(t, Invalid, kind)
}

func TestKind_MapOfIntSlice(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), mapOfIntSlice{})
	assert.Equal(t, SimpleMap, kind)

}

func TestKind_MapOfSliceOfStruct(t *testing.T) {
	type emptyStruct struct{}
	type mapOfSlice map[string][]emptyStruct
	kind := getKind(newEcoContext(os.Stdout), mapOfSlice{})
	assert.Equal(t, Invalid, kind)
}

func TestKind_MapSimple(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), VAL_MapSimple)
	assert.Equal(t, SimpleMap, kind)
}

func TestKind_MapUnsimple(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), VAL_MapUnsimple)
	assert.Equal(t, Invalid, kind)
}

func TestKind_MapUnsimpleKey(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), VAL_MapUnsimpleKey)
	assert.Equal(t, Invalid, kind)
}

func TestKind_MapUnsimpleValue(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), VAL_MapUnsimpleValue)
	assert.Equal(t, SimpleMap, kind)
}

func TestKind_PointerToInt(t *testing.T) {
	var pint pointerToInt
	kind := getKind(newEcoContext(os.Stdout), pint)
	assert.Equal(t, SimplePointer, kind)
}

func TestKind_PointerToIntSlice(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), new(sliceOfInt))
	assert.Equal(t, SimplePointer, kind)
}

func TestKind_PointerToStructSlice(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), new(sliceOfStruct))
	assert.Equal(t, Invalid, kind)
}

func TestKind_SliceSimple(t *testing.T) {
	type intSlice []int
	i := intSlice{1, 2, 3}
	kind := getKind(newEcoContext(os.Stdout), i)
	assert.Equal(t, SimpleSlice, kind)
}

func TestKind_SliceOfMap(t *testing.T) {
	type simpleMap map[int]int
	type sliceOfMap []simpleMap
	i := sliceOfMap{}
	kind := getKind(newEcoContext(os.Stdout), i)
	assert.Equal(t, Invalid, kind)
}

func TestKind_SliceUnsimple(t *testing.T) {
	type emptyStruct struct{}
	type emptyStructSlice []emptyStruct
	i := emptyStructSlice{}
	kind := getKind(newEcoContext(os.Stdout), i)
	assert.Equal(t, Invalid, kind)
}

func TestKind_StatSlice (t *testing.T) {
	val := StatSlice{}
	kind := getKind(newEcoContext(os.Stdout), val)
	assert.Equal(t, SimpleSlice, kind, fmt.Sprintf("Expected %s, got %s", SimpleSlice.String(), kind.String()))
}

func TestKind_StructSimple(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), EcoTest{})
	assert.Equal(t, SimpleStruct, kind)
}

func TestKind_StructUnsimple(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), UnsimpleStruct{})
	assert.Equal(t, Invalid, kind)
}

func TestKind_StructWithMap(t *testing.T) {
	kind := getKind(newEcoContext(os.Stdout), structWithMap{})
	assert.Equal(t, SimpleStruct, kind)
}

func TestKind_StructWithSlice(t *testing.T) {
	type structWithSlice struct{ Slice []int }
	kind := getKind(newEcoContext(os.Stdout), structWithSlice{})
	assert.Equal(t, SimpleStruct, kind)
}

func TestKind_StructWithMapWithSlice(t *testing.T) {
	type mapWithSlice map[string][]int
	type structWithMapWithSlice struct{ M mapWithSlice }
	kind := getKind(newEcoContext(os.Stdout), structWithMapWithSlice{})
	assert.Equal(t, SimpleStruct, kind)
}

func TestKind_StructWithMapWithStruct(t *testing.T) {
	type innerStruct struct{}
	type mapWithStruct map[string]innerStruct
	type structWithMapWithSlice struct{ MapField mapWithStruct }
	kind := getKind(newEcoContext(os.Stdout), structWithMapWithSlice{})
	assert.Equal(t, SimpleStruct, kind)
}

func TestPutObject0(t *testing.T) {
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)
	obj := NewEcoTest("Me", 13)

	prefix := "/_test_/echotest"

	ops := []etcd.Op{}
	opName := etcd.OpPut(filepath.Join(prefix, "Name"), obj.Name, etcd.WithPrevKV())
	ops = append(ops, opName)
	opLuckyNumber := etcd.OpPut(filepath.Join(prefix, "LuckyNumber"), strconv.FormatFloat(obj.LuckyNumber, 'f', 8, 64), etcd.WithPrevKV())
	ops = append(ops, opLuckyNumber)

	txn := etcdClient.Txn(context.Background())
	resp, err := txn.Then(ops...).Commit()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	for _, respOp := range resp.Responses {
		t.Logf("respOp.GetResponsePut().PrevKV: %s => %s", respOp.GetResponsePut().PrevKv.Key, respOp.GetResponsePut().PrevKv.Value)
	}
}

func TestPutObject1(t *testing.T) {
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)
	obj := NewEcoTest("Me", 13)

	prefix := "/_test_/echotest"

	ops := []etcd.Op{}
	opName := etcd.OpPut(filepath.Join(prefix, "Name"), obj.Name, etcd.WithPrevKV())
	ops = append(ops, opName)
	opLuckyNumber := etcd.OpPut(filepath.Join(prefix, "LuckyNumber"), strconv.FormatFloat(obj.LuckyNumber, 'f', 8, 64), etcd.WithPrevKV())
	ops = append(ops, opLuckyNumber)

	txn := etcdClient.Txn(context.Background())
	resp, err := txn.Then(ops...).Commit()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	for _, respOp := range resp.Responses {
		// t.Logf("respOp=%#v", respOp)
		// t.Logf("respOp.GetResponsePut()=%#v", respOp.GetResponsePut())
		// t.Logf("respOp.GetResponsePut()=%#v", respOp.GetResponsePut())
		// t.Logf("respOp.GetResponsePut()=%#v", respOp.GetResponsePut())
		// t.Logf("respOp.GetResponsePut().PrevKv=%#v", respOp.GetResponsePut().PrevKv)
		t.Logf("respOp.GetResponsePut().PrevKV: %s => %s", respOp.GetResponsePut().PrevKv.Key, respOp.GetResponsePut().PrevKv.Value)

	}
}

func TestReflection0(t *testing.T) {
	obj := NewEcoTest("Me", 13)

	var ty reflect.Type
	var val reflect.Value
	ty = reflect.TypeOf(obj)
	if ty.Kind() == reflect.Pointer {
		t.Log("var is a pointer - dereferencing")
		ty = reflect.TypeOf(*obj)
		val = reflect.ValueOf(*obj)
	} else {
		val = reflect.ValueOf(obj)
	}
	assert.Equal(t, 3, ty.NumField())

	for i := range ty.NumField() {
		sf := ty.Field(i)
		var s string
		switch sf.Type.Kind() {
		case reflect.Float32:
			s = strconv.FormatFloat(val.Field(i).Float(), 'f', -1, 32)
		case reflect.Float64:
			s = strconv.FormatFloat(val.Field(i).Float(), 'f', -1, 64)
		case reflect.String:
			s = string(val.Field(i).String())
		default:
			s = "N/A"
		}
		t.Logf("%s=%s", getFieldKey(sf), s)
	}
}

func TestReflection1(t *testing.T) {
	obj := NewEcoTest("Me", 13)

	ty, val, err := reflectStruct(obj)
	require.NoError(t, err)
	require.NotNil(t, ty)
	require.NotEmpty(t, val)
	dumpStruct(t, ty, val)
}

func TestReflection2(t *testing.T) {
	obj := NewEcoTest("Me", 13)
	var ptr **EcoTest = &obj
	var pptr ***EcoTest = &ptr
	var ppptr ****EcoTest = &pptr

	ty, val, err := reflectStruct(ppptr)
	require.NoError(t, err)
	require.NotNil(t, ty)
	require.NotEmpty(t, val)
	require.Equal(t, 2, ty.NumField())
	dumpStruct(t, ty, val)
}

func TestReflection3(t *testing.T) {
	type selfish *selfish
	var obj selfish
	obj = &obj
	ty, val, err := reflectStruct(obj)
	require.Nil(t, ty)
	require.Empty(t, val)
	require.Error(t, err)
}

func reflectStruct(obj any) (reflect.Type, reflect.Value, error) {
	ty := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	var lastPtr uintptr = uintptr(unsafe.Pointer(&val))
	return reflectStructNoCycle(ty, val, lastPtr)
}

func reflectStructNoCycle(ty reflect.Type, val reflect.Value, lastPtr uintptr) (reflect.Type, reflect.Value, error) {
	switch ty.Kind() {
	case reflect.Struct:
		return ty, val, nil
	case reflect.Pointer:
		fmt.Printf("val: %#x\n", val.Pointer())
		ptr := val.Pointer()
		val = reflect.Indirect(val)
		ty = val.Type()
		if ptr == lastPtr {
			return nil, reflect.Value{}, errors.New("cycle detected")
		}
		return reflectStructNoCycle(ty, val, ptr)
	case reflect.UnsafePointer:
		fmt.Printf("val: %p\n", val.Pointer)
		val = reflect.NewAt(ty, unsafe.Pointer(&val))
		fmt.Printf("ty: %s.%s\n", ty.PkgPath(), ty.Name())
		fmt.Printf("val: %p", val.Pointer)
		return ty, val, nil
	default:
		return nil, reflect.Value{}, fmt.Errorf("not struct or pointer to struct: %d (%s.%s)", ty.Kind(), ty.PkgPath(), ty.Name())
	}
}

func dumpOps(t *testing.T, ops []etcd.Op) {
	for _, op := range ops {
		if op.IsGet() {
			key := string(op.KeyBytes())
			s := fmt.Sprintf("%s %s", "GET", key)
			fmt.Println(s)
		} else if op.IsPut() {
			key := string(op.KeyBytes())
			val := string(op.ValueBytes())
			s := fmt.Sprintf("%s %s %s", "PUT", lowlight(key), val)
			fmt.Println(s)
		} else {
			fmt.Printf("%#v\n", op)
		}
	}
}

func dumpStruct(t *testing.T, ty reflect.Type, val reflect.Value) {
	for i := range ty.NumField() {
		var sf reflect.StructField = ty.Field(i)
		var sfv reflect.Value = val.Field(i)
		fieldName := getFieldKey(sf)
		fieldValue, err := getFieldValue(sfv)
		assert.NoError(t, err)
		t.Logf("%s=%s", fieldName, fieldValue)
	}
}

func getSliceKeys (ctx *ecoContext, cli *EtcdClient, prefix string) ([]int, error) {
	ctx.logger.signature("getSliceKeys", prefix)
	childKeys, err := cli.GetKeys(prefix)
	if err != nil { return nil, err }
	srx := fmt.Sprintf(`^%s/(\d)`, prefix)
	rx := regexp.MustCompile(srx)
	matchMap := map[int]struct{}{}
	for _, key := range childKeys {
		if rx.MatchString(key) {
			matches := rx.FindStringSubmatch(key)
			i, err := strconv.Atoi(matches[1])
			if err != nil { return nil, err }
			matchMap[i] = struct{}{}
		}
	}

	var mapKeys []int = slices.Collect(maps.Keys(matchMap))

	return mapKeys, nil
}


func testEncodeBool(t *testing.T, key string, b any) {
	valExpected, err := json.Marshal(b)
	require.NoError(t, err)
	ops, err := Encode(newEcoContext(os.Stdout), key, b)
	require.NoError(t, err)
	dumpOps(t, ops)
	require.NotEmpty(t, ops)
	require.Equal(t, 1, len(ops))
	op := ops[0]
	require.NotEmpty(t, op)
	require.True(t, op.IsPut())
	assert.Equal(t, []byte(key), op.KeyBytes())
	assert.Equal(t, valExpected, op.ValueBytes())
}

func testEncodeNumber(t *testing.T, key string, n any) {
	valExpected, err := json.Marshal(n)
	require.NoError(t, err)
	ops, err := Encode(newEcoContext(os.Stdout), key, n)
	require.NoError(t, err)
	dumpOps(t, ops)
	require.NotEmpty(t, ops)
	require.Equal(t, 1, len(ops))
	op := ops[0]
	require.NotEmpty(t, op)
	require.True(t, op.IsPut())
	assert.Equal(t, []byte(key), op.KeyBytes())
	assert.Equal(t, valExpected, op.ValueBytes())
}

func testEncodeString(t *testing.T, key string, s string) {
	ops, err := Encode(newEcoContext(os.Stdout), key, s)
	dumpOps(t, ops)
	// valExpected := fmt.Sprintf(`"%s"`, s)
	require.NoError(t, err)
	require.NotEmpty(t, ops)
	require.Equal(t, 1, len(ops))
	// op := ops[0]
	// require.NotEmpty(t, op)
	// require.True(t, op.IsPut())
	// assert.Equal(t, []byte(key), op.KeyBytes())
	// assert.Equal(t, []byte(valExpected), op.ValueBytes())
}

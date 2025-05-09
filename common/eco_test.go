package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
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
var VAL_MapWithStructValue = map_string_struct{}

type EcoTest struct {
	Name        string  `eco:"name"`
	LuckyNumber float64 `eco:"lucky_number"`
	Anon        string
}
type pEcoTest *EcoTest

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

func TestGetObject(t *testing.T) {

}

func TestKind_ArrayOfInt(t *testing.T) {
	kind := getKind(newEcoContext(), arrayOfInt{})
	assert.Equal(t, SimpleArray, kind)
}

func TestKind_ArrayOfStruct(t *testing.T) {
	kind := getKind(newEcoContext(), arrayOfStruct{})
	assert.Equal(t, Invalid, kind)
}

func TestKind_MapOfIntSlice(t *testing.T) {
	kind := getKind(newEcoContext(), mapOfIntSlice{})
	assert.Equal(t, SimpleMap, kind)

}

func TestKind_MapOfSliceOfStruct(t *testing.T) {
	type emptyStruct struct{}
	type mapOfSlice map[string][]emptyStruct
	kind := getKind(newEcoContext(), mapOfSlice{})
	assert.Equal(t, Invalid, kind)
}

func TestKind_MapSimple(t *testing.T) {
	kind := getKind(newEcoContext(), VAL_MapSimple)
	assert.Equal(t, SimpleMap, kind)
}

func TestKind_MapUnsimple(t *testing.T) {
	kind := getKind(newEcoContext(), VAL_MapUnsimple)
	assert.Equal(t, Invalid, kind)
}

func TestKind_MapUnsimpleKey(t *testing.T) {
	kind := getKind(newEcoContext(), VAL_MapUnsimpleKey)
	assert.Equal(t, Invalid, kind)
}

func TestKind_MapUnsimpleValue(t *testing.T) {
	kind := getKind(newEcoContext(), VAL_MapUnsimpleValue)
	assert.Equal(t, SimpleMap, kind)
}

func TestKind_PointerToInt(t *testing.T) {
	var pint pointerToInt
	kind := getKind(newEcoContext(), pint)
	assert.Equal(t, SimplePointer, kind)
}

func TestKind_PointerToIntSlice(t *testing.T) {
	kind := getKind(newEcoContext(), new(sliceOfInt))
	assert.Equal(t, SimplePointer, kind)
}

func TestKind_PointerToStructSlice(t *testing.T) {
	kind := getKind(newEcoContext(), new(sliceOfStruct))
	assert.Equal(t, Invalid, kind)
}

func TestKind_SliceSimple(t *testing.T) {
	type intSlice []int
	i := intSlice{1, 2, 3}
	kind := getKind(newEcoContext(), i)
	assert.Equal(t, SimpleSlice, kind)
}

func TestKind_SliceOfMap(t *testing.T) {
	type simpleMap map[int]int
	type sliceOfMap []simpleMap
	i := sliceOfMap{}
	kind := getKind(newEcoContext(), i)
	assert.Equal(t, Invalid, kind)
}

func TestKind_SliceUnsimple(t *testing.T) {
	type emptyStruct struct{}
	type emptyStructSlice []emptyStruct
	i := emptyStructSlice{}
	kind := getKind(newEcoContext(), i)
	assert.Equal(t, Invalid, kind)
}

func TestKind_StructSimple(t *testing.T) {
	kind := getKind(newEcoContext(), EcoTest{})
	assert.Equal(t, SimpleStruct, kind)
}

func TestKind_StructUnsimple(t *testing.T) {
	kind := getKind(newEcoContext(), UnsimpleStruct{})
	assert.Equal(t, Invalid, kind)
}

func TestKind_StructWithMap(t *testing.T) {
	kind := getKind(newEcoContext(), structWithMap{})
	assert.Equal(t, SimpleStruct, kind)
}

func TestKind_StructWithSlice(t *testing.T) {
	type structWithSlice struct{ Slice []int }
	kind := getKind(newEcoContext(), structWithSlice{})
	assert.Equal(t, SimpleStruct, kind)
}

func TestKind_StructWithMapWithSlice(t *testing.T) {
	type mapWithSlice map[string][]int
	type structWithMapWithSlice struct{ M mapWithSlice }
	kind := getKind(newEcoContext(), structWithMapWithSlice{})
	assert.Equal(t, SimpleStruct, kind)
}

func TestKind_StructWithMapWithStruct(t *testing.T) {
	type innerStruct struct{}
	type mapWithStruct map[string]innerStruct
	type structWithMapWithSlice struct{ MapField mapWithStruct }
	kind := getKind(newEcoContext(), structWithMapWithSlice{})
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
			t.Log(s)
		} else if op.IsPut() {
			key := string(op.KeyBytes())
			val := string(op.ValueBytes())
			s := fmt.Sprintf("%s %s %s", "PUT", key, val)
			t.Log(s)
		} else {
			t.Logf("%#v", op)
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

func testEncodeBool(t *testing.T, key string, b any) {
	valExpected, err := json.Marshal(b)
	require.NoError(t, err)
	ops, err := Encode(newEcoContext(), key, b)
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
	ops, err := Encode(newEcoContext(), key, n)
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
	ops, err := Encode(newEcoContext(), key, s)
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

package eco

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
)

func TestEncodeAstros(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	ops, err := Encode(ctx, "/test/astros", VAL_Astros)
	require.NoError(t, err)
	fmt.Println()
	dumpOps(t, ops)
}

func TestEncode_Bool(t *testing.T) {
	ctx, _ := initAndTest(t)
	key := "/test/bool"
	i := true
	testEncodeScalar(t, ctx, key, i)
}

func TestEncode_BoolSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	slice := []bool{true, false, true, true}
	key := "/test/boolSlice"
	ops := createSliceOps(t, ctx, slice, key)

	// test that each Op vaalue is as expected
	testSliceOps(t, slice, key, ops)
}

func TestEncode_EcoTest(t *testing.T) {
	key := "key"
	i := EcoTest{Name: "MEAT", LuckyNumber: 13}
	ctx := newEcoContext(os.Stdout)

	ops, err := Encode(ctx, key, i)
	assert.NoError(t, err)
	dumpOps(t, ops)
}

func TestEncode_Float(t *testing.T) {
	ctx, _ := initAndTest(t)
	key := "key"
	val := 42.0
	testEncodeScalar(t, ctx, key, val)
}

func TestEncode_FloatSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	slice := []float64{13.0, 169.13, -42.0}
	key := "/test/floatSlice"
	ops := createSliceOps(t, ctx, slice, key)

	testSliceOps(t, slice, key, ops)
}


func TestEncode_Int(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	key := "key"
	i := 13
	testEncodeScalar(t, ctx, key, i)
}

func TestEncode_Interface(t *testing.T) {
	type inf interface{}
	var infy = new(inf)
	_, err := Encode(newEcoContext(os.Stdout), "key", infy)
	assert.Error(t, err)
}

func TestEncode_IntSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	slice := []int{5, 8, 13}
	key := "/test/intSlice"
	ops := createSliceOps(t, ctx, slice, key)

	// encode slice + test # of Ops
	testSliceOps(t, slice, key, ops)
}

func TestEncode_MapOfMaps(t *testing.T) {
	key := "stros"
	map0 := map[string]string{"Name": "Altuve", "Position": "LF"}
	map1 := map[string]string{"Name": "Pena", "Position": "SS"}
	map2 := map[string]string{"Name": "Javier", "Position": "P"}
	mapStros := map[int]map[string]string{27: map0, 3: map1, 53: map2}
	ops, err := Encode(newEcoContext(os.Stdout), key, mapStros)
	assert.NoError(t, err)
	dumpOps(t, ops)
}

func TestEncode_MapWithIntKeys(t *testing.T) {
	key := "key"
	i := map[int]string{10: "print 'daylight is great'", 20: "print 'say it again'", 30: "goto 10"}
	ctx := newEcoContext(os.Stdout)

	ops, err := Encode(ctx, key, i)
	assert.NoError(t, err)
	dumpOps(t, ops)
}
func TestEncode_SimpleMap(t *testing.T) {
	key := "key"
	i := map[string]string{"foo": "13", "bar": "thirteen", "bum": "th1rt33n"}
	ctx := newEcoContext(os.Stdout)

	ops, err := Encode(ctx, key, i)
	assert.NoError(t, err)
	dumpOps(t, ops)
}

func TestEncode_Map_String_Struct(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	ops, err := Encode(ctx, "/test/map_string_struct", VAL_Map_String_Struct)
	require.NoError(t, err)
	fmt.Println()
	dumpOps(t, ops)
}

func TestEncode_String(t *testing.T) {
	ctx, _ := initAndTest(t)
	key := "key"
	i := "foo"
	testEncodeScalar(t, ctx, key, i)
}

func TestEncode_StringSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	slice := []string{"foo", "bar", "bum"}
	key := "/test/stringSlice"
	ops := createSliceOps(t, ctx, slice, key)

	// encode slice + test # of Ops
	testSliceOps(t, slice, key, ops)
}

func TestEncoding0(t *testing.T) {
	var s = `"8 is < g but > 13"`
	var buf []byte
	var err error
	buf, err = json.Marshal(s)
	assert.NoError(t, err)
	assert.NotNil(t, buf)
	t.Logf("%-20s %s", "Marshalled s", string(buf))
	bb := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(bb)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(s)
	assert.NoError(t, err)
	t.Logf("%-20s %s", "Encoded s", bb.String())
}


func TestPut_Bool(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := "/test/bool"
	val := true
	testPutScalar(t, ctx, cli, key, val)
}


func TestPut_BoolSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	slice := []bool{true, false, true, true}
	key := "/test/boolSlice"
	ops := createSliceOps(t, ctx, slice, key)
	testSliceOps(t, slice, key, ops)

	cli, err := CreateEtcdClientFromConfig()
	require.NoError(t, err)
	putSlice(t, ctx, cli, ops)
	
	testSliceValuesInEtcd(t, slice, cli, key)
}

func TestPut_Float(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := "/test/float"
	val := 42.0
	testPutScalar(t, ctx, cli, key, val)
}


func TestPut_FloatSlice(t *testing.T) {
	ctx, cli := initAndTest(t)
	slice := []float64{13.0, 169.13, -42.0}
	key := "/test/floatSlice"
	ops := createSliceOps(t, ctx, slice, key)
	testSliceOps(t, slice, key, ops)

	putSlice(t, ctx, cli, ops)
	testSliceValuesInEtcd(t, slice, cli, key)
}

func TestPut_Int(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := "/test/int"
	val := 13
	testPutScalar(t, ctx, cli, key, val)
}


func TestPut_IntSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	slice := []int{5, 8, 13}
	key := "/test/intSlice"
	ops := createSliceOps(t, ctx, slice, key)

	cli, err := CreateEtcdClientFromConfig()
	require.NoError(t, err)
	putSlice(t, ctx, cli, ops)

	testSliceValuesInEtcd(t, slice, cli, key)
}

func TestPut_String(t *testing.T) {
	ctx, cli := initAndTest(t)
	key := "/test/string"
	val := "hello"
	testPutString(t, ctx, cli, key, val)
}

func TestPut_StringSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	slice := []string{"foo", "bar", "bum"}
	key := "/test/stringSlice"
	ops := createSliceOps(t, ctx, slice, key)

	cli, err := CreateEtcdClientFromConfig()
	require.NoError(t, err)
	putSlice(t, ctx, cli, ops)

	testSliceValuesInEtcd(t, slice, cli, key)
}

func createSliceOps[U any](t *testing.T, ctx *ecoContext, slice []U, key string) []etcd.Op {
	ops, err := encodeSlice(ctx, key, reflect.ValueOf(slice))
	require.NoError(t, err)
	require.NotNil(t, ops)
	require.Equal(t, len(slice), len(ops))
	return ops
}


func putSlice(t *testing.T, ctx *ecoContext, cli *EtcdClient, ops []etcd.Op) {
	resp, err := cli.Txn(ctx).Then(ops...).Commit()
	require.NoError(t, err)
	ctx.logger.Infof("%#v", resp)
}


func dumpAndTestEncodeOps(t *testing.T, key string, val any) {
	ctx := newEcoContext(os.Stdout)
	ops, err := Encode(ctx, key, val)
	require.NoError(t, err)
	fmt.Println()
	dumpOps(t, ops)
}


func testSliceValuesInEtcd[U any](t *testing.T, expectedVals []U, cli *EtcdClient, key string) {
// Check each individual slice element in etcd
	for i, val := range expectedVals {
		subKey := fmt.Sprintf("%s/%d", key, i)
		buf, err := cli.Get(subKey)
		require.NoError(t, err)
		require.NotNil(t, buf)
		var etcdVal U
		err = json.Unmarshal(buf, &etcdVal)
		require.NoError(t, err)
		t.Logf("Testing %s=%#v ...", subKey, val)
		require.Equal(t, val, etcdVal)
	}
}

func testSliceOps[U any](t *testing.T, expectedVals []U, key string, ops []etcd.Op) {
	var err error
	// test that each Op vaalue is as expected
	for i, op := range ops {
		require.True(t, op.IsPut())
		elKey := fmt.Sprintf("%s/%d", key, i)
		assert.Equal(t, elKey, string(op.KeyBytes()))
		valExpected := expectedVals[i]
		var val U
		t.Logf("Checking if Op value for %s matches %#v", elKey, valExpected)
		err = json.Unmarshal(op.ValueBytes(), &val)
		require.NoError(t, err)
		assert.Equal(t, expectedVals[i], valExpected)
	}
}
	
package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
	"github.com/dylt-dev/dylt/color"
)

func decode (ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.inc()
	defer ctx.dec()

	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer { return fmt.Errorf("expected pointer; got %s", fullTypeName(ty))}
	if isSimple(ty.Elem().Kind()) {
		resp, err := etcdClient.Client.Get(ctx, key)
		if err != nil { return err }
		if len(resp.Kvs) != 1 {
			return fmt.Errorf("expected one key; got %d", len(resp.Kvs))
		}

		getVal := resp.Kvs[0].Value
		ctx.printf("getVal()=%v (%s)\n", getVal, getVal)
		err = json.Unmarshal(getVal, i)
		if err != nil { return err }
	}

	kindElem := getTypeKind(ctx, ty.Elem())

	if kindElem == SimpleMap {
		return decodeMap(ctx, etcdClient, key, i)

	} else if kindElem == SimpleSlice {
		return decodeSlice(ctx, etcdClient, key, i)
	
	} else if kindElem == SimpleStruct {
		return decodeStruct(ctx, etcdClient, key, i)
	
	} else {
		return errors.New("unsupported type")
	}
}

func decodeMap (ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.printf("decodeMap() - key=%s\n", key)
	ctx.inc(); defer ctx.dec()

	ty := reflect.TypeOf(i)
	var ss  = color.Styledstring(fullTypeName(ty)).FgBg(color.Ansi256.Color194, color.Ansi256.Color201)
	ctx.printf("ty=%s\n", ss)
	if ty.Kind() != reflect.Pointer { return fmt.Errorf("unsupported type (%s)", fullTypeName(ty))}
	kind := getTypeKind(ctx, ty.Elem())
	if kind != SimpleMap { return fmt.Errorf("unsupported type (%s)", fullTypeName(ty.Elem())) }

	if !strings.HasSuffix(key, string(filepath.Separator)) {
		key += string(filepath.Separator)
	}

	_, err := etcdClient.Client.Get(ctx, key, etcd.WithPrefix())
	if err != nil { return err }

	return nil
}

func decodeSlice (ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.inc(); defer ctx.dec()
	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer { return fmt.Errorf("unsupported type (%s)", fullTypeName(ty))}
	kind := getTypeKind(ctx, ty.Elem())
	if kind != SimpleSlice { return fmt.Errorf("unsupported type (%s)", fullTypeName(ty.Elem())) }

	resp, err := etcdClient.Client.Get(ctx, key)
	if err != nil { return err }
	if len(resp.Kvs) != 1 {
		return fmt.Errorf("expected one key; got %d", len(resp.Kvs))
	}

	getVal := resp.Kvs[0].Value
	ctx.printf("getVal()=%v (%s)\n", getVal, getVal)
	err = json.Unmarshal(getVal, i)
	if err != nil { return err }

	return nil
}

func decodeStruct (ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.inc(); defer ctx.dec()
	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer { return fmt.Errorf("unsupported type (%s)", fullTypeName(ty))}
	kind := getTypeKind(ctx, ty.Elem())
	if kind != SimpleStruct { return fmt.Errorf("unsupported type (%s)", fullTypeName(ty.Elem())) }

	if !strings.HasSuffix(key, string(filepath.Separator)) {
		key += string(filepath.Separator)
	}

	_, err := etcdClient.Client.Get(ctx, key, etcd.WithPrefix())
	if err != nil { return err }

	return nil
}


func TestMisc (t *testing.T) {
	etcdClient, err := NewEtcdClientFromConfig()
	ctx := newEcoContext()

	key1 := "/test/f"
	key2 := "/test/f"
	opGet1 := etcd.OpGet(key1)
	opGet2 := etcd.OpGet(key2, etcd.WithPrefix())
	require.NoError(t, err)
	txn := etcdClient.Txn(ctx)
	resp, err := txn.Then(opGet1, opGet2).Commit()
	assert.NoError(t, err)
	for _, resp2 := range resp.Responses {
		t.Logf("%d", resp2.GetResponseRange().Count)
	}
}

func TestBool (t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/flag"
	val := bool(false)
	putAndTest(t, etcdClient, key, val)

	var decodedVal bool
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	assert.NoError(t, err)
	assert.Equal(t, val, decodedVal)
	t.Log(decodedVal)
}

func TestBoolSlice (t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/boolslice"
	val := []bool{ true, true, false }
	putAndTest(t, etcdClient, key, val)

	type boolslice []bool
	var decodedVal boolslice
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, boolslice(val), decodedVal)
	t.Log(decodedVal)
}

func TestFloat (t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/f"
	val := float32(42.0)
	putAndTest(t, etcdClient, key, val)

	var decodedVal float32
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	assert.NoError(t, err)
	assert.Equal(t, val, decodedVal)
	t.Log(decodedVal)
}

func TestFloatSlice (t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/float32slice"
	val := []float32{ 42.0, 1764.0, 6.54321 }
	putAndTest(t, etcdClient, key, val)

	type float32slice []float32
	var decodedVal float32slice
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, float32slice(val), decodedVal)
	t.Log(decodedVal)
}

func TestInt (t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/n"
	val := int(-13)
	putAndTest(t, etcdClient, key, val)

	var decodedVal int
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	assert.NoError(t, err)
	assert.Equal(t, val, decodedVal)
	t.Log(decodedVal)
}

func TestIntSlice (t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/intslice"
	val := []int{5, 8, 13}
	putAndTest(t, etcdClient, key, val)

	type intslice []int
	var decodedVal intslice
	// var decodedVal []int
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, intslice(val), decodedVal)
	t.Log(decodedVal)
}

func TestMapStringString (t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/map/stringstring"
	val1 := "meat"
	val2 := "Meat"
	val3 := "MEEEEAT"
	val := map[string]string {
		"foo": val1,
		"bar": val2,
		"bum": val3,
	}
	putAndTest(t, etcdClient, filepath.Join(key, "foo"), val1)
	putAndTest(t, etcdClient, filepath.Join(key, "bar"), val2)
	putAndTest(t, etcdClient, filepath.Join(key, "baz"), val3)

	type mapstringstring map[string]string
	var decodedVal mapstringstring
	type pmapstringstring *mapstringstring
	var i pmapstringstring = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, (val), decodedVal)
	t.Log(decodedVal)
}

func TestString (t *testing.T) {
	ctx := newEcoContext()
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	key := "/test/s"
	val := `This\nis\a\<difficult>\nstring\n\to\n\e"s'c"a'p"e\n`
	putAndTest(t, etcdClient, key, val)

	var decodedVal string
	var i = &decodedVal
	err = decode(ctx, etcdClient, key, i)
	assert.NoError(t, err)
	assert.Equal(t, val, decodedVal)
	t.Log(decodedVal)
}

func TestStringSlice (t *testing.T) {
	ctx := newEcoContext()
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	key := "/test/stringslice"
	val := []string{"foo", "bar", "bum"}
	putAndTest(t, etcdClient, key, val)

	type stringslice []string
	var decodedVal stringslice
	var i = &decodedVal
	err = decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, stringslice(val), decodedVal)
	t.Log(decodedVal)
}

func TestStructEcoTest (t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/struct/ecotest"
	name := "Me"
	luckyNumber := 13
	val := NewEcoTest(name, float64(luckyNumber))
	putAndTest(t, etcdClient, filepath.Join(key, "name"), val.Name)
	putAndTest(t, etcdClient, filepath.Join(key, "lucky_number"), val.LuckyNumber)
	putAndTest(t, etcdClient, filepath.Join(key, "Anon"), val.Anon)

	var decodedVal EcoTest
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, (val), decodedVal)
	t.Log(decodedVal)
}

func TestUint (t *testing.T) {
	ctx := newEcoContext()
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	key := "/test/n"
	val := uint(13)
	putAndTest(t, etcdClient, key, val)

	var decodedVal uint
	var i = &decodedVal
	err = decode(ctx, etcdClient, key, i)
	assert.NoError(t, err)
	assert.Equal(t, val, decodedVal)
	t.Log(decodedVal)
}

func TestUintSlice (t *testing.T) {
	ctx := newEcoContext()
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	key := "/test/uintslice"
	val := []uint{ 5, 12, 13 }
	putAndTest(t, etcdClient, key, val)

	type uintslice []uint
	var decodedVal uintslice
	var i = &decodedVal
	err = decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, uintslice(val), decodedVal)
	t.Log(decodedVal)
}

func initAndTest (t *testing.T) (*ecoContext, *EtcdClient) {
	ctx := newEcoContext()
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	return ctx, etcdClient
}


func putAndTest (t *testing.T, etcdClient *EtcdClient, key string, i any) {
	// resp, err := etcdClient.Put(context.Background(), key, val)
	j, err := json.Marshal(i)
	require.NoError(t, err)
	resp, err := etcdClient.Put(context.Background(), key, string(j))
	require.NoError(t, err)
	require.NotNil(t, resp)
	// t.Logf("%#v", resp)
}
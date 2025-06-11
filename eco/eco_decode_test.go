package eco

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
)

func decode(ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.logger.signature("decode", key, reflect.TypeOf(i).Elem())
	ctx.inc()
	defer ctx.dec()

	// decode() only works if it's passed a pointer. The stdlin json Decoder has
	// the same constraint.
	// @note I'm not sure why this is better than just returning the object. Stack v heap?
	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer {
		return fmt.Errorf("expected pointer; got %s", common.FullTypeName(ty))
	}

	// Simple objects are easy to deal with. Just use json.Unmarhsal()
	if isScalar(ty.Elem().Kind()) {
		// Get object from etcd + make sure there's only 1
		resp, err := etcdClient.Client.Get(ctx, key)
		if err != nil {
			return err
		}
		if len(resp.Kvs) != 1 {
			return fmt.Errorf("expected one key; got %d", len(resp.Kvs))
		}

		// Unmarshal the result
		getVal := resp.Kvs[0].Value
		ctx.logger.Infof("getVal()=%v (%s)", getVal, getVal)
		err = json.Unmarshal(getVal, i)
		if err != nil {
			return err
		}
		// @note - should we return here?
	}

	// Some non-simple type are supported. The rest of the function checks for them.
	// Note - we want the type of the underlying element, not the type of the pointer
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

// eco stores maps as a number of sub-KVs with a common prefix. Go requires all
// KV values in a map are the same type, but etcd has no way of enforcing this. To
// etcd they're all just KVs.
func decodeMap(ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.logger.signature("decodeMap", "-etcdClient-", key, reflect.TypeOf(i))
	ctx.inc()
	defer ctx.dec()

	ctx.logger.Infof("i=%v ValueOf(i)=%v Elem()=%v ValueOf(Elem())=%v", i, reflect.ValueOf(i), reflect.ValueOf(i).Elem(), reflect.ValueOf(reflect.ValueOf(i).Elem()))

	ty := reflect.TypeOf(i)
	// ctx.println(subtle(fmt.Sprintf("ty=%s", fullTypeName(ty))))
	// Only pointers are supported
	if ty.Kind() != reflect.Pointer {
		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty))
	}

	// Only simple maps are supported
	kind := getTypeKind(ctx, ty.Elem())
	if kind != SimpleMap {
		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty.Elem()))
	}

	// add trailing slash. a/key => a/key/
	if !strings.HasSuffix(key, string(filepath.Separator)) {
		key += string(filepath.Separator)
	}

	// Get entire object tree
	// @note this might be quite large. ideally pagination would avoid issues with huge maps
	resp, err := etcdClient.Client.Get(ctx, key, etcd.WithPrefix())
	ctx.logger.info(common.Highlight("Keys"))
	var valMap reflect.Value
	// The caller may have specified a nil map, or an existing map
	// If nil, create a new map. If not, use the existing map
	if reflect.ValueOf(i).Elem().IsNil() {
		ctx.logger.info("map is nil; initializing new map")
		valMap = reflect.MakeMap(ty.Elem())
		reflect.ValueOf(i).Elem().Set(valMap)
	} else {
		ctx.logger.info("pointer is not nil; using existing map")
		valMap = reflect.Indirect(reflect.ValueOf(i))
	}
	for _, kv := range resp.Kvs {
		// Print a nice log statement
		// @note this is a lot of clutter for logging, esp when the real code
		// is a simple json.Unmarshal()
		skey := strings.TrimPrefix(string(kv.Key), key)
		skeyQuoted := fmt.Sprintf("\"%s\"", skey)
		ctx.logger.Infof("%-16s %-16s", skeyQuoted, kv.Value)
		// (*i)[skey] = kv.Value
		// simple json.Unmarshal() of value
		// @note this only supports maps of scalars. it needs to support nested maps since those are allowed. I think.
		var sval string
		err = json.Unmarshal(kv.Value, &sval)
		if err != nil {
			return err
		}
		// set a map value, reflection-style
		valMap.SetMapIndex(reflect.ValueOf(skey), reflect.ValueOf(sval))
	}

	if err != nil {
		return err
	}

	ctx.logger.info(common.Highlight("returning nil"))
	return nil
}

func decodeSlice(ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.logger.signature("decodeSlice", key, reflect.TypeOf(i).Elem())
	ctx.inc()
	defer ctx.dec()

	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer {
		return fmt.Errorf("unsupported type (%s) -  must be pointer", common.FullTypeName(ty))
	}
	kind := getTypeKind(ctx, ty.Elem())
	if kind != SimpleSlice {
		return fmt.Errorf("unsupported kind (%s) - must be SimpleSlice", common.FullTypeName(ty.Elem()))
	}

	kindElem := getTypeKind(ctx, ty.Elem().Elem())
	ctx.logger.Infof("kindElem=%s", kindElem.String())
	// Get slice keys
	sliceKeys, err := getSliceKeys(ctx, etcdClient, key)
	if err != nil {
		return err
	}
	ctx.logger.Infof("sliceKeys=%v", sliceKeys)
	// slice := reflect.MakeSlice(ty.Elem(), len(sliceKeys), len(sliceKeys))

	// for _, sliceKey := range sliceKeys {
	// 	elKey := path.Join(key, elKey)
	// 	// I have an element type.
	// 	// How do I create a variable to hold that type, and then decode a byte string into it?
	// }

	// Dynamically allocate array
	// For each slice key
	//	get index thingee
	//	somehow do a decode to the reflect.Value, even though I don't know how to do that

	// ctx.logger.Appendf("Getting key %s ...", key)
	// resp, err := etcdClient.Client.Get(ctx, key)
	// if err != nil {
	// 	ctx.logger.Flush(slog.LevelError, err.Error())
	// 	return err
	// }
	// if len(resp.Kvs) != 1 {
	// 	ctx.logger.Flush(slog.LevelError, "error")
	// 	return fmt.Errorf("expected one key; got %d", len(resp.Kvs))
	// }
	// ctx.logger.Flush(slog.LevelInfo, "ok")

	// getVal := resp.Kvs[0].Value
	// ctx.logger.Infof("getVal()=%v (%s)", getVal, getVal)
	// err = json.Unmarshal(getVal, i)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func decodeStruct(ctx *ecoContext, etcdClient *EtcdClient, key string, i any) error {
	ctx.logger.signature("decodeStruct", "-etcdClient", key, reflect.TypeOf(i))
	ctx.inc()
	defer ctx.dec()

	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer {
		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty))
	}
	tyElem := ty.Elem()
	kind := getTypeKind(ctx, tyElem)
	if kind != SimpleStruct {
		return fmt.Errorf("unsupported type (%s)", common.FullTypeName(ty.Elem()))
	}
	nFields := tyElem.NumField()
	ctx.logger.Infof("%-16s %-16d", "nFields", nFields)
	for iField := range nFields {
		field := tyElem.Field(iField)
		ctx.logger.info(string(common.Lowlight(fmt.Sprintf("%-16d %-16s %-16s", iField, field.Name, field.Tag.Get("eco")))))
	}

	if !strings.HasSuffix(key, string(filepath.Separator)) {
		key += string(filepath.Separator)
	}
	resp, err := etcdClient.Client.Get(ctx, key, etcd.WithPrefix())
	if err != nil {
		return err
	}

	fieldNameMap, err := fieldNameMap(i)
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		skey := strings.TrimPrefix(string(kv.Key), key)
		skeyQuoted := fmt.Sprintf("\"%s\"", skey)
		var sval any
		err = json.Unmarshal(kv.Value, &sval)
		if err != nil {
			return err
		}
		field := fieldNameMap[skey]
		field.Set(reflect.ValueOf(sval))
		ctx.logger.Infof("%-16s %-16v", skeyQuoted, sval)
	}

	return nil
}

func TestBool(t *testing.T) {
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

func TestBoolSlice(t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/boolslice"
	val := []bool{true, true, false}
	putAndTest(t, etcdClient, key, val)

	type boolslice []bool
	var decodedVal boolslice
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, boolslice(val), decodedVal)
	t.Log(decodedVal)
}

func TestFloat(t *testing.T) {
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

func TestFieldNameMap(t *testing.T) {
	var ecoTest = EcoTest{}
	var p *EcoTest = &ecoTest

	fieldNameMap, err := fieldNameMap(p)
	require.NoError(t, err)
	require.NotNil(t, fieldNameMap)
	t.Logf("%#v", fieldNameMap)

	fieldNameMap["Anon"].Set(reflect.ValueOf("(...)"))
	fieldNameMap["name"].Set(reflect.ValueOf("Me"))
	fieldNameMap["lucky_number"].Set(reflect.ValueOf(13.0))
	t.Logf("%#v", p)
}

func TestFloatSlice(t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/float32slice"
	val := []float32{42.0, 1764.0, 6.54321}
	putAndTest(t, etcdClient, key, val)

	type float32slice []float32
	var decodedVal float32slice
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, float32slice(val), decodedVal)
	t.Log(decodedVal)
}

func TestInt(t *testing.T) {
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

func TestDecode_IntSlice(t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/intSlice"
	val := []int{5, 8, 13}

	type intslice []int
	var decodedVal intslice
	var i = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, intslice(val), decodedVal)
	t.Log(decodedVal)
}

func TestMapStringString(t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/map/stringstring"
	val1 := "meat"
	val2 := "Meat"
	val3 := "MEEEEAT"
	type mapstringstring map[string]string
	val := mapstringstring{
		"foo": val1,
		"bar": val2,
		"bum": val3,
	}
	for k, v := range val {
		putAndTest(t, etcdClient, filepath.Join(key, k), v)
	}

	var decodedVal mapstringstring = nil
	type pmapstringstring *mapstringstring
	var i pmapstringstring = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, (val), decodedVal)
	t.Log(decodedVal)
}

func TestMisc(t *testing.T) {
	etcdClient, err := NewEtcdClientFromConfig()
	ctx := newEcoContext(os.Stdout)

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

func TestNilMap(t *testing.T) {
	var m map[string]string = nil
	var pm *map[string]string = &m
	t.Logf("reflect.ValueOf(pm).IsNil()=%v", reflect.ValueOf(pm).IsNil())
	t.Logf("reflect.ValueOf(pm).Elem().IsNil()=%v", reflect.ValueOf(pm).Elem().IsNil())

	// Get type of underlying object from pointer
	tyElem := reflect.TypeOf(pm).Elem()
	t.Logf("tyElem=%s", common.FullTypeName(tyElem))

	// Create map and assign a value
	valMap := reflect.MakeMap(tyElem)
	reflect.ValueOf(pm).Elem().Set(valMap)
	valMap.SetMapIndex(reflect.ValueOf("foo"), reflect.ValueOf("13"))

	//
	t.Logf("%#v", *pm)
}

func TestNilMapPointer(t *testing.T) {
	var pm *map[string]string = nil
	pmValue := reflect.ValueOf(pm)
	t.Logf("pmValue.IsNil()=%v", pmValue.IsNil())
}

func TestNilPointer(t *testing.T) {
	var m map[string]string = nil
	var pm *map[string]string = &m
	t.Logf("reflect.ValueOf(pm).IsNil()=%v", reflect.ValueOf(pm).IsNil())
	t.Logf("reflect.ValueOf(pm).Elem().IsNil()=%v", reflect.ValueOf(pm).Elem().IsNil())
}

func TestString(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
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

func TestStringSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
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

func TestStructEcoTest(t *testing.T) {
	ctx, etcdClient := initAndTest(t)

	key := "/test/struct/ecotest"
	name := "Me"
	luckyNumber := 13
	val := NewEcoTest(name, float64(luckyNumber))
	putAndTest(t, etcdClient, filepath.Join(key, "name"), val.Name)
	putAndTest(t, etcdClient, filepath.Join(key, "lucky_number"), val.LuckyNumber)
	putAndTest(t, etcdClient, filepath.Join(key, "Anon"), val.Anon)

	var decodedVal EcoTest
	type pEcoTest *EcoTest
	var i pEcoTest = &decodedVal
	err := decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, *val, *i)
	t.Log(decodedVal)
}

func TestStructSetField0(t *testing.T) {
	var st = EcoTest{}
	type pEcoTest *EcoTest
	var pst pEcoTest = &st

	val := reflect.ValueOf(pst).Elem()
	val.Field(0).Set(reflect.ValueOf("me"))
	val.Field(1).Set(reflect.ValueOf(13.0))
	t.Logf("%#v", *pst)

}

func TestStructSetField1(t *testing.T) {
	var st = EcoTest{}
	p := reflect.ValueOf(&st)
	val := p.Elem()

	t.Logf("fullTypeName(val.Type())=%s", common.FullTypeName(val.Type()))
	t.Logf("val.Type().Kind()=%s", val.Type().Kind().String())

	require.True(t, val.CanSet())
	require.True(t, val.Field(0).CanSet())

	val.Field(0).SetString("me")
	val.Field(1).Set(reflect.ValueOf(13.0))
	t.Logf("%#v", st)
}

func TestUint(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
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

func TestUintSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	key := "/test/uintslice"
	val := []uint{5, 12, 13}
	putAndTest(t, etcdClient, key, val)

	type uintslice []uint
	var decodedVal uintslice
	var i = &decodedVal
	err = decode(ctx, etcdClient, key, i)
	require.NoError(t, err)
	assert.Equal(t, uintslice(val), decodedVal)
	t.Log(decodedVal)
}

func initAndTest(t *testing.T) (*ecoContext, *EtcdClient) {
	ctx := newEcoContext(os.Stdout)
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	return ctx, etcdClient
}

func putAndTest(t *testing.T, etcdClient *EtcdClient, key string, i any) {
	// resp, err := etcdClient.Put(context.Background(), key, val)
	j, err := json.Marshal(i)
	require.NoError(t, err)
	resp, err := etcdClient.Put(context.Background(), key, string(j))
	require.NoError(t, err)
	require.NotNil(t, resp)
	// t.Logf("%#v", resp)
}

/*
ap[string]reflect.Value
{
	"Anon":reflect.Value{
		typ_:(*abi.Type)(0x1014f6460),
		ptr:(unsafe.Pointer)(0x140002f0b28),
		flag:0x98
	},
	"lucky_number":reflect.Value{
		typ_:(*abi.Type)(0x1014f6820),
		ptr:(unsafe.Pointer)(0x140002f0b20),
		flag:0x8e
	},
	"name":reflect.Value{
		typ_:(*abi.Type)(0x1014f6460),
		ptr:(unsafe.Pointer)(0x140002f0b10),
		flag:0x98}
	}

	reflect.Value{
		"Anon":reflect.Value{
			typ_:(*abi.Type)(0x104afa460),
			ptr:(unsafe.Pointer)(0x140002f0b28),
			flag:0x198
		},
		"lucky_number":reflect.Value{
			typ_:(*abi.Type)(0x104afa820),
			ptr:(unsafe.Pointer)(0x140002f0b20),
			flag:0x18e
		},
		"name":reflect.Value{
			typ_:(*abi.Type)(0x104afa460),
			ptr:(unsafe.Pointer)(0x140002f0b10),
			flag:0x198}
		}
*/

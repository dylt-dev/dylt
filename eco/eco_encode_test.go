package eco

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeAstros (t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	ops, err := Encode(ctx, "/test/astros", VAL_Astros)
	require.NoError(t, err)
	fmt.Println()
	dumpOps(t, ops)
}

func TestEncode_Bool(t *testing.T) {
	key := "key"
	i := false
	testEncodeBool(t, key, i)
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
	key := "key"
	i := 42.0
	testEncodeNumber(t, key, i)
}
func TestEncode_Int(t *testing.T) {
	key := "key"
	i := 13
	testEncodeNumber(t, key, i)
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
	// j, err := json.Marshal(slice)
	// assert.NoError(t, err)
	// sJson := string(j)

	key := "/test/intSlice"
	ops, err := encodeSlice(ctx, key, reflect.ValueOf(slice))
	dumpAndTestEncodeOps(t, key, slice)
	require.NoError(t, err)
	for i, op := range ops {
		elKey := fmt.Sprintf("%s/%d", key, i)
		assert.Equal(t, elKey, string(op.KeyBytes()))
		var val int
		err = json.Unmarshal(op.ValueBytes(), &val)
		require.NoError(t, err)
		assert.Equal(t, slice[i], val)
	}
}

func TestPut_IntSlice (t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	slice := []int{5, 8, 13}

	key := "/test/intSlice"
	ops, err := encodeSlice(ctx, key, reflect.ValueOf(slice))
	require.NoError(t, err)
	cli, err := CreateEtcdClientFromConfig()
	require.NoError(t, err)
	resp, err := cli.Txn(context.Background()).Then(ops...).Commit()
	require.NoError(t, err)
	t.Logf("%#v", resp)
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
	key := "key"
	i := "foo"
	testEncodeString(t, key, i)
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

func dumpAndTestEncodeOps (t *testing.T, key string, val any) {
	ctx := newEcoContext(os.Stdout)
	ops, err := Encode(ctx, key, val)
	require.NoError(t, err)
	fmt.Println()
	dumpOps(t, ops)
}
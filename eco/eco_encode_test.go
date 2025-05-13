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
)

func TestEncodeEtcd_Bool(t *testing.T) {
	key := "key"
	i := false
	testEncodeBool(t, key, i)
}

func TestEncodeEtcd_EcoTest(t *testing.T) {
	key := "key"
	i := EcoTest{Name: "MEAT", LuckyNumber: 13}
	ctx := newEcoContext(os.Stdout)

	ops, err := Encode(ctx, key, i)
	assert.NoError(t, err)
	dumpOps(t, ops)
}

func TestEncodeEtcd_Float(t *testing.T) {
	key := "key"
	i := 42.0
	testEncodeNumber(t, key, i)
}
func TestEncodeEtcd_Int(t *testing.T) {
	key := "key"
	i := 13
	testEncodeNumber(t, key, i)
}

func TestEncodeEtcd_Interface(t *testing.T) {
	type inf interface{}
	var infy = new(inf)
	_, err := Encode(newEcoContext(os.Stdout), "key", infy)
	assert.Error(t, err)
}

func TestEncodeEtcd_IntSlice(t *testing.T) {
	ctx := newEcoContext(os.Stdout)
	slice := []int{5, 8, 13}
	j, err := json.Marshal(slice)
	assert.NoError(t, err)
	sJson := string(j)

	key := "key"
	ops, err := encodeSlice(ctx, key, reflect.ValueOf(slice))
	dumpOps(t, ops)
	assert.NoError(t, err)
	assert.Equal(t, key, string(ops[0].KeyBytes()))
	assert.Equal(t, sJson, string(ops[0].ValueBytes()))
}

func TestEncodeEtcd_MapOfMaps(t *testing.T) {
	key := "stros"
	map0 := map[string]string{"Name": "Altuve", "Position": "LF"}
	map1 := map[string]string{"Name": "Pena", "Position": "SS"}
	map2 := map[string]string{"Name": "Javier", "Position": "P"}
	mapStros := map[int]map[string]string{27: map0, 3: map1, 53: map2}
	ops, err := Encode(newEcoContext(os.Stdout), key, mapStros)
	assert.NoError(t, err)
	dumpOps(t, ops)
}

func TestEncodeEtcd_MapWithIntKeys(t *testing.T) {
	key := "key"
	i := map[int]string{10: "print 'daylight is great'", 20: "print 'say it again'", 30: "goto 10"}
	ctx := newEcoContext(os.Stdout)

	ops, err := Encode(ctx, key, i)
	assert.NoError(t, err)
	dumpOps(t, ops)
}
func TestEncodeEtcd_SimpleMap(t *testing.T) {
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

func TestEncodeEtcd_String(t *testing.T) {
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

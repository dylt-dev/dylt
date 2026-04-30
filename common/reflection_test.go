package common

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateMapKey1(t *testing.T) {
	expectedData := "foo"
	var key string
	pkey := &key
	UnmarshalMapKey(expectedData, pkey)
	require.Equal(t, expectedData, key)
}

func TestCreateMapKey2(t *testing.T) {
	expectedData := 13
	var key int
	pkey := &key
	UnmarshalMapKey(fmt.Sprint(expectedData), pkey)
	require.Equal(t, expectedData, key)
}

func TestCreateMapKey3(t *testing.T) {
	expectedData := false
	var key bool
	pkey := &key
	UnmarshalMapKey(fmt.Sprint(expectedData), pkey)
	require.Equal(t, expectedData, key)
}

func TestCreateStructAndSetField(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	expectedName := "foo"
	var pst *TestStruct = nil
	ppst := &pst

	pNewSt, is := CreateOrGetStruct(ctx, reflect.ValueOf(ppst))
	require.True(t, is)
	require.NotNil(t, pNewSt)
	pEco, is := pNewSt.(*TestStruct)
	require.True(t, is)
	require.NotNil(t, pEco)

	err := SetStructField(pst, "Name", expectedName)
	require.NoError(t, err)
	require.Equal(t, expectedName, (*pst).Name)
}

func TestCreateStructAndSetInvalidField(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	var pst *TestStruct = nil
	ppst := &pst

	pNewSt, is := CreateOrGetStruct(ctx, reflect.ValueOf(ppst))
	require.True(t, is)
	require.NotNil(t, pNewSt)
	pEco, is := pNewSt.(*TestStruct)
	require.True(t, is)
	require.NotNil(t, pEco)

	err := SetStructField(pst, "INVALID_FIELD_NAME", nil)
	require.Error(t, err)
}

func TestCreateStructAndUnmarshalField(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	expectedName := "foo"
	expectedNameBuf, err := json.Marshal(expectedName)
	require.NoError(t, err)
	var pst *TestStruct = nil
	ppst := &pst

	pNewSt, is := CreateOrGetStruct(ctx, reflect.ValueOf(ppst))
	require.True(t, is)
	require.NotNil(t, pNewSt)
	pEco, is := pNewSt.(*TestStruct)
	require.True(t, is)
	require.NotNil(t, pEco)

	err = UnmarshalStructField(pEco, "Name", expectedNameBuf)
	require.NoError(t, err)
	require.Equal(t, expectedName, (*pst).Name)
}

func TestCreateStructAndUnmarshalInvalidField(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)
	var pst *TestStruct = nil
	ppst := &pst

	pNewSt, is := CreateOrGetStruct(ctx, reflect.ValueOf(ppst))
	require.True(t, is)
	require.NotNil(t, pNewSt)
	pEco, is := pNewSt.(*TestStruct)
	require.True(t, is)
	require.NotNil(t, pEco)

	err := UnmarshalStructField(pEco, "INVALID_FIELD_NAME", []byte{})
	require.Error(t, err)
}

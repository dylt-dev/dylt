package common

import (
	"errors"
	"io"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConfigFile_Init (t *testing.T) {
	path := "/tmp/config.yaml"
	var fi os.FileInfo
	var err error

	// If os.Stat() returns an error make sure it's because the file doesn't exist
	fi, err = os.Stat(path)
	if err != nil {
		assert.True(t, errors.Is(err, os.ErrNotExist))
		assert.Nil(t, fi)
	}

	// If the file exists, remove it
	if err == nil {
		err = os.Remove(path)
		assert.Nil(t, err)
	}

	// Initialize the file & make sure it exists
	cf := ConfigFile{Path: path}
	err = cf.Init()	
	assert.NoError(t, err)
	fi, err = os.Stat(path)
	assert.Nil(t, err)
	assert.NotNil(t, fi)
}


func Test_ConfigFile_InitIfExists (t *testing.T) {
	path := "/tmp/config.yaml"
	var err error

	// Create the file, & make sure it exists
	f, err := os.OpenFile(path, os.O_CREATE, 0644)
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)
	
	// Confirm calling Init() with an existing file still works
	cf := ConfigFile{Path: path}
	err = cf.Init()	
	assert.NoError(t, err)
}
func TestCreateConfigFile(t *testing.T) {
	err := CreateConfigFile()
	assert.Nil(t, err)
	cfgFilePath := GetConfigFilePath()
	f, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 0400)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	fi, err := f.Stat()
	assert.Nil(t, err)
	assert.NotNil(t, fi)
	assert.Equal(t, "config.yaml", fi.Name())
	assert.False(t, fi.IsDir())
}

func TestConfigMapGet0(t *testing.T) {
	// Simple key, simple map
	key := "foo"
	valExpected := "bar"
	data := ConfigMap{key: valExpected}
	val := data.Get(key)
	assert.NotNil(t, val)
	assert.Equal(t, valExpected, val)
}

func TestConfigMapGet1(t *testing.T) {
	// Composite key
	key := "a.b.c"
	valExpected := "foo"
	data2 := ConfigMap{"c": "foo"}
	data1 := ConfigMap{"b": data2}
	data := ConfigMap{"a": data1}
	val := data.Get(key)
	assert.NotNil(t, val)
	assert.Equal(t, valExpected, val)
}

func TestConfigMapSet0(t *testing.T) {
	// Simple key, empty map
	key := "foo"
	val := "bar"
	data := ConfigMap{}
	assert.Nil(t, data[key])
	data.Set(key, val)
	assert.Equal(t, val, data[key])
}

func TestConfigMapSet1(t *testing.T) {
	// Simple key, value already exists
	// Simple key, empty map
	key := "foo"
	val := "bar"
	newVal := "bum"
	data := ConfigMap{key: val}
	assert.Equal(t, val, data[key])
	data.Set(key, newVal)
	assert.Equal(t, newVal, data[key])
}

func TestConfigMapSet2(t *testing.T) {
	// Composite key, empty map
	key := "a.b.c"
	val := "bar"
	data := ConfigMap{}
	data.Set(key, val)
	data1 := data["a"]
	assert.NotNil(t, data1)
	assert.IsType(t, map[string]any{}, data1)
	data2 := (data1.(map[string]any))["b"]
	assert.NotNil(t, data2)
	assert.IsType(t, map[string]any{}, data2)
	data3 := (data2.(map[string]any))["c"]
	assert.NotNil(t, data3)
	assert.Equal(t, val, data3)
}

func TestConfigMapSet3(t *testing.T) {
	// Composite key, some intermediate keys exist
	key := "a.b.c"
	val := "bar"
	// var initData2 map[string]any
	initData2 := ConfigMap{"c": "MEEEEEEEEAAAAAAAAAAT"}
	initData1 := ConfigMap{"b": initData2}
	data := ConfigMap{"a": initData1}
	data.Set(key, val)
	data1 := data["a"]
	assert.NotNil(t, data1)
	assert.IsType(t, ConfigMap{}, data1)
	data2 := (data1.(ConfigMap))["b"]
	assert.NotNil(t, data2)
	assert.IsType(t, ConfigMap{}, data2)
	data3 := (data2.(ConfigMap))["c"]
	assert.NotNil(t, data3)
	assert.Equal(t, val, data3)
}

/*
func TestLoadConfig (t *testing.T) {
	homeFolder := os.Getenv("HOME")
	configHome := ".config/dylt"
	configFile := "dylt.yaml"
	configFolder := path.Join(homeFolder, configHome)
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(configFolder)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err.Error()))
	}
	field := "etcd-domain"
	etcdDomain := viper.GetString(field)
	assert.NotNil(t, etcdDomain)
	assert.Equal(t, "hello.dylt.dev", etcdDomain)
}
*/

func TestOpenConfigFile(t *testing.T) {
	f, err := OpenConfigFile()
	assert.Nil(t, err)
	assert.NotNil(t, f)
	buf, err := io.ReadAll(f)
	assert.Nil(t, err)
	assert.NotNil(t, buf)
	assert.Greater(t, len(buf), 0)
	t.Logf("%s", string(buf))
}

func TestGetConfigValue(t *testing.T) {
	expected := "hello.dylt.dev"
	key := "etcd-domain"
	val, err := GetConfigValue(key)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)
}

func TestGetConfigValueNoKey(t *testing.T) {
	key := "XXX"
	val, err := GetConfigValue(key)
	assert.NotNil(t, err)
	assert.Nil(t, val)
}

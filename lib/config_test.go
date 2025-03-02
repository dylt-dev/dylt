package lib

import (
	//	"fmt"

	"io"
	"os"

	//	"path"
	"testing"

	//	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func Test (t *testing.T) {
	t.Log("hello")
}

func TestCreateConfigFile (t *testing.T) {
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

func TestGetConfigValue (t *testing.T) {
	expected := "hello.dylt.dev"
	key := "etcd-domain"
	val, err := GetConfigValue(key)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)
}

func TestGetConfigValueNoKey (t *testing.T) {
	key := "XXX"
	val, err := GetConfigValue(key)
	assert.NotNil(t, err)
	assert.Nil(t, val)
}
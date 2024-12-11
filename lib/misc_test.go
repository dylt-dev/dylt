package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSimplePut (t *testing.T) {
	vmClient, err := CreateVmClientFromConfig()
	if err != nil { t.Fatal(err, debug.Stack()) }
	vm := VmInfo{
		Address: "hosty toasty host",
		Name: "ovh-vps0",
	}
	buf, _ := json.Marshal(vm)
	s := string(buf)
	t.Logf("s=%s\n", s)
	name := "ovh-vps0"
	key := fmt.Sprintf("/vm/%s", name)
	vmClient.Client.Put(context.Background(), key, s)
}

func TestLoadConfig (t *testing.T) {
	cfg := Config{}
	err := cfg.Load()
	assert.Nil(t, err)
	domain, _ := cfg.GetEtcDomain()
	assert.Empty(t, domain)
}

func TestLoadConfig2 (t *testing.T) {
	cfg := Config{}
	err := cfg.Load()
	assert.Nil(t, err)
	domain, _ := cfg.GetEtcDomain()
	assert.NotEmpty(t, domain)
	assert.Equal(t, "hello.dylt.dev", domain)
}

func TestSaveConfig (t *testing.T) {
	cfg := Config{}
	err := cfg.Load()
	assert.Nil(t, err)
	err = cfg.SetEtcDomain("hello.dylt.dev")
	assert.Nil(t, err)
	err = cfg.Save()
	assert.Nil(t, err)
}

func TestInitConfig (t *testing.T) {
	viper.SetConfigName(CFG_Filename)
	viper.SetConfigType(CFG_Type)
	viper.AddConfigPath(".")
	cfgFolder := GetConfigFolderPath()
	viper.AddConfigPath(cfgFolder)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err.Error()))
	}
}


func TestClearConfig (t *testing.T) {
	err := ClearConfigFile()
	assert.Nil(t, err)
	cfgFilePath := GetConfigFilePath()
	f, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 400)
	assert.Nil(t, err)
	defer f.Close()
	fi, err := f.Stat()
	assert.Nil(t, err)
	assert.NotNil(t, fi)
	assert.Equal(t, int64(0), fi.Size())
}


func TestInit (t *testing.T) {
	// Init the config
	const etcdDomain = "hello.dylt.dev"
	initInfo := InitInfo{
		EtcdDomain: etcdDomain,
	}
	err := Init(&initInfo)
	assert.Nil(t, err)
	// Test the file exists
	cfgFilePath := GetConfigFilePath()
	cfgFile, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 0400)
	assert.Nil(t, err)
	defer cfgFile.Close()
	// Read file as yaml
	decoder := yaml.NewDecoder(cfgFile)
	cfgStruct := ConfigStruct{}
	err = decoder.Decode(&cfgStruct)
	t.Logf("%#v", cfgStruct)
	assert.Nil(t, err)
	// Test the file contains the expected domain
	assert.Equal(t, etcdDomain, cfgStruct.EtcdDomain)
}

func TestShowConfig (t *testing.T) {
	err := ShowConfig(os.Stdout)
	assert.Nil(t, err)
}

func TestStrings (t *testing.T) {
	str0 := "> GET /etcd-io/etcd/releases/download/v3.5.16/etcd-v3.5.16-linux_amd64.tar.gz HTTP/2"
	str1 := "> GET /etcd-io/etcd/releases/download/v3.5.16/etcd-v3.5.16-linux-amd64.tar.gz HTTP/2"
	len0 := len(str0)
	len1 := len(str1)
	assert.Equal(t, len0, len1)
	for i := range len0 {
		t.Logf("Checking index #%d ...\n", i)
		assert.Equal(t, str0[i], str1[i])
	}
	assert.Equal(t, str0, str1)
}
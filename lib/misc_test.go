package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSimplePut (t *testing.T) {
	cli, err := NewVmClient(viper.GetString("etcd_domain"))
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
	cli.Client.Put(context.Background(), key, s)
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
	_ = cfg.SetEtcDomain("hello.dylt.dev")
	err = cfg.Viper.WriteConfig()
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

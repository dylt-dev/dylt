package lib

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test (t *testing.T) {
	t.Log("hello")
}

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
	field := "etcd_domain"
	etcdDomain := viper.GetString(field)
	assert.NotNil(t, etcdDomain)
	assert.Equal(t, "hello.dylt.dev", etcdDomain)
}


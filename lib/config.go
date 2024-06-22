package lib

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {

}


func (o *Config) Get (name string) (interface{}, error) {
	return nil, nil
}


func (o *Config) Set (name string, value interface{}) error {
	return nil
}


func (o *Config) Load () error {
	return nil
}


func (o *Config) Save () error {
	return nil
}


func InitConfig () {
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
}


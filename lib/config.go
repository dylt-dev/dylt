package lib

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config = struct {

}


func (o *Config) Get (name string) (interface{}, error) {

}


func (o *Config) Set (name string, value interface{}) error {

}


func (o *Config) Load () error {

}


func (o *Config) Save () error {

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


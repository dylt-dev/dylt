package lib

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const CFG_Filename = "config.yaml"
const CFG_Folder = ".config/dylt"
const CFG_Type = "yaml"


type Config struct {
	Viper *viper.Viper
}


func (o *Config) Get (name string) (interface{}, error) {
	return nil, nil
}

func (o *Config) GetEtcDomain () (string, error) {
	key := "etcd_domain"
	isSet := o.Viper.IsSet(key)
	if !isSet { return "", nil }
	domain := o.Viper.GetString(key)
	return domain, nil
}

func (o *Config) SetEtcDomain (domain string) (error) {
	o.Viper.Set("etcd_domain", domain)
	return nil
}


func (o *Config) Set (name string, value interface{}) error {
	return nil
}


func (o *Config) Load () error {
	o.Viper = viper.New()
	o.Viper.SetConfigName(CFG_Filename)
	o.Viper.SetConfigType(CFG_Type)
	o.Viper.AddConfigPath(".")
	cfgFolder := GetConfigFolderPath()
	o.Viper.AddConfigPath(cfgFolder)
	err := o.Viper.ReadInConfig()
	return err
}


func (o *Config) Save () error {
	return nil
}


func CreateConfigFile () error {
	err := CreateConfigFolder()
	if err != nil { return err }
	cfgFilePath := GetConfigFilePath()
	f, err := os.OpenFile(cfgFilePath, os.O_RDONLY|os.O_CREATE, 0644)
	defer f.Close()
	return err
}

func CreateConfigFolder () error {
	cfgFolder := GetConfigFolderPath()
	err := os.MkdirAll(cfgFolder, 0744)
	return err
}

func GetConfigFilePath () string {
	cfgFolder := GetConfigFolderPath()
	cfgFilePath := path.Join(cfgFolder, CFG_Filename)
	return cfgFilePath
}

func GetConfigFolderPath () string {
	homeFolder := os.Getenv("HOME")
	cfgFolder := path.Join(homeFolder, CFG_Folder)
	return cfgFolder
}


func InitConfig () {
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


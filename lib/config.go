package lib

import (
	"io"
	"os"
	"path"
	"gopkg.in/yaml.v3"
)

// Hardcoding this info here is a little sketchy, but it's a convenient location
// for the info until a final scheme is devised
const CFG_Filename = "config.yaml"
const CFG_Folder = ".config/dylt"
const CFG_Type = "yaml"

// type Config struct {
// 	Viper      *viper.Viper
// 	EtcdDomain string `yaml:"etcd_domain"`
// }

type ConfigStruct struct {
	EtcdDomain string `yaml:"etcd_domain"`
}

// func (o *Config) Get (name string) (interface{}, error) {
// 	return nil, nil
// }

// func (o *Config) GetEtcDomain () (string, error) {
// 	key := "etcd_domain"
// 	isSet := o.Viper.IsSet(key)
// 	if !isSet { return "", nil }
// 	domain := o.Viper.GetString(key)
// 	return domain, nil
// }

// func (o *Config) SetEtcDomain (domain string) (error) {
// 	o.Viper.Set("etcd_domain", domain)
// 	return nil
// }

// func (o *Config) SetEtcDomainAndSave (domain string) (error) {
// 	o.Viper.Set("etcd_domain", domain)
// 	err := o.Save()
// 	return err
// }

// func (o *Config) Set (name string, value interface{}) error {
// 	return nil
// }

// func (o *Config) Load () error {
// 	o.Viper = viper.New()
// 	o.Viper.SetConfigName(CFG_Filename)
// 	o.Viper.SetConfigType(CFG_Type)
// 	o.Viper.AddConfigPath(".")
// 	cfgFolder := GetConfigFolderPath()
// 	o.Viper.AddConfigPath(cfgFolder)
// 	err := o.Viper.ReadInConfig()
// 	return err
// }

// func (o *Config) Save () error {
// 	err := o.Viper.WriteConfig()
// 	return err
// }

// func ClearConfigFile () error {
// 	cfgFilePath := GetConfigFilePath()
// 	err := os.Truncate(cfgFilePath, 0)
// 	return err
// }

func CreateConfigFile () error {
	err := CreateConfigFolder()
	if err != nil { return err }
	cfgFilePath := GetConfigFilePath()
	f, err := os.OpenFile(cfgFilePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil { return err }
	defer f.Close()
	return err
}

func CreateConfigFolder() error {
	cfgFolder := GetConfigFolderPath()
	err := os.MkdirAll(cfgFolder, 0744)
	return err
}

func GetConfigFilePath() string {
	cfgFolder := GetConfigFolderPath()
	cfgFilePath := path.Join(cfgFolder, CFG_Filename)
	return cfgFilePath
}

func GetConfigFolderPath() string {
	homeFolder := os.Getenv("HOME")
	cfgFolder := path.Join(homeFolder, CFG_Folder)
	return cfgFolder
}

// func InitConfig() {
// 	viper.SetConfigName(CFG_Filename)
// 	viper.SetConfigType(CFG_Type)
// 	viper.AddConfigPath(".")
// 	cfgFolder := GetConfigFolderPath()
// 	viper.AddConfigPath(cfgFolder)
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		panic(fmt.Errorf("fatal error config file: %s", err.Error()))
// 	}
// }

func LoadConfig() (ConfigStruct, error) {
	cfg := ConfigStruct{}
	cfgFile, err := OpenConfigFile()
	if err != nil { return cfg, err }
	decoder := yaml.NewDecoder(cfgFile)
	err = decoder.Decode(&cfg)
	return cfg, err
}

func OpenConfigFile () (*os.File, error) {
	cfgFilePath := GetConfigFilePath()
	f, err := os.Open(cfgFilePath)
	return f, err
}

func SaveConfig (cfg ConfigStruct) error {
	cfgFilePath := GetConfigFilePath()
	f, err := os.Create(cfgFilePath)
	if err != nil { return err }
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(cfg)
	if err != nil { return err }
	return nil
}

func ShowConfig(out io.Writer) error {
	cfgFilePath := GetConfigFilePath()
	cfgFile, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 0400)
	if err != nil {
		return err
	}
	defer cfgFile.Close()
	// Read file as yaml
	decoder := yaml.NewDecoder(cfgFile)
	cfgStruct := ConfigStruct{}
	err = decoder.Decode(&cfgStruct)
	if err != nil {
		return err
	}
	// Marshal indented output
	buf, err := yaml.Marshal(cfgStruct)
	if err != nil {
		return err
	}
	_, err = out.Write(buf)
	return err
}

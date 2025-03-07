package lib

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Hardcoding this info here is a little sketchy, but it's a convenient location
// for the info until a final scheme is devised
const CFG_Filename = "config.yaml"
const CFG_Folder = ".config/dylt"
const CFG_Type = "yaml"

// type Config struct {
// 	Viper      *viper.Viper
// 	EtcdDomain string `yaml:"etcd-domain"`
// }

type ConfigStruct struct {
	EtcdDomain string `yaml:"etcd-domain"`
}

// func (o *Config) Get (name string) (interface{}, error) {
// 	return nil, nil
// }

// func (o *Config) GetEtcDomain () (string, error) {
// 	key := "etcd-domain"
// 	isSet := o.Viper.IsSet(key)
// 	if !isSet { return "", nil }
// 	domain := o.Viper.GetString(key)
// 	return domain, nil
// }

// func (o *Config) SetEtcDomain (domain string) (error) {
// 	o.Viper.Set("etcd-domain", domain)
// 	return nil
// }

// func (o *Config) SetEtcDomainAndSave (domain string) (error) {
// 	o.Viper.Set("etcd-domain", domain)
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

func CreateConfigFile() error {
	err := CreateConfigFolder()
	if err != nil {
		return err
	}
	cfgFilePath := GetConfigFilePath()
	f, err := os.OpenFile(cfgFilePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
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

func GetConfigValue(key string) (any, error) {
	f, err := OpenConfigFile()
	if err != nil { return nil, err }
	val, err := GetYamlValue(key, f)
	slog.Debug("GetConfigValue()", "key", key, "val", val)
	if err != nil { return nil, err }

	return val, nil
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
	if err != nil {
		return cfg, NewError(err)
	}
	decoder := yaml.NewDecoder(cfgFile)
	err = decoder.Decode(&cfg)
	if err != nil { return cfg, NewError(err) }
	
	return cfg, nil
}

func OpenConfigFile() (*os.File, error) {
	cfgFilePath := GetConfigFilePath()
	slog.Debug("Opening config file", "cfgFilePath", cfgFilePath)
	f, err := os.Open(cfgFilePath)
	if err != nil { return f, NewError(err) }

	return f, nil
}

func SaveConfig(cfg ConfigStruct) error {
	cfgFilePath := GetConfigFilePath()
	cfgFileFolder := filepath.Dir(cfgFilePath)
	err := os.MkdirAll(cfgFileFolder, 0755)
	if err != nil { return NewError(err) }
	f, err := os.Create(cfgFilePath)
	if err != nil { return NewError(err) }
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func WriteConfig(data configMap) error {
	cfgFilePath := GetConfigFilePath()
	slog.Debug("WriteConfig", "cfgFilePath", cfgFilePath)
	f, err := os.Create(cfgFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = WriteYaml(data, f)
	if err != nil {
		return err
	}
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
	data, err := ReadYaml(cfgFile)
	if err != nil {
		return err
	}
	err = WriteYaml(data, out)
	if err != nil {
		return err
	}
	return err
}

func GetYamlValue(key string, in io.Reader) (any, error) {
	keyParts := strings.Split(key, ".")
	data, err := ReadYaml(in)
	if err != nil {
		return nil, err
	}
	var curr any = data
	for _, currKey := range keyParts {
		currMap := curr.(configMap)
		var ok bool
		curr, ok = currMap[currKey]
		if !ok { return nil, fmt.Errorf("missing key: %s", key) }
	}
	value := curr
	return value, nil
}

func SetKey(data configMap, key string, val string) (configMap, error) {
	dataOrig := data
	keyParts := strings.Split(key, ".")
	for i := range len(keyParts) - 1 {
		keyPart := keyParts[i]
		if data[keyPart] == nil {
			data[keyPart] = map[string]any{}
			data = data[keyPart].(map[string]any)
		}
	}
	lastKey := keyParts[len(keyParts)-1]
	data[lastKey] = val
	return dataOrig, nil
}

type configMap map[string]any

func ReadYaml(in io.Reader) (configMap, error) {
	decoder := yaml.NewDecoder(in)
	var data configMap
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteYaml(data configMap, out io.Writer) error {
	encoder := yaml.NewEncoder(out)
	err := encoder.Encode(data)
	return err
}

package api

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dylt-dev/dylt/common"
)

func RunConfigGet(key string) error {
	slog.Debug("RunConfigGet()", "key", key)
	val, err := common.GetConfigValue(key)
	if err != nil {
		return err
	}
	fmt.Printf("\n%s\n", val)

	return nil
}


func RunConfigSet(key string, val string) error {
	slog.Debug("RunConfigSet()", "key", key, "val", val)

	// Open the dylt config file for read+write. Create if necessasry.
	cfgFilePath := common.GetConfigFilePath()
	slog.Debug("Opening config file", "cfgFilePath", cfgFilePath)
	f, err := os.OpenFile(cfgFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return common.NewError(err)
	}
	defer f.Close()

	// Read the dylt config file as YAML
	data, err := common.ReadYaml(f)
	if err != nil {
		return err
	}

	// Truncate the file to 0 and rewrite
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	// Set the config map value and write the updated config map
	data.Set(key, val)
	err = common.WriteConfig(data)
	if err != nil {
		return err
	}
	err = common.WriteYaml(data, f)
	if err != nil {
		return err
	}

	return nil
}

func RunConfigShow() error {
	err := common.ShowConfig(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}

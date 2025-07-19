package lib

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/dylt-dev/dylt/common"
)

func RunStatus () error {
	var err error
	status := new(statusInfo)

	status.isColima, err = checkColima()
	if err != nil { return err }

	status.isConfigFile, err = isExistConfigFile()
	if err != nil { return err }

	status.isIncus = checkIncus()

	fmt.Printf("%#v\n", status)

	isShellAvailable := isShellAvailable()
	common.Logger.Debugf("isShellAvilable: %t", isShellAvailable)

	return nil
}

func isExistConfigFile () (bool, error) {
	cfgFilePath := common.GetConfigFilePath()
	fi, err := os.Stat(cfgFilePath)
	if err != nil { return false, err }
	if fi.IsDir() {
		return false, fmt.Errorf("config file path exists, but is a directory (%s)", cfgFilePath)
	}
	if fi.Mode() & fs.ModeType > 0 {
		return false, fmt.Errorf("config file exists, but its mode is invalid (%d)", fi.Mode())
	}
	return true, nil
}
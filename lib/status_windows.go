package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dylt-dev/dylt/color"
	"github.com/dylt-dev/dylt/common"
)

func RunStatus() error {
	var err error
	status := new(statusInfo)

	status.isConfigFile, err = isExistConfigFile()
	if err != nil {
		common.Logger.Debug(err.Error())
		status.isConfigFile = false
	}

	fmt.Printf("%-42s %s\n", string(common.Highlight("is config file exist")), color.StyleBool(status.isConfigFile))

	return nil
}

func getShellPath() string {
	var shellPaths = getShellPaths()

	for _, shellPath := range shellPaths {
		common.Logger.Debugf("shellPath=%s\n", shellPath)
		_, err := os.Stat(shellPath)
		if err == nil {
			return shellPath
		}
	}

	return ""
}

func getShellPaths() []string {
	windir := os.Getenv("WINDIR")
	if windir == "" {
		windir = `C:\Windows`
	}
	return []string{
		filepath.Join(windir, "System32", "cmd.exe"),
		filepath.Join(windir, "System32", "WindowsPowerShell", "v1.0", "powershell.exe"),
		filepath.Join(windir, "System32", "WindowsPowerShell", "v1.0", "pwsh.exe"),
	}
}

func getIncusSocketPath() string {
	return ""
}

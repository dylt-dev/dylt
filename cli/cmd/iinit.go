package cmd

import (
	"os"
)

var Logger *cliLogger

func init () {
	Logger = NewLogger(os.Stdout)
}


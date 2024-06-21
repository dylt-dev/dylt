package main

import (
	"os"
	
	"github.com/dylt-dev/dylt/cli"
)

func main () {
	// Here is some random text
	os.Exit(cli.Run())
}

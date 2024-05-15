package lib

import (
	"fmt"
	"os"
)

const PATH_DaylightScript = "/opt/bin/daylight.sh"


func CallDaylightScript (args []string) (int, error) {
	fi, err := os.Stat(PATH_DaylightScript)
	if err != nil { return 1, err }
	fmt.Printf("fi=%s\n", fi)
	return 0, nil
}
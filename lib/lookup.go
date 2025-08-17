package lib

import (
	"fmt"
	"os"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/dns"
)

func RunLookupCommand (hostname string) error {
	var addrs []string
	var err error
	addrs, err = dns.GetA(hostname)
	if err != nil { return err }

	common.PrintBlankIfTerminal(os.Stdout)
	for _, addr  := range addrs {
		fmt.Println("uh oh")
		fmt.Println(addr)
	}
	common.PrintBlankIfTerminal(os.Stdout)

	return nil
}
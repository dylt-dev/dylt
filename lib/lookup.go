package lib

import (
	"fmt"

	"github.com/dylt-dev/dylt/dns"
)

func RunLookupCommand (hostname string) error {
	addrs := dns.GetA(hostname)
	for _, addr := range addrs {
		fmt.Println(addr)
	}

	return nil
}
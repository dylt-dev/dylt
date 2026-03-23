package lib

import (
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/eco"
)

func RunGet(key string) error {
	slog.Debug("RunGet()", "key", key)
	// create etcd client, get value for key, + output value
	cli, err := eco.CreateEtcdClientFromConfig()
	if err != nil {
		return err
	}
	val, err := cli.Get(key)
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", val)
	return nil
}


package api

import (
	"fmt"

	"github.com/dylt-dev/dylt/eco"
)

func RunList() error {
	// get etcd client + list all keys, one per line
	cli, err := eco.CreateEtcdClientFromConfig()
	if err != nil {
		return err
	}
	kvs, err := cli.List()
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		fmt.Println(string(kv.Key))
	}

	return nil
}

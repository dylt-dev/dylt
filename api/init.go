package api

import (
	"errors"
	"log/slog"

	"github.com/dylt-dev/dylt/common"
)

func RunInit(etcdDomain string) error {
	slog.Debug("RunInit()", "etcDomain", etcdDomain)
	// create a new config file using the etcdDomain
	if etcdDomain == "" {
		return errors.New("etcd-domain must be set")
	}
	cfg := common.ConfigStruct{EtcdDomain: etcdDomain}
	err := common.SaveConfig(cfg)
	if err != nil {
		return err
	}

	return nil
}

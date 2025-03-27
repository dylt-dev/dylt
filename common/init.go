package common


// All data specified in flags to the `init` subcommand
type InitStruct struct {
	EtcdDomain string
}

func Init(initData *InitStruct) error {
	cfg := ConfigStruct{
		EtcdDomain: initData.EtcdDomain,
	}
	err := SaveConfig(cfg)
	if err != nil { return err }

	return nil
}

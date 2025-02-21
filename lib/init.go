package lib

type InitInfo struct {
	EtcdDomain string
}

func Init (initInfo *InitInfo) error {
	err := CreateConfigFile()
	if err != nil { return err }
	err = ClearConfigFile()
	if err != nil { return err }
	cfg, err := LoadConfig()
	if err != nil { return err }
	err = cfg.SetEtcDomain(initInfo.EtcdDomain)
	if err != nil { return err }
	err = cfg.Save()
	if err != nil { return err }
	return nil
}
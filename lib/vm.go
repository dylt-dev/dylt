package lib

import (
	// "encoding/json"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientV3 "go.etcd.io/etcd/client/v3"
)

const PRE_vm = "/vm"

type VmClient struct {
	EtcdClient
}

type VmInfo struct {
	Address string
	Name string
}

func (o *VmInfo) Get (field string) (any, error) {
	switch field {
	case "Address":
		return o.Address, nil
	case "Name":
		return o.Name, nil
	default:
		errmsg := fmt.Sprintf("Unknown field: %s", field)
		return nil, errors.New(errmsg)
	}
}

func (o *VmInfo) Set (field string, value any) error {
	switch field {
	case "Address":
	{
		o.Address = value.(string)
	}
	case "Name":
		o.Name = value.(string)
	default:
		errmsg := fmt.Sprintf("Unknown field: %s", field)
		return errors.New(errmsg)
	}
	return nil
}

func (o *VmInfo) String () string {
	buf, _ := json.MarshalIndent(o, "", "\t")
	s := string(buf)
	return s
}

type VmInfoMap map[string]*VmInfo

type VmApi interface {
	Get(name string) (*VmInfo, error)
	Names() ([]string, error)
}


func (cli* VmClient) Add (name string, address string) (*VmInfo, error) {
	key := getKeyFromName(name)
	vm := VmInfo {
		Name: name,
		Address: address,
	}
	value, err := json.Marshal(vm)
	if err != nil { return nil, err }
	ctx := context.Background()
	_, err = cli.KV.Put(ctx, key, string(value))
	if err != nil { return nil, err }
	return &vm, nil
}


func (cli *VmClient) All () (VmInfoMap, error) {
	// Use prefix to get all VM entries from etcd
	all := make(VmInfoMap)
	resp, err := cli.Client.Get(context.Background(), PRE_vm, clientV3.WithPrefix())
	if err != nil { return nil, err }
	for _, kv := range resp.Kvs {
		name := GetVmName(kv)
		vmInfo, err := GetValue(kv)
		all[name] = vmInfo
		if err != nil { return nil, err }
	}
	return all, nil
}


// func FilterOnShr (origVmInfoMap VmInfoMap, shr bool) VmInfoMap {
// 	vmInfoMap := make(VmInfoMap)
// 	for name, info := range origVmInfoMap {
// 		if info.Shr == shr {
// 			vmInfoMap[name] = info
// 		}
// 	}

// 	return vmInfoMap
// }


func (cli *VmClient) Get(name string) (*VmInfo, error) {
	key := getKeyFromName(name)
	data, err := cli.EtcdClient.Get(key)
	if err != nil { return nil, err }
	vm := VmInfo{}
	err = json.Unmarshal(data, &vm)
	if err != nil { return nil, err }
	return &vm, nil
}


func GetVmName (kv *mvccpb.KeyValue) string {
	name := getNameFromKey(string(kv.Key))
	return name
}


func GetValue (kv *mvccpb.KeyValue) (*VmInfo, error) {
	rawVal := kv.Value
	vm := VmInfo{}
	err := json.Unmarshal(rawVal, &vm)
	if err != nil { return nil, err }
	return &vm, nil
}


func (cli *VmClient) Names() ([]string, error) {
	resp, err := cli.Client.Get(context.Background(), PRE_vm, clientV3.WithPrefix())
	if err != nil { return nil, err }
	var names []string
	for _, kv := range resp.Kvs {
		name := getNameFromKey(string(kv.Key))
		names = append(names, name)
	}
	return names, nil
}

func CreateVmClientFromConfig () (*VmClient, error) {
	cfg, err := LoadConfig()
	if err != nil { return nil, err }
	vmClient, err := NewVmClient(cfg.EtcdDomain)
	return vmClient, err
}


func NewVmClient(domain string) (*VmClient, error) {
	etcdClient, err := NewEtcdClient(domain)
	if err != nil {
		return nil, err
	}
	cli := VmClient{EtcdClient: *etcdClient}
	return &cli, err
}


func NewVmInfo (kv *mvccpb.KeyValue) (*VmInfo, error) {
	jsonData, err := json.Marshal(kv)
	if err != nil { return nil, err }
	fmt.Print(string(jsonData))
	return nil, nil
}


func getKeyFromName (name string ) string {
	s := fmt.Sprintf("%s/%s", PRE_vm, name)
	return s
}


func getNameFromKey (key string) string {
	prefix := fmt.Sprintf("%s/", PRE_vm)
	name := strings.TrimPrefix(key, prefix)
	return name
}
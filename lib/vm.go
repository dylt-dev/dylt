package lib

import (
	// "encoding/json"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientV3 "go.etcd.io/etcd/client/v3"
)

const PRE_vm = "/vm/"

type VmClient struct {
	EtcdClient
}

type VmInfo struct {
	Host string
	Shr  bool
}

type VmInfoMap map[string]*VmInfo

type VmApi interface {
	Get(name string) (*VmInfo, error)
	Names() ([]string, error)
}


func (cli *VmClient) All () (VmInfoMap, error) {
	// Use prefix to get all VM entries from etcd
	all := make(VmInfoMap)
	resp, err := cli.Client.Get(context.Background(), PRE_vm, clientV3.WithPrefix())
	if err != nil { return nil, err }
	for _, kv := range resp.Kvs {
		key := GetKey(kv)
		vmInfo, err := GetValue(kv)
		all[key] = vmInfo
		if err != nil { return nil, err }
	}
	return all, nil
}


func FilterOnShr (origVmInfoMap VmInfoMap, shr bool) VmInfoMap {
	vmInfoMap := make(VmInfoMap)
	for name, info := range origVmInfoMap {
		if info.Shr == shr {
			vmInfoMap[name] = info
		}
	}

	return vmInfoMap
}


func (cli *VmClient) Get(name string) (*VmInfo, error) {
	key := fmt.Sprintf("%s%s", PRE_vm, name)
	data, err := cli.EtcdClient.Get(key)
	if err != nil { return nil, err }
	vm := VmInfo{}
	err = json.Unmarshal(data, &vm)
	if err != nil { return nil, err }
	return &vm, nil
}


func GetKey (kv *mvccpb.KeyValue) string {
	key := strings.TrimPrefix(string(kv.Key), PRE_vm)
	return key
}


func GetValue (kv *mvccpb.KeyValue) (*VmInfo, error) {
	rawVal := kv.Value
	vm := VmInfo{}
	err := json.Unmarshal(rawVal, &vm)
	if err != nil { return nil, err }
	return &vm, nil
}


func (cli *VmClient) Names() ([]string, error) {
	etcdClient, err := NewEtcdClient("hello.dylt.dev")
	if err != nil { return nil, err }
	resp, err := etcdClient.Client.Get(context.Background(), PRE_vm, clientV3.WithPrefix())
	if err != nil { return nil, err }
	var names []string
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		name, found := strings.CutPrefix(key, PRE_vm)
		if !found {
			panic("This shouldn't happen")
		}
		names = append(names, name)
	}
	return names, nil
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
	fmt.Printf(string(jsonData))
	return nil, nil
}
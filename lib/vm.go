package lib

import (
	// "encoding/json"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	clientV3 "go.etcd.io/etcd/client/v3"
)

const prefix = "/vm/"

type VmClient struct{
	EtcdClient
}

type VmInfo struct {
	Host string
	Shr bool
}

type VmApi interface {
	Get (name string) (*VmInfo, error)
	Names () ([]string, error)
}

func (cli* VmClient) Get (name string) (*VmInfo, error) {
	key := fmt.Sprintf("%s%s", prefix, name)
	data, err := cli.EtcdClient.Get(key)
	if err != nil { return nil, err }
	vm := VmInfo{}
	err = json.Unmarshal(data, &vm)
	if err != nil { return nil, err }
	return &vm, nil
}

func (cli* VmClient) Names () ([]string, error) {
	etcdClient, err := NewEtcdClient("hello.dylt.dev")
	if err != nil { return nil, err }
	resp, err := etcdClient.Client.Get(context.Background(), prefix, clientV3.WithPrefix())
	if err != nil { return nil, err }
	var names []string
	for _, kv := range(resp.Kvs) {
		key := string(kv.Key)
		name, found := strings.CutPrefix(key, prefix)
		if !found { panic("This shouldn't happen") }
		names = append(names, name)
	}
	return names, nil
}

func NewVmClient (domain string) (*VmClient, error) {
	etcdClient, err := NewEtcdClient(domain)
	if err != nil { return nil, err }
	cli := VmClient{EtcdClient: *etcdClient}
	return &cli, err
}
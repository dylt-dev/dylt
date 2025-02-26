package lib

import (
	"context"
	"fmt"
	"os"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientV3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	*clientV3.Client
}

func (cli *EtcdClient) Delete(key string) ([]byte, error) {
	ctx := context.Background()
	oldVal, err := cli.Get(key)
	if err != nil {
		fmt.Println("Error getting old value")
	}
	_, err = cli.Client.Delete(ctx, key)
	if err != nil {
		return nil, err
	}
	return oldVal, nil
}

func (cli *EtcdClient) Get(key string) ([]byte, error) {
	ctx := context.Background()
	resp, err := cli.Client.Get(ctx, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err type=%T\n", err)
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	val := (*resp).Kvs[0].Value
	return val, err
}

func (cli *EtcdClient) GetKeys(prefix string) ([]string, error) {
	ctx := context.Background()
	resp, err := cli.Client.Get(ctx, prefix, clientV3.WithPrefix())
	if err != nil {
		return nil, err
	}
	keys := []string{}
	for _, kv := range resp.Kvs {
		keys = append(keys, string(kv.Key))
	}
	return keys, nil
}

func (cli *EtcdClient) List() ([]*mvccpb.KeyValue, error) {
	ctx := context.Background()
	resp, err := cli.Client.Get(ctx, "", clientV3.WithPrefix())
	if err != nil {
		return nil, err
	}
	return resp.Kvs, nil
}

func CreateEtcdClientFromConfig() (*EtcdClient, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	domain := cfg.EtcdDomain
	cli, err := NewEtcdClient(domain)
	return cli, err
}

func NewEtcdClient(domain string) (*EtcdClient, error) {
	endpoints, err := getEndpoints(domain)
	if err != nil {
		return nil, err
	}
	cfg := clientV3.Config{Endpoints: endpoints}
	cli, err := clientV3.New(cfg)
	if err != nil {
		return nil, err
	}
	client := &EtcdClient{Client: cli}
	return client, err
}

func getEndpoints(domain string) ([]string, error) {
	endpoints := []string{}
	srvs, err := GetSrvs(domain, "etcd-server", "tcp", true)
	if err != nil {
		return nil, err
	}
	for _, srv := range srvs {
		ip := srv.Ips[0]
		port := srv.Port
		endpoint := fmt.Sprintf("http://%s:%d", ip, port)
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, nil
}

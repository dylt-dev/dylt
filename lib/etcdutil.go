package lib

import (
	"context"
	"fmt"
	"log/slog"

	clientV3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	*clientV3.Client
}


func (cli *EtcdClient) Get (key string) ([]byte, error) {
	ctx := context.Background()
	kv, err := cli.Client.Get(ctx, key)
	if err != nil {
		slog.Warn("Error calling cli.Client.Get()", err)
	}
	val := (*kv).Kvs[0].Value
	return val, err
}


func NewEtcdClient (domain string) (*EtcdClient, error) {
	endpoints, err := getEndpoints(domain)
	if err != nil {
		slog.Warn("Error calling NewEtcdClient", err)
	}
	cfg := clientV3.Config{Endpoints: endpoints}
	cli, err := clientV3.New(cfg)
	if err != nil {
		slog.Warn("Error calling NewEtcdClient", err)
	}
	client := &EtcdClient{Client: cli}
	return client, err
}


func getEndpoints (domain string) ([]string, error) {
	endpoints := []string{}
	srvs, err := GetSrvs(domain, true)
	if err != nil {
		slog.Warn("Failure calling GetSrvs()", err)
		return nil, err
	}
	for _, srv := range(srvs) {
		ip := srv.Ips[0]
		port := srv.Port
		endpoint := fmt.Sprintf("http://%s:%d", ip, port)
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, nil
}
package lib

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestGetServiceKey (t *testing.T) {
	svcKey := "/#/svc/yellowrose/bin/foo"
	serviceName, err := GetServiceName(svcKey)
	assert.NoError(t, err)
	assert.Equal(t, "yellowrose", serviceName)
}

func TestGetServiceKeyErr1 (t *testing.T) {
	svcKey := "svc/yellowrose/bin/foo"
	serviceName, err := GetServiceName(svcKey)
	assert.Error(t, err)
	assert.Equal(t, "", serviceName)
}

func TestGetServiceKeyErr2 (t *testing.T) {
	svcKey := "/#/SVC/yellowrose/bin/foo"
	serviceName, err := GetServiceName(svcKey)
	assert.Error(t, err)
	assert.Equal(t, "", serviceName)
}

func TestGetServiceKeyErr3 (t *testing.T) {
	svcKey := "/#/svc/"
	serviceName, err := GetServiceName(svcKey)
	assert.Error(t, err)
	assert.Equal(t, "", serviceName)
}

func TestWatchHello (t *testing.T) {
	defer func () { if x := recover(); x != nil { t.Logf("Panic in the streets of %#v", x)}}()
	
	ctx := clientv3.WithRequireLeader(context.Background())
	assert.NotNil(t, ctx)
	cli, err := CreateEtcdClientFromConfig()
	assert.Nil(t, err)
	chWatch := cli.Watch(ctx, "/hello", clientv3.WithKeysOnly()) 
	var resp clientv3.WatchResponse
	for resp = range chWatch {
		for _, ev := range resp.Events {
			t.Logf("Update detected: %s", ev.Kv.Key)
		}
	}
}

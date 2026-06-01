package lib

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dylt-dev/dylt/eco"
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
	envvar, is := os.LookupEnv("DYLT_SYSTEST")
	if !is || envvar != "Y" {
		t.Skip("sys test only")
	}
	
	defer func () { if x := recover(); x != nil { t.Logf("Panic in the streets of %#v", x)}}()
	
	ctxTimeout, fnCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer fnCancel()
	ctx := clientv3.WithRequireLeader(ctxTimeout)
	assert.NotNil(t, ctx)
	cli, err := eco.CreateEtcdClientFromConfig()
	assert.Nil(t, err)
	chWatch := cli.Watch(ctx, "/hello", clientv3.WithKeysOnly()) 
	t.Logf("%#v", chWatch)
	var resp clientv3.WatchResponse
	for resp = range chWatch {
		for _, ev := range resp.Events {
			t.Logf("Update detected: %s", ev.Kv.Key)
		}
	}
}

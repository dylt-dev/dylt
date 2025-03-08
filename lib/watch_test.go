package lib

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

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

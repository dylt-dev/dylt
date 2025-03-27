package common

import (
	"context"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	etcd3 "go.etcd.io/etcd/client/v3"
)

type EcoTest struct {
	Name string `eco:"Name"`
	LuckyNumber float64 `echo:"LuckyNumber"`
}

func NewEcoTest (name string, luckyNumber float64) *EcoTest {
	return &EcoTest{Name: name, LuckyNumber: luckyNumber}
}

func TestGetObject (t *testing.T) {
	
}

func TestPutObject (t *testing.T) {
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)
	obj := NewEcoTest("Me", 13)

	prefix := "/_test_/echotest"
	opName := etcd3.OpPut("Name", filepath.Join(prefix, obj.Name))
	opLuckyNumber := etcd3.OpPut("LuckyNumber", filepath.Join(prefix, strconv.FormatFloat(obj.LuckyNumber, 'f', 8, 64)))
	txn := etcdClient.Txn(context.Background())
	txn.Then(opName, opLuckyNumber).Commit()
}

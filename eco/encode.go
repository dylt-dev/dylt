package eco

import (
	"reflect"

	"github.com/dylt-dev/dylt/common"
	etcd "go.etcd.io/etcd/client/v3"
)



func Encode(ctx* common.EcoContext, cli *EtcdClient, a any, key string) error {
	ctx.Signature("Encode", reflect.TypeOf(a), key)
	ctx.Inc()
	defer ctx.Dec()
	
	// get kvs	
	kvs := encode(ctx, a)

	// convert to etcd Ops, translating the key along the way
	rootKey := KeyString(key).WithEndSlash()
	ops := make([]etcd.Op, len(kvs))
	for i, kv := range kvs {
		opKey := string(rootKey.Add(kv.Key))
		opVal := string(kv.Value)
		op := etcd.OpPut(opKey, opVal)
		ops[i] = op
	}

	// Write the Ops to the cluster
	txn := cli.Txn(ctx)
	_, err := txn.Then(ops...).Commit()

	return err
}

func encode(ctx *common.EcoContext, a any) []KeyValue {

	ctx.Signature("Encode", reflect.TypeOf(a))
	ctx.Inc()
	defer ctx.Dec()

	rve := NewRvEncodable(a)
	kvs := rve.Encode(ctx, KeyString(""))
	return kvs
}

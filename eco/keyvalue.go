package eco

import (
	"slices"
	"strings"

	"github.com/dylt-dev/dylt/common"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type KeyValue struct {
	Key   KeyString
	Value []byte
}

func NewKeyValue(etcdKv *mvccpb.KeyValue) KeyValue {
	kv := KeyValue{Key: KeyString(etcdKv.Key), Value: etcdKv.Value}
	return kv
}

func createKvSlice(etcdKvs []*mvccpb.KeyValue) []*KeyValue {
	kvSlice := make([]*KeyValue, len(etcdKvs))
	for i, etcdKv := range etcdKvs {
		kvSlice[i] = newKvFromEtcd(etcdKv)
	}

	return kvSlice
}

func deleteKeyFromSlice(ctx *common.EcoContext, kvs []*KeyValue, key string) []*KeyValue {
	ctx.Logger.Signature("deleteKeyFromSlice", len(kvs), key)
	ctx.Inc()
	defer ctx.Dec()

	ctx.Logger.Infof("Before: len(kvs)=%d", len(kvs))
	ctx.Logger.Commentf("Getting index of %s ...", key)
	iKv := slices.IndexFunc(kvs, func(kv *KeyValue) bool { return key == string(kv.Key) })
	ctx.Logger.Infof("iKv=%d", iKv)
	if iKv != -1 {
		ctx.Logger.Comment("Deleting element from slice")
		kvs = slices.Delete(kvs, iKv, iKv+1)
	}
	ctx.Logger.Infof("After: len(kvs)=%d", len(kvs))
	return kvs
}

func newKv(k string, v string) *KeyValue {
	kv := new(KeyValue)
	kv.Key = KeyString(k)
	kv.Value = []byte(v)

	return kv
}

func newKvFromEtcd(etcdKv *mvccpb.KeyValue) *KeyValue {
	kv := new(KeyValue)
	kv.Key = KeyString(string(etcdKv.Key))
	kv.Value = etcdKv.Value

	return kv
}

/*
findKv() is a little tricky, because it will find logical keys that are
present in the KV list, as well as actual keys.

A logical key is a key that matches the prefix of one more more keys in the KV
list. An example is the parent key of a map, struct, or slice. The the parent
key 'exists' in that its presence is implied, in the list, but the key itself
does not exist in etcd and therefore has no assoiciated value. On the other
hand, a key that represents a struct field etc will be explicitly included in
the KV list and therefore have a name and a value.

Logical keys get a KV of (key, nil)
Physical keys get a KV of (key, val)
Keys that are missing entirely get a KV of nil

@note 'implied key' might be more clear than 'logical key'. Or ... it might not.
*/

func findKv(key string, kvs []*KeyValue) *KeyValue {
	var kv *KeyValue
	var isLogical bool = false
	var isPhysical bool = false

	for _, kv = range kvs {
		s := string(kv.Key)
		// Does at least one key have the incoming key as a prefix?
		is := strings.HasPrefix(s, key)
		if is {
			isLogical = true
		}
		// Is there an exact match?
		if s == key {
			isPhysical = true
			break
		}
	}

	if isPhysical {
		return &KeyValue{Key: KeyString(key), Value: kv.Value}
	}

	if isLogical {
		return &KeyValue{Key: KeyString(key), Value: nil}
	}

	return nil
}

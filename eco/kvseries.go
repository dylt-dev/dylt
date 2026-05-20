package eco

import "go.etcd.io/etcd/api/v3/mvccpb"

type KvSeries struct {
	RootKey KeyString
	Kvs     []KeyValue
}


func NewKvSeries (rootKey KeyString, kvs []*mvccpb.KeyValue) (*KvSeries, error) {
	kvSeries := KvSeries{rootKey, nil}
	for _, kv := range kvs {
		kvSeries.Add(NewKeyValue(kv))	
	}

	return &kvSeries, nil
}

func (this *KvSeries) Add (kv KeyValue) bool {
	if !this.IsOwner(kv.Key) && kv.Key != this.RootKey{
		return false
	}
	this.Kvs = append(this.Kvs, kv)
	return true
}


func (this *KvSeries) Len() int {
	if this.Kvs == nil {
		return 0
	}

	return len(this.Kvs)
}

func (this *KvSeries) MaxIndex() int {
	var maxIndex int = 0
	for _, kv := range this.Kvs {
		keyString, is := kv.Key.CutPrefix(this.RootKey)
		if is {
			index, is := keyString.Index()
			if is && index > maxIndex {
				maxIndex = index
			}
		}
	}

	return maxIndex
}

func (this *KvSeries) IsOwner (keyString KeyString) bool {
	return this.RootKey.IsParent(keyString)
}

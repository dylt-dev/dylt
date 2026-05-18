package eco

type KvSeries struct {
	RootKey KeyString
	Kvs     []KeyValue
}


func (this *KvSeries) Add (kv KeyValue) bool {
	if !this.IsOwner(kv.Key) {
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
		keyString, is := kv.Key.CutPrefix(string(this.RootKey))
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

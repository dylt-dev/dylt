package eco

type KvSeries struct {
	Kvs     []KeyValue
	RootKey KeyString
}


func (this KvSeries) Len() int {
	if this.Kvs == nil {
		return 0
	}

	return len(this.Kvs)
}

func (this KvSeries) MaxIndex() int {
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

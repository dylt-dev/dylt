package eco

type KvSlice struct {
	Kvs []KeyValue
	RootKey KeyString
}


func (this KvSlice) MaxIndex () int {
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
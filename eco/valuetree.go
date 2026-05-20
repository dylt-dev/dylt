package eco

import (
	"strconv"

	"github.com/dylt-dev/dylt/common"
)

type ValueTree struct {
	Value []byte
	ChildMap ValueTreeChildMap
}


type ValueTreeChildMap map[string]*ValueTree


func NewValueTree (ctx *common.EcoContext, kvSeries *KvSeries) (*ValueTree, error) {
	tree := new(ValueTree)
	if kvSeries == nil {
		return tree, nil
	}

	for _, kv := range kvSeries.Kvs {
		key, is := KeyString(kv.Key).CutPrefix(kvSeries.RootKey)
		if is {
			tree.Add(ctx, key, kv.Value)
		}
	}

	return tree, nil
}

func (tree *ValueTree) Add(ctx *common.EcoContext, ks KeyString, val []byte) {
	ctx.Logger.Signature("ValueTree.Add", ks)
	ctx.Logger.Infof("tree=%p", tree)
	ctx.Inc()
	defer ctx.Dec()

	if ks == "" {
		ctx.Logger.Comment("IsLeaf() - setting value")
		tree.Value = val
	} else {
		ctx.Logger.Comment("No leaf - recursively add value to child tree")
		head, body := ks.PopHead()
		ctx.Logger.Infof("head=%s body=%s", head, string(body))
		childTree, is := tree.ChildMap[head]
		if !is {
			ctx.Logger.Infof("Child tree for '%s' not found; creating", head)
			if tree.ChildMap == nil {
				tree.ChildMap = ValueTreeChildMap{}
			}
			childTree = &ValueTree{}
			tree.ChildMap[head] = childTree
		}
		childTree.Add(ctx, body, val)
	}
}


func (vtcm *ValueTreeChildMap) MaxIndex(ctx *common.EcoContext) int {
	maxIndex := -1
	for k := range *vtcm {
		i, err := strconv.Atoi(k)
		if err == nil && i > maxIndex {
			maxIndex = i
		}
	}

	return maxIndex
}


func (tree *ValueTree) addBool(ctx *common.EcoContext, key string, val bool) {
	ks := KeyString(key)
	buf := strconv.AppendBool([]byte{}, val)
	tree.Add(ctx, ks, buf)
}


func (tree *ValueTree) setBool (val bool) {
	buf := strconv.AppendBool([]byte{}, val)
	tree.Value = buf
	
}
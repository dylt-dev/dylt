package eco

import "github.com/dylt-dev/dylt/common"

type ValueTree struct {
	Value []byte
	ChildMap ValueTreeChildMap
}


type ValueTreeChildMap map[string]*ValueTree

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
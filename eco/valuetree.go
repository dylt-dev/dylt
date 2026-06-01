package eco

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/dylt-dev/dylt/common"
)

type ValueTree struct {
	Value    []byte
	ChildMap ValueTreeChildMap
}

type ValueTreeChildMap map[string]*ValueTree

func NewValueTree(ctx *common.EcoContext, childArgs ...any) *ValueTree {
	tree := &ValueTree{ChildMap: ValueTreeChildMap{}}
	if len(childArgs) == 1 {
		tree.set(childArgs[0])
	} else if len(childArgs) > 1 {
		i := 0
		for i < len(childArgs) {
			childKey := unmarsalChildMapKey(childArgs[i])
			childValue := childArgs[i+1]
			tree.add(ctx, childKey, childValue)
			i += 2
		}
	}
	return tree
}

func NewValueTreeFromKvSeries(ctx *common.EcoContext, kvSeries *KvSeries) (*ValueTree, error) {
	tree := new(ValueTree)
	if kvSeries == nil {
		return tree, nil
	}

	for _, kv := range kvSeries.Kvs {
		key, is := KeyString(kv.Key).CutPrefix(kvSeries.RootKey)
		if is {
			tree.Add(ctx, key.WithoutEndSlash(), kv.Value)
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

func (tree *ValueTree) add(ctx *common.EcoContext, key string, a any) {
	switch i := a.(type) {
	case *ValueTree:
		if tree.ChildMap == nil {
			tree.ChildMap = ValueTreeChildMap{}
		}
		tree.ChildMap[key] = i
	default:
		if common.NewFlavor(reflect.TypeOf(a).Kind()) == common.Scalar {
			ks := KeyString(key)
			ksName, _ := ks.PopHead()
			ctx.Logger.Infof("ksName=%s", ksName)
			ctx.Logger.Infof("Scalar (val=%v)", a)
			treeChild := &ValueTree{}
			treeChild.set(a)
			tree.add(ctx, string(ksName), treeChild)
		} else {
			panic(fmt.Sprintf("Unexpected type: %s", reflect.TypeOf(a)))
		}
	}
}

// func (tree *ValueTree) addBool(ctx *common.EcoContext, key string, val bool) {
// 	ks := KeyString(key)
// 	buf := strconv.AppendBool([]byte{}, val)
// 	tree.Add(ctx, ks, buf)
// }

// func (tree *ValueTree) addFloat(ctx *common.EcoContext, key string, val float64) {
// 	ks := KeyString(key)
// 	buf := strconv.AppendFloat([]byte{}, val, 'f', -1, 64)
// 	tree.Add(ctx, ks, buf)
// }

// func (tree *ValueTree) addInt(ctx *common.EcoContext, key string, val int64) {
// 	ks := KeyString(key)
// 	buf := strconv.AppendInt([]byte{}, val, 10)
// 	tree.Add(ctx, ks, buf)
// }

// func (tree *ValueTree) addString(ctx *common.EcoContext, key string, val string) {
// 	ks := KeyString(key)
// 	s := fmt.Sprintf(`"%s"`, val)
// 	buf := []byte(s)
// 	tree.Add(ctx, ks, buf)
// }

func (tree *ValueTree) set(a any) {
	buf, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	tree.Value = buf
}

func unmarsalChildMapKey(a any) string {
	switch i := a.(type) {
	case string:
		return i
	default:
		key, err := json.Marshal(a)
		if err != nil {
			panic(err)
		}
		return string(key)
	}
}

// func (tree *ValueTree) setBool (val bool) {
// 	buf := strconv.AppendBool([]byte{}, val)
// 	tree.Value = buf

// }

// func (tree *ValueTree) setFloat (val float64) {
// 	buf := strconv.AppendFloat([]byte{}, val, 'f', -1, 64)
// 	tree.Value = buf

// }

// func (tree *ValueTree) setInt (val int64) {
// 	buf := strconv.AppendInt([]byte{}, val, 10)
// 	tree.Value = buf

// }

// func (tree *ValueTree) setString (val string) {
// 	s := fmt.Sprintf("%q", val)
// 	tree.Value = []byte(s)
// }

// func (tree *ValueTree) setUint (val uint64) {
// 	buf := strconv.AppendUint([]byte{}, val, 10)
// 	tree.Value = buf

// }

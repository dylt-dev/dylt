package eco

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dylt-dev/dylt/common"
)

type KeyValueChildMap map[KeyString]*KeyValueTree

func (m KeyValueChildMap) MaxIndex() uint64 {
	var maxIndex uint64 = 0
	for key := range m {
		keyString := KeyString(key)
		index, is := keyString.Index()
		if is {
			if index > maxIndex {
				maxIndex = index
			}
		}
	}

	return maxIndex
}

type KeyValueTree struct {
	Value    []byte
	Children KeyValueChildMap
}

func (t *KeyValueTree) String() string {
	ctx := common.NewEcoContext(os.Stdout)
	common.InitLogging()
	return treeToString(ctx, t)
}

func treeToString (ctx *common.EcoContext, kvTree *KeyValueTree) string{
	ctx.Inc()
	defer ctx.Dec()

	sb := strings.Builder{}
	buf, err := json.Marshal(kvTree.Value)
	if err != nil {
		buf = []byte{}
	}
	fmt.Fprintf(&sb, "%sValue: %s\n", ctx.Logger.Indent(), string(buf))
	fmt.Fprintf(&sb, "%slen(Children): %d\n", ctx.Logger.Indent(), len(kvTree.Children))
	for _, child := range(kvTree.Children) {
		sChild := treeToString(ctx, child)
		sb.WriteString(sChild)
	}
	return sb.String()
}

/*
It's hard to keep names, keys, and kvs straight. Let's try and account for them all
- the list of KVs that comprise object data
- the key that represents an element
- the root key that represents where the element was found (maybe this doesn't matter and maintaining it is a hassle)
- the name that serves as a key in a kvMap
- the key of a child map - honestly I'm not sure what this is

Let's figure this out

The root key is the root key. It's the key that was used to fetch the initial
object data.
The key representing the element is the current key being processed.
kvMap names are the first segment that follows the element key, for each child key
the list of kvs is just child keys
the key of a child map is the childMap key appended to the prefix. It's a full key
*/

func createKvTree(ctx *common.EcoContext, key string, kvs []*KeyValue, rootKey string) *KeyValueTree {
	ctx.Logger.Signature("createKvTree", key, len(kvs), rootKey)
	ctx.Inc()
	defer ctx.Dec()

	// Find the kv for the specified key in the kvs list, if present
	ctx.Logger.Commentf("Finding key (%s) in kv list ...", key)
	kv := findKv(key, kvs)
	if kv == nil {
		return nil
	}
	ctx.Logger.Infof("Key found. kv=(%v, %v)", kv.Key, string(kv.Value))

	// Create a new tree + set its simple fields
	tree := new(KeyValueTree)
	ctx.Logger.Commentf("Getting element name of %s under rootKey=%s", key, rootKey)
	var value []byte
	value = kv.Value
	ctx.Logger.Comment("value of element determined.")
	ctx.Logger.Infof("value=%v", value)
	tree.Value = value

	ctx.Logger.Comment("creating child nodes, if present")
	tree.Children = KeyValueChildMap{}
	// Remove the current node from the collection of nodes to process
	kvs = deleteKeyFromSlice(ctx, kvs, string(key))
	// Create KeyValueMap + KeyValueTree Children from child keys, if any
	if len(kvs) == 0 {
		return tree
	}

	// Recursively create child nodes
	kvMap := createKvMap(ctx, key, kvs)
	ctx.Logger.Infof("%d children found", len(kvMap))
	for k, v := range kvMap {
		childKey := fmt.Sprintf("%s/%s", key, k)
		childKvs := v
		tree.Children[KeyString(childKey)] = createKvTree(ctx, childKey, childKvs, rootKey)
	}
	
	return tree
}

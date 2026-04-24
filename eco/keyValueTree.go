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
	Name     string
	Value    []byte
	Children KeyValueChildMap
}

func (t *KeyValueTree) String() string {
	ctx := newEcoContext(os.Stdout)
	common.InitLogging()
	return treeToString(ctx, t)
}

func treeToString (ctx *ecoContext, kvTree *KeyValueTree) string{
	ctx.inc()
	defer ctx.dec()

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%sName: %s\n", ctx.logger.indent(), kvTree.Name))
	buf, err := json.Marshal(kvTree.Value)
	if err != nil {
		buf = []byte{}
	}
	fmt.Fprintf(&sb, "%sValue: %s\n", ctx.logger.indent(), string(buf))
	fmt.Fprintf(&sb, "%slen(Children): %d\n", ctx.logger.indent(), len(kvTree.Children))
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

func createKvTree(ctx *ecoContext, key string, kvs []*KeyValue, rootKey string) *KeyValueTree {
	ctx.logger.signature("createKvTree", key, len(kvs), rootKey)
	ctx.inc()
	defer ctx.dec()

	// Find the kv for the specified key in the kvs list, if present
	ctx.logger.commentf("Finding key (%s) in kv list ...", key)
	kv := findKv(key, kvs)
	if kv == nil {
		return nil
	}
	ctx.logger.Infof("Key found. kv=(%v, %v)", kv.Key, string(kv.Value))

	// Create a new tree + set its simp`le fields
	tree := new(KeyValueTree)
	ctx.logger.commentf("Getting element name of %s under rootKey=%s", key, rootKey)
	name := KeyString(key).ElementName(rootKey)
	ctx.logger.Infof("element name=%s", name)
	var value []byte
	value = kv.Value
	ctx.logger.comment("name + value of element determined.")
	ctx.logger.Infof("name=%v value=%v", name, value)
	tree.Name = name
	tree.Value = value

	// Recursively create child nodes
	ctx.logger.comment("creating child nodes, if present")
	tree.Children = KeyValueChildMap{}
	// Remove the current node from the collection of nodes to process
	kvs = deleteKeyFromSlice(ctx, kvs, string(key))
	// Create KeyValueMap + KeyValueTree Children from child keys, if any
	if len(kvs) > 0 {
		kvMap := createKvMap(ctx, key, kvs)
		ctx.logger.Infof("%d children found", len(kvMap))
		for k, v := range kvMap {
			childKey := fmt.Sprintf("%s/%s", key, k)
			childKvs := v
			tree.Children[KeyString(childKey)] = createKvTree(ctx, childKey, childKvs, rootKey)
		}
	}

	return tree
}

package eco

import (
	"encoding/json"
	"fmt"
	"strings"
)

type KeyValueChildMap map[string]*KeyValueTree

type KeyValueTree struct {
	Name     string
	Value    []byte
	Children KeyValueChildMap
}


func (t *KeyValueTree) String () string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Name: %s\n", t.Name))
	buf, err := json.Marshal(t.Value)
	if err != nil {
		buf = []byte{}
	}
	fmt.Fprintf(&sb, "Value: %s\n", string(buf))
	fmt.Fprintf(&sb, "len(Children): %d\n", len(t.Children))
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

The root key is the root key. It's where the initial object data was fetched from
The key representing the element is the current key being processed.
kvMap names are the first segment that follows the element key, for each child key
the list of kvs is just child keys
the key of a child map is the childMap key appended to the prefix. It's a full key
*/

func createKvTree(ctx *ecoContext, key string, kvs []*KeyValue, rootKey string) *KeyValueTree {
	ctx.inc()
	defer ctx.dec()
	ctx.logger.signature("createKvTree", key, len(kvs), rootKey)

	// Find the kv for the specified key in the kvs list, if present
	kv := findKv(key, kvs)
	if kv == nil {
		return nil
	}
	fmt.Printf("kv=(%v, %v)\n", kv.Key, string(kv.Value))

	// Create a new tree + set its simp`le fields
	tree := new(KeyValueTree)
	ctx.logger.commentf("Getting child name of %s (rootKey=%s)", key, rootKey)
	// name := KeyString(key).ChildName(rootKey)
	name := KeyString(key).ElementName(rootKey)
	ctx.logger.commentf("child name=%s", name)
	var value []byte
	value = kv.Value
	fmt.Printf("name=%v\n", name)
	fmt.Printf("value=%v\n", string(value))
	tree.Name = name
	tree.Value = value

	// Recursively create child nodes
	tree.Children = KeyValueChildMap{}
	kvMap := createKvMap(ctx, key, kvs)
	for k, v := range kvMap {
		childKey := fmt.Sprintf("%s/%s", key, k)
		childKvs := v
		tree.Children[childKey] = createKvTree(ctx, childKey, childKvs, rootKey)
	}

	return tree
}

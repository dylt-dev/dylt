package eco

import (
	"encoding/json"
)

/*

KeyValueMap is a helper data structure to assist with the creation of
KeyValueTrees. To produce a KeyValueTree from a (parentKey, []KeyValue) pair,
it is necessary to do the following:

	Create a unique collection of all the child keys of the parent key
	For each child key, create a collection of all its descendant keys, including
      the child key itself

It would be a fair hassle to implement these steps and maintain these data
collections ... if a map weren't already perfect for both. The idiomatic Go way
to create a unique collection is to create a map, since a map's keys are
guaranteed to be unique.

Storing child key => []KV relationships in a map can be done in a single pass
over the []KV collection. This is a nice efficiency win. Otherwise the []KV
would need to be traversed once per element, which is O(n^2). The performance
win is highly unlikely to matter, but avoiding such inefficiencies is
satisfying.

Example of KeyValueMap

*/

type KeyValueMap map[string][]*KeyValue

func (m KeyValueMap) add(ctx *ecoContext, prefix string, kv *KeyValue) bool {
	ctx.logger.signature("KeyValueMap.add()", prefix, kv.Key)
	ctx.inc()
	defer ctx.dec()
	
	// Find the first segment after the prefix. This will be the map key
	// If there is no segment after the prefix, then return false
	ctx.logger.comment("checking if key is a child of the prefix ...")
	fullKey := KeyString(kv.Key)
	ksPrefix := KeyString(prefix).WithoutEndSlash()
	afterPrefix, is := fullKey.CutPrefix(string(ksPrefix))
	if !is {
		ctx.logger.Info("Not a child. Returning.")
		return false
	}
	segments := KeyString(afterPrefix).Segments()
	if len(segments) == 0 {
		ctx.logger.Info("Not a child. Returning.")
		return false
	}
	key := segments[0]
	ctx.logger.Infof("Key (%s) is child.", key)
	// Get the kvs for this key. If the key doesn't exist, create a new lv list
	kvs := m[key]
	if kvs == nil {
		kvs = []*KeyValue{}
	}

	ctx.logger.commentf("append child key (%s) to to m[%s]", kv.Key, key)
	kvs = append(kvs, kv)
	m[key] = kvs
	// Note this final set might be unnecessary, and conceivably inefficient

	// done :)
	return true
}


func (m KeyValueMap) String () string {
	// We don't actually want to do any string formatting here
	// We're perfectly happy to create a similar data structure to
	// mvccpb.KeyValue, but with strings instead of []byte, and to
	// let the stdlib Unmarshaller handle it
	// @note this is yet another hint at the utility of a helper data structure
	//       or two for mvccpb.KeyValue
	type skv struct { Key string; Value string}
	type skvm map[string][]skv
	smap := skvm{}
	for k, v := range m {
		skey := k
		skvs := []skv{}
		for _, el := range(v) {
			sel := skv{Key: string(el.Key), Value: string(el.Value)}
			skvs = append(skvs, sel)
		}
		smap[skey] = skvs
	}
	buf, err := json.MarshalIndent(smap, "", "\t")
	if err != nil {
		return ""
	}
	return string(buf)

}


func createKvMap(ctx *ecoContext, key string, kvs []*KeyValue) KeyValueMap {
	ctx.logger.signature("createKvMap", key, len(kvs))
	ctx.inc()
	defer ctx.dec()

	kvMap := KeyValueMap{}

	for _, kv := range kvs {
		kvMap.add(ctx, key, kv)
	}

	return kvMap
}

package eco

import (
	"encoding/json"
	"fmt"
	"reflect"

	etcd "go.etcd.io/etcd/client/v3"
)

type Decoder interface {
	Decode(*ecoContext, *KeyValueTree, string, reflect.Value) error
}
type MapDecoder struct{}
type MainDecoder struct{}
type ScalarDecoder[U any] struct{}
type SliceDecoder struct{}
type StructDecoder struct{}
type DecoderMap map[reflect.Kind]Decoder

type DecoderMapData map[string][]byte
type DecoderSliceData map[int][]byte

var decoderMap DecoderMap = DecoderMap{
	reflect.Bool:    &ScalarDecoder[bool]{},
	reflect.Int:     &ScalarDecoder[int]{},
	reflect.Int8:    &ScalarDecoder[int8]{},
	reflect.Int16:   &ScalarDecoder[int16]{},
	reflect.Int32:   &ScalarDecoder[int32]{},
	reflect.Int64:   &ScalarDecoder[int64]{},
	reflect.Uint:    &ScalarDecoder[uint]{},
	reflect.Uint8:   &ScalarDecoder[uint8]{},
	reflect.Uint16:  &ScalarDecoder[uint16]{},
	reflect.Uint32:  &ScalarDecoder[uint32]{},
	reflect.Uint64:  &ScalarDecoder[uint64]{},
	reflect.Float32: &ScalarDecoder[float32]{},
	reflect.Float64: &ScalarDecoder[float64]{},
	reflect.String:  &ScalarDecoder[string]{},
	reflect.Array:   &SliceDecoder{},
	reflect.Slice:   &SliceDecoder{},
	reflect.Map:     &MapDecoder{},
	reflect.Struct:  &StructDecoder{},
}

func (d *MainDecoder) Decode(ctx *ecoContext, kvs *KeyValueTree, key string, rv reflect.Value) error {
	// Get the decoder from the decoder map, if it exists
	pKind, err := getUnderlyingPointerKind(rv)
	if err != nil {
		return err
	}
	decoder, is := decoderMap[pKind]
	if !is {
		return fmt.Errorf("Unsupported pointer type (kind=%s)", pKind.String())
	}

	return decoder.Decode(ctx, kvs, key, rv)
}

// Decode the kvs at the key into a map
// The map is specified by a pointer-to-a-pointer-to-map (ppm). The pointer to
// the map (pm) is assumed to be nil, and it is this function's job to allocate
// the map. and then assign the address of the allocated map to the ppm. This is
// how a function can allocation a value and 'return' it via an incoming
// parameter.
//
// @note it might not make a lot of sense to deal with double indirection just
// to support allocating a new value to an incoming parameter. It's what
// json.Unmarshal() does but that's because json.Unmarshal() also supports
// unmarshalling into an existing data structure, as well as allocation. If a
// a function is always allocating and unmarshalling into a new object, it might
// make sense to just return it.
//
// ctx	Context for logging+etcd client
// kvs  key-value pairs which comprise the data
// key  key that prefixes all map keys
// rv   reflection pointer-to-pointer-to-map
func (d *MapDecoder) Decode(ctx *ecoContext, kvTree *KeyValueTree, key string, ppMap reflect.Value) error {
	ctx.logger.signature("MapDecoder.Decode()", key, reflect.TypeOf(ppMap))
	ctx.inc()
	defer ctx.dec()

	// get the reflect.Type for the map to allocate
	typ, err := getUnderlyingMapType(ppMap.Type())
	if err != nil {
		return err
	}

	// allocate the new map + save the value type
	rMap := reflect.MakeMap(typ)
	typValue := rMap.Type().Elem()

	// // get the map data from the kvs+key
	// mapData := getMapData(kvs, key)

	// populate the new map with the data
	for k := range kvTree.Children {
		ctx.logger.Infof("Decoding %s ...", k)
		// create a new map item
		pnew := reflect.New(typValue)
		decoder := decoderMap[typValue.Kind()]
		subkey := fmt.Sprintf("%s/%s", key, k)
		fmt.Printf("subkey=%v\n", subkey)
		err := decoder.Decode(ctx, kvTree, subkey, pnew)
		if err != nil {
			return err
		}
		// // get the address of the new element and unmarshal the mapData value
		// addr := pnew.Elem().Addr()
		// i := addr.Interface()
		// err := json.Unmarshal(v, i)
		// if err != nil {
		// 	return err
		// }

		// Create reflect.Value for mapData key and add key+val to new map
		rk := reflect.ValueOf(k)
		rMap.SetMapIndex(rk, pnew.Elem())
	}

	// Create a new map pointer and assign the new map to it
	pMap := reflect.New(typ)
	pMap.Elem().Set(rMap)

	// assign the new map to the rv
	ppMap.Elem().Set(pMap)

	// done :)
	return nil
}

func (d *ScalarDecoder[U]) Decode(ctx *ecoContext, kvTree *KeyValueTree, key string, rv reflect.Value) error {
	ctx.logger.signature("ScalarDecoder.Decode()", kvTree.Name, key, rv.Kind().String())
	ctx.inc()
	defer ctx.dec()

	if kvTree == nil {
		return fmt.Errorf("no data to decode")
	}
	data := kvTree.Value
	ctx.logger.Infof("data=%#v", data)
	i := rv.Interface()
	err := json.Unmarshal(data, i)
	return err
}

func (d *SliceDecoder) Decode(ctx *ecoContext, kvTree *KeyValueTree, key string, rv reflect.Value) error {
	ctx.logger.signature("SliceDecoder.Decode()", kvTree.Name, key, rv.Kind().String())
	ctx.inc()
	defer ctx.dec()

	// sliceData := getSliceData(kvs, key)
	maxIndex := kvTree.Children.MaxIndex()
	ctx.logger.Infof("%#v", kvTree)
	ctx.logger.Infof("maxIndex=%d", maxIndex)
	typSlice, err := getUnderlyingSliceType(rv)
	if err != nil {
		return err
	}
	len := maxIndex + 1
	cap := maxIndex + 1
	ctx.logger.Infof("Making slice: len=%d cap=%d", len, cap)
	rvSlice := reflect.MakeSlice(typSlice, int(len), int(cap))

	// Unmarshal all the elements
	for childKey, childTree := range kvTree.Children {
		ctx.logger.Infof("Decoding %s ...", childKey)
		// Check if the childKey is a uint
		keyString := KeyString(childKey)
		i, is := keyString.Index()
		if !is {
			ctx.logger.Infof("Key not a uint - skipping key (%s)", childKey)
			continue
		}
		// Get a pointer to the slice element at the specified index
		el := rvSlice.Index(int(i))
		addr := el.Addr()
		subkey := fmt.Sprintf("%s/%d", key, i)
		ctx.logger.Infof("subkey=%s ...", subkey)
		decoder := decoderMap[el.Kind()]
		ctx.logger.Infof("delgating to decoder: type=%s", reflect.TypeOf(decoder))
		decoder.Decode(ctx, childTree, subkey, addr)
		ctx.logger.commentf("subkey (%s) decoded", subkey)
		// pEl := addr.Interface()

		// // Unmarshal the specified data into the element pointer
		// err := json.Unmarshal(data, pEl)
		// if err != nil {
		// 	return err
		// }
	}

	// Make a slice pointer + assign the new slice to the pointer's Elem()
	rvNew := reflect.New(typSlice)
	rvNew.Elem().Set(rvSlice)

	// Assign the new slice pointer to the incoming rv
	rv.Elem().Set(rvNew)

	return nil
}

func (d *StructDecoder) Decode(ctx *ecoContext, kvs *KeyValueTree, key string, rv reflect.Value) error {
	return nil
}

func Decode(ctx *ecoContext, cli *EtcdClient, key string, pp any) error {
	ctx.logger.signature("decode", key, reflect.TypeOf(pp).Elem())
	ctx.inc()
	defer ctx.dec()

	// Confirm p is a 'normal pointer', ie a pointer that is not a pointer-to-a-pointer
	if !isValidPointer(pp) {
		return fmt.Errorf("p must be a pointer-to-a-pointer, with an element type that is not a pointer (kind=%s)",
			reflect.TypeOf(pp).Kind().String())
	}

	// Get kvs from etcd
	op := etcd.OpGet(key, etcd.WithPrefix())
	txn := cli.Txn(ctx)
	resp, err := txn.Then(op).Commit()
	if err != nil {
		return nil
	}

	rangeResp := resp.Responses[0].GetResponseRange()
	etcdKvs := rangeResp.Kvs
	kvs := createKvSlice(etcdKvs)
	kvTree := createKvTree(ctx, key, kvs, key)
	decoder := MainDecoder{}
	rv := reflect.ValueOf(pp)
	err = decoder.Decode(ctx, kvTree, key, rv)
	return err

	// // Simple objects are easy to deal with. Just use json.Unmarhsal()
	// if isScalar(ty.Elem().Kind()) {
	// 	// Get object from etcd + make sure there's only 1
	// 	resp, err := cli.Client.Get(ctx, key)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(resp.Kvs) != 1 {
	// 		return fmt.Errorf("expected one key; got %d", len(resp.Kvs))
	// 	}

	// 	// Unmarshal the result
	// 	getVal := resp.Kvs[0].Value
	// 	ctx.logger.Infof("getVal()=%v (%s)", getVal, getVal)
	// 	err = json.Unmarshal(getVal, i)
	// 	if err != nil {
	// 		ctx.logger.Errorf("Unmarshalling error: %s (%#v)", err.Error(), getVal)
	// 		return err
	// 	}
	// 	// @note - should we return here?
	// 	return nil
	// }

	// Some non-simple type are supported. The rest of the function checks for them.
	// Note - we want the type of the underlying element, not the type of the pointer
	// kindElem := getTypeKind(ctx, ty.Elem())

	// switch kindElem {
	// case SimpleMap: return decodeMap(ctx, cli, key, i)
	// case SimpleSlice: return decodeSlice(ctx, cli, key, i)
	// case SimpleStruct: return decodeStruct(ctx, cli, key, i)

	// default:
	// 	return errors.New("unsupported type")
	// }
}

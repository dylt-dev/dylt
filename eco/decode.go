package eco

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/dylt-dev/dylt/common"
)

type Decoder interface {
	Decode(*common.EcoContext, *KeyValueTree, string, reflect.Value) error
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

func (d *MainDecoder) Decode(ctx *common.EcoContext, kvs *KeyValueTree, key string, rv reflect.Value) error {
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
func (d *MapDecoder) Decode(ctx *common.EcoContext, kvTree *KeyValueTree, key string, rv reflect.Value) error {
	ctx.Logger.Signature("MapDecoder.Decode()", key, reflect.TypeOf(rv))
	ctx.Inc()
	defer ctx.Dec()

	pMap, is := common.CreateOrGetMap(ctx, rv)
	if !is {
		return fmt.Errorf("Unable to create or get map for ... reasons")	
	}
	if pMap == nil {
		return fmt.Errorf("nil map -- this shouldn't happen")
	}
	rvpMap := reflect.ValueOf(pMap)
	if rvpMap.IsNil() {
		return fmt.Errorf("rvpMap.IsNil() == true")
	}

	// get the reflect.Type for the underlying map
	typMap, err := common.GetUnderlyingMapType(ctx, rv)
	if err != nil {
		return err
	}
	ctx.Logger.Infof("map type=map[%s]%s", typMap.Key(), typMap.Elem())

	typKey := typMap.Key()
	typValue := typMap.Elem()
	// populate the new map with the data
	for keyString, childTree := range kvTree.Children {
		ctx.Logger.Infof("Decoding %s ...", keyString)
		ctx.Logger.Infof("childTree.Value=%#v", childTree.Value)
		// // create a new map item
		pnew := reflect.New(typValue)
		decoder := decoderMap[typValue.Kind()]
		ctx.Logger.Infof("decoder type=%v\n", reflect.TypeOf(decoder))
		elemName := keyString.ElementName(key)
		subkey := fmt.Sprintf("%s/%s", key, elemName)
		ctx.Logger.Infof("subkey=%v\n", subkey)
		err := decoder.Decode(ctx, childTree, subkey, pnew)
		if err != nil {
			return err
		}
		// // get the address of the new element and unmarshal the mapData value
		// addr := pnew.Elem().Addr()
		// i := addr.Interface()
		// err := json.Unmarshal(v, i)
		// if err != nil {
		// 	return err
		// // }

		// // Create reflect.Value for mapData key and add key+val to new map
		prvKey := reflect.New(typKey)
		pKey := prvKey.Interface()
		err = common.UnmarshalMapKey(elemName, pKey)
		if err != nil {
			return err
		}
		ctx.Logger.Infof("rvpMap type=map[%s]%s", rvpMap.Elem().Type().Key(), rvpMap.Elem().Type().Elem())
		rvpMap.Elem().SetMapIndex(prvKey.Elem(), pnew.Elem())
	}

	// Create a new map pointer and assign the new map to it
	// pMap := reflect.New(typ)
	// pMap.Elem().Set(rMap)

	// assign the new map to the rv
	// ppMap.Elem().Set(pMap)

	// done :)
	return nil
}

func (d *ScalarDecoder[U]) Decode(ctx *common.EcoContext, kvTree *KeyValueTree, key string, rv reflect.Value) error {
	ctx.Logger.Signature("ScalarDecoder.Decode()", key, rv.Kind().String())
	ctx.Inc()
	defer ctx.Dec()

	if kvTree == nil {
		return fmt.Errorf("no data to decode")
	}
	data := kvTree.Value
	ctx.Logger.Infof("data=%#v", data)
	i := rv.Interface()
	err := json.Unmarshal(data, i)
	return err
}

func (d *SliceDecoder) Decode(ctx *common.EcoContext, kvTree *KeyValueTree, key string, rv reflect.Value) error {
	ctx.Logger.Signature("SliceDecoder.Decode()", key, rv.Kind().String())
	ctx.Inc()
	defer ctx.Dec()

	// sliceData := getSliceData(kvs, key)
	maxIndex := kvTree.Children.MaxIndex()
	ctx.Logger.Infof("%#v", kvTree)
	ctx.Logger.Infof("maxIndex=%d", maxIndex)
	typSlice, err := getUnderlyingSliceType(rv)
	if err != nil {
		return err
	}
	len := maxIndex + 1
	cap := maxIndex + 1
	ctx.Logger.Infof("Making slice: len=%d cap=%d", len, cap)
	rvSlice := reflect.MakeSlice(typSlice, int(len), int(cap))

	// Unmarshal all the elements
	for childKey, childTree := range kvTree.Children {
		ctx.Logger.Infof("Decoding %s ...", childKey)
		// Check if the childKey is a uint
		keyString := KeyString(childKey)
		i, is := keyString.Index()
		if !is {
			ctx.Logger.Infof("Key not a uint - skipping key (%s)", childKey)
			continue
		}
		// Get a pointer to the slice element at the specified index
		el := rvSlice.Index(int(i))
		addr := el.Addr()
		subkey := fmt.Sprintf("%s/%d", key, i)
		ctx.Logger.Infof("subkey=%s ...", subkey)
		decoder := decoderMap[el.Kind()]
		ctx.Logger.Infof("delgating to decoder: type=%s", reflect.TypeOf(decoder))
		decoder.Decode(ctx, childTree, subkey, addr)
		ctx.Logger.Commentf("subkey (%s) decoded", subkey)
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

func (d *StructDecoder) Decode(ctx *common.EcoContext, kvTree *KeyValueTree, key string, rv reflect.Value) error {
	ctx.Logger.Signature("StructDecoder.Decode()", key, reflect.TypeOf(rv))
	ctx.Inc()
	defer ctx.Dec()

	pStruct, is := common.CreateOrGetStruct(ctx, rv)
	if !is {
		return fmt.Errorf("Unable to create or get struct for ... reasons")	
	}

	// get the reflect.Type for the underlying struct to iterate over
	ctx.Logger.Comment("Getting underlying struct type ...")
	typStruct, err := common.GetUnderlyingStructType(ctx, rv)
	if err != nil {
		return err
	}
	ctx.Logger.Infof("struct type=%s", typStruct.Name())

	ctx.Logger.Comment("Dumping child keys ...")
	for childKey, _ := range kvTree.Children {
		ctx.Logger.Infof("childKey=%s", childKey)
	}
	
	ctx.Logger.Comment("Iterating over struct fields")
	rvpStruct := common.ToRv(pStruct)
	for field := range typStruct.Fields() {
		if field.Tag != "" {
			ctx.Logger.Infof("%-20s %-20s", field.Name, field.Type)
		} else {
			ctx.Logger.Infof("%-20s %-20s (%s)", field.Name, field.Type, field.Tag)
		}
		childKeyName := getFieldKey(field)
		childKey := createKeyString(key, childKeyName)
		ctx.Logger.Infof("childKey=%s", childKey)
		kvChildTree, is := kvTree.Children[childKey]
		if !is {
			ctx.Logger.Infof("Field not found in KVs: %s", childKey)
		}
		// Decode LV value into struct field
		decoder := decoderMap[field.Type.Kind()]
		structField := rvpStruct.Elem().FieldByName(field.Name)
		addr := structField.Addr()
		decoder.Decode(ctx, kvChildTree, childKeyName, addr) 
		// common.UnmarshalStructField(pStruct, field.Name, kv.Value)
		// else {
		// 	addr := rvStruct.FieldByName(string(childKey)).Addr()
		// 	decoder := decoderMap[field.Type.Kind()]
		// 	decoder.Decode(ctx, kvTree, kv.Name, addr)
		// }
	}


	// populate the new map with the data
	// for k, childTree := range kvTree.Children {
	// 	ctx.Logger.Infof("Decoding %s ...", k)
	// }

	return nil
}


func Decode(ctx *common.EcoContext, cli *EtcdClient, key string, i any) error {
	ctx.Logger.Signature("decode", key, reflect.TypeOf(i))
	ctx.Inc()
	defer ctx.Dec()

	// Confirm p is a valid pointer
	if !isValidPointer(i) {
		return fmt.Errorf("p must be a non-nil pointer, not %s",
			reflect.TypeOf(i).Kind().String())
	}

	// Get etcd KVs from cluster
	ctx.Logger.Comment("Getting KVs from cluster ...")
	etcdKvs, err := getEtcdKvs(ctx, cli, key)
	if err != nil {
		return err
	}
	ctx.Logger.Info("Done")
	ctx.Logger.Info()

	// Create kvTree
	ctx.Logger.Comment("Creating KV slice and KV tree ...")
	kvs := createKvSlice(etcdKvs)
	kvTree := createKvTree(ctx, key, kvs)
	ctx.Logger.Info("Done")
	ctx.Logger.Info()

	// Create reflect.Value for i
	rv := common.ToRv(i)
	
	// Decode using the top-level Decoder
	ctx.Logger.Comment("Decoding ...")
	decoder := MainDecoder{}
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
	// 	ctx.Logger.Infof("getVal()=%v (%s)", getVal, getVal)
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


func GetStructFieldKey (fld reflect.StructField) string {
	var key  string
	tag, is := fld.Tag.Lookup("eco")
	if is {
		key = tag
	} else {
		key = fld.Name
	}

	return key
}

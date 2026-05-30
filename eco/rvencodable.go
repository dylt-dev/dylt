package eco

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/dylt-dev/dylt/common"
)

type Encodable interface {
	Data() map[string]any
}

type RvEncodable reflect.Value
type RvMap reflect.Value
type RvScalar reflect.Value
type RvSlice reflect.Value
type RvStruct reflect.Value

func NewRvMap(a any) RvMap {
	rv := common.Reflect(a)
	if rv.Kind() != reflect.Map {
		panic(fmt.Sprintf("NewRvMap: expecting map (%s)", rv.Kind()))
	}

	return RvMap(rv)
}

func (rvs RvMap) Data() map[string]any {
	data := map[string]any{}
	rv := reflect.Value(rvs)

	iter := rv.MapRange()
	for iter.Next() {
		k := fmt.Sprint(iter.Key().Interface())
		v := iter.Value().Interface()
		data[k] = v
	}

	return data
}

func NewRvScalar(a any) RvScalar {
	rv := common.Reflect(a)
	if common.NewFlavor(rv.Kind()) != common.Scalar {
		panic(fmt.Sprintf("NewRvScalar: expecting scalar (%s)", rv.Kind()))
	}

	return RvScalar(rv)
}

func (rvs RvScalar) Data() map[string]any {
	data := map[string]any{}
	rv := reflect.Value(rvs)
	data[""] = rv.Interface()

	return data
}

func NewRvSlice(a any) RvSlice {
	rv := common.Reflect(a)
	if rv.Kind() != reflect.Slice {
		panic(fmt.Sprintf("NewRvSlice: expecting slice (%s)", rv.Kind()))
	}

	return RvSlice(rv)
}

func (rvs RvSlice) Data() map[string]any {
	data := map[string]any{}
	rv := reflect.Value(rvs)

	for i := range rv.Len() {
		k := strconv.Itoa(i)
		v := rv.Index(i).Interface()
		data[k] = v
	}

	return data
}

func (rvs RvStruct) Data() map[string]any {
	data := map[string]any{}
	rv := reflect.Value(rvs)

	for structField := range rv.Fields() {
		k := GetStructFieldKey(structField)
		field := rv.FieldByName(structField.Name)
		v := field.Interface()
		data[k] = v
	}

	return data
}

func NewRvEncodable(a any) RvEncodable {
	rv := common.Reflect(a)
	return RvEncodable(rv)
}

func (rve RvEncodable) Data() map[string]any {
	return rve.subType().Data()
}

func (rve RvEncodable) Encode(ctx *common.EcoContext, ks KeyString) []KeyValue {
	ctx.Signature("RvEncodable.Encode", ks)
	ctx.Inc()
	defer ctx.Dec()

	kvs := []KeyValue{}

	if rve.IsNil() {
		kv := KeyValue{ks, []byte(nil)}
		return []KeyValue{kv}
	}

	rve, err := rve.walkPointer()
	if err != nil {
		kv := KeyValue{ks, []byte(nil)}
		return []KeyValue{kv}
	}

	if rve.Flavor() == common.Scalar {
		var val []byte
		if rve.Kind() == reflect.String {
			// val = []byte(reflect.Value(rve).String())
			val, err = json.Marshal(rve.Interface())
			if err != nil {
				panic(err)
			}
		} else {
			val, err = json.Marshal(rve.Interface())
			if err != nil {
				panic(err)
			}
		}
		kv := KeyValue{ks, val}
		return []KeyValue{kv}
	}

	for k, v := range rve.Data() {
		ctx.Infof("ks=%s", ks)
		ksc := ks.AddSegment(k)
		ctx.Infof("ks=%s ksc=%s", ks, ksc)
		if v != nil {
			rvec := NewRvEncodable(v)
			kvsc := rvec.Encode(ctx, ksc)
			kvs = append(kvs, kvsc...)
		}
	}

	return kvs
}

func (rve RvEncodable) Elem() RvEncodable {
	rv := reflect.Value(rve)
	return RvEncodable(rv.Elem())
}

func (rve RvEncodable) Flavor() common.Flavor {
	return common.NewFlavor(rve.Kind())
}

func (rve RvEncodable) Interface() any {
	rv := reflect.Value(rve)
	return rv.Interface()
}

func (rve RvEncodable) IsNil() bool {
	switch rve.Flavor() {
	case common.Scalar,
		common.Struct:
		return false

	case common.Map,
		common.Slice:
		{
			return reflect.Value(rve).IsNil()
		}

	case common.Pointer:
		{
			if reflect.Value(rve).IsNil() {
				return true
			}
			return rve.Elem().IsNil()
		}

	default:
		panic(fmt.Sprintf("unexpected type (%s)", rve.Kind()))
	}
}

func (rve RvEncodable) Kind() reflect.Kind {
	return rve.Type().Kind()
}

func (rve RvEncodable) Type() reflect.Type {
	return reflect.Value(rve).Type()
}

func (rve RvEncodable) subType() Encodable {
	switch rve.Flavor() {
	case common.Map:
		return RvMap(rve)
	case common.Scalar:
		return RvScalar(rve)
	case common.Slice:
		return RvSlice(rve)
	case common.Struct:
		return RvStruct(rve)
	default:
		panic(fmt.Errorf("unsupported type (%s)", rve.Type()))
	}
}

func (rve RvEncodable) walkPointer() (RvEncodable, error) {
	flavor := rve.Flavor()
	switch flavor {
	case common.Map,
		common.Scalar,
		common.Slice,
		common.Struct:
		return rve, nil
	case common.Pointer:
		{
			rv := reflect.Value(rve)
			if rv.IsNil() {
				return rve, nil
			}
			rve = RvEncodable(rv.Elem())
			return rve.walkPointer()
		}
	default:
		{
			return rve, fmt.Errorf("unsupported type (%s)", rve.Kind())
		}
	}
}

// func (rve RvEncodable) zero () []byte {
// 	rv := reflect.Value(rve)
// 	rvZero := reflect.Zero(rv.Type())
// 	buf, err := json.Marshal(rvZero.Interface())
// 	if err != nil {
// 		panic(err)
// 	}
// 	return buf
// }

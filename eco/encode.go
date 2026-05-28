package eco

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dylt-dev/dylt/common"
)


type Encodeable interface {
	Data () map[string]any
}


type RvMap reflect.Value
type RvSlice reflect.Value
type RvStruct reflect.Value


func NewRvMap (a any) RvMap {
	rv := common.Reflect(a)
	if rv.Kind() != reflect.Map {
		panic(fmt.Sprintf("NewRvMap: expecting map (%s)", rv.Kind()))
	}

	return RvMap(rv)
}


func (rvs RvMap) Data () map[string]any {
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


func NewRvSlice (a any) RvSlice {
	rv := common.Reflect(a)
	if rv.Kind() != reflect.Slice {
		panic(fmt.Sprintf("NewRvSlice: expecting slice (%s)", rv.Kind()))
	}

	return RvSlice(rv)
}


func (rvs RvSlice) Data () map[string]any {
	data := map[string]any{}
	rv := reflect.Value(rvs)

	for i := range rv.Len() {
		k := strconv.Itoa(i)
		v := rv.Index(i).Interface()
		data[k] = v
	}

	return data
}


func (rvs RvStruct) Data () map[string]any {
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



func Encode (ctx *common.EcoContext, a any) []KeyValue {
	ctx.Signature("Encode", reflect.TypeOf(a))
	ctx.Inc()
	defer ctx.Dec()
	
	return nil
}
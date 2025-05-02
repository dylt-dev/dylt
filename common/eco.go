package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/dylt-dev/dylt/color"
	etcd "go.etcd.io/etcd/client/v3"
)

func Encode(ctx *ecoContext, key string, i any) ([]etcd.Op, error) {
	ctx.printf("Encoding key=%s (type=%s) ...)\n", key, fullTypeName(reflect.TypeOf(i)))
	if _, ok := i.(reflect.Value); ok {
		ctx.println("* arg i is of type reflect.Value; did you mean to call i.Interface()?")
	}
	ctx.inc()
	defer ctx.dec()

	var ty reflect.Type = reflect.TypeOf(i)
	// var _ reflect.Value = reflect.ValueOf(i)
	var ops = []etcd.Op{}
	var kind reflect.Kind = ty.Kind()
	var val reflect.Value = reflect.ValueOf(i)
	var err error
	ctx.printf("Switching on kind=%s ...\n", kind.String())
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		opsPut, err := encodeDefault(ctx, key, val)
		if err != nil {
			return nil, err
		}
		ops = append(ops, opsPut...)
	case reflect.Array:
		opsArray, err := encodeSlice(ctx, key, val)
		if err != nil {
			return nil, err
		}
		ops = append(ops, opsArray...)
	case reflect.Map:
		opsMap, err := encodeMap(ctx, key, val)
		if err != nil {
			return nil, err
		}
		ops = append(ops, opsMap...)
	case reflect.Slice:
		opsSlice, err := encodeSlice(ctx, key, val)
		if err != nil {
			return nil, err
		}
		ops = append(ops, opsSlice...)
	case reflect.Struct:
		opsStruct, err := encodeStruct(ctx, key, val)
		if err != nil {
			return nil, err
		}
		ops = append(ops, opsStruct...)
	default:
		err = errors.New("unsupported")
		return nil, err
	}
	return ops, nil
}

type ecoContext struct {
	context.Context
	level int
}

func newEcoContext() *ecoContext {
	return &ecoContext{Context: context.Background(), level: 0}
}

func (ctx *ecoContext) dec() *ecoContext {
	ctx.level--
	return ctx
}

func (ctx *ecoContext) inc() *ecoContext {
	ctx.level++
	return ctx
}

func (ctx *ecoContext) indent() string {
	return strings.Repeat("  ", ctx.level)
}

func (ctx *ecoContext) printf(format string, a ...any) (int, error) {
	format = fmt.Sprintf("%s%s", ctx.indent(), format)
	return fmt.Printf(format, a...)
}

func (ctx *ecoContext) println(a ...any) (int, error) {
	args := fmt.Sprintln(a...)
	return fmt.Printf("%s%s", ctx.indent(), args)
}

type kind uint

const (
	Invalid kind = iota
	Bool
	Number
	String
	SimpleArray
	SimpleInterface
	SimpleMap
	SimplePointer
	SimpleSlice
	SimpleStruct
)

func (k kind) String() string {
	switch k {
	case Invalid:
		return "Invalid"
	case Bool:
		return "Bool"
	case Number:
		return "Number"
	case String:
		return "String"
	case SimpleArray:
		return "SimpleArray"
	case SimpleInterface:
		return "SimpleInterface"
	case SimpleMap:
		return "SimpleMap"
	case SimplePointer:
		return "SimplePointer"
	case SimpleSlice:
		return "SimpleSlice"
	case SimpleStruct:
		return "SimpleStruct"
	default:
		return fmt.Sprintf("Unknown kind :%d", k)
	}
}

func arrayKind(ctx *ecoContext, ty reflect.Type) kind {
	ctx.printf("In arrayKind(): ty=%s\n", fullTypeName(ty))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Array {
		ctx.println("type is not a array; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	ctx.printf("Checking element type (%s) ... ", fullTypeName(tyElem))
	if isTypeSimple(ty.Elem()) {
		ctx.println("element type is simple; returning SimpleArray", fullTypeName(tyElem))
		return SimpleArray
	}
	ctx.println("conditions were not met; returning Invalid")
	return Invalid
}

func encodeDefault(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.printf("encodeDefault() - key=%s\n", key)
	ctx.inc()
	defer ctx.dec()

	i := val.Interface()
	j, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	opPut := etcd.OpPut(key, string(j))

	return []etcd.Op{opPut}, nil
}

func encodeMap(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.printf("encodeMap() - key=%s\n", key)
	ctx.inc()
	defer ctx.dec()

	ty := val.Type()
	ctx.printf("Confirming type (%s) is SimpleMap ... ", fullTypeName(ty))
	if getTypeKind(ctx, ty) != SimpleMap {
		ctx.println("incorrect.")
		return nil, fmt.Errorf("expecting SimpleMap; got %s", fullTypeName(ty))
	}

	ctx.println("confirmed.")
	ctx.println("Encoding keys and values ...")
	var ops = []etcd.Op{}
	mapIter := val.MapRange()
	for mapIter.Next() {
		miKey := fmt.Sprintf("%v", mapIter.Key().Interface())
		elKey := filepath.Join(key, string(miKey))
		elVal := mapIter.Value()
		elOps, err := Encode(ctx, elKey, elVal.Interface())
		if err != nil {
			return nil, err
		}
		ops = append(ops, elOps...)
	}

	return ops, nil
}

func encodeSlice(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.printf("encodeSlice() - key=%s", key)
	ctx.inc()
	defer ctx.dec()

	ty := val.Type()
	if getTypeKind(ctx, ty) != SimpleSlice {
		return nil, fmt.Errorf("expecting SimpleSlice; got %s", fullTypeName(ty))
	}

	j, err := json.Marshal(val.Interface())
	if err != nil {
		return nil, err
	}
	op := etcd.OpPut(key, string(j))

	return []etcd.Op{op}, nil
}

func encodeStruct(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	ctx.printf("encodeStruct() - key=%s\n", key)
	ctx.inc()
	defer ctx.dec()

	ty := val.Type()
	ctx.printf("Confirming type (%s) is SimpleStruct ... ", fullTypeName(ty))
	if getTypeKind(ctx, ty) != SimpleStruct {
		ctx.println("incorrect.")
		return nil, fmt.Errorf("expecting SimpleStruct; got %s", fullTypeName(ty))
	}

	ctx.println("confirmed.")
	ctx.println("Encoding fields ...")
	var ops = []etcd.Op{}
	for i := range ty.NumField() {
		sf := ty.Field(i)
		sfName := getFieldKey(sf)
		sfVal := val.Field(i)
		sfKey := filepath.Join(key, sfName)
		sfOps, err := Encode(ctx, sfKey, sfVal.Interface())
		if err != nil {
			return nil, err
		}
		ops = append(ops, sfOps...)
	}

	return ops, nil
}

func fieldNameMap (i any) (map[string]reflect.Value, error) {
	var tyElem reflect.Type
	var valElem reflect.Value
	ty := reflect.TypeOf(i)
	if ty.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("expecting pointer; got %s", fullTypeName(ty))
	}
	tyElem = ty.Elem()
	if tyElem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expecting pointer to struct; got %s", tyElem.Kind().String())
	}
	val := reflect.ValueOf(i)
	valElem = val.Elem()
	
	fieldNameMap := map[string]reflect.Value{}
	for i := range tyElem.NumField() {
		tyField := tyElem.Field(i)
		fieldKey := getFieldKey(tyField)
		fieldName := tyField.Name
		valField := valElem.FieldByName(fieldName)
		fieldNameMap[fieldKey] = valField
	}

	return fieldNameMap, nil
}


func fullTypeName(ty reflect.Type) string {
	pkgPath := ty.PkgPath()
	typeName := ty.Name()
	if typeName == "" {
		typeName = "(anon)"
	}

	if pkgPath == "" {
		return typeName
	}

	if filepath.Dir(pkgPath) == "github.com/dylt-dev/dylt" {
		pkgPath = filepath.Base(pkgPath)
	}

	return fmt.Sprintf("%s.%s", pkgPath, typeName)
}

func getFieldKey(sf reflect.StructField) string {
	tagValue, ok := sf.Tag.Lookup("eco")
	var fieldName string
	if ok {
		fieldName = tagValue
	} else {
		fieldName = sf.Name
	}

	return fieldName
}

func getFieldValue(val reflect.Value) (string, error) {
	var s string
	kind := val.Type().Kind()
	switch kind {
	case reflect.Bool:
		s = strconv.FormatBool(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s = strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32:
		s = strconv.FormatFloat(val.Float(), 'f', -1, 32)
	case reflect.Float64:
		s = strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case reflect.String:
		s = string(val.String())
	default:
		return "", fmt.Errorf("unsupported kind: %s", kind.String())
	}

	return s, nil
}

func getKind(ctx *ecoContext, i any) kind {
	// ty := reflect.TypeOf(i)
	// if fullTypeName(ty) == "reflect.Type" {
	// 	fmt.Println("Warning - GetKind() called with reflect.Type(). Did you mean GetTypeKind()?")
	// }
	return getTypeKind(ctx, reflect.TypeOf(i))
}

func getTypeKind(ctx *ecoContext, ty reflect.Type) kind {
	reflectKind := ty.Kind()
	// fmt.Printf("ty=%s reflectKind=%s\n", fullTypeName(ty), reflectKind.String())
	switch reflectKind {
	case reflect.Bool:
		return Bool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32,
		reflect.Float64:
		return Number
	case reflect.String:
		return String
	case reflect.Array:
		return arrayKind(ctx, ty)
	case reflect.Interface:
		return interfaceKind(ctx, ty)
	case reflect.Map:
		return mapKind(ctx, ty)
	case reflect.Pointer:
		return pointerKind(ctx, ty)
	case reflect.Slice:
		return sliceKind(ctx, ty)
	case reflect.Struct:
		return structKind(ctx, ty)
	case reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.Func, reflect.UnsafePointer:
		return Invalid
	default:
		return Invalid
	}
}

func interfaceKind(ctx *ecoContext, ty reflect.Type) kind {
	ctx.printf("In interfaceKind(): ty=%s\n", fullTypeName(ty))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Interface {
		ctx.println("type is not an interface; returning Invalid")
		return Invalid
	}

	return Invalid
}

func isSimple(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

func isTypeSimple(ty reflect.Type) bool {
	return isSimple(ty.Kind())
}

func mapKind(ctx *ecoContext, ty reflect.Type) kind {
	sig := fullTypeName(ty)
	ctx.printf("%s(%s)\n", highlight("mapKind"), lowlight(sig))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Map {
		ctx.println("type is not a map; returning Invalid")
		return Invalid
	}

	ctx.println(lowlight(fmt.Sprintf("Checking key (%s) ...", fullTypeName(ty.Key()))))
	if !isTypeSimple(ty.Key()) {
		ctx.printf("%-32s: %-32s; %s\n", "key", "non-simple", highlight("returning Invalid"))
		return Invalid
	}
	ctx.printf("%-16s %-16s; continuing\n", "key", "simple")

	tyElem := ty.Elem()
	ctx.println(lowlight(fmt.Sprintf("Checking element type (%s) ...", fullTypeName(tyElem))))
	ctx.inc()
	kindElem := getTypeKind(ctx, tyElem)
	ctx.dec()
	if isTypeSimple(tyElem) {
		ctx.printf("%-16s %-16s; %s\n", "type", "simple", highlight("returning SimpleMap"))
		return SimpleMap
	}
	ctx.printf("type: not simple; continuing\n")
	ctx.printf("%-16s %-16s; continuing\n", "type", "not simple")

	ctx.println(lowlight(fmt.Sprintf("Checking element kind (%s) ...", kindElem.String())))
	if kindElem == SimpleMap ||
		kindElem == SimpleStruct ||
		kindElem == SimpleSlice {
		ctx.printf("kind: simple; returning SimpleMap\n", kindElem.String())
		return SimpleMap
	}
	ctx.printf("type: not simple; continuing\n")

	ctx.println(highlight("conditions were not met; returning Invalid"))
	return Invalid
}

func pointerKind(ctx *ecoContext, ty reflect.Type) kind {
	ctx.printf("In pointerKind(): ty=%s\n", fullTypeName(ty))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Pointer {
		ctx.println("type is not a pointer; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	ctx.printf("Checking pointer type (%s) ... ", fullTypeName(tyElem))
	if isTypeSimple(tyElem) {
		ctx.println("pointer type is simple; returning SimplePointer")
		return SimplePointer
	}

	ctx.println("conditions were not met; returning Invalid")
	return Invalid
}

func sliceKind(ctx *ecoContext, ty reflect.Type) kind {
	ctx.printf("In sliceKind(): ty=%s\n", fullTypeName(ty))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Slice {
		ctx.println("type is not a slice; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	ctx.printf("Checking element type (%s) ... ", fullTypeName(tyElem))
	if isTypeSimple(tyElem) {
		ctx.println("element type is simple; returning SimpleSlice")
		return SimpleSlice
	}

	ctx.println("conditions were not met; returning Invalid")
	return Invalid
}

func structKind(ctx *ecoContext, ty reflect.Type) kind {
	ctx.printf("In structKind(): ty=%s\n", fullTypeName(ty))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Struct {
		ctx.println("type is not a struct; returning Invalid")

		return Invalid
	}

	ctx.printf("%d field(s)\n", ty.NumField())
	for i := range ty.NumField() {
		sf := ty.Field(i)
		sfType := sf.Type
		ctx.printf("checking field '%s': type name = %s ...\n", sf.Name, fullTypeName(sfType))
		sfReflectKind := sfType.Kind()

		if isTypeSimple(sfType) {
			ctx.println("field type is simple; continuing")
			continue
		}

		if sfReflectKind == reflect.Map && mapKind(ctx, sfType) == SimpleMap {
			ctx.println("field type is SimpleMap; continuing")
			continue
		}

		if sfReflectKind == reflect.Slice && sliceKind(ctx, sfType) == SimpleSlice {
			ctx.println("field type is SimpleSlice; continuing")
			continue
		}

		if sfReflectKind == reflect.Struct && structKind(ctx, sf.Type) == SimpleStruct {
			ctx.println("field type is SimpleStruct; continuing")
			continue
		}

		return Invalid
	}

	ctx.println("All fields passed; returning SimpleStruct")
	return SimpleStruct
}

/* log styles */

func highlight(s string) string {

	var ss = color.Styledstring(s).Style("\033[1;97m")

	return string(ss)
}

func lowlight(s string) string {
	var ss = color.Styledstring(s).Fg(color.Sys.BrightBlack)

	return string(ss)
}

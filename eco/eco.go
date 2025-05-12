package eco

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/dylt-dev/dylt/color"
	etcd "go.etcd.io/etcd/client/v3"
)

func Encode(ctx *ecoContext, key string, i any) ([]etcd.Op, error) {
	logger.signature(ctx, "Encode", key, reflect.TypeOf(i))
	if _, ok := i.(reflect.Value); ok {
		logger.info(ctx, "arg i is of type reflect.Value; did you mean to call i.Interface()?")
	}
	ctx.inc()
	defer ctx.dec()

	var ty reflect.Type = reflect.TypeOf(i)
	// var _ reflect.Value = reflect.ValueOf(i)
	var ops = []etcd.Op{}
	var kind reflect.Kind = ty.Kind()
	var val reflect.Value = reflect.ValueOf(i)
	var err error

	// ctx.println(color.Styledstring("Check object type to confirm it can be encoded").Fg(color.X11.CornflowerBlue))
	logger.comment(ctx, "Check object type to confirm it can be encoded")
	logger.infof(ctx, "Switching on kind=%s ...", kind.String())
	switch kind {

	// simple case for simple types
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		ops, err = encodeDefault(ctx, key, val)

	case reflect.Array:
		ops, err = encodeSlice(ctx, key, val)

	case reflect.Map:
		ops, err = encodeMap(ctx, key, val)

	case reflect.Slice:
		ops, err = encodeSlice(ctx, key, val)

	case reflect.Struct:
		ops, err = encodeStruct(ctx, key, val)

	default:
		err = errors.New("unsupported")
	}

	if err != nil {
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

// func (ctx *ecoContext) indent() string {
// 	const tab = "  "
// 	return strings.Repeat(tab, ctx.level)
// }

// func (ctx *ecoContext) printf(format string, a ...any) (int, error) {
// 	format = fmt.Sprintf("%s%s", ctx.indent(), format)
// 	return fmt.Printf(format, a...)
// }

// func (ctx *ecoContext) println(a ...any) (int, error) {
// 	args := fmt.Sprintln(a...)
// 	return fmt.Printf("%s%s", ctx.indent(), args)
// }

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
	logger.signature(ctx, "arrayKind", ty)
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Array {
		logger.info(ctx, "type is not a array; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	logger.infof(ctx, "Checking element type (%s) ... ", fullTypeName(tyElem))
	if isTypeScalar(ty.Elem()) {
		logger.infof(ctx, "element type (%s) is scalar; returning SimpleArray", fullTypeName(tyElem))
		return SimpleArray
	}
	logger.info(ctx, "conditions were not met; returning Invalid")
	return Invalid
}

func encodeDefault(ctx *ecoContext, key string, val reflect.Value) ([]etcd.Op, error) {
	logger.signature(ctx, "encodeDefault", key, val.Type())
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
	logger.signature(ctx, "encodeMap", key, val.Type())
	ctx.inc()
	defer ctx.dec()

	ty := val.Type()
	logger.infof(ctx, "Confirming type (%s) is SimpleMap ...", fullTypeName(ty))
	if getTypeKind(ctx, ty) != SimpleMap {
		logger.comment(ctx, "incorrect.")
		return nil, fmt.Errorf("expecting SimpleMap; got %s", fullTypeName(ty))
	}

	logger.comment(ctx, "confirmed.")
	logger.info(ctx, "Encoding keys and values ...")
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
	logger.signature(ctx, "encodeSlice", key, val.Type())
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
	logger.signature(ctx, "encodeStruct", key, val.Type())
	ctx.inc()
	defer ctx.dec()

	ty := val.Type()
	logger.commentf(ctx, "Confirming type (%s) is SimpleStruct ...", fullTypeName(ty))
	if getTypeKind(ctx, ty) != SimpleStruct {
		logger.comment(ctx, "incorrect.")
		return nil, fmt.Errorf("expecting SimpleStruct; got %s", fullTypeName(ty))
	}

	logger.info(ctx, "confirmed.")
	logger.info(ctx, "Encoding fields ...")
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

func fieldNameMap(i any) (map[string]reflect.Value, error) {
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
	logger.signature(ctx, "interfaceKind", ty)
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Interface {
		logger.info(ctx, "type is not an interface; returning Invalid")
		return Invalid
	}

	return Invalid
}

func isScalar(kind reflect.Kind) bool {
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

func isTypeScalar(ty reflect.Type) bool {
	return isScalar(ty.Kind())
}

func mapKind(ctx *ecoContext, ty reflect.Type) kind {
	sig := fullTypeName(ty)
	logger.infof(ctx, "%s(%s)", highlight("mapKind"), lowlight(sig))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Map {
		logger.info(ctx, "type is not a map; returning Invalid")
		return Invalid
	}

	logger.infof(ctx, "%-70s", lowlight(fmt.Sprintf("checking key (%s) ...", fullTypeName(ty.Key()))))
	if !isTypeScalar(ty.Key()) {
		logger.infof(ctx, "%-32s: %-32s; %s", "key", "non-scalar", highlight("returning Invalid"))
		return Invalid
	}
	logger.infof(ctx, "%-16s %-16s; continuing", "key", "scalar")

	tyElem := ty.Elem()
	logger.info(ctx, lowlight(fmt.Sprintf("checking element type (%s) ...", fullTypeName(tyElem))))
	ctx.inc()
	kindElem := getTypeKind(ctx, tyElem)
	ctx.dec()
	if isTypeScalar(tyElem) {
		logger.infof(ctx, "%-16s %-16s; %s", "type", "scalar", highlight("returning SimpleMap"))
		return SimpleMap
	}
	logger.infof(ctx, "%-16s %-16s; continuing", fullTypeName(tyElem), "not scalar")

	logger.infof(ctx, "%-70s", lowlight(fmt.Sprintf("Checking element kind (%s) ...", kindElem.String())))
	if kindElem == SimpleMap ||
		kindElem == SimpleStruct ||
		kindElem == SimpleSlice {
		logger.infof(ctx, "%s: simple; returning SimpleMap", kindElem.String())
		return SimpleMap
	}
	logger.info(ctx, "type: not simple; continuing")

	logger.info(ctx, highlight("conditions were not met; returning Invalid"))
	return Invalid
}

func pointerKind(ctx *ecoContext, ty reflect.Type) kind {
	logger.signature(ctx, "pointerKind", ty)
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Pointer {
		logger.info(ctx, "type is not a pointer; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	logger.infof(ctx, "Checking pointer type (%s) ... ", fullTypeName(tyElem))
	if isTypeScalar(tyElem) {
		logger.info(ctx, "pointer type is scalar; returning SimplePointer")
		return SimplePointer
	}

	logger.info(ctx, "conditions were not met; returning Invalid")
	return Invalid
}

func sliceKind(ctx *ecoContext, ty reflect.Type) kind {
	logger.signature(ctx, "sliceKind", ty)
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Slice {
		logger.info(ctx, "type is not a slice; returning Invalid")
		return Invalid
	}

	tyElem := ty.Elem()
	logger.infof(ctx, "Checking element type (%s) ... ", fullTypeName(tyElem))
	if isTypeScalar(tyElem) {
		logger.info(ctx, "element type is scalar; returning SimpleSlice")
		return SimpleSlice
	}

	logger.info(ctx, "conditions were not met; returning Invalid")
	return Invalid
}

func structKind(ctx *ecoContext, ty reflect.Type) kind {
	// ctx.printf("%s(%s)\n", highlight("structKind"), lowlight(fullTypeName(ty)))
	logger.signature(ctx, "structKind", fullTypeName(ty))
	ctx.inc()
	defer ctx.dec()

	if ty.Kind() != reflect.Struct {
		logger.info(ctx, "type is not a struct; returning Invalid")

		return Invalid
	}

	logger.infof(ctx, "%d field(s)", ty.NumField())
	for i := range ty.NumField() {
		sf := ty.Field(i)
		sfType := sf.Type
		logger.infof(ctx, "%-70s", lowlight(fmt.Sprintf("checking field '%s' (%s) ...", sf.Name, fullTypeName(sfType))))
		sfReflectKind := sfType.Kind()

		if isTypeScalar(sfType) {
			logger.infof(ctx, "%-16s %-16s; %s", sfType, "scalar", "continuing")
			continue
		}

		if sfReflectKind == reflect.Map && mapKind(ctx, sfType) == SimpleMap {
			logger.info(ctx, "field type is SimpleMap; continuing")
			continue
		}

		if sfReflectKind == reflect.Slice && sliceKind(ctx, sfType) == SimpleSlice {
			logger.info(ctx, "field type is SimpleSlice; continuing")
			continue
		}

		if sfReflectKind == reflect.Struct && structKind(ctx, sf.Type) == SimpleStruct {
			logger.info(ctx, "field type is SimpleStruct; continuing")
			continue
		}

		return Invalid
	}

	logger.infof(ctx, "%s; %s", "All fields passed", highlight("returning SimpleStruct"))
	return SimpleStruct
}

/* log styles */

func highlight(s string) string {

	var ss = color.Styledstring(s).Style("\033[1;97m")

	return string(ss)
}

func lowlight(s string) string {
	var ss = color.Styledstring(s).Fg(color.X11.Gray50)

	return string(ss)
}

type ecoLogger struct {
	*slog.Logger
}

func newEcoLogger (w io.Writer) *ecoLogger{
	options := color.ColorOptions{Level: slog.LevelDebug}
	handler := color.NewColorHandler(w, options)
	return &ecoLogger {
		Logger: slog.New(handler),
	}
}

func (l *ecoLogger) comment(ctx *ecoContext, msg string) {
	l.Logger.InfoContext(ctx, l.indent(ctx) + string(color.Styledstring(msg).Fg(color.X11.CornflowerBlue)))
}

func (l *ecoLogger) commentf(ctx *ecoContext, sFmt string, args ...any) {
	msg := fmt.Sprintf(sFmt, args...)
	l.comment(ctx, l.indent(ctx) + msg)
}

func (l *ecoLogger) indent(ctx *ecoContext) string {
	const tab = "  "
	return strings.Repeat(tab, ctx.level)
}

func (l *ecoLogger) info(ctx *ecoContext, s string) {
	l.Logger.InfoContext(ctx, l.indent(ctx) + s)
}

func (l *ecoLogger) infof(ctx *ecoContext, sfmt string, args ...any) {
	s := fmt.Sprintf(sfmt, args...)
	l.Logger.InfoContext(ctx, l.indent(ctx) + s)
}

func (l *ecoLogger) signature (ctx *ecoContext, name string, args ...any) {
	sig := createSignature(name, args...)
	l.Logger.InfoContext(ctx, l.indent(ctx) + sig)
}

func createSignature (name string, args ...any) string {
	// highlight, concat, all that good stuff
	sFmt := fmt.Sprintf("%%s(%s)", strings.Repeat("%v, ", len(args)-1) + "%v")
	args2 := make([]any, len(args)+1)
	args2[0] = highlight(name)
	for i, arg := range args {
		ty, is := arg.(reflect.Type) 
		var sArg string
		if is {
			sArg = fmt.Sprintf("-%s-", fullTypeName(ty))
		} else {
			sArg = fmt.Sprintf("%v", arg)
		}
		args2[i+1] = lowlight(sArg)
	}
	s := fmt.Sprintf(sFmt, args2...)

	return s
}

var logger *ecoLogger

func init () {
	logger = newEcoLogger(os.Stdout)
}
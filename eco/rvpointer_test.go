package eco

import (
	"os"
	"reflect"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestRvpCreateOrGetMap1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/map"}

	x := map[string]int{}
	p := &x
	rvp, err := NewRvPointer(p)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*map[string]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotNil(t, *ptr)
	require.Equal(t, p, ptr)
	require.Equal(t, 0, len(*ptr))
}

func TestRvpCreateOrGetMap2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/map"}

	var p *map[string]int = nil
	pp := &p
	rvp, err := NewRvPointer(pp)
	require.NoError(t, err)
	require.NotNil(t, rvp)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*map[string]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotNil(t, *ptr)
	require.NotEqual(t, p, ptr)
	require.Equal(t, 0, len(*ptr))
}

func TestRvpCreateOrGetMap3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/map"}

	var x map[string]int = nil
	p := &x
	pp := &p
	rvPtr, err := NewRvPointer(pp)
	require.NoError(t, err)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*map[string]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotNil(t, *ptr)
	require.NotEqual(t, p, ptr)
	require.Equal(t, 0, len(*ptr))
}

func TestRvpCreateOrGetMap4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/map"}

	x := map[string]int{}
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp
	rv := reflect.ValueOf(pppp)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*map[string]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotNil(t, *ptr)
	require.Equal(t, p, ptr)
	require.Equal(t, 0, len(*ptr))
}

func TestRvpCreateOrGetMap5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/map"}

	var p *map[string]int = nil
	pp := &p
	ppp := &pp
	pppp := &ppp
	rv := reflect.ValueOf(pppp)
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*map[string]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotNil(t, *ptr)
	require.NotEqual(t, p, ptr)
	require.Equal(t, 0, len(*ptr))
}

func TestRvpCreateOrGetMap6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/map"}

	var x map[string]int = nil
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp
	rv := reflect.ValueOf(pppp)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*map[string]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotNil(t, *ptr)
	require.NotEqual(t, p, ptr)
	require.Equal(t, 0, len(*ptr))
}

func TestRvpCreateOrGetScalar1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/scalar/n", []byte("13")},
	}
	kvSeries := KvSeries{kvs, "/scalar"}

	var x int
	p := &x
	rvp, err := NewRvPointer(p)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.Equal(t, p, ptr)
}

func TestRvpCreateOrGetScalar2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/scalar/n", []byte("13")},
	}
	kvSeries := KvSeries{kvs, "/scalar"}

	var p *int = nil
	pp := &p
	rvp, err := NewRvPointer(pp)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, p, ptr)
}

func TestRvpCreateOrGetScalar3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/scalar/n", []byte("13")},
	}
	kvSeries := KvSeries{kvs, "/scalar"}

	var x int
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp
	rvp, err := NewRvPointer(pppp)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.Equal(t, p, ptr)
}

func TestRvpCreateOrGetScalar4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/scalar/n", []byte("13")},
	}
	kvSeries := KvSeries{kvs, "/scalar"}

	var p *int = nil
	pp := &p
	ppp := &pp
	pppp := &ppp
	rvp, err := NewRvPointer(pppp)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, p, ptr)
}

func TestRvpCreateOrGetStruct1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/struct"}

	x := struct{}{}
	p := &x
	rvp, err := NewRvPointer(p)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*struct{})
	require.True(t, is)
	require.NotNil(t, ptr)
	require.Equal(t, p, ptr)
}

func TestRvpCreateOrGetStruct2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/scalar"}

	var p *struct{} = nil
	pp := &p
	rvp, err := NewRvPointer(pp)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*struct{})
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, p, ptr)
}

func TestRvpCreateOrGetStruct3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/struct"}

	x := struct{}{}
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp
	rvp, err := NewRvPointer(pppp)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*struct{})
	require.True(t, is)
	require.NotNil(t, ptr)
	require.Equal(t, p, ptr)
}

func TestRvpCreateOrGetStruct4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{}
	kvSeries := KvSeries{kvs, "/struct"}

	var p *struct{} = nil
	pp := &p
	ppp := &pp
	pppp := &ppp
	rvp, err := NewRvPointer(pppp)
	require.NoError(t, err)
	normPtr, err := rvp.CreateOrGet(ctx, kvSeries)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*struct{})
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, p, ptr)
}

func TestRvpCreateOrGetSlice1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/3", []byte("3")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	slice := make([]int, 5)
	rv := reflect.ValueOf(&slice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGetSlice(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.Equal(t, &slice, ptr)
}

func TestRvpCreateOrGetSlice2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/13", []byte("13")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	slice := make([]int, 5)
	rv := reflect.ValueOf(&slice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGetSlice(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, &slice, ptr)
	require.Equal(t, 14, len(*ptr))
}

func TestRvpCreateOrGetSlice3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/13", []byte("13")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	var p *[]int = nil
	pp := &p
	rvPtr, err := NewRvPointer(pp)
	require.NoError(t, err)
	normPtr, err := rvPtr.CreateOrGetSlice(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, p, ptr)
	require.Equal(t, 14, len(*ptr))
}

func TestRvpCreateOrGetSlice4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/3", []byte("3")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	slice := make([]int, 5)
	pslice := &slice
	ppslice := &pslice
	pppslice := &ppslice
	ppppslice := &pppslice
	rv := reflect.ValueOf(ppppslice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGetSlice(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.Equal(t, &slice, ptr)
}

func TestRvpCreateOrGetSlice5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/13", []byte("13")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	slice := make([]int, 5)
	pslice := &slice
	ppslice := &pslice
	pppslice := &ppslice
	ppppslice := &pppslice
	rv := reflect.ValueOf(ppppslice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGetSlice(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, &slice, ptr)
	require.Equal(t, 14, len(*ptr))
}

func TestRvpCreateOrGetSlice6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/13", []byte("13")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	var pslice *[]int = nil
	ppslice := &pslice
	pppslice := &ppslice
	ppppslice := &pppslice
	rv := reflect.ValueOf(ppppslice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGetSlice(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, pslice, ptr)
	require.Equal(t, 14, len(*ptr))
}

func TestRvpCreateOrGet_Slice1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/3", []byte("3")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	slice := make([]int, 5)
	rv := reflect.ValueOf(&slice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.Equal(t, &slice, ptr)
}

func TestRvpCreateOrGet_Slice2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/13", []byte("13")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	slice := make([]int, 5)
	rv := reflect.ValueOf(&slice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, &slice, ptr)
	require.Equal(t, 14, len(*ptr))
}

func TestRvpCreateOrGet_Slice3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/13", []byte("13")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	var p *[]int = nil
	pp := &p
	rvPtr, err := NewRvPointer(pp)
	require.NoError(t, err)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, p, ptr)
	require.Equal(t, 14, len(*ptr))
}

func TestRvpCreateOrGet_Slice4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/3", []byte("3")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	slice := make([]int, 5)
	pslice := &slice
	ppslice := &pslice
	pppslice := &ppslice
	ppppslice := &pppslice
	rv := reflect.ValueOf(ppppslice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.Equal(t, &slice, ptr)
}

func TestRvpCreateOrGet_Slice5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/13", []byte("13")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	slice := make([]int, 5)
	pslice := &slice
	ppslice := &pslice
	pppslice := &ppslice
	ppppslice := &pppslice
	rv := reflect.ValueOf(ppppslice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, &slice, ptr)
	require.Equal(t, 14, len(*ptr))
}

func TestRvpCreateOrGet_Slice6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	kvs := []KeyValue{
		{"/slice/1", []byte("1")},
		{"/slice/2", []byte("2")},
		{"/slice/13", []byte("13")},
	}
	kvSlice := KvSeries{kvs, "/slice"}

	var pslice *[]int = nil
	ppslice := &pslice
	pppslice := &ppslice
	ppppslice := &pppslice
	rv := reflect.ValueOf(ppppslice)
	t.Logf("rv.Kind()=%s", rv.Kind().String())
	rvPtr := RvPointer(rv)
	normPtr, err := rvPtr.CreateOrGet(ctx, kvSlice)
	require.NoError(t, err)
	require.NotNil(t, normPtr)
	ptr, is := normPtr.Value.(*[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	require.NotEqual(t, pslice, ptr)
	require.Equal(t, 14, len(*ptr))
}

func TestRvpElemTypeInt1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedData := reflect.TypeFor[int]()
	var pn *int = nil
	var ppn **int = &pn
	rv := reflect.ValueOf(ppn)

	typ, err := RvPointer(rv).ElemType(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

func TestRvpElemTypeInt2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedData := reflect.TypeFor[int]()
	var n int = 13
	var pn *int = &n
	rv := reflect.ValueOf(pn)

	typ, err := RvPointer(rv).ElemType(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// RvPointer.ElemType() - pointer chain to nil pointer
func TestRvpElemTypeInt3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedData := reflect.TypeFor[int]()
	var pn *int = nil
	var ppn **int = &pn
	var pppn ***int = &ppn
	var ppppn ****int = &pppn
	var pppppn *****int = &ppppn
	rv := reflect.ValueOf(pppppn)

	typ, err := RvPointer(rv).ElemType(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// RvPointer.ElemType() - pointer chain to non-nil pointer
func TestRvpElemTypeInt4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedData := reflect.TypeFor[int]()
	var n int = 13
	var pn *int = &n
	var ppn **int = &pn
	var pppn ***int = &ppn
	var ppppn ****int = &pppn
	var pppppn *****int = &ppppn
	rv := reflect.ValueOf(pppppn)

	typ, err := RvPointer(rv).ElemType(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// pointer to nil map
func TestRvpElemTypeMap1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedData := reflect.TypeFor[map[int]string]()
	var x map[int]string = nil
	var p = &x
	var pp = &p
	rv := reflect.ValueOf(pp)

	typ, err := RvPointer(rv).ElemType(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// pointer to non-nil map
func TestRvpElemTypeMap2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedData := reflect.TypeFor[map[int]string]()
	var m map[int]string = map[int]string{13: "meat"}
	var pm = &m
	rv := reflect.ValueOf(pm)

	typ, err := RvPointer(rv).ElemType(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// pointer chain to nil map
func TestRvpElemTypeMap3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedData := reflect.TypeFor[map[int]string]()
	var m map[int]string = nil
	var pm = &m
	var ppm = &pm
	var pppm = &ppm
	var ppppm = &pppm
	rv := reflect.ValueOf(ppppm)

	typ, err := RvPointer(rv).ElemType(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

// pointer chain to non-nil map
func TestRvpElemTypeMap4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	expectedData := reflect.TypeFor[map[int]string]()
	var m map[int]string = map[int]string{13: "meat"}
	var pm = &m
	var ppm = &pm
	var pppm = &ppm
	var ppppm = &pppm
	rv := reflect.ValueOf(ppppm)

	typ, err := RvPointer(rv).ElemType(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedData, typ)
}

func TestRvpIsNil1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var a *int = nil
	rvp, err := NewRvPointer(a)
	require.NoError(t, err)
	require.True(t, rvp.IsNil(ctx))
}

func TestRvpIsNil2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var a *[]int = nil
	rvp, err := NewRvPointer(a)
	require.NoError(t, err)
	require.True(t, rvp.IsNil(ctx))
}

func TestRvpIsNil3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x []int = nil
	a := &x
	rvp, err := NewRvPointer(a)
	require.NoError(t, err)
	require.True(t, rvp.IsNil(ctx))
}

func TestRvpIsNil4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x *[]int = nil
	a := &x
	rvp, err := NewRvPointer(a)
	require.NoError(t, err)
	require.False(t, rvp.IsNil(ctx))
}

func TestRvpIsNil5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	x := []int{}
	a := &x
	rvp, err := NewRvPointer(a)
	require.NoError(t, err)
	require.False(t, rvp.IsNil(ctx))
}

func TestRvpIsNil6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x int = 0
	a := &x
	rvp, err := NewRvPointer(a)
	require.NoError(t, err)
	require.False(t, rvp.IsNil(ctx))
}

func TestRvpIsNil7(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[string]int = nil
	rvp, err := NewRvPointer(p)
	require.NoError(t, err)
	require.True(t, rvp.IsNil(ctx))
}

func TestRvpIsReference1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *[]int = &[]int{}
	rv := reflect.ValueOf(p)
	rvp := RvPointer(rv)
	is := rvp.IsReference(ctx)
	require.True(t, is)
}

func TestRvpIsReference2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[string]int = &(map[string]int{})
	rv := reflect.ValueOf(p)
	rvp := RvPointer(rv)
	is := rvp.IsReference(ctx)
	require.True(t, is)
}

func TestRvpIsReference3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *[]int = nil
	rv := reflect.ValueOf(p)
	rvp := RvPointer(rv)
	is := rvp.IsReference(ctx)
	require.True(t, is)
}

func TestRvpIsReference4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *int = nil
	rv := reflect.ValueOf(p)
	rvp := RvPointer(rv)
	is := rvp.IsReference(ctx)
	require.False(t, is)
}

func TestRvpIsReference5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	rvp := RvPointer(reflect.ValueOf(nil))
	is := rvp.IsReference(ctx)
	require.False(t, is)
}

func TestRvpIsValue1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x int = 13
	p := &x
	rv := reflect.ValueOf(p)
	rvp := RvPointer(rv)
	is := rvp.IsValue(ctx)
	require.True(t, is)
}

func TestRvpIsValue2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	p := &struct{}{}
	rv := reflect.ValueOf(p)
	rvp := RvPointer(rv)
	is := rvp.IsValue(ctx)
	require.True(t, is)
}

func TestRvpIsValue3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *int = nil
	rv := reflect.ValueOf(p)
	rvp := RvPointer(rv)
	is := rvp.IsValue(ctx)
	require.True(t, is)
}

func TestRvpIsValue4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *[]int = nil
	rv := reflect.ValueOf(p)
	rvp := RvPointer(rv)
	is := rvp.IsValue(ctx)
	require.False(t, is)
}

func TestRvpIsValue5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[string]int = nil
	rv := reflect.ValueOf(p)
	rvp := RvPointer(rv)
	is := rvp.IsValue(ctx)
	require.False(t, is)
}

func TestRvpIsValue6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	rvp := RvPointer(reflect.ValueOf(nil))
	is := rvp.IsValue(ctx)
	require.False(t, is)
}

func TestRvpNew1(t *testing.T) {
	var a *int
	rvp, err := NewRvPointer(a)
	require.NoError(t, err)
	require.NotNil(t, rvp)
	rv := reflect.Value(*rvp)
	require.Equal(t, reflect.Pointer, rv.Kind())
	require.Equal(t, reflect.Int, rv.Type().Elem().Kind())
	require.Equal(t, a, rv.Interface())
}

func TestRvpNew2(t *testing.T) {
	var a *****************int = nil
	rvp, err := NewRvPointer(a)
	require.NoError(t, err)
	require.NotNil(t, rvp)
	rv := reflect.Value(*rvp)
	require.Equal(t, reflect.Pointer, rv.Kind())
	require.Equal(t, a, rv.Interface())
}

func TestRvpNew3(t *testing.T) {
	var a int
	rvp, err := NewRvPointer(a)
	require.Error(t, err)
	require.Nil(t, rvp)
}

func TestRvpNew4(t *testing.T) {
	rvp, err := NewRvPointer(nil)
	require.Error(t, err)
	require.Nil(t, rvp)
}

// RvPointer.Walk() - pointer to nil pointer
func TestRvpWalkInt1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *int = nil
	var pp **int = &p

	rv := reflect.ValueOf(pp)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	require.Equal(t, pp, normPtr.Value)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	pptr, is := normPtr.Value.(**int)
	require.NotNil(t, pptr)
	require.True(t, is)
	require.Nil(t, *pptr)
}

// RvPointer.Walk() - non-nil pointer
func TestRvpWalkInt2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var n int = 13
	var pn *int = &n

	rv := reflect.ValueOf(pn)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.True(t, isAllocated)
	ptr, is := normPtr.Value.(*int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, *ptr, n)
}

// RvPointer.Walk() - pointer chain to nil pointer
func TestRvpWalkInt3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var pn *int = nil
	var ppn **int = &pn
	var pppn ***int = &ppn
	var ppppn ****int = &pppn
	var pppppn *****int = &ppppn

	rv := reflect.ValueOf(pppppn)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	pptr, is := normPtr.Value.(**int)
	require.NotNil(t, pptr)
	require.True(t, is)
	require.Nil(t, *pptr)
	require.Equal(t, pptr, ppn)
}

// RvPointer.Walk() - pointer chain to non-nil pointer
func TestRvpWalkInt4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var n int = 13
	var pn *int = &n
	var ppn **int = &pn
	var pppn ***int = &ppn
	var ppppn ****int = &pppn
	var pppppn *****int = &ppppn

	rv := reflect.ValueOf(pppppn)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.True(t, isAllocated)
	ptr, is := normPtr.Value.(*int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, *ptr, n)
}

// pointer to nil map
func TestRvpWalkMap1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x map[int]string = nil
	p := &x
	pp := &p
	rv := reflect.ValueOf(pp)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	ctx.Logger.Infof("normPtr.Value type=%v", reflect.TypeOf(normPtr.Value))
	ptr, is := normPtr.Value.(**map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.NotNil(t, *ptr)
	require.Nil(t, **ptr)
}

// pointer to non-nil map
func TestRvpWalkMap2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	x := map[int]string{13: "meat"}
	var p = &x

	rv := reflect.ValueOf(p)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.True(t, isAllocated)
	ptr, is := normPtr.Value.(*map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, "meat", (*ptr)[13])
}

// pointer chain to nil map
func TestRvpWalkMap3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x map[int]string = nil
	var p = &x
	var pp = &p
	var ppp = &pp
	var pppp = &ppp

	rv := reflect.ValueOf(pppp)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	ptr, is := normPtr.Value.(**map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, pp, ptr)
	require.NotNil(t, *ptr)
	require.Nil(t, **ptr)
}

// pointer chain to non-nil map
func TestRvpWalkMap4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var m map[int]string = map[int]string{13: "meat"}
	var pm = &m
	var ppm = &pm
	var pppm = &ppm
	var ppppm = &pppm

	rv := reflect.ValueOf(ppppm)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.True(t, isAllocated)
	ptr, is := normPtr.Value.(*map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, pm, ptr)
	require.Equal(t, "meat", (*ptr)[13])
}


// pointer chain to nil map pointer
func TestRvpWalkMap5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[int]string = nil
	var pp = &p

	rvp, err := NewRvPointer(pp)
	require.NoError(t, err)
	require.NotNil(t, rvp)
	normPtr, err := rvp.Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	elemType, err := normPtr.ElemType(ctx)
	require.NoError(t, err)
	require.Equal(t, reflect.Map, elemType.Kind())
	ptr, is := normPtr.Value.(**map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, pp, ptr)
}


// pointer chain to nil map pointer
func TestRvpWalkMap6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[int]string = nil
	var pp = &p
	var ppp = &pp
	var pppp = &ppp

	rv := reflect.ValueOf(pppp)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	ptr, is := normPtr.Value.(**map[int]string)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, pp, ptr)
}


// pointer to nil slice
func TestRvpWalkSlice1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x []int = nil
	p := &x
	pp := &p

	rv := reflect.ValueOf(pp)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	ptr, is := normPtr.Value.(**[]int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, pp, ptr)
	require.NotNil(t, *ptr)
	require.Nil(t, **ptr)
}

// pointer to non-nil slice
func TestRvpWalkSlice2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x []int = []int{1, 2, 3, 4, 5}
	p := &x

	rv := reflect.ValueOf(p)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.True(t, isAllocated)
	ptr, is := normPtr.Value.(*[]int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, p, ptr)
	require.Equal(t, 3, (*ptr)[2])
}

// pointer chain to nil slice
func TestRvpWalkSlice3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var x []int = nil
	p := &x
	pp := &p
	ppp := &pp
	pppp := &ppp

	rv := reflect.ValueOf(pppp)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	ptr, is := normPtr.Value.(**[]int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, pp, ptr)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	require.NotNil(t, *ptr)
	require.Nil(t, **ptr)
}

// pointer chain to non-nil slice
func TestRvpWalkSlice4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var slice []int = []int{1, 2, 3, 4, 5}
	var pslice = &slice
	var ppslice = &pslice
	var pppslice = &ppslice
	var ppppslice = &pppslice

	rv := reflect.ValueOf(ppppslice)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.True(t, isAllocated)
	ptr, is := normPtr.Value.(*[]int)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, 3, (*ptr)[2])
}

// pointer chain to nil slice pointer
func TestRvpWalkSlice5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var pslice *[]int = nil
	var ppslice = &pslice
	var pppslice = &ppslice
	var ppppslice = &pppslice

	rv := reflect.ValueOf(ppppslice)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	ptr, is := normPtr.Value.(**[]int)
	require.True(t, is)
	require.NotNil(t, ptr)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	require.Nil(t, *ptr)

}

// RvPointer.Walk() - nil struct pointer
func TestRvpWalkStruct1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var pst *common.TestStruct = nil
	var ppst **common.TestStruct = &pst

	rv := reflect.ValueOf(ppst)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	pptr, is := normPtr.Value.(**common.TestStruct)
	require.NotNil(t, pptr)
	require.True(t, is)
	require.Nil(t, *pptr)
}

// RvPointer.Walk() - non-nil pointer to struct
func TestRvpWalkStruct2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	name := "meat"
	var st common.TestStruct = common.TestStruct{Name: name}
	var pst *common.TestStruct = &st

	rv := reflect.ValueOf(pst)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.True(t, isAllocated)
	ptr, is := normPtr.Value.(*common.TestStruct)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, *ptr, st)
	require.Equal(t, name, (*ptr).Name)
}

// RvPointer.Walk() - nil struct pointer chain
func TestRvpWalkStruct3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var pst *common.TestStruct = nil
	var ppst **common.TestStruct = &pst
	var pppst ***common.TestStruct = &ppst
	var ppppst ****common.TestStruct = &pppst

	rv := reflect.ValueOf(ppppst)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.False(t, isAllocated)
	pptr, is := normPtr.Value.(**common.TestStruct)
	require.NotNil(t, pptr)
	require.True(t, is)
	require.Nil(t, *pptr)
}

// RvPointer.Walk() - non-nil pointer chain to struct
func TestRvpWalkStruct4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	name := "meat"
	var st common.TestStruct = common.TestStruct{Name: name}
	var pst *common.TestStruct = &st
	var ppst **common.TestStruct = &pst
	var pppst ***common.TestStruct = &ppst
	var ppppst ****common.TestStruct = &pppst

	rv := reflect.ValueOf(ppppst)
	normPtr, err := RvPointer(rv).Walk(ctx)
	require.NoError(t, err)
	isAllocated, err := normPtr.IsAllocated(ctx)
	require.NoError(t, err)
	require.True(t, isAllocated)
	ptr, is := normPtr.Value.(*common.TestStruct)
	require.NotNil(t, ptr)
	require.True(t, is)
	require.Equal(t, *ptr, st)
	require.Equal(t, name, (*ptr).Name)
}

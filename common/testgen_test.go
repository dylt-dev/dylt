package common

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

/*
type typ [][]struct { Tempora struct { Eum map[string]struct{ Dolorem map[int][]map[bool][]string } } }
type typ0 string
type typ1 []typ0
type typ2 map[bool]typ1
type typ3 []typ2
type typ4 map[int]typ3
type typ5 struct{Dolorem typ4}
type typ6 map[string]typ5
type typ7 struct{Eum typ6}
type typ8 struct{Tempora typ7}
type typ9 []typ8
type typ10 []typ9

x0 := typ0("meat")
x1 := typ1{x0}
x2 := typ2{true: x1}
x3 := typ3{x2}
x4 := typ4{13: x3}
x5 := typ5{Dolorem: x4}
x6 := typ6{"foo": x5}
x7 := typ7{Eum: x6}
x8 := typ8{Tempora: x7}
x9 := typ9{x8}
x10 := typ10{x9}
x := x10

expected, err := json.Marshal(x0)
require.NoError(t, err)
kvs := encode(ctx, x)
require.NotNil(t, kvs)
require.Equal(t, 1, len(kvs))
require.Equal(t, KeyString("/0/0/Tempora/Eum/foo/Dolorem/13/0/true/0"), kvs[0].Key)
require.Equal(t, expected, kvs[0].Value)
fmt.Fprint(t.Output(), kvs)
*/
func TestGenEncodeTest(t *testing.T) {

}


// Ex: x4 := typ4{13: x3}
func TestCreateEncodeMapAssignStatement1 (t *testing.T) {
	val := 13
	n := 4
	stmt := createEncodeMapAssignStatement(val, n)
	t.Log(stmt)
}

// Ex: type typ4 map[int]typ3
func TestCreateEncodeMapTypeStatement1 (t *testing.T) {
	rtElem := reflect.TypeFor[int]()
	n := 4
	stmt := createEncodeMapTypeStatement(rtElem, n)
	t.Log(stmt)
}


// Ex: x0 := typ0("meat")
func TestCreateEncodeScalarAssignStatement1 (t *testing.T) {
	value := "meat"
	n := 0
	stmt := createEncodeScalarAssignStatement(value, n)
	t.Log(stmt)
}

// Ex: type typ0 string
func TestCreateEncodeScalarTypeStatement1 (t *testing.T) {
	rt := reflect.TypeFor[string]()
	n := 0
	stmt := createEncodeScalarTypeStatement(rt, n)
	t.Log(stmt)
}



// Ex: x3 := typ3{x2}
func TestCreateEncodeSliceAssignStatement1 (t *testing.T) {
	n := 3
	stmt := createEncodeSliceAssignStatement(n)
	t.Log(stmt)
}

// Ex: type typ3 []typ2
func TestCreateEncodeSliceTypeStatement1 (t *testing.T) {
	n := 3
	stmt := createEncodeSliceTypeStatement(n)
	t.Log(stmt)
}


// Ex: x7 := typ7{Eum: x6}
func TestCreateEncodeStructAssignStatement1 (t *testing.T) {
	fieldName := "Eum"
	n := 7
	stmt := createEncodeStructAssignStatement(fieldName, n)
	t.Log(stmt)
}

// Ex: type typ7 struct{Eum typ6}
func TestCreateEncodeStructTypeStatement1 (t *testing.T) {
	fieldName := "Eum"
	n := 7
	stmt := createEncodeStructTypeStatement(fieldName, n)
	t.Log(stmt)
}


func TestGenDeclarations(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	decls := GenDeclarations(ctx, 10, 3)
	for i, decl := range decls {
		fmt.Fprintf(t.Output(), "decls[%d] = %v\n", i, decl)
	}
}


func TestGenSliceDeclaration(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	decl := genSliceDeclaration(ctx, 1)
	t.Log(decl)
}


func TestGenObjCtorStmts1(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type typ map[string]int
	rt := reflect.TypeFor[typ]()
	values := GenScalarValues(ctx, rt)

	stmts := GenObjCtorStmts(ctx, rt, values)
	for _, stmt := range stmts {
		fmt.Fprintln(t.Output(), stmt)
	}
}


func TestGenObjCtorStmts2(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

    type typ struct{Ipsa struct{Placeat []string}}
	rt := reflect.TypeFor[typ]()
	values := GenScalarValues(ctx, rt)

	stmts := GenObjCtorStmts(ctx, rt, values)
	for _, stmt := range stmts {
		fmt.Fprintln(t.Output(), stmt)
	}
}

func TestGenObjCtorStmts3(t *testing.T) {
	ctx := NewEcoContext(os.Stdout)

	type typ struct{N int}
	rt := reflect.TypeFor[typ]()
	values := GenScalarValues(ctx, rt)

	stmts := GenObjCtorStmts(ctx, rt, values)
	for _, stmt := range stmts {
		fmt.Fprintln(t.Output(), stmt)
	}
}


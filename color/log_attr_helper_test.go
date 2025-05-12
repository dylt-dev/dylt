package color

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"testing"
	"unsafe"

	"github.com/dylt-dev/seq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupAttr0 (t *testing.T) {
	var attr slog.Attr
	var gn groupName
	
	attr = slog.String("foo", "thirteen")
	gn = "g"
	var args = metaarglist{}
	args = addMetaargs(args, gn, attr)
	require.Equal(t, 2, len(args))	
	require.IsType(t, *new(groupName), args[0])
	require.IsType(t, *new(slog.Attr), args[1])
}

func TestGroupAttr1 (t *testing.T) {
	var attr, attr0, attr1 slog.Attr
	var gn groupName
	
	gn = "g"
	attr0 = slog.String("foo", "13")
	attr1 = slog.String("bar", "thirteen")
	var args = metaarglist{}
	args = addMetaargs(args, gn, attr0, attr1)
	require.Equal(t, 3, len(args))	
	require.IsType(t, *new(groupName), args[0])
	require.IsType(t, *new(slog.Attr), args[1])
	require.IsType(t, *new(slog.Attr), args[2])

	var attrData = map[string]string{}
	var val string
	var groupNames = []groupName{}
	var sq seq.Seq[any] = newArraySeq(args)
	var arg any
	var err error
	var is bool
	var key string

	arg, err = sq.Next()
	t.Log("Next() ...")
	for true {
		t.Logf("Top o main loop")
		// Loop until a non-groupname is found	
		for true {
			t.Logf("Top o groupname loop")
			t.Log("Checking for io.EOF (normal) or other error (not normal)")
			if errors.Is(err, io.EOF) { t.Log("io.EOF found; breaking"); break }
			require.NoError(t, err)
			t.Logf("arg type=%s", reflect.TypeOf(arg))
			gn, is = arg.(groupName)
			if !is { t.Log("Non-groupname found; breaking"); break }
			t.Logf("groupname found (%s); adding to list", gn)
			groupNames = append(groupNames, gn)
			t.Log("Next() ...")
			arg, err = sq.Next()
		}

		// Loop until a groupname is found
		for true {
			t.Logf("Top o Attr loop")
			t.Log("Checking for io.EOF (normal) or other error (not normal)")
			if errors.Is(err, io.EOF) { t.Log("io.EOF found; breaking"); break }
			require.NoError(t, err)
			attr, is = arg.(slog.Attr)
			if !is { t.Log("Non-Attr found; breaking"); break }
			val = attr.Value.String()
			key = fmt.Sprintf("%s.%s", join(groupNames...), attr.Key)
			t.Logf("Attr found (%s=%s); adding to map", attr.Key, attr.Value.String())
			attrData[key] = val
			t.Log("Next() ...")
			arg, err = sq.Next()
		}

		t.Log("Checking main loop for io.EOF (normal) or other error (not normal)")
		if err != nil {
			if errors.Is(err, io.EOF) { t.Log("io.EOF found; terminating main loop"); break }
		} else {
			t.Logf("Weird error (%s); terminating main loop", err)
		}

	}
	t.Logf("%#v", attrData)
}

func TestArraySeq0 (t *testing.T) {
	a := []int{1, 2, 3}
	var sq seq.Seq[int] = newArraySeq(a)
	var err error
	var el int

	el, err = sq.Next()
	require.NoError(t, err)
	require.Equal(t, 1, el)

	el, err = sq.Next()
	require.NoError(t, err)
	require.Equal(t, 2, el)

	el, err = sq.Next()
	require.NoError(t, err)
	require.Equal(t, 3, el)

	el, err = sq.Next()
	require.Error(t, err)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, 0, el)

	el, err = sq.Next()
	require.Error(t, err)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, 0, el)
}

func TestCopyAndAppend0 (t *testing.T) {
	l := []int64{1, 2, 3}
	var elNew int64 = 4
	// lNew := copyAndAppend(l, elNew)
	lNew := append(l, elNew)
	assert.Equal(t, len(l)+1, len(lNew))
	assert.Equal(t, l[:], lNew[:len(l)])
	assert.Equal(t, elNew, lNew[len(lNew)-1])
	oldEl0 := l[0]
	lNew[0] = 13
	assert.Equal(t, oldEl0, l[0])
	t.Logf("#l: %p\n", unsafe.Pointer(&l))
	t.Logf("#l[0]: %p\n", unsafe.Pointer(&(l[0])))
	t.Logf("#lNew: %p\n", unsafe.Pointer(&lNew))
}

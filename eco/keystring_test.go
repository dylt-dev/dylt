package eco

import (
	"testing"

	"github.com/stretchr/testify/require"
)


func TestChildName1(t *testing.T) {
	expectedData := "bum"
	keyString := KeyString("/foo/bar/bum")
	prefix := KeyString("/foo/bar")
	data, is := keyString.ChildName(prefix)
	require.True(t, is)
	require.Equal(t, expectedData, data)
}

func TestChildName2(t *testing.T) {
	expectedData := ""
	keyString := KeyString("/foo/bar")
	prefix := KeyString("/foo/bar")
	data, is := keyString.ChildName(prefix)
	require.False(t, is)
	require.Equal(t, expectedData, data)
}


func TestCutPrefix3(t *testing.T) {
	expectedData := "/1313"
	prefix := "/prefix"
	keyString := KeyString(expectedData)
	s, is := keyString.CutPrefix(prefix)
	require.False(t, is)
	require.Equal(t, expectedData, string(s))
}


func TestIndex(t *testing.T) {
	expectedData := int(13)
	keyString := KeyString("/foo/bar/13")
	index, is := keyString.Index()
	require.True(t, is)
	require.Equal(t, expectedData, index)
}

func TestIndexBad (t *testing.T) {
	expectedData := int(2)
	keyString := KeyString("/test/boolSlice/2")
	index, is := keyString.Index()
	require.True(t, is)
	require.Equal(t, expectedData, index)
}


func TestIndexEmpty(t *testing.T) {
	expectedData := int(0)
	keyString := KeyString("")
	index, is := keyString.Index()
	require.False(t, is)
	require.Equal(t, expectedData, index)
}

func TestIndexNoInt(t *testing.T) {
	expectedData := int(0)
	keyString := KeyString("/foo/bar/bum")
	index, is := keyString.Index()
	require.False(t, is)
	require.Equal(t, expectedData, index)
}

func TestIndexNoPrefix(t *testing.T) {
	expectedData := int(13)
	keyString := KeyString("/13")
	index, is := keyString.Index()
	require.True(t, is)
	require.Equal(t, expectedData, index)
}

func TestIndexNoSlash(t *testing.T) {
	expectedData := int(13)
	keyString := KeyString("13")
	index, is := keyString.Index()
	require.True(t, is)
	require.Equal(t, expectedData, index)
}


func TestIsParent1(t *testing.T) {
	ksParent := KeyString("/foo")
	ks := KeyString("/foo/bar")
	require.True(t, ksParent.IsParent(ks))
}


func TestIsParent2(t *testing.T) {
	ksParent := KeyString("/foo")
	ks := KeyString("/foo/bar/bum")
	require.True(t, ksParent.IsParent(ks))
}


func TestIsParent3(t *testing.T) {
	ksParent := KeyString("/foo")
	ks := KeyString("/bar")
	require.False(t, ksParent.IsParent(ks))
}


func TestIsParent4(t *testing.T) {
	ksParent := KeyString("/foo")
	ks := KeyString("/foo")
	require.False(t, ksParent.IsParent(ks))
}


func TestIsParent5(t *testing.T) {
	ksParent := KeyString("/foo/")
	ks := KeyString("/foo")
	require.False(t, ksParent.IsParent(ks))
}


func TestIsParent6(t *testing.T) {
	ksParent := KeyString("/foo")
	ks := KeyString("/foo/")
	require.False(t, ksParent.IsParent(ks))
}


func TestIsParent7(t *testing.T) {
	ksParent := KeyString("/foo/")
	ks := KeyString("/foo/")
	require.False(t, ksParent.IsParent(ks))
}


func TestKeyStringChild1(t *testing.T) {
	expectedChild := KeyString("/foo")
	ksParent := KeyString("/foo/bar/bum")
	ksPrefix := KeyString("/")
	ksChild, is := ksParent.Child(ksPrefix)
	require.True(t, is)
	require.Equal(t, expectedChild, ksChild)
}


func TestKeyStringChild2(t *testing.T) {
	expectedChild := KeyString("/foo/bar")
	ksParent := KeyString("/foo/bar/bum")
	ksPrefix := KeyString("/foo")
	ksChild, is := ksParent.Child(ksPrefix)
	require.True(t, is)
	require.Equal(t, expectedChild, ksChild)
}


func TestKeyStringChild3(t *testing.T) {
	expectedChild := KeyString("/foo/bar/bum")
	ksParent := KeyString("/foo/bar/bum")
	ksPrefix := KeyString("/foo/bar")
	ksChild, is := ksParent.Child(ksPrefix)
	require.True(t, is)
	require.Equal(t, expectedChild, ksChild)
}


func TestKeyStringChild4(t *testing.T) {
	expectedChild := KeyString("")
	ksParent := KeyString("/foo/bar/bum")
	ksPrefix := KeyString("/foo/bar/bum")
	ksChild, is := ksParent.Child(ksPrefix)
	require.False(t, is)
	require.Equal(t, expectedChild, ksChild)
}

func TestKeyStringIsLeaf1 (t *testing.T) {
	var ks KeyString = "/foo"
	is := ks.IsLeaf()
	require.True(t, is)
}

func TestKeyStringIsLeaf2 (t *testing.T) {
	var ks KeyString = "/foo/bar"
	is := ks.IsLeaf()
	require.False(t, is)
}


func TestKeyStringIsLeaf3 (t *testing.T) {
	var ks KeyString = ""
	is := ks.IsLeaf()
	require.False(t, is)
}

func TestKeyStringPopHead1(t *testing.T) {
	var expectedHead string = "foo"
	var expectedBody KeyString = "/bar/bum"
	var ks KeyString = "/foo/bar/bum/"
	head, body := ks.PopHead()
	require.Equal(t, expectedHead, head)
	require.Equal(t, expectedBody, body)
}

func TestKeyStringPopHead2(t *testing.T) {
	var expectedHead string = "foo"
	var expectedBody KeyString = ""
	var ks KeyString = "/foo"
	head, body := ks.PopHead()
	require.Equal(t, expectedHead, head)
	require.Equal(t, expectedBody, body)
}

func TestKeyStringPopHead3(t *testing.T) {
	var expectedHead string = ""
	var expectedBody KeyString = ""
	var ks KeyString = ""
	head, body := ks.PopHead()
	require.Equal(t, expectedHead, head)
	require.Equal(t, expectedBody, body)
}

func TestKeyStringTrimHead1(t *testing.T) {
	var expected string = "/bar/bum"
	var ks KeyString = "/foo/bar/bum/"
	ks2 := ks.TrimHead()
	require.Equal(t, expected, string(ks2))
}

func TestKeyStringTrimHead2(t *testing.T) {
	var expected string = ""
	var ks KeyString = "/foo"
	ks2 := ks.TrimHead()
	require.Equal(t, expected, string(ks2))		
}

func TestKeyStringTrimHead3(t *testing.T) {
	var expected string = ""
	var ks KeyString = ""
	ks2 := ks.TrimHead()
	require.Equal(t, expected, string(ks2))		
}

func TestSegments1(t *testing.T) {
	var ks KeyString = "/foo/bar/bum/"
	segments := ks.Segments()
	require.Equal(t, 3, len(segments))
	require.Equal(t, "foo", segments[0])
	require.Equal(t, "bar", segments[1])
	require.Equal(t, "bum", segments[2])
}

func TestSegments2(t *testing.T) {
	var s KeyString = "foo/bar/bum"
	segments := s.Segments()
	require.Equal(t, 3, len(segments))
	require.Equal(t, "foo", segments[0])
	require.Equal(t, "bar", segments[1])
	require.Equal(t, "bum", segments[2])
}

func TestSegments3(t *testing.T) {
	var s KeyString = "//"
	segments := s.Segments()
	require.Equal(t, 1, len(segments))
	require.Equal(t, "", segments[0])
}

func TestSegments4(t *testing.T) {
	var s KeyString = "////"
	segments := s.Segments()
	require.Equal(t, 3, len(segments))
	require.Equal(t, "", segments[0])
	require.Equal(t, "", segments[1])
	require.Equal(t, "", segments[2])
}

func TestWithEndSlash1(t *testing.T) {
	var s KeyString = "/foo/bar"
	var expectedVal KeyString = "/foo/bar/"
	require.Equal(t, expectedVal, s.WithEndSlash())
}

func TestWithEndSlash2(t *testing.T) {
	var s KeyString = "/foo/bar/"
	var expectedVal KeyString = "/foo/bar/"
	require.Equal(t, expectedVal, s.WithEndSlash())
}

func TestWithEndSlash3(t *testing.T) {
	var s KeyString = ""
	var expectedVal KeyString = "/"
	require.Equal(t, expectedVal, s.WithEndSlash())
}

func TestWithEndSlash4(t *testing.T) {
	var s KeyString = "/"
	var expectedVal KeyString = "/"
	require.Equal(t, expectedVal, s.WithEndSlash())
}

func TestWithoutEndSlash1(t *testing.T) {
	var s KeyString = "/foo/bar/"
	var expectedVal KeyString = "/foo/bar"
	require.Equal(t, expectedVal, s.WithoutEndSlash())
}

func TestWithoutEndSlash2(t *testing.T) {
	var s KeyString = "/foo/bar/"
	var expectedVal KeyString = "/foo/bar"
	require.Equal(t, expectedVal, s.WithoutEndSlash())
}

func TestWithoutEndSlash3(t *testing.T) {
	var s KeyString = "/"
	var expectedVal KeyString = ""
	require.Equal(t, expectedVal, s.WithoutEndSlash())
}

func TestWithoutEndSlash4(t *testing.T) {
	var s KeyString = ""
	var expectedVal KeyString = ""
	require.Equal(t, expectedVal, s.WithoutEndSlash())
}

func TestWithStartSlash1(t *testing.T) {
	var s KeyString = "foo/bar/"
	var expectedVal KeyString = "/foo/bar/"
	require.Equal(t, expectedVal, s.WithStartSlash())
}

func TestWithStartSlash2(t *testing.T) {
	var s KeyString = "/foo/bar/"
	var expectedVal KeyString = "/foo/bar/"
	require.Equal(t, expectedVal, s.WithStartSlash())
}

func TestWithStartSlash3(t *testing.T) {
	var s KeyString = ""
	var expectedVal KeyString = "/"
	require.Equal(t, expectedVal, s.WithStartSlash())
}

func TestWithStartSlash4(t *testing.T) {
	var s KeyString = "/"
	var expectedVal KeyString = "/"
	require.Equal(t, expectedVal, s.WithStartSlash())
}

func TestWithoutStartSlash1(t *testing.T) {
	var s KeyString = "/foo/bar/"
	var expectedVal KeyString = "foo/bar/"
	require.Equal(t, expectedVal, s.WithoutStartSlash())
}

func TestWithoutStartSlash2(t *testing.T) {
	var s KeyString = "/foo/bar"
	var expectedVal KeyString = "foo/bar"
	require.Equal(t, expectedVal, s.WithoutStartSlash())
}

func TestWithoutStartSlash3(t *testing.T) {
	var s KeyString = "/"
	var expectedVal KeyString = ""
	require.Equal(t, expectedVal, s.WithoutStartSlash())
}

func TestWithoutStartSlash4(t *testing.T) {
	var s KeyString = ""
	var expectedVal KeyString = ""
	require.Equal(t, expectedVal, s.WithoutStartSlash())
}

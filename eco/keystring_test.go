package eco

import (
	"testing"

	"github.com/stretchr/testify/require"
)


func TestChildName1(t *testing.T) {
	expectedData := "bum"
	keyString := KeyString("/foo/bar/bum")
	prefix := "/foo/bar"
	data := keyString.ChildName(prefix)
	require.Equal(t, expectedData, data)
}

func TestChildName2(t *testing.T) {
	expectedData := ""
	keyString := KeyString("/foo/bar")
	prefix := "/foo/bar"
	data := keyString.ChildName(prefix)
	require.Equal(t, expectedData, data)
}

func TestIndex(t *testing.T) {
	expectedData := uint64(13)
	keyString := KeyString("/foo/bar/13")
	index, is := keyString.Index()
	require.True(t, is)
	require.Equal(t, expectedData, index)
}

func TestIndexBad (t *testing.T) {
	expectedData := uint64(2)
	keyString := KeyString("/test/boolSlice/2")
	index, is := keyString.Index()
	require.True(t, is)
	require.Equal(t, expectedData, index)
}


func TestIndexEmpty(t *testing.T) {
	expectedData := uint64(0)
	keyString := KeyString("")
	index, is := keyString.Index()
	require.False(t, is)
	require.Equal(t, expectedData, index)
}

func TestIndexNoInt(t *testing.T) {
	expectedData := uint64(0)
	keyString := KeyString("/foo/bar/bum")
	index, is := keyString.Index()
	require.False(t, is)
	require.Equal(t, expectedData, index)
}

func TestIndexNoPrefix(t *testing.T) {
	expectedData := uint64(13)
	keyString := KeyString("/13")
	index, is := keyString.Index()
	require.True(t, is)
	require.Equal(t, expectedData, index)
}

func TestIndexNoSlash(t *testing.T) {
	expectedData := uint64(13)
	keyString := KeyString("13")
	index, is := keyString.Index()
	require.True(t, is)
	require.Equal(t, expectedData, index)
}


func TestSegments1(t *testing.T) {
	var s KeyString = "/foo/bar/bum/"
	segments := s.Segments()
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

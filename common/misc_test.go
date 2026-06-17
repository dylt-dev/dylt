package common

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestLev1 (t *testing.T) {
	Lev("aaa", "bbbbb")
}


func TestLev2 (t *testing.T) {
	n := Lev("abcde", "12345")
	t.Logf("n = %d", n)
}


func TestLev3 (t *testing.T) {
	n := Lev("cat", "dog")
	t.Logf("n = %d", n)
}


func TestLev4 (t *testing.T) {
	n := Lev("cat", "rat")
	t.Logf("n = %d", n)
}


func TestLev5 (t *testing.T) {
	n := Lev("cat", "cab")
	t.Logf("n = %d", n)
}


func TestLev6 (t *testing.T) {
	n := Lev("cat", "cot")
	t.Logf("n = %d", n)
}


func TestLev7 (t *testing.T) {
	n := Lev("cat", "cat")
	t.Logf("n = %d", n)
}


func TestWriteDeclaration(t *testing.T) {
	envWriteDecl, is := os.LookupEnv("ECOGEN")
	if !is || (envWriteDecl != "1" && strings.ToLower(envWriteDecl) != "y") {
		t.Skipf("%s not set or not set to 1 or Y/y", "ECOGEN")
	}
	
	ctx := NewEcoContext(os.Stdout)
	ctx.Mute()
	sDepth := os.Args[len(os.Args)-1]
	depth, err := strconv.ParseUint(sDepth, 10, 16)
	require.NoError(t, err)

	bbDecl := bytes.Buffer{}
	WriteDeclaration(ctx, int(depth), &bbDecl)
	decl := bbDecl.String()
	t.Output().Write([]byte(decl))
}

func TestWriteScalarValues(t *testing.T) {
	envWriteDecl, is := os.LookupEnv("ECOGEN")
	if !is || (envWriteDecl != "1" && strings.ToLower(envWriteDecl) != "y") {
		t.Skipf("%s not set or not set to 1 or Y/y", "ECOGEN")
	}
	
	ctx := NewEcoContext(os.Stdout)
	ctx.Mute()
	sDepth := os.Args[len(os.Args)-1]
	depth, err := strconv.ParseUint(sDepth, 10, 16)
	require.NoError(t, err)

	bbDecl := bytes.Buffer{}
	WriteDeclaration(ctx, int(depth), &bbDecl)
	decl := bbDecl.String()
	t.Output().Write([]byte(decl))
}

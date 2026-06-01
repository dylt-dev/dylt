package common

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
	WriteDeclaration(ctx, int(depth), nil, &bbDecl)
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
	WriteDeclaration(ctx, int(depth), nil, &bbDecl)
	decl := bbDecl.String()
	t.Output().Write([]byte(decl))
}

package eco

import (
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestIsPointerAllocated1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *int = nil
	is, err := IsPointerAllocated(ctx, &p)
	require.NoError(t, err)
	require.False(t, is)
}

func TestIsPointerAllocated2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var n int = 13
	var p *int = &n
	is, err := IsPointerAllocated(ctx, p)
	require.NoError(t, err)
	require.True(t, is)
}

func TestIsPointerAllocated3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *[]int = nil
	is, err := IsPointerAllocated(ctx, p)
	require.NoError(t, err)
	require.False(t, is)
}

func TestIsPointerAllocated4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var slice []int = []int{ 13 }
	var p *[]int = &slice
	is, err := IsPointerAllocated(ctx, p)
	require.NoError(t, err)
	require.True(t, is)
}

func TestIsPointerAllocated5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *map[string]int = nil
	is, err := IsPointerAllocated(ctx, p)
	require.NoError(t, err)
	require.False(t, is)
}

func TestIsPointerAllocated6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var slice map[string]int = map[string]int{}
	var p *map[string]int = &slice
	is, err := IsPointerAllocated(ctx, p)
	require.NoError(t, err)
	require.True(t, is)
}

func TestIsPointerAllocated7(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *struct{} = nil
	is, err := IsPointerAllocated(ctx, &p)
	require.NoError(t, err)
	require.False(t, is)
}

func TestIsPointerAllocated8(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	st := struct{}{}
	var p *struct{} = &st
	is, err := IsPointerAllocated(ctx, p)
	require.NoError(t, err)
	require.True(t, is)
}

// non-pointer - error
func TestIsPointerAllocatedError1(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	st := struct{}{}
	is, err := IsPointerAllocated(ctx, st)
	require.Error(t, err)
	require.False(t, is)
}

// nil - error
func TestIsPointerAllocatedError2(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	is, err := IsPointerAllocated(ctx, nil)
	require.Error(t, err)
	require.False(t, is)
}

// scalar nil pointer - error
func TestIsPointerAllocatedError3(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *int = nil
	is, err := IsPointerAllocated(ctx, p)
	require.Error(t, err)
	require.False(t, is)
}

// struct nil pointer - error
func TestIsPointerAllocatedError4(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	var p *struct{} = nil
	is, err := IsPointerAllocated(ctx, p)
	require.Error(t, err)
	require.False(t, is)
}

// pointer to non-nil pointer - error
func TestIsPointerAllocatedError5(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := 13
	p := &n
	is, err := IsPointerAllocated(ctx, &p)
	require.Error(t, err)
	require.False(t, is)
}

// pointer to pointer to pointer -- error
func TestIsPointerAllocatedError6(t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	n := 13
	p := &n
	pp := &p
	ppp := &pp
	is, err := IsPointerAllocated(ctx, ppp)
	require.Error(t, err)
	require.False(t, is)
}

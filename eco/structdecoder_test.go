package eco

import (
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestStructDecoder1 (t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	
	decoder := StructDecoder{}
	tree := &ValueTree{}
	tree.addString(ctx, "/name", "Smitty")
	tree.addInt(ctx, "/lucky_number", 13)
	tree.addString(ctx, "/NoTag", "tagless")
	
	var x common.TestStruct
	p := &x
	err := decoder.Decode(ctx, tree, p)
	
	require.NoError(t, err)
	require.Equal(t, "Smitty", x.Name)
	require.Equal(t, float64(13), x.LuckyNumber)
	require.Equal(t, "tagless", x.NoTag)
}


func TestStructDecoder2 (t *testing.T) {
	ctx := common.NewEcoContext(os.Stdout)
	
	decoder := StructDecoder{}
	tree := &ValueTree{}
	tree.addString(ctx, "/name", "Smitty")
	tree.addInt(ctx, "/lucky_number", 13)
	tree.addString(ctx, "/NoTag", "tagless")
	
	var p *common.TestStruct = nil
	pp := &p
	err := decoder.Decode(ctx, tree, pp)
	x := *p
		
	require.NoError(t, err)
	require.Equal(t, "Smitty", x.Name)
	require.Equal(t, float64(13), x.LuckyNumber)
	require.Equal(t, "tagless", x.NoTag)
}
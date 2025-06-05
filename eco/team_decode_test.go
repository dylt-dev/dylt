package eco

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeTeam_Misc (t *testing.T) {
	var val MiscMap
	decodeAndTest(t, "/test/team/astros/Players/altuve/Misc", &val)
	require.NotEmpty(t, val)
	t.Logf("%#v", val)
}

func TestDecodeTeam_Stats(t *testing.T) {
	var val StatSlice
	decodeAndTest(t, "/test/team/astros/Players/altuve/Stats", &val)
	require.NotEmpty(t, val)
	t.Logf("%#v", val)

}

func decodeAndTest (t *testing.T, key string, i any) {
	etcd, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	ctx := newEcoContext(os.Stdout)
	err = decode(ctx, etcd, key, i)
	require.NoError(t, err)
	el := reflect.ValueOf(i).Elem()
	require.NotNil(t, el)
	require.NotEmpty(t, el)
}

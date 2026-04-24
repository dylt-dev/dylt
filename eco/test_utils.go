package eco

import (
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func initAndTest(t *testing.T) (*ecoContext, *EtcdClient) {
	common.InitLogging()
	ctx := newEcoContext(os.Stdout)
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	return ctx, etcdClient
}

package eco

import (
	"os"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)


func initAndTest(t *testing.T) (*common.EcoContext, *EtcdClient) {
	common.InitLogging()
	ctx := common.NewEcoContext(os.Stdout)
	etcdClient, err := NewEtcdClientFromConfig()
	require.NoError(t, err)
	require.NotNil(t, etcdClient)

	return ctx, etcdClient
}

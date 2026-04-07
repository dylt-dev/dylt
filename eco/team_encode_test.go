package eco

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
)

func TestEncodeTeam (t *testing.T) {
	encodeAndTest(t, "/test/team/astros", VAL_Astros)	
}

func TestEncodeTeam_Misc (t *testing.T) {
	encodeAndTest(t, "/test/team/astros/Players/altuve/Misc", VAL_AltuveMisc)
}

func TestEncodeTeam_Player (t *testing.T) {
	encodeAndTest(t, "/test/team/astros/Players/altuve", VAL_Altuve)
}

func TestEncodeTeam_Players (t *testing.T) {
	encodeAndTest(t, "/test/team/astros/Players", VAL_Players)
}

func TestEncodeTeam_Stats (t *testing.T) {
	encodeAndTest(t, "/test/team/astros/Players/altuve/Stats", VAL_AltuveStats)
}

func createTxn (t *testing.T, cli *EtcdClient) etcd.Txn{
	var err error
	if cli == nil {
		cli, err = CreateEtcdClientFromConfig()
		require.NoError(t, err)
	}
	require.NotNil(t, cli)
	txn := cli.Txn(context.Background())
	require.NotEmpty(t, txn)
	
	return txn
}

func encodeAndTest (t *testing.T, key string, val any) {
	ctx, cli := initAndTest(t)
	ops, err := Encode(ctx, key, val)
	require.NoError(t, err)
	txn := createTxn(t, cli)
	resp, err := txn.Then(ops...).Commit()
	require.NoError(t, err)
	t.Logf("%#v", resp)
}
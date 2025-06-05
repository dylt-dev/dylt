package eco

import (
	"context"
	"os"
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

func createTxn (t *testing.T) etcd.Txn{
	etcd, err := CreateEtcdClientFromConfig()
	require.NoError(t, err)
	txn := etcd.Txn(context.Background())
	require.NotEmpty(t, txn)
	
	return txn
}

func encodeAndTest (t *testing.T, key string, val any) {
	ctx := newEcoContext(os.Stdout)
	ops, err := Encode(ctx, key, val)
	require.NoError(t, err)
	txn := createTxn(t)
	resp, err := txn.Then(ops...).Commit()
	require.NoError(t, err)
	t.Logf("%#v", resp)
}
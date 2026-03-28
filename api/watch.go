package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dylt-dev/dylt/eco"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func RunWatchScript(scriptKey string, targetPath string) error {
	slog.Debug("RunWatchScript()", "scriptKey", scriptKey, targetPath, "targetPath")
	// Get etcd client
	slog.Info("Creating etcd client ...")
	cli, err := eco.CreateEtcdClientFromConfig()
	if err != nil {
		return err
	}

	// Create watch
	ctx := clientv3.WithRequireLeader(context.Background())
	fmt.Printf("Watching %s ...", scriptKey)
	chWatch := cli.Watch(ctx, scriptKey, clientv3.WithKeysOnly())

	// Loop over watch channel
	var resp clientv3.WatchResponse
	for resp = range chWatch {
		for _, ev := range resp.Events {
			key := string(ev.Kv.Key)
			fmt.Printf("%s updated\n", key)
			slog.Debug("Update detected", "key", key)

			if key != scriptKey {
				return fmt.Errorf("key mismatch: execpted %s, got %s", scriptKey, key)
			}

			ctx := context.Background()
			resp, err := cli.Client.Get(ctx, key)
			if err != nil {
				return err
			}
			if len(resp.Kvs) == 0 {
				slog.Warn("No KVs found for watch")
				return nil
			}

			val := (*resp).Kvs[0].Value
			slog.Debug("Value found", "len(val)", len(val))
			fmt.Printf("Writing value to %s ...\n", targetPath)
			err = os.WriteFile(targetPath, val, 0755)
			if err != nil {
				return err
			}

			fmt.Printf("Watching %s ...\n", scriptKey)
		}
	}

	fmt.Println("Done.")
	return nil
}

func RunWatchSvc() error {
	slog.Debug("RunWatchSvc()")

	slog.Info("Creating etcd client ...")
	cli, err := eco.CreateEtcdClientFromConfig()
	if err != nil {
		return err
	}

	// Create watch
	prefix := "svc/"
	ctx := clientv3.WithRequireLeader(context.Background())
	fmt.Printf("Watching %s ...", prefix)
	chWatch := cli.Watch(ctx, prefix, clientv3.WithKeysOnly(), clientv3.WithPrefix())

	// Loop over watch channel
	var resp clientv3.WatchResponse
	for resp = range chWatch {
		for _, ev := range resp.Events {
			key := string(ev.Kv.Key)
			fmt.Printf("%s updated\n", key)
			slog.Debug("Update detected", "key", key)

			fmt.Printf("Watching %s ...\n", prefix)
		}
	}

	fmt.Println("Done.")
	return nil
}


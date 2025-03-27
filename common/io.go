package common

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

func ChownR(folderPath string, uid int, gid int) error {
	var fnWalk fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		if err == nil {
			slog.Debug("lib.ChownR.fnWalk", "path", path, "d.Name()", d.Name())
			fullPath := filepath.Join(folderPath, path)
			err = os.Chown(fullPath, uid, gid)
			if err != nil {
				return err
			}
		}
		return nil
	}

	slog.Debug("lib.ChownR()", "folderPath", folderPath, "uid", uid, "gid", gid)
	dir := os.DirFS(folderPath)
	err := fs.WalkDir(dir, ".", fnWalk)
	if err != nil {
		return err
	}

	return nil
}


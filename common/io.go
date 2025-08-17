package common

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"golang.org/x/term"
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

func IsStderrTerminal () bool {
	var fd = int(os.Stderr.Fd())
	var isTerm = term.IsTerminal(fd)

	return isTerm
}

func IsStdinTerminal () bool {
	var fd = int(os.Stdin.Fd())
	var isTerm = term.IsTerminal(fd)

	return isTerm
}

func IsStdoutTerminal () bool {
	var fd = int(os.Stdout.Fd())
	var isTerm = term.IsTerminal(fd)

	return isTerm
}

func PrintBlankIfTerminal (f *os.File) {
	var fd = int(f.Fd())
	var isTerm bool = term.IsTerminal(fd)
	switch f {
	case os.Stdout: if isTerm { fmt.Fprintln(os.Stdout) }
	case os.Stderr: if isTerm { fmt.Fprintln(os.Stderr) }
	}
}
package lib

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
)

//go:embed svc/*
var svc embed.FS


func CreateWatchDaylightService () error {
	filename := "svc/watch-daylight/watch-daylight.service"

	fs.WalkDir(svc, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fmt.Printf("%s\n", p)
		}
		return nil
	})

	f, err := svc.Open(filename)
	if err != nil { return err }
	s, err := io.ReadAll(f)
	if err != nil { return err }
	fmt.Println(string(s))

	return nil
}
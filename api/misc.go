package api

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/dylt-dev/dylt/common"
	// "github.com/dylt-dev/dylt/template"
)

//go:embed content/*
var content embed.FS

func RunCreateTwoNodeCluster() error {
	common.Logger.Debug("RunCreateTwoNodeCluster()")
	var err error

	r := bufio.NewReader(os.Stdin)

	fmt.Println("Two node cluster time!")
	fmt.Println()

	fmt.Printf("Get your two node's IP addresses or hostnames, and whatever ssh private keys are necessary to connect to them. ")
	_, err = r.ReadBytes('\n')
	fmt.Println()

	fmt.Print("Done! (hit <Enter>) ")
	_, err = r.ReadBytes('\n')
	return err
}


func RunGenEtcdRunScript() error {
	common.Logger.Debug("RunGenEtcdRunScript()")

	fmt.Println("I'm gennin a script!")

	buf, err := content.ReadFile("content/hello.tmpl")
	if err != nil {
		return err
	}
	var tmpl *template.Template = template.New("hello")
	tmpl, err = tmpl.Parse(string(buf))
	tmpl.Execute(os.Stdout, nil)
	return nil
}

package api

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/eco"
)

func RunVmAdd(name string, fqdn string) error {
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// execute command
	vm, err := cli.Add(name, fqdn)
	if err != nil {
		return err
	}
	fmt.Println(vm)

	return nil
}


func RunVmAll() error {
	// get vm-specific etcd client, get all vm data, + show it
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	vmInfoMap, err := cli.All()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(vmInfoMap)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))

	return nil
}


func RunVmDel(name string) error {
	slog.Debug("RunVmDel()", "name", name)
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// delete vm data from cluster
	prevVal, err := cli.Del(name)
	if err != nil {
		return err
	}
	// log deleted vm, if it existede
	if prevVal == nil {
		fmt.Printf("vm '%s' not found\n", name)
	} else {
		fmt.Printf("%s\n", string(prevVal))
	}

	return nil
}


func RunVmGet(name string) error {
	slog.Debug("RunVmGet()", "name", name)
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// get vm data from cluster
	vm, err := cli.Get(name)
	if err != nil {
		return err
	}
	// pritn vm data if vm was found
	if vm == nil {
		fmt.Printf("\nvm '%s' not found.\n\n", name)
	} else {
		fmt.Println(vm)
	}

	return nil
}


func RunVmList() error {
	slog.Debug("RunVmList()")
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// List all vm names, one per line
	names, err := cli.Names()
	if err != nil {
		return err
	}
	fmt.Println()
	for _, name := range names {
		fmt.Println(name)
	}
	fmt.Println()

	return nil
}


func RunVmSet(name string, key string, val string) error {
	slog.Debug("RunVmSet()", "name", name, "key", key, "val", val)
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// get the vm data from the cluster, set the field (if it exists), and save updated object
	vm, err := cli.Get(name)
	if err != nil {
		return err
	}
	err = vm.Set(key, val)
	if err != nil {
		return err
	}
	vm, err = cli.Put(name, vm)
	if err != nil {
		return err
	}
	// print the updated vm if it exists
	if vm == nil {
		fmt.Printf("vm '%s' not found", name)
	} else {
		fmt.Println(vm)
	}

	return nil
}

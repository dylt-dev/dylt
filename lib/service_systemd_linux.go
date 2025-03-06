// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

package lib

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
)

var logger log.Logger
func init () {
	logger = *log.Default()
}
	

func DisableSystemdService (svcName string) error {
	// Check if service exists. If not, leave peacefully.
	does, err := DoesSystemdServiceExist(svcName)
	if err != nil { return err }
	if !does { logger.Printf("%s doesn't exist.", svcName); return nil }

	// Stop service if running
	err = StopSystemdService(svcName)
	if err != nil { return err }

	// Check if service is enabled - if not, leave peacefully
	is, err := IsSystemdServiceEnabled(svcName)
	logger.Printf("is=%t err=%s", is, err)
	if err != nil { return err }
	if !is { fmt.Printf("%s is already disabled", svcName); return nil }

	// Build and run command
	logger.Printf("Disabling %s ...", svcName)
	cmd := BuildDisableServiceCommand(svcName)
	err = cmd.Run()
	if err != nil { return err }

	return nil
}


func DoesSystemdServiceExist (svcName string) (bool, error) {
	// Build and run command. Non-ExitErrors are bad. ExitError.ErrorCode()=4 means service doesn't exist.
	cmd := BuildDoesServiceExistCommand(svcName)
	err := cmd.Run()
	if err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if ok && exiterr.ExitCode() == 4 {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func EnableSystemdService (svcName string, svcFS *ServiceFS) error {
	// Check if service is enabled. If so, leave peacefully.
	is, err := IsSystemdServiceEnabled(svcName)
	if err != nil { return err }
	if is { logger.Printf("%s is already enabled", svcName); return nil }
	
	// Build and run command
	logger.Printf("Enable %s ...", svcName)
	cmd := BuildEnableServiceCommand(svcName, svcFS)
	err = cmd.Run()
	if err != nil { return err }

	return nil
}


func RemoveSystemdService (svcName string) error {
	// Check if service exists. If so, try and disable
	does, err := DoesSystemdServiceExist(svcName)
	if err != nil { return err }
	if does { 
		err := DisableSystemdService(svcName)
		if err != nil { return nil }
	}	

	// Remove service folder
	folderPath := fmt.Sprintf("/opt/svc/%s", svcName)
	err = os.RemoveAll(folderPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logger.Printf("%s already removed.", folderPath)
			return nil
		}
		return err
	}
	logger.Printf("%s successfully removed.", folderPath)

	return nil
}

func StartSystemdService (svcName string) error {
	// Check if service exists. If not, leave peacefully.
	does, err := DoesSystemdServiceExist(svcName)
	logger.Printf("does=%t err=%s", does, err)
	if err != nil { return err }
	if !does { logger.Printf("%s doesn't exist.", svcName); return nil }

	// Check if service is running. If so, leave peacefully. 
	is, err := IsSystemdServiceRunning(svcName)
	logger.Printf("is=%t err%s", is, err)
	if err != nil { return err; }
	if is { logger.Printf("%s is already running", svcName); return nil }

	// Start service
	logger.Printf("Starting %s ...", svcName)
	cmd := BuildStartServiceCommand(svcName) 
	err = cmd.Run()
	if err != nil { return err }

	return nil
}

func StopSystemdService (svcName string) error {
	// Check if service exists. If not, leave peacefully.
	does, err := DoesSystemdServiceExist(svcName)
	logger.Printf("does=%t err=%s", does, err)
	if err != nil { return err }
	if !does { logger.Printf("%s doesn't exist.", svcName); return nil }

	// Check if service is running. If so, stop it.
	is, err := IsSystemdServiceRunning(svcName)
	logger.Printf("is=%t err%s", is, err)
	if err != nil { return err; }
	if !is { logger.Printf("%s is already stopped.", svcName); return nil }
	
	// Stop service
	logger.Printf("Stopping %s ...", svcName)
	cmd := BuildStopServiceCommand(svcName) 
	err = cmd.Run()
	if err != nil { return err }

	return nil
}


// Answers the question 'should we disable this service?'
func IsSystemdServiceEnabled (svcName string) (bool, error) {
	cmd := BuildIsServiceEnabledCommand(svcName)
	err := cmd.Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}


// Answers the question 'should we stop this service?'
func IsSystemdServiceRunning (svcName string) (bool, error) {
	cmd := BuildIsServiceActiveCommand(svcName)
	err := cmd.Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

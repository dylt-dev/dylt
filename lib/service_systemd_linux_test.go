package lib

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const svcName = "watch-daylight"

func TestDisableService (t *testing.T) {
	err := DisableSystemdService(svcName)
	assert.Nil(t, err)
}

func TestEnableService (t *testing.T) {
	var svcFS ServiceFS = ServiceFS{RootPath: PATH_SvcFolderRoot}
	err := EnableSystemdService(svcName, &svcFS)
	assert.Nil(t, err)
}


func TestInstallService (t *testing.T) {
}

func TestStartService (t *testing.T) {
	err := StartSystemdService(svcName)
	assert.Nil(t, err)
}

func TestRemoveService (t *testing.T) {
	err := RemoveSystemdService(svcName)
	assert.Nil(t, err)
}

func TestStopService (t *testing.T) {
	err := StopSystemdService(svcName)
	assert.Nil(t, err)
}

func TestSystemdServiceExistsBadName (t *testing.T) {
	exists, err := DoesSystemdServiceExist("XXXXX")
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestSystemdServiceExistsOk (t *testing.T) {
	exists, err := DoesSystemdServiceExist("watch-daylight")
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestSystemdServiceIsEnabledBadName (t *testing.T) {
	exists, err := IsSystemdServiceEnabled("XXXXX")
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestSystemdServiceIsEnabledOk (t *testing.T) {
	exists, err := IsSystemdServiceEnabled("watch-daylight")
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestSystemdServiceIsRunningBadName (t *testing.T) {
	exists, err := IsSystemdServiceRunning("XXXXX")
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestSystemdServiceIsRunningOk (t *testing.T) {
	exists, err := IsSystemdServiceRunning("watch-daylight")
	assert.Nil(t, err)
	assert.True(t, exists)
}


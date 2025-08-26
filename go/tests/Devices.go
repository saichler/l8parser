package tests

import "github.com/saichler/probler/go/types"
import "github.com/saichler/probler/go/tests"

func Devices() *types.NetworkDeviceList {
	return tests.GenerateExactDeviceTableMockData()
}

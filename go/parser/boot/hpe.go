package boot

import (
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8pollaris/go/types"
)

// CreateHPEServerBootPolls creates collection and parsing Pollaris model for HPE servers
func CreateHPEServerBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "hpe-server"
	polaris.Groups = []string{common.BOOT_GROUP}
	polaris.Polling = make(map[string]*types.Poll)
	createHPESystemPoll(polaris)
	createHPEStoragePoll(polaris)
	createHPEPowerThermalPoll(polaris)
	return polaris
}

// HPE server-specific polling functions
func createHPESystemPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("hpeSystem")
	poll.What = ".1.3.6.1.4.1.232.2"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createHPEVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	poll.Attributes = append(poll.Attributes, createHPEVersion())
	p.Polling[poll.Name] = poll
}

func createHPEStoragePoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("hpeStorage")
	poll.What = ".1.3.6.1.4.1.232.3.2.5.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createDiskStatus())
	p.Polling[poll.Name] = poll
}

func createHPEPowerThermalPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("hpePowerThermal")
	poll.What = ".1.3.6.1.4.1.232.6.2.6.7.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createPowerSupplyStatus())
	poll.Attributes = append(poll.Attributes, createTemperatureSensors())
	p.Polling[poll.Name] = poll
}

// HPE-specific attribute creation functions
func createHPEVendor() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("hpe", ".1.3.6.1.2.1.1.1.0", "Hewlett Packard Enterprise"))
	return attr
}

func createHPEVersion() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.232.2.2.4.2.0"))
	return attr
}
package boot

import (
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8pollaris/go/types"
)

// CreateAristaSwitchBootPolls creates collection and parsing Pollaris model for Arista switches
func CreateAristaSwitchBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "arista-switch"
	polaris.Groups = []string{common.BOOT_GROUP}
	polaris.Polling = make(map[string]*types.Poll)
	createAristaSystemPoll(polaris)
	createAristaInterfacesPoll(polaris)
	createAristaEnvironmentalPoll(polaris)
	return polaris
}

// Arista device-specific polling functions
func createAristaSystemPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("aristaSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createAristaVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	poll.Attributes = append(poll.Attributes, createAristaVersion())
	p.Polling[poll.Name] = poll
}

func createAristaInterfacesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("aristaInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	poll.Attributes = append(poll.Attributes, createInterfaceSpeed())
	p.Polling[poll.Name] = poll
}

func createAristaEnvironmentalPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("aristaEnvironmental")
	poll.What = ".1.3.6.1.4.1.30065.3.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createTemperatureSensors())
	p.Polling[poll.Name] = poll
}

// Arista-specific attribute creation functions
func createAristaVendor() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("arista", ".1.3.6.1.2.1.1.1.0", "Arista"))
	return attr
}

func createAristaVersion() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.30065.1.3.1.1.0"))
	return attr
}
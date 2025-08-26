package boot

import (
	"github.com/saichler/l8pollaris/go/types"
)

// CreateHuaweiRouterBootPolls creates collection and parsing Pollaris model for Huawei routers
func CreateHuaweiRouterBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "huawei-router"
	polaris.Groups = []string{"huawei", "huawei-router"}
	polaris.Polling = make(map[string]*types.Poll)
	createHuaweiSystemPoll(polaris)
	createHuaweiInterfacesPoll(polaris)
	createHuaweiEnvironmentalPoll(polaris)
	return polaris
}

// Huawei device-specific polling functions
func createHuaweiSystemPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("huaweiSystem")
	poll.What = ".1.3.6.1.4.1.2011.5.25"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createHuaweiVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	poll.Attributes = append(poll.Attributes, createHuaweiVersion())
	p.Polling[poll.Name] = poll
}

func createHuaweiInterfacesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("huaweiInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	p.Polling[poll.Name] = poll
}

func createHuaweiEnvironmentalPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("huaweiEnvironmental")
	poll.What = ".1.3.6.1.4.1.2011.5.25.31.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createTemperatureSensors())
	p.Polling[poll.Name] = poll
}

// Huawei-specific attribute creation functions
func createHuaweiVendor() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("huawei", ".1.3.6.1.2.1.1.1.0", "Huawei"))
	return attr
}

func createHuaweiVersion() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.2011.5.25.1.1.1.0"))
	return attr
}
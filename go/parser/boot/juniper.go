package boot

import (
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8pollaris/go/types"
)

// CreateJuniperRouterBootPolls creates collection and parsing Pollaris model for Juniper routers
func CreateJuniperRouterBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "juniper-router"
	polaris.Groups = []string{common.BOOT_GROUP}
	polaris.Polling = make(map[string]*types.Poll)
	createJuniperSystemPoll(polaris)
	createJuniperInterfacesPoll(polaris)
	createJuniperChassisComponentsPoll(polaris)
	createJuniperRoutingEnginePoll(polaris)
	return polaris
}

// Juniper device-specific polling functions
func createJuniperSystemPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("juniperSystem")
	poll.What = ".1.3.6.1.4.1.2636.3.1.2"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createJuniperVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	poll.Attributes = append(poll.Attributes, createJuniperVersion())
	p.Polling[poll.Name] = poll
}

func createJuniperInterfacesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("juniperInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	poll.Attributes = append(poll.Attributes, createInterfaceSpeed())
	p.Polling[poll.Name] = poll
}

func createJuniperChassisComponentsPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("juniperChassis")
	poll.What = ".1.3.6.1.4.1.2636.3.1.13.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createChassisComponentStatus())
	p.Polling[poll.Name] = poll
}

func createJuniperRoutingEnginePoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("juniperRoutingEngine")
	poll.What = ".1.3.6.1.4.1.2636.3.1.13.1.8"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createRoutingEngineUtilization())
	p.Polling[poll.Name] = poll
}

// Juniper-specific attribute creation functions
func createJuniperVendor() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("juniper", ".1.3.6.1.2.1.1.1.0", "Juniper"))
	return attr
}

func createJuniperVersion() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.2636.3.1.2.0"))
	return attr
}
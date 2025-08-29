package boot

import (
	"github.com/saichler/l8pollaris/go/types"
)

// CreateCiscoSwitchBootPolls creates collection and parsing Pollaris model for Cisco switches
func CreateCiscoSwitchBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "cisco-switch"
	polaris.Groups = []string{"cisco", "cisco-switch"}
	polaris.Polling = make(map[string]*types.Poll)
	createCiscoSystemPoll(polaris)
	createCiscoVersionPoll(polaris)
	createCiscoSerialPoll(polaris)
	createCiscoInterfacesPoll(polaris)
	createCiscoModulesPoll(polaris)
	createCiscoPowerSupplyPoll(polaris)
	createCiscoFanPoll(polaris)
	createCiscoCpuPoll(polaris)
	createCiscoMemoryPoll(polaris)
	return polaris
}

// CreateCiscoRouterBootPolls creates collection and parsing Pollaris model for Cisco routers
func CreateCiscoRouterBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "cisco-router"
	polaris.Groups = []string{"cisco", "cisco-router"}
	polaris.Polling = make(map[string]*types.Poll)
	createCiscoSystemPoll(polaris)
	createCiscoVersionPoll(polaris)
	createCiscoSerialPoll(polaris)
	createCiscoInterfacesPoll(polaris)
	createCiscoRouterModulesPoll(polaris)
	createCiscoPowerSupplyPoll(polaris)
	createCiscoCpuPoll(polaris)
	createCiscoMemoryPoll(polaris)
	createCiscoRoutingPoll(polaris)
	return polaris
}

// Cisco device-specific polling functions
func createCiscoSystemPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createCiscoVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	poll.Attributes = append(poll.Attributes, createSysOid())
	p.Polling[poll.Name] = poll
}

func createCiscoVersionPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoVersion")
	poll.What = ".1.3.6.1.4.1.9.9.25.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createCiscoVersion())
	p.Polling[poll.Name] = poll
}

func createCiscoSerialPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoSerial")
	poll.What = ".1.3.6.1.4.1.9.3.6"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createCiscoSerial())
	p.Polling[poll.Name] = poll
}

func createCiscoInterfacesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	poll.Attributes = append(poll.Attributes, createInterfaceSpeed())
	poll.Attributes = append(poll.Attributes, createInterfaceMtu())
	p.Polling[poll.Name] = poll
}

func createCiscoModulesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoModules")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createModuleName())
	poll.Attributes = append(poll.Attributes, createModuleModel())
	poll.Attributes = append(poll.Attributes, createModuleStatus())
	p.Polling[poll.Name] = poll
}

func createCiscoPowerSupplyPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoPowerSupply")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createPowerSupplyStatus())
	poll.Attributes = append(poll.Attributes, createPowerSupplyModel())
	p.Polling[poll.Name] = poll
}

func createCiscoFanPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoFans")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createFanStatus())
	p.Polling[poll.Name] = poll
}

func createCiscoCpuPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoCpu")
	poll.What = ".1.3.6.1.4.1.9.9.109.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createCpuUtilization())
	p.Polling[poll.Name] = poll
}

func createCiscoMemoryPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoMemory")
	poll.What = ".1.3.6.1.4.1.9.9.48.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createMemoryUtilization())
	p.Polling[poll.Name] = poll
}

func createCiscoRouterModulesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoRouterModules")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createRouteProcessorStatus())
	p.Polling[poll.Name] = poll
}

func createCiscoRoutingPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoRouting")
	poll.What = ".1.3.6.1.2.1.4.21.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createRoutingTableEntry())
	p.Polling[poll.Name] = poll
}

// Cisco-specific attribute creation functions
func createCiscoVendor() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("cisco", ".1.3.6.1.2.1.1.1.0", "Cisco"))
	return attr
}

func createCiscoVersion() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.9.25.1.1.1.2.2"))
	return attr
}

func createCiscoSerial() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.serialnumber"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.3.6.3.0"))
	return attr
}
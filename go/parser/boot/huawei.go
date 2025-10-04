package boot

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

// CreateHuaweiRouterBootPolls creates collection and parsing Pollaris model for Huawei routers
func CreateHuaweiRouterBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "huawei-router"
	polaris.Groups = []string{"huawei", "huawei-router"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createHuaweiSystemPoll(polaris)
	createHuaweiMibSystemPoll(polaris)
	createHuaweiInterfacesPoll(polaris)
	createHuaweiEnvironmentalPoll(polaris)
	return polaris
}

// Huawei device-specific polling functions
func createHuaweiSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("huaweiSystem")
	poll.What = ".1.3.6.1.4.1.2011.5.25.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createHuaweiVersion())
	p.Polling[poll.Name] = poll
}

func createHuaweiMibSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("huaweiMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createHuaweiVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createHuaweiInterfacesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("huaweiInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	p.Polling[poll.Name] = poll
}

func createHuaweiEnvironmentalPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("huaweiEnvironmental")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createTemperatureSensors())
	p.Polling[poll.Name] = poll
}

// Huawei-specific attribute creation functions
func createHuaweiVendor() *l8tpollaris.L8P_Attribute {
	attr := &l8tpollaris.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("huawei", ".1.3.6.1.2.1.1.1.0", "Huawei"))
	return attr
}

func createHuaweiVersion() *l8tpollaris.L8P_Attribute {
	attr := &l8tpollaris.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.2011.5.25.1.1.1.0"))
	return attr
}

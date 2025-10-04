package boot

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

// CreateHPEServerBootPolls creates collection and parsing Pollaris model for HPE servers
func CreateHPEServerBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "hpe-server"
	polaris.Groups = []string{"hpe", "hpe-server"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createHPESystemPoll(polaris)
	createHPEMibSystemPoll(polaris)
	createHPEStoragePoll(polaris)
	createHPEPowerThermalPoll(polaris)
	return polaris
}

// HPE server-specific polling functions
func createHPESystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("hpeSystem")
	poll.What = ".1.3.6.1.4.1.232.2.2.4"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createHPEVersion())
	p.Polling[poll.Name] = poll
}

func createHPEMibSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("hpeMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createHPEVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createHPEStoragePoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("hpeStorage")
	poll.What = ".1.3.6.1.2.1.25.2.3.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createDiskStatus())
	p.Polling[poll.Name] = poll
}

func createHPEPowerThermalPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("hpePowerThermal")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createPowerSupplyStatus())
	poll.Attributes = append(poll.Attributes, createTemperatureSensors())
	p.Polling[poll.Name] = poll
}

// HPE-specific attribute creation functions
func createHPEVendor() *l8tpollaris.L8P_Attribute {
	attr := &l8tpollaris.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("hpe", ".1.3.6.1.2.1.1.1.0", "Hewlett Packard Enterprise"))
	return attr
}

func createHPEVersion() *l8tpollaris.L8P_Attribute {
	attr := &l8tpollaris.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.232.2.2.4.2.0"))
	return attr
}

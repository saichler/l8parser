package boot

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

// CreateDellServerBootPolls creates collection and parsing Pollaris model for Dell servers
func CreateDellServerBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "dell-server"
	polaris.Groups = []string{"dell", "dell-server"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createDellSystemPoll(polaris)
	createDellMibSystemPoll(polaris)
	createDellStoragePoll(polaris)
	createDellPowerThermalPoll(polaris)
	return polaris
}

// Dell server-specific polling functions
func createDellSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("dellSystem")
	poll.What = ".1.3.6.1.4.1.674.10892.5.1.3"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createDellVersion())
	p.Polling[poll.Name] = poll
}

func createDellMibSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("dellMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createDellVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createDellStoragePoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("dellStorage")
	poll.What = ".1.3.6.1.2.1.25.2.3.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createDiskStatus())
	p.Polling[poll.Name] = poll
}

func createDellPowerThermalPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("dellPowerThermal")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createPowerSupplyStatus())
	poll.Attributes = append(poll.Attributes, createTemperatureSensors())
	p.Polling[poll.Name] = poll
}

// Dell-specific attribute creation functions
func createDellVendor() *l8tpollaris.L8P_Attribute {
	attr := &l8tpollaris.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("dell", ".1.3.6.1.2.1.1.1.0", "Dell"))
	return attr
}

func createDellVersion() *l8tpollaris.L8P_Attribute {
	attr := &l8tpollaris.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.674.10892.5.1.3.1.6.0"))
	return attr
}

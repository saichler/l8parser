package boot

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

// CreateNokiaRouterBootPolls creates collection and parsing Pollaris model for Nokia routers
func CreateNokiaRouterBootPolls() *l8poll.L8Pollaris {
	polaris := &l8poll.L8Pollaris{}
	polaris.Name = "nokia-router"
	polaris.Groups = []string{"nokia", "nokia-router"}
	polaris.Polling = make(map[string]*l8poll.L8Poll)
	createNokiaSystemPoll(polaris)
	createNokiaMibSystemPoll(polaris)
	createNokiaInterfacesPoll(polaris)
	createNokiaCardsPoll(polaris)
	return polaris
}

// Nokia device-specific polling functions
func createNokiaSystemPoll(p *l8poll.L8Pollaris) {
	poll := createBaseSNMPPoll("nokiaSystem")
	poll.What = ".1.3.6.1.4.1.6527.3.1.2.2.1"
	poll.Operation = l8poll.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createNokiaVersion())
	p.Polling[poll.Name] = poll
}

func createNokiaMibSystemPoll(p *l8poll.L8Pollaris) {
	poll := createBaseSNMPPoll("nokiaMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8poll.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createNokiaVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createNokiaInterfacesPoll(p *l8poll.L8Pollaris) {
	poll := createBaseSNMPPoll("nokiaInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = l8poll.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	p.Polling[poll.Name] = poll
}

func createNokiaCardsPoll(p *l8poll.L8Pollaris) {
	poll := createBaseSNMPPoll("nokiaCards")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = l8poll.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createCardStatus())
	p.Polling[poll.Name] = poll
}

// Nokia-specific attribute creation functions
func createNokiaVendor() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("nokia", ".1.3.6.1.2.1.1.1.0", "Nokia"))
	return attr
}

func createNokiaVersion() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.6527.3.1.2.2.1.4.0"))
	return attr
}

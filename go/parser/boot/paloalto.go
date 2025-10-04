package boot

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

// CreatePaloAltoFirewallBootPolls creates collection and parsing Pollaris model for Palo Alto firewalls
func CreatePaloAltoFirewallBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "paloalto-firewall"
	polaris.Groups = []string{"paloalto", "paloalto-firewall"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createPaloAltoSystemPoll(polaris)
	createPaloAltoMibSystemPoll(polaris)
	createPaloAltoSessionsPoll(polaris)
	createPaloAltoThreatsPoll(polaris)
	return polaris
}

// Palo Alto Networks device-specific polling functions
func createPaloAltoSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("paloAltoSystem")
	poll.What = ".1.3.6.1.4.1.25461.2.1.2.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createPaloAltoVersion())
	p.Polling[poll.Name] = poll
}

func createPaloAltoMibSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("paloAltoMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createPaloAltoVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createPaloAltoInterfacesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("ifTable")
	poll.What = ".1.3.6.1.2.1.2.2"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createPaloAltoIfTableRule())
	p.Polling[poll.Name] = poll
}

func createPaloAltoSessionsPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("paloAltoSessions")
	poll.What = ".1.3.6.1.4.1.25461.2.1.2.3"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createActiveSessions())
	p.Polling[poll.Name] = poll
}

func createPaloAltoThreatsPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("paloAltoThreats")
	poll.What = ".1.3.6.1.4.1.25461.2.1.2.2"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createThreatCount())
	p.Polling[poll.Name] = poll
}

func createPaloAltoIfTableRule() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.physicals"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)

	// Use custom rule to translate ifTable CTable to NetworkDevice.physicals
	rule := &l8tpollaris.L8PRule{}
	rule.Name = "IfTableToPhysicals"
	rule.Params = make(map[string]*l8tpollaris.L8PParameter)
	attr.Rules = append(attr.Rules, rule)

	return attr
}

// Palo Alto-specific attribute creation functions
func createPaloAltoVendor() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("paloalto", ".1.3.6.1.2.1.1.1.0", "Palo Alto Networks"))
	return attr
}

func createPaloAltoVersion() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.25461.2.1.2.1.1.0"))
	return attr
}

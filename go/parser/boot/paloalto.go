package boot

import (
	"github.com/saichler/l8pollaris/go/types"
)

// CreatePaloAltoFirewallBootPolls creates collection and parsing Pollaris model for Palo Alto firewalls
func CreatePaloAltoFirewallBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "paloalto-firewall"
	polaris.Groups = []string{"paloalto", "paloalto-firewall"}
	polaris.Polling = make(map[string]*types.Poll)
	createPaloAltoSystemPoll(polaris)
	createPaloAltoMibSystemPoll(polaris)
	createPaloAltoInterfacesPoll(polaris)
	createPaloAltoSessionsPoll(polaris)
	createPaloAltoThreatsPoll(polaris)
	return polaris
}

// Palo Alto Networks device-specific polling functions
func createPaloAltoSystemPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("paloAltoSystem")
	poll.What = ".1.3.6.1.4.1.25461.2.1.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createPaloAltoVersion())
	p.Polling[poll.Name] = poll
}

func createPaloAltoMibSystemPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("paloAltoMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createPaloAltoVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createPaloAltoInterfacesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("paloAltoInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	p.Polling[poll.Name] = poll
}

func createPaloAltoSessionsPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("paloAltoSessions")
	poll.What = ".1.3.6.1.4.1.25461.2.1.2.3.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createActiveSessions())
	p.Polling[poll.Name] = poll
}

func createPaloAltoThreatsPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("paloAltoThreats")
	poll.What = ".1.3.6.1.4.1.25461.2.1.2.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createThreatCount())
	p.Polling[poll.Name] = poll
}

// Palo Alto-specific attribute creation functions
func createPaloAltoVendor() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("paloalto", ".1.3.6.1.2.1.1.1.0", "Palo Alto Networks"))
	return attr
}

func createPaloAltoVersion() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.25461.2.1.2.1.1.0"))
	return attr
}
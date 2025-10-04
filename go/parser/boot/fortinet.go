package boot

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

// CreateFortinetFirewallBootPolls creates collection and parsing Pollaris model for Fortinet firewalls
func CreateFortinetFirewallBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "fortinet-firewall"
	polaris.Groups = []string{"fortinet", "fortinet-firewall"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createFortinetSystemPoll(polaris)
	createFortinetMibSystemPoll(polaris)
	createFortinetInterfacesPoll(polaris)
	createFortinetSessionsPoll(polaris)
	createFortinetVpnPoll(polaris)
	return polaris
}

// Fortinet device-specific polling functions
func createFortinetSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("fortinetSystem")
	poll.What = ".1.3.6.1.4.1.12356.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createFortinetVersion())
	p.Polling[poll.Name] = poll
}

func createFortinetMibSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("fortinetMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createFortinetVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createFortinetInterfacesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("fortinetInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	p.Polling[poll.Name] = poll
}

func createFortinetSessionsPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("fortinetSessions")
	poll.What = ".1.3.6.1.4.1.12356.101.4.1.8"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createFortinetActiveSessions())
	p.Polling[poll.Name] = poll
}

func createFortinetVpnPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("fortinetVpn")
	poll.What = ".1.3.6.1.4.1.12356.101.12.2.3.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createFortinetVpnTunnelStatus())
	p.Polling[poll.Name] = poll
}

// Fortinet-specific attribute creation functions
func createFortinetVendor() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("fortinet", ".1.3.6.1.2.1.1.1.0", "Fortinet"))
	return attr
}

func createFortinetVersion() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.12356.1.1.0"))
	return attr
}

func createFortinetActiveSessions() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.physicals.performance.activeconnections"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.12356.101.4.1.8.0"))
	return attr
}

func createFortinetVpnTunnelStatus() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.networklinks.linkstatus"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.12356.101.12.2.3.1.3"))
	return attr
}

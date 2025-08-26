package boot

import (
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8pollaris/go/types"
)

// CreateNokiaRouterBootPolls creates collection and parsing Pollaris model for Nokia routers
func CreateNokiaRouterBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "nokia-router"
	polaris.Groups = []string{common.BOOT_GROUP}
	polaris.Polling = make(map[string]*types.Poll)
	createNokiaSystemPoll(polaris)
	createNokiaInterfacesPoll(polaris)
	createNokiaCardsPoll(polaris)
	return polaris
}

// Nokia device-specific polling functions
func createNokiaSystemPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("nokiaSystem")
	poll.What = ".1.3.6.1.4.1.6527.3.1.2.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createNokiaVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	poll.Attributes = append(poll.Attributes, createNokiaVersion())
	p.Polling[poll.Name] = poll
}

func createNokiaInterfacesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("nokiaInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	p.Polling[poll.Name] = poll
}

func createNokiaCardsPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("nokiaCards")
	poll.What = ".1.3.6.1.4.1.6527.3.1.2.2.3.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createCardStatus())
	p.Polling[poll.Name] = poll
}

// Nokia-specific attribute creation functions
func createNokiaVendor() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("nokia", ".1.3.6.1.2.1.1.1.0", "Nokia"))
	return attr
}

func createNokiaVersion() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.6527.3.1.2.2.1.4.0"))
	return attr
}
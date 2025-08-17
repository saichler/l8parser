package boot

import (
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8parser/go/parser/rules"
	"github.com/saichler/l8pollaris/go/types"
	"strconv"
)

var DEFAULT_CADENCE int64 = 300
var DEFAULT_TIMEOUT int64 = 30

func CreateSNMPBootPolls() *types.Pollaris {
	snmpPolaris := &types.Pollaris{}
	snmpPolaris.Name = "mib2"
	snmpPolaris.Groups = []string{common.BOOT_GROUP}
	snmpPolaris.Polling = make(map[string]*types.Poll)
	createSystemMibPoll(snmpPolaris)
	return snmpPolaris
}

func createSystemMibPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("systemMib")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createVendor() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("cisco", ".1.3.6.1.2.1.1.1.0", "Cisco"))
	attr.Rules = append(attr.Rules, createContainsRule("ubuntu", ".1.3.6.1.2.1.1.1.0", "Ubuntu Linux"))
	return attr
}

func createSysName() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.systemname"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.5.0"))
	return attr
}

func createVersion() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("cisco", ".1.3.6.1.2.1.1.1.0", "Cisco"))
	attr.Rules = append(attr.Rules, createContainsRule("ubuntu", ".1.3.6.1.2.1.1.1.0", "Ubuntu Linux"))
	return attr
}

func addParameter(name, value string, rule *types.Rule) {
	param := &types.Parameter{}
	param.Name = name
	param.Value = value
	rule.Params[name] = param
}

func createContainsRule(what, from, output string) *types.Rule {
	rule := &types.Rule{}
	rule.Name = "Contains"
	rule.Params = make(map[string]*types.Parameter)
	addParameter("what", what, rule)
	addParameter("from", from, rule)
	addParameter("output", output, rule)
	return rule
}

func createToTable(columns, keycolumn int) *types.Rule {
	rule := &types.Rule{}
	rule.Name = "ToTable"
	rule.Params = make(map[string]*types.Parameter)
	rule.Params[rules.Columns] = &types.Parameter{Name: rules.Columns, Value: strconv.Itoa(columns)}
	rule.Params[rules.KeyColumn] = &types.Parameter{Name: rules.KeyColumn, Value: strconv.Itoa(keycolumn)}
	return rule
}

func createTableToMap() *types.Rule {
	rule := &types.Rule{}
	rule.Name = "TableToMap"
	rule.Params = make(map[string]*types.Parameter)
	return rule
}

func createSetRule(from string) *types.Rule {
	rule := &types.Rule{}
	rule.Name = "Set"
	rule.Params = make(map[string]*types.Parameter)
	addParameter("from", from, rule)
	return rule
}

func createBaseSNMPPoll(jobName string) *types.Poll {
	poll := &types.Poll{}
	poll.Name = jobName
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Cadence = DEFAULT_CADENCE
	poll.Protocol = types.Protocol_PSNMPV2
	return poll
}

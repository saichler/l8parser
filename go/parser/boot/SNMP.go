package boot

import (
	"strconv"
	"strings"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8parser/go/parser/rules"
	"github.com/saichler/l8pollaris/go/types"
)

var DEFAULT_CADENCE int64 = 300
var DEFAULT_TIMEOUT int64 = 30

// CreateSNMPBootPolls creates generic SNMP collection and parsing Pollaris model
func CreateSNMPBootPolls() *types.Pollaris {
	snmpPolaris := &types.Pollaris{}
	snmpPolaris.Name = "mib2"
	snmpPolaris.Groups = []string{common.BOOT_GROUP}
	snmpPolaris.Polling = make(map[string]*types.Poll)
	createSystemMibPoll(snmpPolaris)
	return snmpPolaris
}

// GetPollarisByOid returns the appropriate vendor-specific Pollaris model based on sysOID
func GetPollarisByOid(sysOid string) *types.Pollaris {
	// Cisco devices
	if isCiscoOid(sysOid) {
		if isCiscoSwitchOid(sysOid) {
			return CreateCiscoSwitchBootPolls()
		}
		return CreateCiscoRouterBootPolls() // Default to router for Cisco
	}

	// Juniper devices
	if isJuniperOid(sysOid) {
		return CreateJuniperRouterBootPolls()
	}

	// Palo Alto Networks devices
	if isPaloAltoOid(sysOid) {
		return CreatePaloAltoFirewallBootPolls()
	}

	// Fortinet devices
	if isFortinetOid(sysOid) {
		return CreateFortinetFirewallBootPolls()
	}

	// Arista devices
	if isAristaOid(sysOid) {
		return CreateAristaSwitchBootPolls()
	}

	// Nokia devices
	if isNokiaOid(sysOid) {
		return CreateNokiaRouterBootPolls()
	}

	// Huawei devices
	if isHuaweiOid(sysOid) {
		return CreateHuaweiRouterBootPolls()
	}

	// Dell servers
	if isDellOid(sysOid) {
		return CreateDellServerBootPolls()
	}

	// HPE servers
	if isHPEOid(sysOid) {
		return CreateHPEServerBootPolls()
	}

	// Default to generic SNMP polling if no vendor match
	return CreateSNMPBootPolls()
}

// GetAllPolarisModels returns a slice of all available Pollaris models
func GetAllPolarisModels() []*types.Pollaris {
	models := make([]*types.Pollaris, 0)

	// Generic SNMP
	models = append(models, CreateSNMPBootPolls())

	// Cisco devices
	models = append(models, CreateCiscoSwitchBootPolls())
	models = append(models, CreateCiscoRouterBootPolls())

	// Juniper devices
	models = append(models, CreateJuniperRouterBootPolls())

	// Palo Alto devices
	models = append(models, CreatePaloAltoFirewallBootPolls())

	// Fortinet devices
	models = append(models, CreateFortinetFirewallBootPolls())

	// Arista devices
	models = append(models, CreateAristaSwitchBootPolls())

	// Nokia devices
	models = append(models, CreateNokiaRouterBootPolls())

	// Huawei devices
	models = append(models, CreateHuaweiRouterBootPolls())

	// Dell devices
	models = append(models, CreateDellServerBootPolls())

	// HPE devices
	models = append(models, CreateHPEServerBootPolls())

	return models
}

// OID matching helper functions
func isCiscoOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}
	return strings.HasPrefix(normalizedOid, ".1.3.6.1.4.1.9.")
}

func isCiscoSwitchOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}

	// Cisco switch specific OIDs (Catalyst series)
	switchOids := []string{
		".1.3.6.1.4.1.9.1.122",  // Catalyst 2960
		".1.3.6.1.4.1.9.1.616",  // Catalyst 3560
		".1.3.6.1.4.1.9.1.717",  // Catalyst 3750
		".1.3.6.1.4.1.9.1.1208", // Catalyst 4500
		".1.3.6.1.4.1.9.1.1146", // Catalyst 6500
	}
	for _, switchOid := range switchOids {
		if strings.HasPrefix(normalizedOid, switchOid) {
			return true
		}
	}
	return false
}

func isJuniperOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}
	return strings.HasPrefix(normalizedOid, ".1.3.6.1.4.1.2636.")
}

func isPaloAltoOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}
	return strings.HasPrefix(normalizedOid, ".1.3.6.1.4.1.25461.")
}

func isFortinetOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}
	return strings.HasPrefix(normalizedOid, ".1.3.6.1.4.1.12356.")
}

func isAristaOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}
	return strings.HasPrefix(normalizedOid, ".1.3.6.1.4.1.30065.")
}

func isNokiaOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}
	return strings.HasPrefix(normalizedOid, ".1.3.6.1.4.1.6527.")
}

func isHuaweiOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}
	return strings.HasPrefix(normalizedOid, ".1.3.6.1.4.1.2011.")
}

func isDellOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}
	return strings.HasPrefix(normalizedOid, ".1.3.6.1.4.1.674.")
}

func isHPEOid(sysOid string) bool {
	// Normalize OID by ensuring it starts with a dot
	normalizedOid := sysOid
	if !strings.HasPrefix(normalizedOid, ".") {
		normalizedOid = "." + normalizedOid
	}
	return strings.HasPrefix(normalizedOid, ".1.3.6.1.4.1.232.")
}

func createSystemMibPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("systemMib")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	poll.Attributes = append(poll.Attributes, createSysOid())
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

func createSysOid() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.sysoid"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.2.0"))
	return attr
}

func createSysName() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.sysname"
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

// Common utility functions for creating rules and polls
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

// Common network attribute functions (used by vendor-specific files)
func createInterfaceName() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.name"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.2"))
	return attr
}

func createInterfaceStatus() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.status"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.8"))
	return attr
}

func createInterfaceSpeed() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.speed"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.5"))
	return attr
}

func createInterfaceMtu() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.mtu"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.4"))
	return attr
}

// Module and chassis attribute functions
func createModuleName() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.name"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2"))
	return attr
}

func createModuleModel() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.model"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.13"))
	return attr
}

func createModuleStatus() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.status"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5"))
	return attr
}

func createChassisComponentStatus() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.status"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.6"))
	return attr
}

// Power and environmental attribute functions
func createPowerSupplyStatus() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.status"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.3"))
	return attr
}

func createPowerSupplyModel() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.model"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2"))
	return attr
}

func createFanStatus() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.status"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.3"))
	return attr
}

func createTemperatureSensors() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.temperaturecelsius"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.4"))
	return attr
}

// Performance attribute functions
func createCpuUtilization() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.cpuusagepercent"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.9.109.1.1.1.1.5"))
	return attr
}

func createMemoryUtilization() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.memoryusagepercent"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.9.48.1.1.1.6"))
	return attr
}

func createRoutingEngineUtilization() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.utilizationpercent"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.2636.3.1.13.1.8"))
	return attr
}

// Router-specific attribute functions
func createRouteProcessorStatus() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.status"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5"))
	return attr
}

func createRoutingTableEntry() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.ipaddress"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.4.21.1.1"))
	return attr
}

func createCardStatus() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.status"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.9"))
	return attr
}

// Firewall-specific attribute functions
func createActiveSessions() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.activeconnections"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.25461.2.1.2.3.1.0"))
	return attr
}

func createThreatCount() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alerts.count"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.25461.2.1.2.2.1.0"))
	return attr
}

func createVpnTunnelStatus() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.networklinks.linkstatus"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.25461.2.1.2.4.1.3"))
	return attr
}

// Server-specific attribute functions
func createDiskStatus() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.processes.status"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.25.2.3.1.4"))
	return attr
}

package boot

import (
	"strconv"
	"strings"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8parser/go/parser/rules"
	"github.com/saichler/l8pollaris/go/types/l8poll"
)

var DEFAULT_CADENCE = &l8poll.L8C_CadencePlan{Cadences: []int64{900, 3600, 7200}, Enabled: true}
var EVERY_5_MINUTES = &l8poll.L8C_CadencePlan{Cadences: []int64{300, 3600, 7200}, Enabled: true}
var EVERY_5_MINUTES_ALWAYS = &l8poll.L8C_CadencePlan{Cadences: []int64{300}, Enabled: true}
var DISABLED = &l8poll.L8C_CadencePlan{Cadences: []int64{7200}, Enabled: false}
var DEFAULT_TIMEOUT int64 = 30

func CreateBoot00() *l8poll.L8Pollaris {
	boot00 := &l8poll.L8Pollaris{}
	boot00.Name = "boot00"
	boot00.Groups = []string{common.BOOT_STAGE_00}
	boot00.Polling = make(map[string]*l8poll.L8Poll)
	createIpAddressPoll(boot00)
	return boot00
}

func CreateBoot01() *l8poll.L8Pollaris {
	boot01 := &l8poll.L8Pollaris{}
	boot01.Name = "boot01"
	boot01.Groups = []string{common.BOOT_STAGE_01}
	boot01.Polling = make(map[string]*l8poll.L8Poll)
	createSystemMibPoll(boot01)
	return boot01
}

func CreateBoot02() *l8poll.L8Pollaris {
	boot02 := &l8poll.L8Pollaris{}
	boot02.Name = "boot02"
	boot02.Groups = []string{common.BOOT_STAGE_02}
	boot02.Polling = make(map[string]*l8poll.L8Poll)
	createDeviceStatusPoll(boot02)
	return boot02
}

// CreateSNMPBootPolls creates generic SNMP collection and parsing Pollaris model
func CreateBoot03() *l8poll.L8Pollaris {
	boot03 := &l8poll.L8Pollaris{}
	boot03.Name = "boot03"
	boot03.Groups = []string{common.BOOT_STAGE_03}
	boot03.Polling = make(map[string]*l8poll.L8Poll)
	createIfTable(boot03)
	createEntityMibPoll(boot03)
	return boot03
}

// GetPollarisByOid returns the appropriate vendor-specific Pollaris model based on sysOID
func GetPollarisByOid(sysOid string) *l8poll.L8Pollaris {
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
	return CreateBoot03()
}

// GetAllPolarisModels returns a slice of all available Pollaris models
func GetAllPolarisModels() []*l8poll.L8Pollaris {
	models := make([]*l8poll.L8Pollaris, 0)

	//Generic K8s
	models = append(models, CreateK8sBootPolls())

	// Generic Pre Boot
	models = append(models, CreateBoot00())
	models = append(models, CreateBoot01())
	models = append(models, CreateBoot02())

	// Generic SNMP
	models = append(models, CreateBoot03())

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

func createSystemMibPoll(p *l8poll.L8Pollaris) {
	poll := createBaseSNMPPoll("systemMib")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8poll.L8C_Operation_L8C_Map
	poll.Cadence = EVERY_5_MINUTES
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createVendor())            // networkdevice.equipmentinfo.vendor
	poll.Attributes = append(poll.Attributes, createSysName())           // networkdevice.equipmentinfo.sys_name
	poll.Attributes = append(poll.Attributes, createSysOid())            // networkdevice.equipmentinfo.sys_oid
	poll.Attributes = append(poll.Attributes, createSystemDescription()) // networkdevice.equipmentinfo.hardware
	poll.Attributes = append(poll.Attributes, createSystemSoftware())    // networkdevice.equipmentinfo.software
	poll.Attributes = append(poll.Attributes, createSystemVersion())     // networkdevice.equipmentinfo.version
	poll.Attributes = append(poll.Attributes, createSystemModel())       // networkdevice.equipmentinfo.model
	poll.Attributes = append(poll.Attributes, createSystemUptime())      // networkdevice.equipmentinfo.uptime
	poll.Attributes = append(poll.Attributes, createSystemLocation())    // networkdevice.equipmentinfo.location
	poll.Attributes = append(poll.Attributes, createSystemDeviceType())  // networkdevice.equipmentinfo.device_type
	p.Polling[poll.Name] = poll
}

func createIfTable(p *l8poll.L8Pollaris) {
	poll := createBaseSNMPPoll("ifTable")
	poll.What = ".1.3.6.1.2.1.2.2"
	poll.Operation = l8poll.L8C_Operation_L8C_Table
	poll.Cadence = DISABLED // Disable ifTable polling
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createIfTableRule())
	p.Polling[poll.Name] = poll
}

func createEntityMibPoll(p *l8poll.L8Pollaris) {
	poll := createBaseSNMPPoll("entityMib")
	poll.What = ".1.3.6.1.2.1.47.1.1"
	poll.Operation = l8poll.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createEntityMibRule())
	p.Polling[poll.Name] = poll
}

func createIpAddressPoll(p *l8poll.L8Pollaris) {
	poll := createBaseSNMPPoll("ipAddress")
	poll.What = "ipaddress" // Static value instead of SNMP OID
	poll.Operation = l8poll.L8C_Operation_L8C_Map
	poll.Cadence = EVERY_5_MINUTES
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createIpAddress())
	p.Polling[poll.Name] = poll
}

func createDeviceStatusPoll(p *l8poll.L8Pollaris) {
	poll := createBaseSNMPPoll("deviceStatus")
	poll.What = "devicestatus" // Static value instead of SNMP OID
	poll.Operation = l8poll.L8C_Operation_L8C_Map
	poll.Cadence = EVERY_5_MINUTES_ALWAYS
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createDeviceStatus())
	p.Polling[poll.Name] = poll
}

func createVendor() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("cisco", ".1.3.6.1.2.1.1.1.0", "Cisco"))
	attr.Rules = append(attr.Rules, createContainsRule("ubuntu", ".1.3.6.1.2.1.1.1.0", "Ubuntu Linux"))
	return attr
}

func createIfTableRule() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)

	// Use custom rule to translate ifTable CTable to NetworkDevice.physicals
	rule := &l8poll.L8P_Rule{}
	rule.Name = "IfTableToPhysicals"
	rule.Params = make(map[string]*l8poll.L8P_Parameter)
	attr.Rules = append(attr.Rules, rule)

	return attr
}

func createEntityMibRule() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)

	// Use custom rule to translate Entity MIB CTable to NetworkDevice.physicals
	rule := &l8poll.L8P_Rule{}
	rule.Name = "EntityMibToPhysicals"
	rule.Params = make(map[string]*l8poll.L8P_Parameter)
	attr.Rules = append(attr.Rules, rule)

	return attr
}

func createSysOid() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.sysoid"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.2.0"))
	return attr
}

func createSysName() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.sysname"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.5.0"))
	return attr
}

func createVersion() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("cisco", ".1.3.6.1.2.1.1.1.0", "Cisco"))
	attr.Rules = append(attr.Rules, createContainsRule("ubuntu", ".1.3.6.1.2.1.1.1.0", "Ubuntu Linux"))
	return attr
}

// Common utility functions for creating rules and polls
func addParameter(name, value string, rule *l8poll.L8P_Rule) {
	param := &l8poll.L8P_Parameter{}
	param.Name = name
	param.Value = value
	rule.Params[name] = param
}

func createContainsRule(what, from, output string) *l8poll.L8P_Rule {
	rule := &l8poll.L8P_Rule{}
	rule.Name = "Contains"
	rule.Params = make(map[string]*l8poll.L8P_Parameter)
	addParameter("what", what, rule)
	addParameter("from", from, rule)
	addParameter("output", output, rule)
	return rule
}

func createToTable(columns, keycolumn int) *l8poll.L8P_Rule {
	rule := &l8poll.L8P_Rule{}
	rule.Name = "StringToCTable"
	rule.Params = make(map[string]*l8poll.L8P_Parameter)
	rule.Params[rules.Columns] = &l8poll.L8P_Parameter{Name: rules.Columns, Value: strconv.Itoa(columns)}
	rule.Params[rules.KeyColumn] = &l8poll.L8P_Parameter{Name: rules.KeyColumn, Value: strconv.Itoa(keycolumn)}
	return rule
}

func createTableToMap() *l8poll.L8P_Rule {
	rule := &l8poll.L8P_Rule{}
	rule.Name = "CTableToMapProperty"
	rule.Params = make(map[string]*l8poll.L8P_Parameter)
	return rule
}

func createSetRule(from string) *l8poll.L8P_Rule {
	rule := &l8poll.L8P_Rule{}
	rule.Name = "Set"
	rule.Params = make(map[string]*l8poll.L8P_Parameter)
	addParameter("from", from, rule)
	return rule
}

func createDeviceStatusRule() *l8poll.L8P_Rule {
	rule := &l8poll.L8P_Rule{}
	rule.Name = "MapToDeviceStatus"
	rule.Params = make(map[string]*l8poll.L8P_Parameter)
	addParameter("from", "devicestatus", rule)
	return rule
}

func createBaseSNMPPoll(jobName string) *l8poll.L8Poll {
	poll := &l8poll.L8Poll{}
	poll.Name = jobName
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Cadence = DEFAULT_CADENCE
	poll.Protocol = l8poll.L8C_Protocol_L8P_PSNMPV2
	return poll
}

// Common network attribute functions (used by vendor-specific files)
func createInterfaceName() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.2"))
	return attr
}

func createInterfaceStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.8"))
	return attr
}

func createInterfaceSpeed() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.speed"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.5"))
	return attr
}

func createInterfaceMtu() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.mtu"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.4"))
	return attr
}

// Module and chassis attribute functions
func createModuleName() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2"))
	return attr
}

func createModuleModel() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.model"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.13"))
	return attr
}

func createModuleStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5"))
	return attr
}

func createChassisComponentStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.6"))
	return attr
}

// Power and environmental attribute functions
func createPowerSupplyStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.3"))
	return attr
}

func createPowerSupplyModel() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.model"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2"))
	return attr
}

func createFanStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.3"))
	return attr
}

func createTemperatureSensors() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.temperaturecelsius"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.4"))
	return attr
}

// Performance attribute functions
func createCpuUtilization() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.cpuusagepercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.9.109.1.1.1.1.5"))
	return attr
}

func createMemoryUtilization() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.memoryusagepercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.9.48.1.1.1.6"))
	return attr
}

func createRoutingEngineUtilization() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.utilizationpercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.2636.3.1.13.1.8"))
	return attr
}

// Router-specific attribute functions
func createRouteProcessorStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5"))
	return attr
}

func createRoutingTableEntry() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.ipaddress"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.4.21.1.1"))
	return attr
}

func createCardStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.9"))
	return attr
}

// Firewall-specific attribute functions
func createActiveSessions() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.activeconnections"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.25461.2.1.2.3.1.0"))
	return attr
}

func createThreatCount() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alertcount"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.25461.2.1.2.2.1.0"))
	return attr
}

func createVpnTunnelStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.linkstatus"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.25461.2.1.2.4.1.3"))
	return attr
}

// Server-specific attribute functions
func createDiskStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.processes.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.25.2.3.1.4"))
	return attr
}

// =====================================
// COMPREHENSIVE POLLING CONFIGURATIONS
// Supporting all NetworkDevice model attributes
// =====================================

// Equipment Info Extended Attributes
func createSysNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.sysname"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.5.0"))
	return attr
}

func createSysOidAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.sysoid"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.2.0"))
	return attr
}

func createTargetIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.deviceid"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: TargetId typically derived from inventory system, not directly available via SNMP
	// Could be mapped from sysName or other identifier
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.5.0")) // Use sysName as fallback
	return attr
}

func createLocationAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.location"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.6.0")) // sysLocation
	return attr
}

func createUptimeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.uptime"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.3.0")) // sysUpTime
	return attr
}

func createLastSeenAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.lastseen"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: LastSeen is typically managed by polling system, not available via SNMP
	// This would be updated by the collector based on successful polling
	return attr
}

func createHardwareAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.hardware"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.1.0")) // sysDescr (contains hardware info)
	return attr
}

func createSoftwareAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.software"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.1.0")) // sysDescr (contains software info)
	return attr
}

func createSeriesAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.series"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Series typically derived from model parsing, not directly available
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.1.0")) // Parse from sysDescr
	return attr
}

func createFamilyAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.family"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Family typically derived from model parsing, not directly available
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.1.0")) // Parse from sysDescr
	return attr
}

func createInterfaceCountAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.interfacecount"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.1.0")) // ifNumber
	return attr
}

// Physical Component Attributes
func createChassisSerialAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.serialnumber"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.11.1")) // entPhysicalSerialNum
	return attr
}

func createChassisModelAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.model"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.13.1")) // entPhysicalModelName
	return attr
}

func createChassisDescriptionAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.description"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2.1")) // entPhysicalDescr
	return attr
}

func createChassisTemperatureAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.temperature"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.99.1.1.1.4.1")) // entSensorValue (temperature)
	return attr
}

// Module Attributes
func createModuleNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.7")) // entPhysicalName (table)
	return attr
}

func createModuleModelAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.model"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.13")) // entPhysicalModelName (table)
	return attr
}

func createModuleDescriptionAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.description"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2")) // entPhysicalDescr (table)
	return attr
}

func createModuleStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5")) // entPhysicalClass (table)
	return attr
}

func createModuleTypeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.moduletype"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: ModuleType derived from entPhysicalClass and description parsing
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5")) // entPhysicalClass
	return attr
}

func createModuleTemperatureAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.temperature"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.99.1.1.1.4")) // entSensorValue (table)
	return attr
}

// CPU Attributes
func createCpuIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.id"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.1")) // entPhysicalIndex for CPU
	return attr
}

func createCpuNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.7")) // entPhysicalName for CPU
	return attr
}

func createCpuModelAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.model"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.13")) // entPhysicalModelName for CPU
	return attr
}

func createCpuArchitectureAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.architecture"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Architecture not typically available via standard SNMP, would need vendor-specific MIBs
	// Placeholder using description field
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2"))
	return attr
}

func createCpuCoresAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.cores"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Core count not typically available via standard SNMP
	// Would require vendor-specific MIBs or parsing from description
	return attr
}

func createCpuFrequencyAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.frequencymhz"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: CPU frequency not typically available via standard SNMP
	// Would require vendor-specific MIBs
	return attr
}

func createCpuStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5")) // entPhysicalClass
	return attr
}

func createCpuTemperatureAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.cpus.temperature"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.99.1.1.1.4")) // entSensorValue for CPU temp
	return attr
}

// Memory Attributes
func createMemoryIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.memorymodules.id"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.1")) // entPhysicalIndex for Memory
	return attr
}

func createMemoryNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.memorymodules.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.7")) // entPhysicalName for Memory
	return attr
}

func createMemoryTypeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.memorymodules.type"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Memory type not typically available via standard SNMP
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2")) // Parse from description
	return attr
}

func createMemorySizeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.memorymodules.sizebytes"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Memory size not typically available via standard SNMP for modules
	// Would use HOST-RESOURCES-MIB for total memory
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.25.2.2.0")) // hrMemorySize
	return attr
}

func createMemoryFrequencyAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.memorymodules.frequencymhz"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Memory frequency not typically available via standard SNMP
	return attr
}

func createMemoryStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.chassis.modules.memorymodules.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5")) // entPhysicalClass
	return attr
}

// Port Attributes
func createPortIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.id"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.1")) // ifIndex (table)
	return attr
}

// Interface Attributes (nested in ports)
func createInterfaceIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.id"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.1")) // ifIndex
	return attr
}

func createInterfaceNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.2")) // ifDescr
	return attr
}

func createInterfaceStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.8")) // ifOperStatus
	return attr
}

func createInterfaceDescriptionAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.description"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.31.1.1.1.18")) // ifAlias
	return attr
}

func createInterfaceTypeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.interfacetype"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.3")) // ifType
	return attr
}

func createInterfaceSpeedAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.speed"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.5")) // ifSpeed
	return attr
}

func createInterfaceMacAddressAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.macaddress"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.6")) // ifPhysAddress
	return attr
}

func createInterfaceIpAddressAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.ipaddress"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.4.20.1.1")) // ipAdEntAddr (table)
	return attr
}

func createInterfaceMtuAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.mtu"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.4")) // ifMtu
	return attr
}

func createInterfaceAdminStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.adminstatus"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.7")) // ifAdminStatus
	return attr
}

// Interface Statistics
func createInterfaceRxPacketsAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.statistics.rxpackets"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.11")) // ifInUcastPkts
	return attr
}

func createInterfaceTxPacketsAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.statistics.txpackets"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.17")) // ifOutUcastPkts
	return attr
}

func createInterfaceRxBytesAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.statistics.rxbytes"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.10")) // ifInOctets
	return attr
}

func createInterfaceTxBytesAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.statistics.txbytes"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.16")) // ifOutOctets
	return attr
}

func createInterfaceRxErrorsAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.statistics.rxerrors"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.14")) // ifInErrors
	return attr
}

func createInterfaceTxErrorsAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.statistics.txerrors"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.20")) // ifOutErrors
	return attr
}

func createInterfaceRxDropsAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.statistics.rxdrops"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.13")) // ifInDiscards
	return attr
}

func createInterfaceTxDropsAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.statistics.txdrops"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.19")) // ifOutDiscards
	return attr
}

func createInterfaceCollisionsAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.ports.interfaces.statistics.collisions"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Collisions not available in standard IF-MIB, would need EtherLike-MIB
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.10.7.2.1.4")) // dot3StatsLateCollisions
	return attr
}

// Power Supply Attributes
func createPowerSupplyIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.id"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.1")) // entPhysicalIndex for PSU
	return attr
}

func createPowerSupplyNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.7")) // entPhysicalName
	return attr
}

func createPowerSupplyModelAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.model"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.13")) // entPhysicalModelName
	return attr
}

func createPowerSupplySerialNumberAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.serialnumber"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.11")) // entPhysicalSerialNum
	return attr
}

func createPowerSupplyWattageAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.wattage"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Wattage not typically available via standard SNMP, vendor-specific required
	return attr
}

func createPowerSupplyPowerTypeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.powertype"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Power type (AC/DC) not typically available via standard SNMP
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2")) // Parse from description
	return attr
}

func createPowerSupplyStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5")) // entPhysicalClass/status
	return attr
}

func createPowerSupplyTemperatureAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.temperature"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.99.1.1.1.4")) // entSensorValue for PSU temp
	return attr
}

func createPowerSupplyLoadPercentAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.loadpercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Load percentage not typically available via standard SNMP
	return attr
}

func createPowerSupplyVoltageAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.voltage"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Voltage not typically available via standard SNMP
	return attr
}

func createPowerSupplyCurrentAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.powersupplies.current"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Current not typically available via standard SNMP
	return attr
}

// Fan Attributes
func createFanIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.id"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.1")) // entPhysicalIndex for Fan
	return attr
}

func createFanNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.7")) // entPhysicalName
	return attr
}

func createFanDescriptionAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.description"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.2")) // entPhysicalDescr
	return attr
}

func createFanStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.47.1.1.1.1.5")) // entPhysicalClass/status
	return attr
}

func createFanSpeedRpmAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.speedrpm"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.99.1.1.1.4")) // entSensorValue for fan speed
	return attr
}

func createFanMaxSpeedRpmAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.maxspeedrpm"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Max speed not typically available via standard SNMP
	return attr
}

func createFanTemperatureAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.temperature"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.99.1.1.1.4")) // entSensorValue for fan temp
	return attr
}

func createFanVariableSpeedAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.fans.variablespeed"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Variable speed capability not typically available via standard SNMP
	return attr
}

// Performance Metrics Attributes
func createPerformanceCpuUsageAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.cpuusagepercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.25.3.3.1.2.1")) // hrProcessorLoad
	return attr
}

func createPerformanceMemoryUsageAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.memoryusagepercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Memory usage percentage calculation required from hrStorageUsed/hrStorageSize
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.25.2.3.1.6.1")) // hrStorageUsed
	return attr
}

func createPerformanceTemperatureAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.temperaturecelsius"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.99.1.1.1.4.1")) // entSensorValue (main temp sensor)
	return attr
}

func createPerformanceUptimeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.uptime"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.3.0")) // sysUpTime
	return attr
}

func createPerformanceLoadAverageAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.loadaverage"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.2021.10.1.3.1")) // laLoad.1 (if UCD-SNMP-MIB available)
	return attr
}

// Process Information Attributes
func createProcessNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.processes.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.25.4.2.1.2")) // hrSWRunName (table)
	return attr
}

func createProcessPidAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.processes.pid"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.25.4.2.1.1")) // hrSWRunIndex (table)
	return attr
}

func createProcessCpuPercentAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.processes.cpupercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Per-process CPU percentage not available in standard HOST-RESOURCES-MIB
	return attr
}

func createProcessMemoryPercentAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.processes.memorypercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Per-process memory percentage not available in standard HOST-RESOURCES-MIB
	return attr
}

func createProcessStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.physicals.performance.processes.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.25.4.2.1.7")) // hrSWRunStatus (table)
	return attr
}

// =====================================
// LOGICAL INTERFACES SECTION
// =====================================

// Logical Interface Attributes
func createLogicalInterfaceIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.id"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.1")) // ifIndex for logical interfaces
	return attr
}

func createLogicalInterfaceNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.2")) // ifDescr for logical interfaces
	return attr
}

func createLogicalInterfaceStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.8")) // ifOperStatus
	return attr
}

func createLogicalInterfaceDescriptionAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.description"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.31.1.1.1.18")) // ifAlias
	return attr
}

func createLogicalInterfaceTypeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.interfacetype"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.3")) // ifType
	return attr
}

func createLogicalInterfaceIpAddressAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.ipaddress"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.4.20.1.1")) // ipAdEntAddr (table)
	return attr
}

func createLogicalInterfaceMtuAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.mtu"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.4")) // ifMtu
	return attr
}

func createLogicalInterfaceAdminStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.logicals.interfaces.adminstatus"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.7")) // ifAdminStatus
	return attr
}

// =====================================
// NETWORK TOPOLOGY SECTION
// =====================================

// Network Topology Attributes
func createTopologyIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.topologyid"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Topology ID not available via SNMP - generated by management system
	return attr
}

func createTopologyNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Topology name not available via SNMP - configured by management system
	return attr
}

func createTopologyTypeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.topologytype"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Topology type not available via SNMP - configured by management system
	return attr
}

func createTopologyLastUpdatedAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.lastupdated"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Last updated timestamp managed by topology discovery system
	return attr
}

// Network Node Attributes
func createNetworkNodeIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.nodes.nodeid"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Node ID typically derived from device identification, not directly from SNMP
	return attr
}

func createNetworkNodeNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.nodes.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.5.0")) // sysName
	return attr
}

func createNetworkNodeTypeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.nodes.nodetype"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Node type derived from device classification, not directly from SNMP
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.2.0")) // Parse from sysObjectID
	return attr
}

func createNetworkNodeStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.nodes.status"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Node status derived from polling success/failure and device operational state
	return attr
}

func createNetworkNodeLocationAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.nodes.location"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.6.0")) // sysLocation
	return attr
}

func createNetworkNodeRegionAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.nodes.region"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Region typically configured or derived from location parsing
	return attr
}

func createNetworkNodeTierAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.topology.nodes.tier"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Tier classification not available via SNMP - network design concept
	return attr
}

// =====================================
// NETWORK LINKS SECTION
// =====================================

// Network Link Attributes
func createNetworkLinkIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.linkid"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Link ID generated by topology discovery, not available via SNMP
	return attr
}

func createNetworkLinkNameAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.name"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Link name typically configured or generated by management system
	return attr
}

func createNetworkLinkFromNodeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.fromnode"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: From node determined by topology discovery
	return attr
}

func createNetworkLinkToNodeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.tonode"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: To node determined by topology discovery
	return attr
}

func createNetworkLinkStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.linkstatus"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Link status derived from interface operational states
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.8")) // ifOperStatus
	return attr
}

func createNetworkLinkBandwidthAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.bandwidth"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.5")) // ifSpeed
	return attr
}

func createNetworkLinkTypeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.linktype"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.3")) // ifType
	return attr
}

func createNetworkLinkUtilizationAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.utilizationpercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Utilization calculated from ifInOctets/ifOutOctets vs ifSpeed over time
	return attr
}

func createNetworkLinkLatencyAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.latencyms"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Latency not available via SNMP - requires active measurement
	return attr
}

func createNetworkLinkDistanceAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.distancekm"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Distance not available via SNMP - geographic calculation or configuration
	return attr
}

func createNetworkLinkUptimeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.uptime"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Link uptime calculated from interface uptime tracking
	return attr
}

func createNetworkLinkErrorRateAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.errorrate"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Error rate calculated from ifInErrors/ifOutErrors over time
	return attr
}

func createNetworkLinkAvailabilityAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.availabilitypercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Availability calculated from uptime tracking
	return attr
}

// Network Link Metrics
func createLinkMetricsBytesTransmittedAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.bytestransmitted"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.16")) // ifOutOctets
	return attr
}

func createLinkMetricsBytesReceivedAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.bytesreceived"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.10")) // ifInOctets
	return attr
}

func createLinkMetricsPacketsTransmittedAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.packetstransmitted"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.17")) // ifOutUcastPkts
	return attr
}

func createLinkMetricsPacketsReceivedAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.packetsreceived"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.11")) // ifInUcastPkts
	return attr
}

func createLinkMetricsErrorCountAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.errorcount"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Combined error count from ifInErrors + ifOutErrors
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.14")) // ifInErrors
	return attr
}

func createLinkMetricsDropCountAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.dropcount"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Combined drop count from ifInDiscards + ifOutDiscards
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.2.2.1.13")) // ifInDiscards
	return attr
}

func createLinkMetricsJitterAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.jitterms"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Jitter not available via SNMP - requires active measurement
	return attr
}

func createLinkMetricsPacketLossAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.packetlosspercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Packet loss calculated from error/drop statistics over time
	return attr
}

func createLinkMetricsThroughputAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.throughputbps"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Throughput calculated from octets counters over time
	return attr
}

func createLinkMetricsLastMeasurementAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networklinks.metrics.lastmeasurement"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Last measurement timestamp managed by polling system
	return attr
}

// =====================================
// NETWORK HEALTH SECTION
// =====================================

// Network Health Attributes
func createHealthOverallStatusAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.overallstatus"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Overall status calculated from device and interface states
	return attr
}

func createHealthTotalDevicesAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.totaldevices"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Total devices managed by inventory system, not available via single device SNMP
	return attr
}

func createHealthOnlineDevicesAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.onlinedevices"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Online devices count managed by monitoring system
	return attr
}

func createHealthOfflineDevicesAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.offlinedevices"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Offline devices count managed by monitoring system
	return attr
}

func createHealthWarningDevicesAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.warningdevices"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Warning devices count calculated from threshold monitoring
	return attr
}

func createHealthCriticalDevicesAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.criticaldevices"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Critical devices count calculated from threshold monitoring
	return attr
}

func createHealthTotalLinksAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.totallinks"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Total links count from topology discovery
	return attr
}

func createHealthActiveLinksAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.activelinks"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Active links count from interface status monitoring
	return attr
}

func createHealthInactiveLinksAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.inactivelinks"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Inactive links count from interface status monitoring
	return attr
}

func createHealthWarningLinksAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.warninglinks"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Warning links count from utilization/error threshold monitoring
	return attr
}

func createHealthNetworkAvailabilityAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.networkavailabilitypercent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Network availability calculated from uptime statistics
	return attr
}

func createHealthLastHealthCheckAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.lasthealthcheck"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Last health check timestamp managed by monitoring system
	return attr
}

// Health Alert Attributes
func createHealthAlertIdAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alerts.alertid"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Alert ID generated by alerting system, not available via SNMP
	return attr
}

func createHealthAlertSeverityAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alerts.severity"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Alert severity calculated from threshold breaches
	return attr
}

func createHealthAlertTitleAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alerts.title"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Alert title generated by alerting rules
	return attr
}

func createHealthAlertDescriptionAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alerts.description"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Alert description generated by alerting rules
	return attr
}

func createHealthAlertAffectedComponentAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alerts.affectedcomponent"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Affected component identified by monitoring system
	return attr
}

func createHealthAlertComponentTypeAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alerts.componenttype"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Component type classification by monitoring system
	return attr
}

func createHealthAlertTimestampAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alerts.timestamp"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Alert timestamp when threshold breach detected
	return attr
}

func createHealthAlertAcknowledgedAttribute() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.networkhealth.alerts.acknowledged"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// NOTE: Acknowledgement status managed by alerting system
	return attr
}

// System MIB attribute functions for EquipmentInfo

func createSystemDescription() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.hardware"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.1.0")) // sysDescr
	return attr
}

func createSystemUptime() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.uptime"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.3.0")) // sysUpTime
	return attr
}

func createSystemLocation() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.location"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.6.0")) // sysLocation
	return attr
}

func createSystemDeviceType() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.device_type"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	// Use InferDeviceType rule that analyzes sysObjectID to determine device type
	rule := &l8poll.L8P_Rule{}
	rule.Name = "InferDeviceType"
	rule.Params = make(map[string]*l8poll.L8P_Parameter)
	// Pass the sysObjectID OID for analysis
	rule.Params["from"] = &l8poll.L8P_Parameter{Value: ".1.3.6.1.2.1.1.2.0"} // sysObjectID
	attr.Rules = append(attr.Rules, rule)
	return attr
}

func createSystemSoftware() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.software"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.1.0")) // sysDescr (extract software info)
	return attr
}

func createSystemVersion() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.1.0")) // sysDescr (extract version info)
	return attr
}

func createSystemModel() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.model"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.2.1.1.1.0")) // sysDescr (extract model info)
	return attr
}

func createIpAddress() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.ipaddress"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule("ipaddress"))
	return attr
}

func createDeviceStatus() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.devicestatus"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createDeviceStatusRule())
	return attr
}

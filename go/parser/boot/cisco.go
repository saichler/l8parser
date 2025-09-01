package boot

import (
	"github.com/saichler/l8pollaris/go/types"
)

// CreateCiscoSwitchBootPolls creates collection and parsing Pollaris model for Cisco switches
func CreateCiscoSwitchBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "cisco-switch"
	polaris.Groups = []string{"cisco", "cisco-switch"}
	polaris.Polling = make(map[string]*types.Poll)
	createCiscoSystemPoll(polaris)
	createCiscoVersionPoll(polaris)
	createCiscoSerialPoll(polaris)
	createCiscoInterfacesPoll(polaris)
	createCiscoModulesPoll(polaris)
	createCiscoPowerSupplyPoll(polaris)
	createCiscoFanPoll(polaris)
	createCiscoCpuPoll(polaris)
	createCiscoMemoryPoll(polaris)
	return polaris
}

// CreateCiscoRouterBootPolls creates collection and parsing Pollaris model for Cisco routers
func CreateCiscoRouterBootPolls() *types.Pollaris {
	polaris := &types.Pollaris{}
	polaris.Name = "cisco-router"
	polaris.Groups = []string{"cisco", "cisco-router"}
	polaris.Polling = make(map[string]*types.Poll)
	createCiscoSystemPoll(polaris)
	createCiscoVersionPoll(polaris)
	createCiscoSerialPoll(polaris)
	createCiscoInterfacesPoll(polaris)
	createCiscoRouterModulesPoll(polaris)
	createCiscoPowerSupplyPoll(polaris)
	createCiscoCpuPoll(polaris)
	createCiscoMemoryPoll(polaris)
	createCiscoRoutingPoll(polaris)
	
	// Enhanced Comprehensive Polling
	createCiscoPhysicalComponentsPoll(polaris)
	createCiscoPerformanceMetricsPoll(polaris) 
	createCiscoNetworkHealthPoll(polaris)
	createCiscoNetworkLinksPoll(polaris)
	return polaris
}

// Cisco device-specific polling functions
func createCiscoSystemPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	// Equipment Info Attributes
	poll.Attributes = append(poll.Attributes, createCiscoVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	poll.Attributes = append(poll.Attributes, createSysOid())
	poll.Attributes = append(poll.Attributes, createSysNameAttribute())
	poll.Attributes = append(poll.Attributes, createSysOidAttribute())
	poll.Attributes = append(poll.Attributes, createLocationAttribute())
	poll.Attributes = append(poll.Attributes, createUptimeAttribute())
	poll.Attributes = append(poll.Attributes, createHardwareAttribute())
	poll.Attributes = append(poll.Attributes, createSoftwareAttribute())
	poll.Attributes = append(poll.Attributes, createSeriesAttribute())
	poll.Attributes = append(poll.Attributes, createFamilyAttribute())
	p.Polling[poll.Name] = poll
}

func createCiscoVersionPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoVersion")
	poll.What = ".1.3.6.1.4.1.9.9.25.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createCiscoVersion())
	p.Polling[poll.Name] = poll
}

func createCiscoSerialPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoSerial")
	poll.What = ".1.3.6.1.4.1.9.3.6"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createCiscoSerial())
	p.Polling[poll.Name] = poll
}

func createCiscoInterfacesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	// Original Interface Attributes
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	poll.Attributes = append(poll.Attributes, createInterfaceSpeed())
	poll.Attributes = append(poll.Attributes, createInterfaceMtu())
	
	// Enhanced Physical Interface Attributes (nested in ports)
	poll.Attributes = append(poll.Attributes, createInterfaceIdAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceNameAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceStatusAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceDescriptionAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceTypeAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceSpeedAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceMacAddressAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceIpAddressAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceMtuAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceAdminStatusAttribute())
	
	// Interface Statistics
	poll.Attributes = append(poll.Attributes, createInterfaceRxPacketsAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceTxPacketsAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceRxBytesAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceTxBytesAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceRxErrorsAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceTxErrorsAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceRxDropsAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceTxDropsAttribute())
	poll.Attributes = append(poll.Attributes, createInterfaceCollisionsAttribute())
	
	// Logical Interfaces (for VLANs, Loopbacks, etc.)
	poll.Attributes = append(poll.Attributes, createLogicalInterfaceIdAttribute())
	poll.Attributes = append(poll.Attributes, createLogicalInterfaceNameAttribute())
	poll.Attributes = append(poll.Attributes, createLogicalInterfaceStatusAttribute())
	poll.Attributes = append(poll.Attributes, createLogicalInterfaceDescriptionAttribute())
	poll.Attributes = append(poll.Attributes, createLogicalInterfaceTypeAttribute())
	poll.Attributes = append(poll.Attributes, createLogicalInterfaceIpAddressAttribute())
	poll.Attributes = append(poll.Attributes, createLogicalInterfaceMtuAttribute())
	poll.Attributes = append(poll.Attributes, createLogicalInterfaceAdminStatusAttribute())
	
	p.Polling[poll.Name] = poll
}

func createCiscoModulesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoModules")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	// Original Module Attributes
	poll.Attributes = append(poll.Attributes, createModuleName())
	poll.Attributes = append(poll.Attributes, createModuleModel())
	poll.Attributes = append(poll.Attributes, createModuleStatus())
	
	// Enhanced Chassis Attributes
	poll.Attributes = append(poll.Attributes, createChassisSerialAttribute())
	poll.Attributes = append(poll.Attributes, createChassisModelAttribute())
	poll.Attributes = append(poll.Attributes, createChassisDescriptionAttribute())
	poll.Attributes = append(poll.Attributes, createChassisTemperatureAttribute())
	
	// Enhanced Module Attributes
	poll.Attributes = append(poll.Attributes, createModuleNameAttribute())
	poll.Attributes = append(poll.Attributes, createModuleModelAttribute())
	poll.Attributes = append(poll.Attributes, createModuleDescriptionAttribute())
	poll.Attributes = append(poll.Attributes, createModuleStatusAttribute())
	poll.Attributes = append(poll.Attributes, createModuleTypeAttribute())
	poll.Attributes = append(poll.Attributes, createModuleTemperatureAttribute())
	
	// CPU Attributes (nested in modules)
	poll.Attributes = append(poll.Attributes, createCpuIdAttribute())
	poll.Attributes = append(poll.Attributes, createCpuNameAttribute())
	poll.Attributes = append(poll.Attributes, createCpuModelAttribute())
	poll.Attributes = append(poll.Attributes, createCpuArchitectureAttribute())
	poll.Attributes = append(poll.Attributes, createCpuCoresAttribute())
	poll.Attributes = append(poll.Attributes, createCpuFrequencyAttribute())
	poll.Attributes = append(poll.Attributes, createCpuStatusAttribute())
	poll.Attributes = append(poll.Attributes, createCpuTemperatureAttribute())
	
	// Memory Attributes (nested in modules)
	poll.Attributes = append(poll.Attributes, createMemoryIdAttribute())
	poll.Attributes = append(poll.Attributes, createMemoryNameAttribute())
	poll.Attributes = append(poll.Attributes, createMemoryTypeAttribute())
	poll.Attributes = append(poll.Attributes, createMemorySizeAttribute())
	poll.Attributes = append(poll.Attributes, createMemoryFrequencyAttribute())
	poll.Attributes = append(poll.Attributes, createMemoryStatusAttribute())
	
	p.Polling[poll.Name] = poll
}

func createCiscoPowerSupplyPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoPowerSupply")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createPowerSupplyStatus())
	poll.Attributes = append(poll.Attributes, createPowerSupplyModel())
	p.Polling[poll.Name] = poll
}

func createCiscoFanPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoFans")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createFanStatus())
	p.Polling[poll.Name] = poll
}

func createCiscoCpuPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoCpu")
	poll.What = ".1.3.6.1.4.1.9.9.109.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createCpuUtilization())
	p.Polling[poll.Name] = poll
}

func createCiscoMemoryPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoMemory")
	poll.What = ".1.3.6.1.4.1.9.9.48.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createMemoryUtilization())
	p.Polling[poll.Name] = poll
}

func createCiscoRouterModulesPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoRouterModules")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createRouteProcessorStatus())
	p.Polling[poll.Name] = poll
}

func createCiscoRoutingPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoRouting")
	poll.What = ".1.3.6.1.2.1.4.21.1"
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createRoutingTableEntry())
	p.Polling[poll.Name] = poll
}

// Cisco-specific attribute creation functions
func createCiscoVendor() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("cisco", ".1.3.6.1.2.1.1.1.0", "Cisco"))
	return attr
}

func createCiscoVersion() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.9.25.1.1.1.2.2"))
	return attr
}

func createCiscoSerial() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.serialnumber"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.3.6.3.0"))
	return attr
}

// =====================================
// COMPREHENSIVE CISCO POLLING FUNCTIONS
// Supporting all NetworkDevice model attributes with Cisco-specific optimizations
// =====================================

// Cisco Physical Components Polling
func createCiscoPhysicalComponentsPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoPhysicalComponents")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1" // ENTITY-MIB
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	
	// Power Supply Attributes (comprehensive)
	poll.Attributes = append(poll.Attributes, createPowerSupplyIdAttribute())
	poll.Attributes = append(poll.Attributes, createPowerSupplyNameAttribute())
	poll.Attributes = append(poll.Attributes, createPowerSupplyModelAttribute())
	poll.Attributes = append(poll.Attributes, createPowerSupplySerialNumberAttribute())
	poll.Attributes = append(poll.Attributes, createPowerSupplyStatusAttribute())
	poll.Attributes = append(poll.Attributes, createPowerSupplyTemperatureAttribute())
	poll.Attributes = append(poll.Attributes, createPowerSupplyWattageAttribute()) // NOTE: Vendor-specific required
	poll.Attributes = append(poll.Attributes, createPowerSupplyPowerTypeAttribute())
	poll.Attributes = append(poll.Attributes, createPowerSupplyLoadPercentAttribute()) // NOTE: Vendor-specific required
	poll.Attributes = append(poll.Attributes, createPowerSupplyVoltageAttribute()) // NOTE: Vendor-specific required
	poll.Attributes = append(poll.Attributes, createPowerSupplyCurrentAttribute()) // NOTE: Vendor-specific required
	
	// Fan Attributes (comprehensive)  
	poll.Attributes = append(poll.Attributes, createFanIdAttribute())
	poll.Attributes = append(poll.Attributes, createFanNameAttribute())
	poll.Attributes = append(poll.Attributes, createFanDescriptionAttribute())
	poll.Attributes = append(poll.Attributes, createFanStatusAttribute())
	poll.Attributes = append(poll.Attributes, createFanSpeedRpmAttribute())
	poll.Attributes = append(poll.Attributes, createFanMaxSpeedRpmAttribute()) // NOTE: Vendor-specific required
	poll.Attributes = append(poll.Attributes, createFanTemperatureAttribute())
	poll.Attributes = append(poll.Attributes, createFanVariableSpeedAttribute()) // NOTE: Vendor-specific required
	
	// Port Attributes (comprehensive)
	poll.Attributes = append(poll.Attributes, createPortIdAttribute())
	
	p.Polling[poll.Name] = poll
}

// Cisco Performance Metrics Polling  
func createCiscoPerformanceMetricsPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoPerformanceMetrics")
	poll.What = ".1.3.6.1.2.1.25" // HOST-RESOURCES-MIB
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	
	// Performance Metrics
	poll.Attributes = append(poll.Attributes, createPerformanceCpuUsageAttribute())
	poll.Attributes = append(poll.Attributes, createPerformanceMemoryUsageAttribute())
	poll.Attributes = append(poll.Attributes, createPerformanceTemperatureAttribute())
	poll.Attributes = append(poll.Attributes, createPerformanceUptimeAttribute())
	poll.Attributes = append(poll.Attributes, createPerformanceLoadAverageAttribute()) // NOTE: UCD-SNMP-MIB if available
	poll.Attributes = append(poll.Attributes, createInterfaceCountAttribute())
	
	// Process Information
	poll.Attributes = append(poll.Attributes, createProcessNameAttribute())
	poll.Attributes = append(poll.Attributes, createProcessPidAttribute())
	poll.Attributes = append(poll.Attributes, createProcessCpuPercentAttribute()) // NOTE: Not available in standard MIB
	poll.Attributes = append(poll.Attributes, createProcessMemoryPercentAttribute()) // NOTE: Not available in standard MIB
	poll.Attributes = append(poll.Attributes, createProcessStatusAttribute())
	
	p.Polling[poll.Name] = poll
}

// Cisco Network Health Polling
func createCiscoNetworkHealthPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoNetworkHealth")
	poll.What = ".1.3.6.1.2.1.2" // IF-MIB for interface health
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	
	// Network Health Attributes (NOTE: Many calculated by management system)
	poll.Attributes = append(poll.Attributes, createHealthOverallStatusAttribute())
	poll.Attributes = append(poll.Attributes, createHealthTotalDevicesAttribute()) // NOTE: Management system
	poll.Attributes = append(poll.Attributes, createHealthOnlineDevicesAttribute()) // NOTE: Management system  
	poll.Attributes = append(poll.Attributes, createHealthOfflineDevicesAttribute()) // NOTE: Management system
	poll.Attributes = append(poll.Attributes, createHealthWarningDevicesAttribute()) // NOTE: Management system
	poll.Attributes = append(poll.Attributes, createHealthCriticalDevicesAttribute()) // NOTE: Management system
	poll.Attributes = append(poll.Attributes, createHealthTotalLinksAttribute()) // NOTE: Topology discovery
	poll.Attributes = append(poll.Attributes, createHealthActiveLinksAttribute()) // NOTE: Interface monitoring
	poll.Attributes = append(poll.Attributes, createHealthInactiveLinksAttribute()) // NOTE: Interface monitoring
	poll.Attributes = append(poll.Attributes, createHealthWarningLinksAttribute()) // NOTE: Threshold monitoring
	poll.Attributes = append(poll.Attributes, createHealthNetworkAvailabilityAttribute()) // NOTE: Calculated
	poll.Attributes = append(poll.Attributes, createHealthLastHealthCheckAttribute()) // NOTE: Management system
	
	// Health Alert Attributes (NOTE: Generated by alerting system)
	poll.Attributes = append(poll.Attributes, createHealthAlertIdAttribute())
	poll.Attributes = append(poll.Attributes, createHealthAlertSeverityAttribute())
	poll.Attributes = append(poll.Attributes, createHealthAlertTitleAttribute())
	poll.Attributes = append(poll.Attributes, createHealthAlertDescriptionAttribute())
	poll.Attributes = append(poll.Attributes, createHealthAlertAffectedComponentAttribute())
	poll.Attributes = append(poll.Attributes, createHealthAlertComponentTypeAttribute())
	poll.Attributes = append(poll.Attributes, createHealthAlertTimestampAttribute())
	poll.Attributes = append(poll.Attributes, createHealthAlertAcknowledgedAttribute())
	
	p.Polling[poll.Name] = poll
}

// Cisco Network Links Polling
func createCiscoNetworkLinksPoll(p *types.Pollaris) {
	poll := createBaseSNMPPoll("ciscoNetworkLinks")
	poll.What = ".1.3.6.1.2.1.2" // IF-MIB for link information
	poll.Operation = types.Operation_OMap
	poll.Attributes = make([]*types.Attribute, 0)
	
	// Network Link Attributes (NOTE: Many derived from topology discovery)
	poll.Attributes = append(poll.Attributes, createNetworkLinkIdAttribute()) // NOTE: Topology discovery
	poll.Attributes = append(poll.Attributes, createNetworkLinkNameAttribute()) // NOTE: Management system
	poll.Attributes = append(poll.Attributes, createNetworkLinkFromNodeAttribute()) // NOTE: Topology discovery
	poll.Attributes = append(poll.Attributes, createNetworkLinkToNodeAttribute()) // NOTE: Topology discovery
	poll.Attributes = append(poll.Attributes, createNetworkLinkStatusAttribute()) // From ifOperStatus
	poll.Attributes = append(poll.Attributes, createNetworkLinkBandwidthAttribute()) // From ifSpeed
	poll.Attributes = append(poll.Attributes, createNetworkLinkTypeAttribute()) // From ifType
	poll.Attributes = append(poll.Attributes, createNetworkLinkUtilizationAttribute()) // NOTE: Calculated
	poll.Attributes = append(poll.Attributes, createNetworkLinkLatencyAttribute()) // NOTE: Active measurement
	poll.Attributes = append(poll.Attributes, createNetworkLinkDistanceAttribute()) // NOTE: Geographic calculation
	poll.Attributes = append(poll.Attributes, createNetworkLinkUptimeAttribute()) // NOTE: Calculated
	poll.Attributes = append(poll.Attributes, createNetworkLinkErrorRateAttribute()) // NOTE: Calculated
	poll.Attributes = append(poll.Attributes, createNetworkLinkAvailabilityAttribute()) // NOTE: Calculated
	
	// Network Link Metrics (from interface statistics)
	poll.Attributes = append(poll.Attributes, createLinkMetricsBytesTransmittedAttribute()) // ifOutOctets
	poll.Attributes = append(poll.Attributes, createLinkMetricsBytesReceivedAttribute()) // ifInOctets
	poll.Attributes = append(poll.Attributes, createLinkMetricsPacketsTransmittedAttribute()) // ifOutUcastPkts
	poll.Attributes = append(poll.Attributes, createLinkMetricsPacketsReceivedAttribute()) // ifInUcastPkts
	poll.Attributes = append(poll.Attributes, createLinkMetricsErrorCountAttribute()) // ifInErrors + ifOutErrors
	poll.Attributes = append(poll.Attributes, createLinkMetricsDropCountAttribute()) // ifInDiscards + ifOutDiscards
	poll.Attributes = append(poll.Attributes, createLinkMetricsJitterAttribute()) // NOTE: Active measurement
	poll.Attributes = append(poll.Attributes, createLinkMetricsPacketLossAttribute()) // NOTE: Calculated
	poll.Attributes = append(poll.Attributes, createLinkMetricsThroughputAttribute()) // NOTE: Calculated
	poll.Attributes = append(poll.Attributes, createLinkMetricsLastMeasurementAttribute()) // NOTE: Management system
	
	p.Polling[poll.Name] = poll
}
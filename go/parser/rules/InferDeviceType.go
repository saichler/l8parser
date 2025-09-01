package rules

import (
	"fmt"
	"strings"

	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
)

type InferDeviceType struct{}

func (this *InferDeviceType) Name() string {
	return "InferDeviceType"
}

func (this *InferDeviceType) ParamNames() []string {
	return []string{}
}

// DeviceType enum values from proto
const (
	DEVICE_TYPE_UNKNOWN       = 0
	DEVICE_TYPE_ROUTER        = 1
	DEVICE_TYPE_SWITCH        = 2
	DEVICE_TYPE_FIREWALL      = 3
	DEVICE_TYPE_LOAD_BALANCER = 4
	DEVICE_TYPE_ACCESS_POINT  = 5
	DEVICE_TYPE_SERVER        = 6
	DEVICE_TYPE_STORAGE       = 7
	DEVICE_TYPE_GATEWAY       = 8
)

func (this *InferDeviceType) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*types.Parameter, any interface{}, pollWhat string) error {
	// Get the NetworkDevice
	networkDevice, ok := any.(*types2.NetworkDevice)
	if !ok {
		return resources.Logger().Error("Target object is not a NetworkDevice")
	}

	// Ensure Equipmentinfo exists (note: lowercase 'i')
	if networkDevice.Equipmentinfo == nil {
		networkDevice.Equipmentinfo = &types2.EquipmentInfo{}
	}

	// Get sysObjectID from input using GetValueInput (like other parsing rules)
	input := workSpace[Input]
	if input == nil {
		return resources.Logger().Error("nil input for InferDeviceType")
	}

	// Get the sysObjectID value using the same pattern as Set rule
	value, _, err := GetValueInput(resources, input, params, pollWhat)
	if err != nil {
		return resources.Logger().Error("Error getting sysObjectID:", err)
	}

	if value == nil {
		return resources.Logger().Error("nil sysObjectID value")
	}

	sysObjectID := convertInterfaceToString(value)
	fmt.Printf("DEBUG InferDeviceType: sysObjectID='%s'\n", sysObjectID)

	// Infer device type based on sysObjectID patterns
	deviceType := inferDeviceTypeFromOID(sysObjectID)
	networkDevice.Equipmentinfo.DeviceType = types2.DeviceType(deviceType)

	fmt.Printf("DEBUG InferDeviceType: Inferred device type: %d\n", deviceType)
	return nil
}

func inferDeviceTypeFromOID(sysObjectID string) int32 {
	sysObjectIDLower := strings.ToLower(sysObjectID)

	// OID-based detection (vendor-specific OIDs)
	// Cisco (1.3.6.1.4.1.9) - Most common enterprise network vendor
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9") {
		// Cisco routers typically have specific sub-OIDs
		if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9.1.1") ||  // Old Cisco routers
		   strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9.12") ||   // ISR/ASR series
		   strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9.1.222") { // Newer routers
			return DEVICE_TYPE_ROUTER
		}
		// Cisco switches
		if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9.1.23") ||  // Catalyst switches
		   strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9.1.516") || // Catalyst 2960
		   strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9.1.717") {  // Catalyst 3750
			return DEVICE_TYPE_SWITCH
		}
		// Cisco ASA firewalls
		if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9.1.745") || // ASA 5500 series
		   strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9.1.1069") { // ASA 5585
			return DEVICE_TYPE_FIREWALL
		}
		// Default Cisco to router (most common)
		return DEVICE_TYPE_ROUTER
	}

	// Juniper (1.3.6.1.4.1.2636) - Primarily routing/switching
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.2636") {
		// Juniper is primarily routing equipment
		return DEVICE_TYPE_ROUTER
	}

	// Palo Alto (1.3.6.1.4.1.25461) - Firewalls
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.25461") {
		return DEVICE_TYPE_FIREWALL
	}

	// Fortinet (1.3.6.1.4.1.12356) - Firewalls
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.12356") {
		return DEVICE_TYPE_FIREWALL
	}

	// F5 (1.3.6.1.4.1.3375) - Load Balancers
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.3375") {
		return DEVICE_TYPE_LOAD_BALANCER
	}

	// Arista (1.3.6.1.4.1.30065) - Switches
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.30065") {
		return DEVICE_TYPE_SWITCH
	}

	// Dell servers (1.3.6.1.4.1.674)
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.674") {
		return DEVICE_TYPE_SERVER
	}

	// HP/HPE (1.3.6.1.4.1.232)
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.232") {
		// HP makes both servers and network equipment
		// Default to server for HP as that's more common for SNMP management
		return DEVICE_TYPE_SERVER
	}

	// Ubiquiti (1.3.6.1.4.1.41112) - Access Points
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.41112") {
		return DEVICE_TYPE_ACCESS_POINT
	}

	// NetApp (1.3.6.1.4.1.789) - Storage
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.789") {
		return DEVICE_TYPE_STORAGE
	}

	// EMC (1.3.6.1.4.1.1139) - Storage
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.1139") {
		return DEVICE_TYPE_STORAGE
	}

	// Default to unknown if no patterns match
	return DEVICE_TYPE_UNKNOWN
}

func inferDeviceTypeFromData(sysDescr string, sysObjectID string) int32 {
	sysDescrLower := strings.ToLower(sysDescr)
	sysObjectIDLower := strings.ToLower(sysObjectID)

	// Router patterns (highest priority for routers)
	routerPatterns := []string{
		"router", "asr", "isr", "7200", "7300", "7400", "7500", "7600", 
		"nexus", "mx", "ex", "qfx", "srx", "vmx", "vqfx", "routing",
		"ospf", "bgp", "mpls", "cisco ios", "junos", "nx-os",
	}

	// Switch patterns  
	switchPatterns := []string{
		"switch", "catalyst", "3750", "3560", "2960", "4500", "6500",
		"switching", "vlan", "spanning-tree", "stp", "ethernet switch",
	}

	// Firewall patterns
	firewallPatterns := []string{
		"firewall", "asa", "pix", "fortigate", "palo alto", "checkpoint", 
		"fortios", "pan-os", "security", "utm", "ngfw", "threat",
	}

	// Load balancer patterns
	loadBalancerPatterns := []string{
		"load balancer", "f5", "big-ip", "netscaler", "citrix", "a10", 
		"loadbalancer", "application delivery", "adc",
	}

	// Access point patterns
	accessPointPatterns := []string{
		"access point", "wireless", "wifi", "ap", "wap", "aironet", 
		"aruba", "ubiquiti", "unifi",
	}

	// Server patterns
	serverPatterns := []string{
		"server", "dell", "hp", "hpe", "supermicro", "lenovo", "ibm",
		"proliant", "poweredge", "system x", "blade", "rack",
	}

	// Storage patterns
	storagePatterns := []string{
		"storage", "san", "nas", "netapp", "emc", "dell emc", "pure storage",
		"hitachi", "hds", "vnx", "unity", "isilon", "compellent",
	}

	// Gateway patterns
	gatewayPatterns := []string{
		"gateway", "edge", "branch", "wan optimizer", "sd-wan", "viptela",
	}

	// OID-based detection (vendor-specific OIDs)
	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.9") { // Cisco
		if containsAnyPattern(sysDescrLower, routerPatterns) {
			return DEVICE_TYPE_ROUTER
		}
		if containsAnyPattern(sysDescrLower, switchPatterns) {
			return DEVICE_TYPE_SWITCH  
		}
		if containsAnyPattern(sysDescrLower, firewallPatterns) {
			return DEVICE_TYPE_FIREWALL
		}
	}

	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.2636") { // Juniper
		return DEVICE_TYPE_ROUTER // Juniper is primarily routing/switching
	}

	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.25461") { // Palo Alto
		return DEVICE_TYPE_FIREWALL
	}

	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.12356") { // Fortinet
		return DEVICE_TYPE_FIREWALL
	}

	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.3375") { // F5
		return DEVICE_TYPE_LOAD_BALANCER
	}

	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.674") { // Dell
		return DEVICE_TYPE_SERVER
	}

	if strings.Contains(sysObjectIDLower, "1.3.6.1.4.1.232") { // HP/HPE  
		return DEVICE_TYPE_SERVER
	}

	// Pattern-based detection (fallback to description analysis)
	if containsAnyPattern(sysDescrLower, routerPatterns) {
		return DEVICE_TYPE_ROUTER
	}

	if containsAnyPattern(sysDescrLower, switchPatterns) {
		return DEVICE_TYPE_SWITCH
	}

	if containsAnyPattern(sysDescrLower, firewallPatterns) {
		return DEVICE_TYPE_FIREWALL
	}

	if containsAnyPattern(sysDescrLower, loadBalancerPatterns) {
		return DEVICE_TYPE_LOAD_BALANCER
	}

	if containsAnyPattern(sysDescrLower, accessPointPatterns) {
		return DEVICE_TYPE_ACCESS_POINT
	}

	if containsAnyPattern(sysDescrLower, serverPatterns) {
		return DEVICE_TYPE_SERVER
	}

	if containsAnyPattern(sysDescrLower, storagePatterns) {
		return DEVICE_TYPE_STORAGE
	}

	if containsAnyPattern(sysDescrLower, gatewayPatterns) {
		return DEVICE_TYPE_GATEWAY
	}

	// Default to unknown if no patterns match
	return DEVICE_TYPE_UNKNOWN
}

func containsAnyPattern(text string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}
	return false
}

func convertInterfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}
	
	// Handle byte array to string conversion
	if byteArray, ok := value.([]uint8); ok {
		return strings.TrimSpace(string(byteArray))
	}
	
	return strings.TrimSpace(fmt.Sprintf("%v", value))
}
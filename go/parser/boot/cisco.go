/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package boot

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

// CreateCiscoSwitchBootPolls creates collection and parsing Pollaris model for Cisco switches
func CreateCiscoSwitchBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "cisco-switch"
	polaris.Groups = []string{"cisco", "cisco-switch"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createCiscoSystemPoll(polaris)
	createCiscoVersionPoll(polaris)
	createCiscoSerialPoll(polaris)
	createCiscoInterfacesPoll(polaris)
	createCiscoCpuPoll(polaris)
	createCiscoMemoryPoll(polaris)
	return polaris
}

// CreateCiscoRouterBootPolls creates collection and parsing Pollaris model for Cisco routers
func CreateCiscoRouterBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "cisco-router"
	polaris.Groups = []string{"cisco", "cisco-router"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createCiscoSystemPoll(polaris)
	createCiscoVersionPoll(polaris)
	createCiscoSerialPoll(polaris)
	createCiscoInterfacesPoll(polaris)
	createCiscoCpuPoll(polaris)
	createCiscoMemoryPoll(polaris)
	createCiscoRoutingPoll(polaris)
	return polaris
}

// Cisco device-specific polling functions
func createCiscoSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("ciscoSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
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

func createCiscoVersionPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("ciscoVersion")
	poll.What = ".1.3.6.1.4.1.9.9.25.1.1.1.2.2"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createCiscoVersion())
	p.Polling[poll.Name] = poll
}

func createCiscoSerialPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("ciscoSerial")
	poll.What = ".1.3.6.1.4.1.9.3.6.3.0"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createCiscoSerial())
	p.Polling[poll.Name] = poll
}

func createCiscoInterfacesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("ciscoInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
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

func createCiscoCpuPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("ciscoCpu")
	poll.What = ".1.3.6.1.4.1.9.9.109.1.1.1.1.5.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createCpuUtilization())
	p.Polling[poll.Name] = poll
}

func createCiscoMemoryPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("ciscoMemory")
	poll.What = ".1.3.6.1.4.1.9.9.48.1.1.1.6.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createMemoryUtilization())
	p.Polling[poll.Name] = poll
}

func createCiscoRoutingPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("ciscoRouting")
	poll.What = ".1.3.6.1.2.1.4.21.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createRoutingTableEntry())
	p.Polling[poll.Name] = poll
}

// Cisco-specific attribute creation functions
func createCiscoVendor() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("cisco", ".1.3.6.1.2.1.1.1.0", "Cisco"))
	return attr
}

func createCiscoVersion() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.9.25.1.1.1.2.2"))
	return attr
}

func createCiscoSerial() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.serialnumber"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.9.3.6.3.0"))
	return attr
}


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

// CreateJuniperRouterBootPolls creates collection and parsing Pollaris model for Juniper routers
func CreateJuniperRouterBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "juniper-router"
	polaris.Groups = []string{"juniper", "juniper-router"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createJuniperSystemPoll(polaris)
	createJuniperMibSystemPoll(polaris)
	createJuniperInterfacesPoll(polaris)
	createJuniperChassisComponentsPoll(polaris)
	createJuniperRoutingEnginePoll(polaris)
	return polaris
}

// Juniper device-specific polling functions
func createJuniperSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("juniperSystem")
	poll.What = ".1.3.6.1.4.1.2636.3.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createJuniperVersion())
	p.Polling[poll.Name] = poll
}

func createJuniperMibSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("juniperMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createJuniperVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createJuniperInterfacesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("juniperInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	poll.Attributes = append(poll.Attributes, createInterfaceSpeed())
	p.Polling[poll.Name] = poll
}

func createJuniperChassisComponentsPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("juniperChassis")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createChassisComponentStatus())
	p.Polling[poll.Name] = poll
}

func createJuniperRoutingEnginePoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("juniperRoutingEngine")
	poll.What = ".1.3.6.1.4.1.2636.3.1.13.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createRoutingEngineUtilization())
	p.Polling[poll.Name] = poll
}

// Juniper-specific attribute creation functions
func createJuniperVendor() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("juniper", ".1.3.6.1.2.1.1.1.0", "Juniper"))
	return attr
}

func createJuniperVersion() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.2636.3.1.2.0"))
	return attr
}

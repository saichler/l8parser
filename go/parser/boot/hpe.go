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

// CreateHPEServerBootPolls creates collection and parsing Pollaris model for HPE servers
func CreateHPEServerBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "hpe-server"
	polaris.Groups = []string{"hpe", "hpe-server"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createHPESystemPoll(polaris)
	createHPEMibSystemPoll(polaris)
	createHPEStoragePoll(polaris)
	createHPEPowerThermalPoll(polaris)
	return polaris
}

// HPE server-specific polling functions
func createHPESystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("hpeSystem")
	poll.What = ".1.3.6.1.4.1.232.2.2.4"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createHPEVersion())
	p.Polling[poll.Name] = poll
}

func createHPEMibSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("hpeMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createHPEVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createHPEStoragePoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("hpeStorage")
	poll.What = ".1.3.6.1.2.1.25.2.3.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createDiskStatus())
	p.Polling[poll.Name] = poll
}

func createHPEPowerThermalPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("hpePowerThermal")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createPowerSupplyStatus())
	poll.Attributes = append(poll.Attributes, createTemperatureSensors())
	p.Polling[poll.Name] = poll
}

// HPE-specific attribute creation functions
func createHPEVendor() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("hpe", ".1.3.6.1.2.1.1.1.0", "Hewlett Packard Enterprise"))
	return attr
}

func createHPEVersion() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.232.2.2.4.2.0"))
	return attr
}

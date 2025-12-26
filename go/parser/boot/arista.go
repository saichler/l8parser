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

// CreateAristaSwitchBootPolls creates collection and parsing Pollaris model for Arista switches
func CreateAristaSwitchBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "arista-switch"
	polaris.Groups = []string{"arista", "arista-switch"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createAristaSystemPoll(polaris)
	createAristaMibSystemPoll(polaris)
	createAristaInterfacesPoll(polaris)
	createAristaEnvironmentalPoll(polaris)
	return polaris
}

// Arista device-specific polling functions
func createAristaSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("aristaSystem")
	poll.What = ".1.3.6.1.4.1.30065.1.3.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createAristaVersion())
	p.Polling[poll.Name] = poll
}

func createAristaMibSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("aristaMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createAristaVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createAristaInterfacesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("aristaInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	poll.Attributes = append(poll.Attributes, createInterfaceSpeed())
	p.Polling[poll.Name] = poll
}

func createAristaEnvironmentalPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("aristaEnvironmental")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createTemperatureSensors())
	p.Polling[poll.Name] = poll
}

// Arista-specific attribute creation functions
func createAristaVendor() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("arista", ".1.3.6.1.2.1.1.1.0", "Arista"))
	return attr
}

func createAristaVersion() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.30065.1.3.1.1.0"))
	return attr
}

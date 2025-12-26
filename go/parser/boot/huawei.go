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

// CreateHuaweiRouterBootPolls creates collection and parsing Pollaris model for Huawei routers
func CreateHuaweiRouterBootPolls() *l8tpollaris.L8Pollaris {
	polaris := &l8tpollaris.L8Pollaris{}
	polaris.Name = "huawei-router"
	polaris.Groups = []string{"huawei", "huawei-router"}
	polaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createHuaweiSystemPoll(polaris)
	createHuaweiMibSystemPoll(polaris)
	createHuaweiInterfacesPoll(polaris)
	createHuaweiEnvironmentalPoll(polaris)
	return polaris
}

// Huawei device-specific polling functions
func createHuaweiSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("huaweiSystem")
	poll.What = ".1.3.6.1.4.1.2011.5.25.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createHuaweiVersion())
	p.Polling[poll.Name] = poll
}

func createHuaweiMibSystemPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("huaweiMibSystem")
	poll.What = ".1.3.6.1.2.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createHuaweiVendor())
	poll.Attributes = append(poll.Attributes, createSysName())
	p.Polling[poll.Name] = poll
}

func createHuaweiInterfacesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("huaweiInterfaces")
	poll.What = ".1.3.6.1.2.1.2.2.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createInterfaceName())
	poll.Attributes = append(poll.Attributes, createInterfaceStatus())
	p.Polling[poll.Name] = poll
}

func createHuaweiEnvironmentalPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseSNMPPoll("huaweiEnvironmental")
	poll.What = ".1.3.6.1.2.1.47.1.1.1.1"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Map
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createTemperatureSensors())
	p.Polling[poll.Name] = poll
}

// Huawei-specific attribute creation functions
func createHuaweiVendor() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.vendor"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createContainsRule("huawei", ".1.3.6.1.2.1.1.1.0", "Huawei"))
	return attr
}

func createHuaweiVersion() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "networkdevice.equipmentinfo.version"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createSetRule(".1.3.6.1.4.1.2011.5.25.1.1.1.0"))
	return attr
}

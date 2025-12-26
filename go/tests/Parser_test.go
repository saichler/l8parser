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

package tests

import (
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	common2 "github.com/saichler/probler/go/prob/common"
	"github.com/saichler/probler/go/prob/common/creates"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/common"

	"github.com/saichler/l8collector/go/collector/service"
	"github.com/saichler/l8inventory/go/tests/utils_inventory"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
)

// TestMain sets up and tears down the test environment for all tests in this package.
func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

// TestParser tests the end-to-end parsing flow with a simulated network device.
// It sets up the Pollaris, Collector, and Parser services and verifies
// that data flows correctly through the parsing pipeline.
func TestParser(t *testing.T) {

	linksId := common2.NetworkDevice_Links_ID

	all := boot.GetAllPolarisModels()
	common.SmoothFirstCollection = true
	for _, snmpPolls := range all {
		for _, poll := range snmpPolls.Polling {
			if poll.Cadence.Enabled {
				poll.Cadence.Cadences[0] = 3
			}
		}
	}

	initData := []interface{}{}
	for _, snmpPolls := range all {
		initData = append(initData, snmpPolls)
	}

	//use opensim to simulate this device with this ip
	//https://github.com/saichler/opensim
	//curl -X POST http://localhost:8080/api/v1/devices -H "Content-Type: application/json" -d '{"start_ip":"10.10.10.1","device_count":3,"netmask":"24"}'
	device := creates.CreateDevice("10.20.30.1", linksId, "sim")

	vnic := topo.VnicByVnetNum(2, 2)
	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, pollaris.ServiceArea, true, nil)
	sla.SetInitItems(initData)
	vnic.Resources().Services().Activate(sla, vnic)

	targets.Activate("postgres", "probler", vnic)

	collSN, collSA := targets.Links.Collector(linksId)
	sla = ifs.NewServiceLevelAgreement(&service.CollectorService{}, collSN, collSA, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	parSN, parSA := targets.Links.Parser(linksId)
	sla = ifs.NewServiceLevelAgreement(&parsing.ParsingService{}, parSN, parSA, true, nil)
	sla.SetServiceItem(&types2.NetworkDevice{})
	sla.SetPrimaryKeys("Id")
	sla.SetArgs(false)
	vnic.Resources().Services().Activate(sla, vnic)

	boxSN, boxSA := targets.Links.Cache(linksId)
	sla = ifs.NewServiceLevelAgreement(&utils_inventory.MockOrmService{}, boxSN, boxSA, false, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	cl := topo.VnicByVnetNum(1, 1)
	err := cl.Multicast(targets.ServiceName, targets.ServiceArea, ifs.POST, device)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 6)

	m, ok := vnic.Resources().Services().ServiceHandler(boxSN, boxSA)
	if !ok {
		vnic.Resources().Logger().Fail(t, "Cannot find mock service")
		return
	}
	mock := m.(*utils_inventory.MockOrmService)
	if mock.PatchCount() < 1 {
		vnic.Resources().Logger().Fail(t, "Expected 1 patch count in mock ", mock.PatchCount())
		return
	}

}

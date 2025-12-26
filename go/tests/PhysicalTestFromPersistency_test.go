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
	"fmt"
	"github.com/saichler/l8collector/go/tests/utils_collector"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	common2 "github.com/saichler/probler/go/prob/common"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/service"
	inventory "github.com/saichler/l8inventory/go/inv/service"
	"github.com/saichler/l8inventory/go/tests/utils_inventory"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
)

// TestPhysicalFromPersistency tests parsing from previously persisted job files.
// This enables replay of historical collection data for debugging and regression testing.
func TestPhysicalFromPersistency(t *testing.T) {

	linksId := common2.NetworkDevice_Links_ID

	//use opensim to simulate this device with this ip
	//https://github.com/saichler/opensim
	//curl -X POST http://localhost:8080/api/v1/devices -H "Content-Type: application/json" -d '{"start_ip":"10.10.10.1","device_count":3,"netmask":"24"}'
	//device := creates.CreateDevice("10.20.30.3", linksId, "sim")

	vnic := topo.VnicByVnetNum(2, 2)
	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, pollaris.ServiceArea, true, nil)
	utils_collector.SetPolls(sla)
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
	sla = ifs.NewServiceLevelAgreement(&utils_inventory.MockOrmService{}, boxSN, boxSA, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	time.Sleep(time.Second)

	job, err := parsing.LoadJob("mib2", "entityMib", "10.20.30.3", "10.20.30.3")
	if err != nil {
		vnic.Resources().Logger().Fail(t, err.Error())
		return

	}
	ps, _ := vnic.Resources().Services().ServiceHandler(parSN, parSA)
	parserService := ps.(*parsing.ParsingService)
	jobElem := object.New(nil, job)
	parserService.Post(jobElem, vnic)

	time.Sleep(time.Second)

	inv := inventory.Inventory(vnic.Resources(), boxSN, boxSA)
	filter := &types2.NetworkDevice{Id: "10.20.30.3"}
	elem := inv.ElementByElement(filter)
	networkDevice := elem.(*types2.NetworkDevice)

	fmt.Printf("DEBUG: NetworkDevice has %d physicals\n", len(networkDevice.Physicals))

	if len(networkDevice.Physicals) == 0 {
		vnic.Resources().Logger().Fail(t, "No physicals found in NetworkDevice")
		return
	}

	for physicalKey, physical := range networkDevice.Physicals {
		fmt.Printf("DEBUG: Physical key '%s' has %d ports\n", physicalKey, len(physical.Ports))

		if physical.Ports == nil || len(physical.Ports) < 2 {
			fmt.Printf("DEBUG: Physical '%s' has insufficient ports. Expected > 2, got %d\n", physicalKey, len(physical.Ports))

			// Let's check what we actually have
			if physical.Ports != nil {
				for portKey, port := range physical.Ports {
					fmt.Printf("DEBUG: Port key '%v': %+v\n", portKey, port)
					if port.Interfaces != nil {
						fmt.Printf("DEBUG: Port '%v' has %d interfaces\n", portKey, len(port.Interfaces))
						for ifKey, iface := range port.Interfaces {
							fmt.Printf("DEBUG: Interface '%v': Name='%s', ID='%s'\n", ifKey, iface.Name, iface.Id)
						}
					} else {
						fmt.Printf("DEBUG: Port '%v' has no interfaces\n", portKey)
					}
				}
			}

			vnic.Resources().Logger().Fail(t, "Expected more ports")
			return
		}
		for portKey, port := range physical.Ports {
			fmt.Printf("Port '%v': %+v\n", portKey, port)
		}
	}
}

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
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	common2 "github.com/saichler/probler/go/prob/common"
	"github.com/saichler/probler/go/prob/common/creates"
	"os"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8collector/go/collector/service"
	"github.com/saichler/l8collector/go/tests/utils_collector"
	inventory "github.com/saichler/l8inventory/go/inv/service"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestPhysical tests the physical inventory parsing capabilities.
// It verifies that Entity MIB and IF-MIB data is correctly parsed
// into NetworkDevice physical structures (chassis, ports, interfaces).
func TestPhysical(t *testing.T) {
	linksId := common2.NetworkDevice_Links_ID
	ip := "10.20.30.3"
	common.SmoothFirstCollection = true

	//use opensim to simulate this device with this ip
	//https://github.com/saichler/opensim
	//curl -X POST http://localhost:8080/api/v1/devices -H "Content-Type: application/json" -d '{"start_ip":"10.10.10.1","device_count":3,"netmask":"24"}'
	device := creates.CreateDevice(ip, linksId, "sim")
	device.State = l8tpollaris.L8PTargetState_Up

	vnic := topo.VnicByVnetNum(2, 2)

	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, pollaris.ServiceArea, true, nil)
	utils_collector.SetPolls(sla)
	vnic.Resources().Services().Activate(sla, vnic)

	targets.Activate("admin", "admin", vnic)

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
	sla = ifs.NewServiceLevelAgreement(&inventory.InventoryService{}, boxSN, boxSA, true, nil)
	sla.SetServiceItem(&types2.NetworkDevice{})
	sla.SetServiceItemList(&types2.NetworkDeviceList{})
	sla.SetPrimaryKeys("Id")
	//sla.SetArgs(forwardInfo)
	vnic.Resources().Services().Activate(sla, vnic)

	time.Sleep(time.Second)

	cl := topo.VnicByVnetNum(1, 1)
	cl.Multicast(targets.ServiceName, targets.ServiceArea, ifs.POST, device)

	time.Sleep(time.Second * 20)

	inv := inventory.Inventory(vnic.Resources(), boxSN, boxSA)
	filter := &types2.NetworkDevice{Id: device.TargetId}
	elem := inv.ElementByElement(filter)
	networkDevice := elem.(*types2.NetworkDevice)

	marshalOptions := protojson.MarshalOptions{
		UseEnumNumbers: true,
	}
	jsn, _ := marshalOptions.Marshal(networkDevice)
	os.WriteFile("/tmp/NetworkDevice.json", jsn, 0644)

	if len(networkDevice.Physicals) == 0 {
		vnic.Resources().Logger().Fail(t, "No physicals found in NetworkDevice")
		return
	}

	if networkDevice.Equipmentinfo.DeviceType == types2.DeviceType_DEVICE_TYPE_UNKNOWN {
		vnic.Resources().Logger().Fail(t, "Unknown device type")
		return
	}

	for _, physical := range networkDevice.Physicals {
		if physical.Performance == nil || physical.Performance.CpuUsagePercent == nil {
			vnic.Resources().Logger().Fail(t, "No performance")
			return
		}
		if physical.Ports == nil || len(physical.Ports) < 2 {
			vnic.Resources().Logger().Fail(t, "Expected more ports")
			return
		}
	}
	if networkDevice.Equipmentinfo.DeviceType == types2.DeviceType_DEVICE_TYPE_UNKNOWN {
		vnic.Resources().Logger().Fail(t, "Unknown device type")
		return
	}
	if networkDevice.Equipmentinfo.IpAddress == "" {
		vnic.Resources().Logger().Fail(t, "Unknown device ip address")
		return
	}
	if networkDevice.Equipmentinfo.DeviceStatus == types2.DeviceStatus_DEVICE_STATUS_UNKNOWN {
		vnic.Resources().Logger().Fail(t, "Unknown device status")
		return
	}
}

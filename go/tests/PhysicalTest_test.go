package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8collector/go/collector/service"
	"github.com/saichler/l8collector/go/collector/targets"
	"github.com/saichler/l8collector/go/tests/utils_collector"
	inventory "github.com/saichler/l8inventory/go/inv/service"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestPhysical(t *testing.T) {

	ip := "10.20.30.9"
	common.SmoothFirstCollection = true
	serviceArea := byte(0)
	allPolls := boot.GetAllPolarisModels()
	for _, snmpPolls := range allPolls {
		for _, poll := range snmpPolls.Polling {
			if poll.Cadence.Enabled {
				poll.Cadence.Cadences[0] = 3
			}
		}
	}

	//use opensim to simulate this device with this ip
	//https://github.com/saichler/opensim
	//curl -X POST http://localhost:8080/api/v1/devices -H "Content-Type: application/json" -d '{"start_ip":"10.10.10.1","device_count":3,"netmask":"24"}'
	device := utils_collector.CreateDevice(ip, serviceArea)

	vnic := topo.VnicByVnetNum(2, 2)

	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, serviceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	p := pollaris.Pollaris(vnic.Resources())
	for _, snmpPolls := range allPolls {
		err := p.Post(snmpPolls, false)
		if err != nil {
			vnic.Resources().Logger().Fail(t, err.Error())
			return
		}
	}

	sla = ifs.NewServiceLevelAgreement(&targets.TargetService{}, targets.ServiceName, serviceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&service.CollectorService{}, common.CollectorService, serviceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&parsing.ParsingService{}, device.LinkParser.ZsideServiceName, byte(device.LinkParser.ZsideServiceArea), true, nil)
	sla.SetServiceItem(&types2.NetworkDevice{})
	sla.SetPrimaryKeys([]string{"Id"})
	sla.SetArgs(false)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&inventory.InventoryService{}, device.LinkData.ZsideServiceName, byte(device.LinkData.ZsideServiceArea), true, nil)
	sla.SetServiceItem(&types2.NetworkDevice{})
	sla.SetServiceItemList(&types2.NetworkDeviceList{})
	sla.SetPrimaryKeys([]string{"Id"})
	//sla.SetArgs(forwardInfo)
	vnic.Resources().Services().Activate(sla, vnic)

	time.Sleep(time.Second)

	cl := topo.VnicByVnetNum(1, 1)
	cl.Multicast(targets.ServiceName, serviceArea, ifs.POST, device)

	time.Sleep(time.Second * 20)

	inv := inventory.Inventory(vnic.Resources(), device.LinkData.ZsideServiceName, byte(device.LinkData.ZsideServiceArea))
	filter := &types2.NetworkDevice{Id: ip}
	elem := inv.ElementByElement(filter)
	networkDevice := elem.(*types2.NetworkDevice)

	marshalOptions := protojson.MarshalOptions{
		UseEnumNumbers: true,
	}
	jsn, _ := marshalOptions.Marshal(networkDevice)
	os.WriteFile("/tmp/NetworkDevice.json", jsn, 0644)

	fmt.Printf("DEBUG: NetworkDevice has %d physicals\n", len(networkDevice.Physicals))

	if len(networkDevice.Physicals) == 0 {
		vnic.Resources().Logger().Fail(t, "No physicals found in NetworkDevice")
		return
	}

	if networkDevice.Equipmentinfo.DeviceType == types2.DeviceType_DEVICE_TYPE_UNKNOWN {
		vnic.Resources().Logger().Fail(t, "Unknown device type")
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

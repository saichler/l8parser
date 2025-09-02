package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8collector/go/collector/devices"
	"github.com/saichler/l8collector/go/collector/service"
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

	serviceArea := byte(0)
	allPolls := boot.GetAllPolarisModels()
	for _, snmpPolls := range allPolls {
		for _, poll := range snmpPolls.Polling {
			if poll.Cadence > 3 {
				poll.Cadence = 3
			}
		}
	}

	//use opensim to simulate this device with this ip
	//https://github.com/saichler/opensim
	//curl -X POST http://localhost:8080/api/v1/devices -H "Content-Type: application/json" -d '{"start_ip":"10.10.10.1","device_count":3,"netmask":"24"}'
	device := utils_collector.CreateDevice("10.20.30.1", serviceArea)

	vnic := topo.VnicByVnetNum(2, 2)
	vnic.Resources().Registry().Register(pollaris.PollarisService{})
	vnic.Resources().Services().Activate(pollaris.ServiceType, pollaris.ServiceName, serviceArea, vnic.Resources(), vnic)
	vnic.Resources().Registry().Register(&types2.NetworkDeviceList{})

	p := pollaris.Pollaris(vnic.Resources())
	for _, snmpPolls := range allPolls {
		err := p.Add(snmpPolls, false)
		if err != nil {
			vnic.Resources().Logger().Fail(t, err.Error())
			return
		}
	}

	vnic.Resources().Registry().Register(devices.DeviceService{})
	vnic.Resources().Services().Activate(devices.ServiceType, devices.ServiceName, serviceArea, vnic.Resources(), vnic)
	vnic.Resources().Registry().Register(service.CollectorService{})
	vnic.Resources().Services().Activate(service.ServiceType, common.CollectorService, serviceArea, vnic.Resources(), vnic)

	vnic.Resources().Registry().Register(&parsing.ParsingService{})
	vnic.Resources().Services().Activate(parsing.ServiceType, device.ParsingService.ServiceName, byte(device.ParsingService.ServiceArea),
		vnic.Resources(), vnic, &types2.NetworkDevice{}, "Id", true)

	vnic.Resources().Registry().Register(&inventory.InventoryService{})
	vnic.Resources().Services().Activate(inventory.ServiceType, device.InventoryService.ServiceName, byte(device.InventoryService.ServiceArea),
		vnic.Resources(), vnic, "Id", &types2.NetworkDevice{})

	time.Sleep(time.Second)

	cl := topo.VnicByVnetNum(1, 1)
	cl.Multicast(devices.ServiceName, serviceArea, ifs.POST, device)

	time.Sleep(time.Second * 10)

	inv := inventory.Inventory(vnic.Resources(), device.InventoryService.ServiceName, byte(device.InventoryService.ServiceArea))
	elem := inv.ElementByKey("10.20.30.1")
	networkDevice := elem.(*types2.NetworkDevice)

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
	marshalOptions := protojson.MarshalOptions{
		UseEnumNumbers: true,
	}
	jsn, _ := marshalOptions.Marshal(networkDevice)
	os.WriteFile("/tmp/NetworkDevice.json", jsn, 0644)
}

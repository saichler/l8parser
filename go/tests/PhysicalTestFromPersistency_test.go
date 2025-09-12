package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/tests/utils_collector"
	inventory "github.com/saichler/l8inventory/go/inv/service"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8srlz/go/serialize/object"
	types2 "github.com/saichler/probler/go/types"
)

func TestPhysicalFromPersistency(t *testing.T) {

	serviceArea := byte(0)
	allPolls := boot.GetAllPolarisModels()

	//use opensim to simulate this device with this ip
	//https://github.com/saichler/opensim
	//curl -X POST http://localhost:8080/api/v1/devices -H "Content-Type: application/json" -d '{"start_ip":"10.10.10.1","device_count":3,"netmask":"24"}'
	device := utils_collector.CreateDevice("10.20.30.3", serviceArea)

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

	vnic.Resources().Registry().Register(&parsing.ParsingService{})
	vnic.Resources().Services().Activate(parsing.ServiceType, device.ParsingService.ServiceName, byte(device.ParsingService.ServiceArea),
		vnic.Resources(), vnic, &types2.NetworkDevice{}, "Id", true)

	forwardInfo := &types.DeviceServiceInfo{}
	forwardInfo.ServiceName = "MockOrm"
	forwardInfo.ServiceArea = 0

	vnic.Resources().Registry().Register(&inventory.InventoryService{})
	vnic.Resources().Services().Activate(inventory.ServiceType, device.InventoryService.ServiceName, byte(device.InventoryService.ServiceArea),
		vnic.Resources(), vnic, "Id", &types2.NetworkDevice{}, forwardInfo)

	time.Sleep(time.Second)

	job, err := parsing.LoadJob("mib2", "entityMib", "10.20.30.3", "10.20.30.3")
	if err != nil {
		vnic.Resources().Logger().Fail(t, err.Error())
		return
	}
	ps, _ := vnic.Resources().Services().ServiceHandler(device.ParsingService.ServiceName, byte(device.ParsingService.ServiceArea))
	parserService := ps.(*parsing.ParsingService)
	jobElem := object.New(nil, job)
	parserService.Post(jobElem, vnic)

	time.Sleep(time.Second)

	inv := inventory.Inventory(vnic.Resources(), device.InventoryService.ServiceName, byte(device.InventoryService.ServiceArea))
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

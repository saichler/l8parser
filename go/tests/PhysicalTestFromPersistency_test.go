package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8collector/go/collector/service"
	"github.com/saichler/l8collector/go/collector/targets"
	"github.com/saichler/l8collector/go/tests/utils_collector"
	inventory "github.com/saichler/l8inventory/go/inv/service"
	"github.com/saichler/l8inventory/go/tests/utils_inventory"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
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
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&utils_inventory.MockOrmService{}, device.LinkData.ZsideServiceName, byte(device.LinkData.ZsideServiceArea), true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	time.Sleep(time.Second)

	job, err := parsing.LoadJob("mib2", "entityMib", "10.20.30.3", "10.20.30.3")
	if err != nil {
		vnic.Resources().Logger().Fail(t, err.Error())
		return
	}
	ps, _ := vnic.Resources().Services().ServiceHandler(device.LinkParser.ZsideServiceName, byte(device.LinkParser.ZsideServiceArea))
	parserService := ps.(*parsing.ParsingService)
	jobElem := object.New(nil, job)
	parserService.Post(jobElem, vnic)

	time.Sleep(time.Second)

	inv := inventory.Inventory(vnic.Resources(), device.LinkData.ZsideServiceName, byte(device.LinkData.ZsideServiceArea))
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

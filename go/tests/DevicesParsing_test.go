package tests

import (
	"strconv"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/common"
	
	"github.com/saichler/l8collector/go/collector/service"
	"github.com/saichler/l8collector/go/tests/utils_collector"
	"github.com/saichler/l8inventory/go/tests/utils_inventory"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
)

func TestFullDevicesParsing(t *testing.T) {

	serviceArea := byte(0)
	all := boot.GetAllPolarisModels()
	for _, snmpPolls := range all {
		for _, poll := range snmpPolls.Polling {
			if poll.Cadence.Enabled {
				poll.Cadence.Cadences[0] = 3
			}
		}
	}
	//use opensim to simulate this device with this ip
	//https://github.com/saichler/opensim
	//curl -X POST http://localhost:8080/api/v1/devices -H "Content-Type: application/json" -d '{"start_ip":"10.10.10.1","device_count":3,"netmask":"24"}'
	netDevices := make([]*types.Device, 0)
	for i := 1; i <= 19; i++ {
		ii := strconv.Itoa(i)
		device := utils_collector.CreateDevice("10.20.30."+ii, serviceArea)
		netDevices = append(netDevices, device)
	}

	vnic := topo.VnicByVnetNum(2, 2)
	vnic.Resources().Registry().Register(pollaris.PollarisService{})
	vnic.Resources().Services().Activate(pollaris.ServiceType, pollaris.ServiceName, serviceArea, vnic.Resources(), vnic)

	p := pollaris.Pollaris(vnic.Resources())
	for _, snmpPolls := range all {
		err := p.Add(snmpPolls, false)
		if err != nil {
			vnic.Resources().Logger().Fail(t, err.Error())
			return
		}
	}

	vnic.Resources().Registry().Register(targets.DeviceService{})
	vnic.Resources().Services().Activate(targets.ServiceType, targets.ServiceName, serviceArea, vnic.Resources(), vnic)
	vnic.Resources().Registry().Register(service.CollectorService{})
	vnic.Resources().Services().Activate(service.ServiceType, common.CollectorService, serviceArea, vnic.Resources(), vnic)

	vnic.Resources().Registry().Register(&parsing.ParsingService{})
	vnic.Resources().Services().Activate(parsing.ServiceType, netDevices[0].ParsingService.ServiceName, byte(netDevices[0].ParsingService.ServiceArea),
		vnic.Resources(), vnic, &types2.NetworkDevice{}, "Id", false)

	vnic.Resources().Registry().Register(&utils_inventory.MockOrmService{})
	vnic.Resources().Services().Activate(utils_inventory.ServiceType, netDevices[0].InventoryService.ServiceName, byte(netDevices[0].InventoryService.ServiceArea),
		vnic.Resources(), vnic)

	time.Sleep(time.Second)

	cl := topo.VnicByVnetNum(1, 1)
	for _, device := range netDevices {
		err := cl.Multicast(targets.ServiceName, serviceArea, ifs.POST, device)
		if err != nil {
			panic(err)
		}
		//time.Sleep(time.Second)
	}

	time.Sleep(time.Second * 10)
}

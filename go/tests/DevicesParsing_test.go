package tests

import (
	"strconv"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8collector/go/collector/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"

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
	netDevices := make([]*l8tpollaris.L8PTarget, 0)
	for i := 1; i <= 19; i++ {
		ii := strconv.Itoa(i)
		device := utils_collector.CreateDevice("10.20.30."+ii, serviceArea)
		netDevices = append(netDevices, device)
	}

	vnic := topo.VnicByVnetNum(2, 2)
	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, serviceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	p := pollaris.Pollaris(vnic.Resources())
	for _, snmpPolls := range all {
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

	sla = ifs.NewServiceLevelAgreement(&parsing.ParsingService{}, netDevices[0].LinkParser.ZsideServiceName,
		byte(netDevices[0].LinkParser.ZsideServiceArea), true, nil)
	sla.SetServiceItem(&types2.NetworkDevice{})
	sla.SetPrimaryKeys([]string{"Id"})
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&utils_inventory.MockOrmService{}, netDevices[0].LinkData.ZsideServiceName,
		byte(netDevices[0].LinkData.ZsideServiceArea), true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

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

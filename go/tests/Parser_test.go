package tests

import (
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8collector/go/collector/targets"

	"github.com/saichler/l8collector/go/collector/service"
	"github.com/saichler/l8collector/go/tests/utils_collector"
	"github.com/saichler/l8inventory/go/tests/utils_inventory"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

func TestParser(t *testing.T) {
	serviceArea := byte(0)
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
	device := utils_collector.CreateDevice("10.20.30.1", serviceArea)

	vnic := topo.VnicByVnetNum(2, 2)
	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, serviceArea, true, nil)
	sla.SetInitItems(initData)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&targets.TargetService{}, targets.ServiceName, serviceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&service.CollectorService{}, common.CollectorService, serviceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&parsing.ParsingService{}, device.LinkParser.ZsideServiceName, byte(device.LinkParser.ZsideServiceArea), true, nil)
	sla.SetServiceItem(&types2.NetworkDevice{})
	sla.SetPrimaryKeys("Id")
	sla.SetArgs(false)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&utils_inventory.MockOrmService{}, device.LinkData.ZsideServiceName, byte(device.LinkData.ZsideServiceArea), false, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	cl := topo.VnicByVnetNum(1, 1)
	err := cl.Multicast(targets.ServiceName, serviceArea, ifs.POST, device)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 6)

	m, ok := vnic.Resources().Services().ServiceHandler(device.LinkData.ZsideServiceName, byte(device.LinkData.ZsideServiceArea))
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

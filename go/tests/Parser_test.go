package tests

import (
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8collector/go/collector/devices"
	"github.com/saichler/l8collector/go/collector/service"
	"github.com/saichler/l8collector/go/tests/utils_collector"
	"github.com/saichler/l8inventory/go/tests/utils_inventory"
	types2 "github.com/saichler/l8inventory/go/types"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

func TestParser(t *testing.T) {
	serviceArea := byte(0)
	snmpPolls := boot.CreateSNMPBootPolls()
	for _, poll := range snmpPolls.Polling {
		poll.Cadence = 3
	}

	device := utils_collector.CreateDevice("192.168.86.179", serviceArea)

	vnic := topo.VnicByVnetNum(2, 2)
	vnic.Resources().Registry().Register(pollaris.PollarisService{})
	vnic.Resources().Services().Activate(pollaris.ServiceType, pollaris.ServiceName, serviceArea, vnic.Resources(), vnic)
	vnic.Resources().Registry().Register(devices.DeviceService{})
	vnic.Resources().Services().Activate(devices.ServiceType, devices.ServiceName, serviceArea, vnic.Resources(), vnic)
	vnic.Resources().Registry().Register(service.CollectorService{})
	vnic.Resources().Services().Activate(service.ServiceType, common.CollectorService, serviceArea, vnic.Resources(), vnic)

	vnic.Resources().Registry().Register(&parsing.ParsingService{})
	vnic.Resources().Services().Activate(parsing.ServiceType, device.ParsingService.ServiceName, byte(device.ParsingService.ServiceArea),
		vnic.Resources(), vnic, &types2.NetworkBox{}, "Id")

	vnic.Resources().Registry().Register(&utils_inventory.MockOrmService{})
	vnic.Resources().Services().Activate(utils_inventory.ServiceType, device.InventoryService.ServiceName, byte(device.InventoryService.ServiceArea),
		vnic.Resources(), vnic)

	time.Sleep(time.Second)

	p := pollaris.Pollaris(vnic.Resources())
	err := p.Add(snmpPolls, false)
	if err != nil {
		vnic.Resources().Logger().Fail(t, err.Error())
		return
	}

	cl := topo.VnicByVnetNum(1, 1)
	err = cl.Multicast(devices.ServiceName, serviceArea, ifs.POST, device)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	m, ok := vnic.Resources().Services().ServiceHandler(device.InventoryService.ServiceName, byte(device.InventoryService.ServiceArea))
	if !ok {
		vnic.Resources().Logger().Fail(t, "Cannot find mock service")
		return
	}
	mock := m.(*utils_inventory.MockOrmService)
	if mock.PatchCount() != 1 {
		vnic.Resources().Logger().Fail(t, "Expected 1 patch count in mock")
		return
	}

	/*
		defer func() {
			deActivateDeviceAndPollConfigServices(cfg, 0)
		}()
	*/
	/*
		cl := topo.VnicByVnetNum(1, 1)
		err = cl.Multicast(devices.ServiceName, serviceArea, ifs.POST, device)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 3)
		mp, ok := vnic.Resources().Services().ServiceHandler(device.ParsingService.ServiceName, byte(device.ParsingService.ServiceArea))
		if !ok {
			panic("No mock service found")
		}
		mock := mp.(*utils_collector.MockParsingService)
		for k, v := range mock.JobsCounts() {
			for k1, v1 := range v {
				if v1 != 1 {
					vnic.Resources().Logger().Fail(t, "Expected 1 but got ", v1, " job ", k, ":", k1)
				}
			}
		}*/
}

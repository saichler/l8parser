package tests

import (
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	common2 "github.com/saichler/probler/go/prob/common"
	"github.com/saichler/probler/go/prob/common/creates"
	"strconv"
	"testing"
	"time"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"

	"github.com/saichler/l8collector/go/collector/service"
	"github.com/saichler/l8inventory/go/tests/utils_inventory"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"

	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
)

func TestFullDevicesParsing(t *testing.T) {

	linksId := common2.NetworkDevice_Links_ID

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
		device := creates.CreateDevice("10.20.30."+ii, linksId, "sim")
		netDevices = append(netDevices, device)
	}

	vnic := topo.VnicByVnetNum(2, 2)
	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, pollaris.ServiceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	p := pollaris.Pollaris(vnic.Resources())
	for _, snmpPolls := range all {
		err := p.Post(snmpPolls, false)
		if err != nil {
			vnic.Resources().Logger().Fail(t, err.Error())
			return
		}
	}

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

	cl := topo.VnicByVnetNum(1, 1)
	for _, device := range netDevices {
		err := cl.Multicast(targets.ServiceName, targets.ServiceArea, ifs.POST, device)
		if err != nil {
			panic(err)
		}
		//time.Sleep(time.Second)
	}

	time.Sleep(time.Second * 10)
}

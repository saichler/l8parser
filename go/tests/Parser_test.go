package tests

import (
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/types"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

func TestParser(t *testing.T) {
	vnic := topo.VnicByVnetNum(2, 2)
	vnic.Resources().Registry().Register(pollaris.PollarisService{})
	vnic.Resources().Services().Activate(pollaris.ServiceType, pollaris.ServiceName, 0, vnic.Resources(), vnic)
	vnic.Resources().Registry().Register(parsing.ParsingService{})
	vnic.Resources().Services().Activate(parsing.ServiceType, "Parser", 0, vnic.Resources(), vnic, &types.Pollaris{}, "Name")

	p := pollaris.Pollaris(vnic.Resources())
	p.Add(boot.CreateSNMPBootPolls(), false)
	p.Add(boot.CreateK8sBootPolls(), false)
}

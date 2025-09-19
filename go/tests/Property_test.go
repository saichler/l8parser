package tests

import (
	"fmt"
	"testing"

	"github.com/saichler/probler/go/types"
	"github.com/saichler/l8reflect/go/reflect/cloning"
	"github.com/saichler/l8reflect/go/reflect/updating"
)

func TestProperty(t *testing.T) {
	aside := &types.NetworkDevice{}
	aside.NetworkLinks = make([]*types.NetworkLink, 1)
	aside.NetworkLinks[0] = &types.NetworkLink{}
	aside.NetworkLinks[0].LinkStatus = types.LinkStatus_LINK_STATUS_ACTIVE
	aside.Physicals = make(map[string]*types.Physical)
	aside.Physicals["physical-1"] = &types.Physical{}
	aside.Physicals["physical-1"].Id = "id5"

	c := cloning.NewCloner()
	zside := c.Clone(aside).(*types.NetworkDevice)
	zside.NetworkLinks[0].LinkStatus = types.LinkStatus_LINK_STATUS_DOWN
	aside.Physicals["physical-1"].Id = "id6"

	vnic := topo.VnicByVnetNum(1, 1)
	vnic.Resources().Introspector().Inspect(aside)

	updater := updating.NewUpdater(vnic.Resources(), false, false)
	err := updater.Update(aside, zside)
	if err != nil {
		vnic.Resources().Logger().Fail(t, err.Error())
		return
	}
	for _, chg := range updater.Changes() {
		fmt.Println("PropertyId=", chg.PropertyId())
	}

}

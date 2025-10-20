package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8collector/go/collector/service"
	"github.com/saichler/l8collector/go/collector/targets"
	"github.com/saichler/l8collector/go/tests/utils_collector"
	inventory "github.com/saichler/l8inventory/go/inv/service"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestCluster(t *testing.T) {
	serviceArea := byte(0)
	k8sPolls := boot.CreateK8sBootPolls()

	//use opensim to simulate this device with this ip
	//https://github.com/saichler/opensim
	//curl -X POST http://localhost:8080/api/v1/devices -H "Content-Type: application/json" -d '{"start_ip":"10.10.10.1","device_count":3,"netmask":"24"}'
	cluster := utils_collector.CreateCluster("./lab.conf", "lab", int32(serviceArea))

	vnic := topo.VnicByVnetNum(2, 2)

	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, serviceArea, true, nil)
	sla.SetInitItems([]interface{}{k8sPolls})
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&targets.TargetService{}, targets.ServiceName, serviceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&service.CollectorService{}, common.CollectorService, serviceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&parsing.ParsingService{}, cluster.LinkParser.ZsideServiceName, byte(cluster.LinkParser.ZsideServiceArea), true, nil)
	sla.SetServiceItem(&types2.K8SCluster{})
	sla.SetPrimaryKeys("Name")
	sla.SetArgs(false)
	vnic.Resources().Services().Activate(sla, vnic)

	sla = ifs.NewServiceLevelAgreement(&inventory.InventoryService{}, cluster.LinkData.ZsideServiceName, byte(cluster.LinkData.ZsideServiceArea), true, nil)
	sla.SetServiceItem(&types2.K8SCluster{})
	sla.SetServiceItemList(&types2.K8SClusterList{})
	sla.SetPrimaryKeys("Name")
	vnic.Resources().Services().Activate(sla, vnic)

	time.Sleep(time.Second)

	cl := topo.VnicByVnetNum(1, 1)
	cl.Multicast(targets.ServiceName, serviceArea, ifs.POST, cluster)

	time.Sleep(time.Second * 10)

	inv := inventory.Inventory(vnic.Resources(), cluster.LinkData.ZsideServiceName, byte(cluster.LinkData.ZsideServiceArea))

	filter := &types2.K8SCluster{Name: "lab"}
	elem := inv.ElementByElement(filter)
	k8sCluster := elem.(*types2.K8SCluster)
	list := &types2.K8SClusterList{List: []*types2.K8SCluster{k8sCluster}}

	fmt.Println(len(k8sCluster.Pods))

	jsn, _ := protojson.Marshal(list)
	os.WriteFile("./clusters.json", jsn, 0777)

}

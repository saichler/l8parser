package tests

import (
	"fmt"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	common2 "github.com/saichler/probler/go/prob/common"
	"os"
	"testing"
	"time"

	"github.com/saichler/l8collector/go/collector/service"
	inventory "github.com/saichler/l8inventory/go/inv/service"
	"github.com/saichler/l8parser/go/parser/boot"
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/probler/go/prob/common/creates"
	types2 "github.com/saichler/probler/go/types"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestCluster(t *testing.T) {
	linksID := common2.K8s_Links_ID
	k8sPolls := boot.CreateK8sBootPolls()
	cluster := creates.CreateCluster("lab")
	cluster.State = l8tpollaris.L8PTargetState_Up

	vnic := topo.VnicByVnetNum(2, 2)

	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, pollaris.ServiceArea, true, nil)
	sla.SetInitItems([]interface{}{k8sPolls})
	vnic.Resources().Services().Activate(sla, vnic)

	targets.Activate("postgres", "probler", vnic)
	targets.Links.Collector(linksID)

	collServiceName, collServiceArea := targets.Links.Collector(linksID)
	sla = ifs.NewServiceLevelAgreement(&service.CollectorService{}, collServiceName, collServiceArea, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)

	parserServiceName, parserServiceArea := targets.Links.Parser(linksID)
	sla = ifs.NewServiceLevelAgreement(&parsing.ParsingService{}, parserServiceName, parserServiceArea, true, nil)
	sla.SetServiceItem(&types2.K8SCluster{})
	sla.SetPrimaryKeys("Name")
	sla.SetArgs(false)
	vnic.Resources().Services().Activate(sla, vnic)

	k8sServiceName, k8sServiceArea := targets.Links.Cache(linksID)
	sla = ifs.NewServiceLevelAgreement(&inventory.InventoryService{}, k8sServiceName, k8sServiceArea, true, nil)
	sla.SetServiceItem(&types2.K8SCluster{})
	sla.SetServiceItemList(&types2.K8SClusterList{})
	sla.SetPrimaryKeys("Name")
	vnic.Resources().Services().Activate(sla, vnic)

	time.Sleep(time.Second)

	cl := topo.VnicByVnetNum(1, 1)
	cl.Multicast(targets.ServiceName, targets.ServiceArea, ifs.POST, cluster)

	time.Sleep(time.Second * 10)

	inv := inventory.Inventory(vnic.Resources(), k8sServiceName, k8sServiceArea)

	filter := &types2.K8SCluster{Name: "lab"}
	elem := inv.ElementByElement(filter)
	k8sCluster := elem.(*types2.K8SCluster)
	list := &types2.K8SClusterList{List: []*types2.K8SCluster{k8sCluster}}

	fmt.Println(len(k8sCluster.Pods))

	jsn, _ := protojson.Marshal(list)
	os.WriteFile("./clusters.json", jsn, 0777)

}

/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	parsing "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/probler/go/prob/common/creates"
	"github.com/saichler/probler/go/serializers"
	types2 "github.com/saichler/probler/go/types"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestCluster tests the Kubernetes cluster parsing capabilities.
// It verifies that kubectl output is correctly parsed into K8sCluster structures
// including nodes, pods, deployments, and other Kubernetes resources.
func TestCluster(t *testing.T) {
	linksID := common2.K8s_Links_ID

	cluster := creates.CreateCluster("lab")
	cluster.State = l8tpollaris.L8PTargetState_Up

	vnic := topo.VnicByVnetNum(2, 2)

	vnic.Resources().Registry().Register(&types2.K8SReadyState{})
	vnic.Resources().Registry().Register(&types2.K8SRestartsState{})

	info, err := vnic.Resources().Registry().Info("K8SReadyState")
	if err != nil {
		vnic.Resources().Logger().Error(err)
	} else {
		info.AddSerializer(&serializers.Ready{})
	}

	info, err = vnic.Resources().Registry().Info("K8SRestartsState")
	if err != nil {
		vnic.Resources().Logger().Error(err)
	} else {
		info.AddSerializer(&serializers.Restarts{})
	}

	vnic.Resources().Registry().RegisterEnums(types2.K8SPodStatus_value)

	pollaris.Activate(vnic)

	targets.Activate(common2.DB_CREDS, common2.DB_NAME, vnic)

	service.Activate(linksID, vnic)

	//Activate Inventory parser
	parsing.Activate(common2.NetworkDevice_Links_ID, &types2.NetworkDevice{}, false, vnic, "Id")
	parsing.Activate(linksID, &types2.K8SCluster{}, false, vnic, "Name")

	k8sServiceName, k8sServiceArea := targets.Links.Cache(linksID)
	sla := ifs.NewServiceLevelAgreement(&inventory.InventoryService{}, k8sServiceName, k8sServiceArea, true, nil)
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

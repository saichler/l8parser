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
	"testing"

	"github.com/saichler/l8reflect/go/reflect/cloning"
	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/probler/go/types"
)

// TestProperty tests the property reflection and update mechanisms
// used by the parser to set values on NetworkDevice objects.
func TestProperty(t *testing.T) {
	aside := &types.NetworkDevice{}
	aside.Physicals = make(map[string]*types.Physical)
	aside.Physicals["physical-1"] = &types.Physical{}
	aside.Physicals["physical-1"].Id = "id5"

	c := cloning.NewCloner()
	zside := c.Clone(aside).(*types.NetworkDevice)
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

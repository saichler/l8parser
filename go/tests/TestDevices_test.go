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
	"strings"
	"testing"

	"github.com/saichler/l8parser/go/parser/boot"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)
import "google.golang.org/protobuf/encoding/protojson"

// TestDevices tests that mock device data can be correctly serialized to JSON.
func TestDevices(t *testing.T) {
	deviceList := Devices()
	/*
		for _, device := range deviceList.List {
			if device.Equipmentinfo.SysName == "NY-CORE-01" {
				for _, phy := range device.Physicals {
					fmt.Println(phy)
				}
				panic("dddd")
			}
		}*/
	_, err := protojson.Marshal(deviceList)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	//fmt.Println(string(devices))
}

// TestPolling validates all Pollaris model configurations for correctness.
func TestPolling(t *testing.T) {
	m := CheckPollaris()
	if len(m) > 0 {
		t.Fail()
	}

}

// CheckPollaris validates all Pollaris models and returns a map of invalid configurations.
// It checks that Set rules have valid "from" parameters that match the poll's What field.
func CheckPollaris() map[string]string {
	plrs := boot.GetAllPolarisModels()
	result := make(map[string]string)
	for _, p := range plrs {
		checkPollaris(p, result)
	}
	return result
}

func checkPollaris(p *l8tpollaris.L8Pollaris, invalid map[string]string) {
	for _, poll := range p.Polling {
		for _, attr := range poll.Attributes {
			for _, rule := range attr.Rules {
				if rule.Name == "Set" {
					from, ok := rule.Params["from"]
					if ok && !strings.HasPrefix(from.Value, poll.What) {
						invalid[attr.PropertyId] = poll.Name
					}
				}
			}
		}
	}
}

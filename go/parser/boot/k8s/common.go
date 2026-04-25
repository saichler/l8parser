/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package k8s contains one Pollaris definition per Kubernetes/Istio prime object.
// Each prime object has its own file declaring a single K8sResourcePollDef value;
// CreateClientBootPolls aggregates all of them into a slice of *L8Pollaris (one
// Pollaris per prime object). Each Pollaris is named after the prime object's
// LinksId (matching probler/prob/common/Links_k8s.go) so the collector's
// BootSequence can match the right Pollaris per target by name.
package k8s

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8parser/go/parser/rules"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
)

// K8sResourcePollDef defines one Kubernetes/Istio prime object collection.
// LinksId is the per-prime-object LinkID (matches probler/prob/common/Links_k8s.go);
// it doubles as the resulting Pollaris.Name so BootSequence can route it to the
// matching target.
type K8sResourcePollDef struct {
	Name      string
	GVR       string
	Fields    []string
	Headers   []string
	ModelName string
	ColCount  int
	KeyIdx    []int
	LinksId   string
}

// every5MinutesAlways and defaultTimeout duplicate the parent boot package's
// constants so this sub-package can stand alone (avoids an import cycle that
// would form if boot imports k8s and k8s imports boot).
var (
	every5MinutesAlways = &l8tpollaris.L8PCadencePlan{Cadences: []int64{300}, Enabled: true}
	defaultTimeout      = int64(60)
	stringConvert       = &strings2.String{TypesPrefix: true}
)

// makeClientPollaris builds one *L8Pollaris per prime object. Pollaris.Name is
// set to def.LinksId so that BootSequence in l8collector matches it exactly to
// the target whose LinksId == Pollaris.Name.
func makeClientPollaris(def K8sResourcePollDef) *l8tpollaris.L8Pollaris {
	p := &l8tpollaris.L8Pollaris{}
	p.Name = def.LinksId
	p.Groups = []string{common.BOOT_STAGE_00}
	p.Polling = make(map[string]*l8tpollaris.L8Poll)

	poll := &l8tpollaris.L8Poll{}
	poll.Name = def.Name
	poll.Timeout = defaultTimeout
	poll.Cadence = every5MinutesAlways
	poll.Protocol = l8tpollaris.L8PProtocol_L8PKubernetesAPI
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.What = clientTableSpec(def.GVR, def.Fields, def.Headers)
	poll.Attributes = []*l8tpollaris.L8PAttribute{clientAttribute(def)}

	p.Polling[poll.Name] = poll
	return p
}

func clientAttribute(def K8sResourcePollDef) *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = map[string]string{def.ModelName: def.ModelName}
	attr.Rules = []*l8tpollaris.L8PRule{
		toTableRule(def.ColCount, def.KeyIdx...),
		tableToInstancesRule(),
	}
	return attr
}

func toTableRule(columns int, keyColumn ...int) *l8tpollaris.L8PRule {
	rule := &l8tpollaris.L8PRule{}
	rule.Name = "StringToCTable"
	rule.Params = make(map[string]*l8tpollaris.L8PParameter)
	rule.Params[rules.Columns] = &l8tpollaris.L8PParameter{Name: rules.Columns, Value: strconv.Itoa(columns)}
	keyStr := stringConvert.ToString(reflect.ValueOf(keyColumn))
	rule.Params[rules.KeyColumn] = &l8tpollaris.L8PParameter{Name: rules.KeyColumn, Value: keyStr}
	return rule
}

func tableToInstancesRule() *l8tpollaris.L8PRule {
	rule := &l8tpollaris.L8PRule{}
	rule.Name = "CTableToInstances"
	rule.Params = make(map[string]*l8tpollaris.L8PParameter)
	return rule
}

func clientTableSpec(gvr string, fields, columnNames []string) string {
	return fmt.Sprintf(`{"result":"table","mode":"list","gvr":"%s","fields":["%s"],"columnNames":["%s"]}`,
		gvr, joinDoubleQuoted(fields), joinDoubleQuoted(columnNames))
}

func joinDoubleQuoted(values []string) string {
	if len(values) == 0 {
		return ""
	}
	result := values[0]
	for i := 1; i < len(values); i++ {
		result += `","` + values[i]
	}
	return result
}

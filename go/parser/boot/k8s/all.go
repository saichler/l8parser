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

package k8s

import "github.com/saichler/l8pollaris/go/types/l8tpollaris"

// allDefs is the list of every K8s/Istio prime object collection definition.
// One Pollaris is produced per entry, named after the LinksId. Order doesn't
// affect routing; it only affects boot-time iteration order.
var allDefs = []K8sResourcePollDef{
	Node,
	Pod,
	Deploy,
	Sts,
	Ds,
	Rs,
	Job,
	Cj,
	Hpa,
	Svc,
	Ing,
	NetPol,
	Ep,
	EpSl,
	IngCl,
	Pv,
	Pvc,
	Scl,
	Cm,
	Sec,
	Rq,
	Lr,
	Pdb,
	Sa,
	Role,
	Cr,
	Rb,
	Crb,
	Ns,
	VCl,
	IstioVs,
	IstioDr,
	IstioGw,
	IstioSe,
	IstioPa,
	IstioAp,
	IstioSc,
	IstioEf,
	Crd,
	Event,
}

// CreateClientBootPolls returns one *L8Pollaris per K8s/Istio prime object.
// Each Pollaris.Name matches the corresponding LinksId so the collector's
// BootSequence routes the matching Pollaris to the target whose LinksId is
// the same string.
func CreateClientBootPolls() []*l8tpollaris.L8Pollaris {
	out := make([]*l8tpollaris.L8Pollaris, 0, len(allDefs))
	for _, def := range allDefs {
		out = append(out, makeClientPollaris(def))
	}
	return out
}

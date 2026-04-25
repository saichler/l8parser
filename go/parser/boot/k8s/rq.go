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

var Rq = K8sResourcePollDef{
	LinksId:   "K8sRq",
	Name:      "resourcequotas",
	GVR:       "v1/resourcequotas",
	ModelName: "k8sresourcequota",
	ColCount:  11,
	KeyIdx:    []int{0, 1},
	Fields:    []string{"metadata.namespace", "metadata.name", "_k.age", "_k.requestcpu", "_k.requestmemory", "_k.limitcpu", "_k.limitmemory", "_k.usedrequestcpu", "_k.usedrequestmemory", "_k.usedlimitcpu", "_k.usedlimitmemory"},
	Headers:   []string{"NAMESPACE", "NAME", "AGE", "REQUEST CPU", "REQUEST MEMORY", "LIMIT CPU", "LIMIT MEMORY", "USED REQUEST CPU", "USED REQUEST MEMORY", "USED LIMIT CPU", "USED LIMIT MEMORY"},
}

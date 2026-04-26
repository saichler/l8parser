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

var Ds = K8sResourcePollDef{
	LinksId:   "K8sDs",
	Name:      "daemonsets",
	GVR:       "apps/v1/daemonsets",
	ModelName: "k8sdaemonset",
	ColCount:  12,
	KeyIdx:    []int{0, 1},
	// All five count columns are now pre-stringified by enrichDaemonSet in
	// l8collector. The previous direct status.*NumberScheduled / numberReady
	// paths produced rune-garbage or "[]" depending on whether K8s omitted
	// the field and how the decoder unmarshalled the missing value.
	Fields:  []string{"metadata.namespace", "metadata.name", "_k.desired", "_k.current", "_k.ready", "_k.uptodate", "_k.available", "_k.nodeselector", "_k.age", "_k.containers", "_k.images", "_k.selector"},
	Headers: []string{"NAMESPACE", "NAME", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE", "NODE SELECTOR", "AGE", "CONTAINERS", "IMAGES", "SELECTOR"},
}

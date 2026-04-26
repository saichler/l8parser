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

var Pod = K8sResourcePollDef{
	LinksId:   "K8sPod",
	Name:      "pods",
	GVR:       "v1/pods",
	ModelName: "k8spod",
	// CONTAINERS_JSON carries the JSON-encoded container array produced by
	// l8collector's enrichPodContainers. The parser stores the raw string
	// in K8sPod.containers_json and the UI's pod detail popup parses it to
	// render image / imagePullPolicy / ports / env / resources / volumeMounts
	// per container.
	ColCount: 11,
	KeyIdx:   []int{0, 1},
	Fields:   []string{"metadata.namespace", "metadata.name", "_k.ready", "status.phase", "_k.restarts", "_k.age", "status.podIP", "spec.nodeName", "_k.nominatednode", "_k.containers_json"},
	Headers:  []string{"NAMESPACE", "NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP", "NODE", "NOMINATED NODE", "CONTAINERS_JSON"},
}

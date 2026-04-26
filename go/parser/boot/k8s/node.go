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

var Node = K8sResourcePollDef{
	LinksId:   "K8sNode",
	Name:      "nodes",
	GVR:       "v1/nodes",
	ModelName: "k8snode",
	// Adding STATUS as the second column to match `kubectl get nodes`
	// layout. The collector's enrichNode() computes the Ready/NotReady
	// string from status.conditions[type=Ready].status; the parser maps
	// it through the EnumRegistry to the K8SNodeStatus enum field.
	ColCount: 11,
	KeyIdx:   []int{0},
	Fields:   []string{"metadata.name", "_k.status", "_k.roles", "_k.age", "status.nodeInfo.kubeletVersion", "_k.internalip", "_k.externalip", "status.nodeInfo.osImage", "status.nodeInfo.kernelVersion", "status.nodeInfo.containerRuntimeVersion"},
	Headers:  []string{"NAME", "STATUS", "ROLES", "AGE", "VERSION", "INTERNAL-IP", "EXTERNAL-IP", "OS-IMAGE", "KERNEL-VERSION", "CONTAINER-RUNTIME"},
}

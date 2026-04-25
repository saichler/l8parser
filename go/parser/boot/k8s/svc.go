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

var Svc = K8sResourcePollDef{
	LinksId:   "K8sSvc",
	Name:      "services",
	GVR:       "v1/services",
	ModelName: "k8sservice",
	ColCount:  8,
	KeyIdx:    []int{0, 1},
	Fields:    []string{"metadata.namespace", "metadata.name", "spec.type", "spec.clusterIP", "_k.externalip", "_k.ports", "_k.age"},
	Headers:   []string{"NAMESPACE", "NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE"},
}

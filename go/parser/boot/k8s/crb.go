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

var Crb = K8sResourcePollDef{
	LinksId:   "K8sCrb",
	Name:      "clusterrolebindings",
	GVR:       "rbac.authorization.k8s.io/v1/clusterrolebindings",
	ModelName: "k8sclusterrolebinding",
	ColCount:  3,
	KeyIdx:    []int{0},
	Fields:    []string{"metadata.name", "_k.roleref", "_k.age"},
	Headers:   []string{"NAME", "ROLE REF", "AGE"},
}

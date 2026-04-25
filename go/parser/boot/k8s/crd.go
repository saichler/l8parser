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

var Crd = K8sResourcePollDef{
	LinksId:   "K8sCrd",
	Name:      "crds",
	GVR:       "apiextensions.k8s.io/v1/customresourcedefinitions",
	ModelName: "k8scrd",
	ColCount:  5,
	KeyIdx:    []int{0},
	Fields:    []string{"metadata.name", "spec.group", "_k.version", "spec.scope", "_k.age"},
	Headers:   []string{"NAME", "GROUP", "VERSION", "SCOPE", "AGE"},
}

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

var Sa = K8sResourcePollDef{
	LinksId:   "K8sSa",
	Name:      "serviceaccounts",
	GVR:       "v1/serviceaccounts",
	ModelName: "k8sserviceaccount",
	ColCount:  4,
	KeyIdx:    []int{0, 1},
	Fields:    []string{"metadata.namespace", "metadata.name", "_k.secrets", "_k.age"},
	Headers:   []string{"NAMESPACE", "NAME", "SECRETS", "AGE"},
}

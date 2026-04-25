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

var Event = K8sResourcePollDef{
	LinksId:   "K8sEvt",
	Name:      "events",
	GVR:       "v1/events",
	ModelName: "k8sevent",
	ColCount:  10,
	KeyIdx:    []int{0, 1},
	Fields:    []string{"metadata.namespace", "metadata.name", "type", "reason", "_k.object", "message", "_k.source", "count", "_k.firstseen", "_k.lastseen"},
	Headers:   []string{"NAMESPACE", "NAME", "TYPE", "REASON", "OBJECT", "MESSAGE", "SOURCE", "COUNT", "FIRST SEEN", "LAST SEEN"},
}

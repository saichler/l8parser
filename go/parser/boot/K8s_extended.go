/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package boot

var k8sExtendedClientPolls = []K8sResourcePollDef{
	{Name: "replicasets", GVR: "apps/v1/replicasets", ModelName: "k8sreplicaset", ColCount: 6, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.desired", "_k.current", "_k.ready", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "DESIRED", "CURRENT", "READY", "AGE"}},
	{Name: "jobs", GVR: "batch/v1/jobs", ModelName: "k8sjob", ColCount: 5, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.completions", "_k.duration", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "COMPLETIONS", "DURATION", "AGE"}},
	{Name: "cronjobs", GVR: "batch/v1/cronjobs", ModelName: "k8scronjob", ColCount: 7, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "spec.schedule", "_k.lastschedule", "spec.suspend", "_k.active", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "SCHEDULE", "LAST SCHEDULE", "SUSPEND", "ACTIVE", "AGE"}},
	{Name: "hpas", GVR: "autoscaling/v2/horizontalpodautoscalers", ModelName: "k8shpa", ColCount: 8, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.reference", "_k.targets", "spec.minReplicas", "spec.maxReplicas", "status.currentReplicas", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "REFERENCE", "TARGETS", "MIN REPLICAS", "MAX REPLICAS", "CURRENT REPLICAS", "AGE"}},
	{Name: "ingresses", GVR: "networking.k8s.io/v1/ingresses", ModelName: "k8singress", ColCount: 7, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "spec.ingressClassName", "_k.hosts", "_k.address", "_k.ports", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "CLASS NAME", "HOSTS", "ADDRESS", "PORTS", "AGE"}},
	{Name: "ingressclasses", GVR: "networking.k8s.io/v1/ingressclasses", ModelName: "k8singressclass", ColCount: 3, KeyIdx: []int{0},
		Fields:  []string{"metadata.name", "spec.controller", "_k.age"},
		Headers: []string{"NAME", "CONTROLLER", "AGE"}},
	{Name: "endpoints", GVR: "v1/endpoints", ModelName: "k8sendpoints", ColCount: 4, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.endpoints", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "ENDPOINTS", "AGE"}},
	{Name: "endpointslices", GVR: "discovery.k8s.io/v1/endpointslices", ModelName: "k8sendpointslice", ColCount: 6, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.addresstype", "_k.ports", "_k.endpoints", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "ADDRESS TYPE", "PORTS", "ENDPOINTS", "AGE"}},
	{Name: "persistentvolumes", GVR: "v1/persistentvolumes", ModelName: "k8spersistentvolume", ColCount: 10, KeyIdx: []int{0},
		Fields:  []string{"metadata.name", "_k.capacity", "_k.accessmodes", "spec.persistentVolumeReclaimPolicy", "_k.status", "_k.claim", "spec.storageClassName", "_k.reason", "_k.age", "spec.volumeMode"},
		Headers: []string{"NAME", "CAPACITY", "ACCESS MODES", "RECLAIM POLICY", "STATUS", "CLAIM", "STORAGE CLASS", "REASON", "AGE", "VOLUME MODE"}},
	{Name: "persistentvolumeclaims", GVR: "v1/persistentvolumeclaims", ModelName: "k8spersistentvolumeclaim", ColCount: 9, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.status", "spec.volumeName", "_k.capacity", "_k.accessmodes", "spec.storageClassName", "spec.volumeMode", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "STATUS", "VOLUME", "CAPACITY", "ACCESS MODES", "STORAGE CLASS", "VOLUME MODE", "AGE"}},
	{Name: "storageclasses", GVR: "storage.k8s.io/v1/storageclasses", ModelName: "k8sstorageclass", ColCount: 6, KeyIdx: []int{0},
		Fields:  []string{"metadata.name", "provisioner", "reclaimPolicy", "volumeBindingMode", "_k.allowexpansion", "_k.age"},
		Headers: []string{"NAME", "PROVISIONER", "RECLAIM POLICY", "VOLUME BINDING MODE", "ALLOW VOLUME EXPANSION", "AGE"}},
	{Name: "configmaps", GVR: "v1/configmaps", ModelName: "k8sconfigmap", ColCount: 4, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.datacount", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "DATA COUNT", "AGE"}},
	{Name: "secrets", GVR: "v1/secrets", ModelName: "k8ssecret", ColCount: 5, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "type", "_k.datacount", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "TYPE", "DATA COUNT", "AGE"}},
	{Name: "serviceaccounts", GVR: "v1/serviceaccounts", ModelName: "k8sserviceaccount", ColCount: 4, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.secrets", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "SECRETS", "AGE"}},
	{Name: "roles", GVR: "rbac.authorization.k8s.io/v1/roles", ModelName: "k8srole", ColCount: 3, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "AGE"}},
	{Name: "clusterroles", GVR: "rbac.authorization.k8s.io/v1/clusterroles", ModelName: "k8sclusterrole", ColCount: 2, KeyIdx: []int{0},
		Fields:  []string{"metadata.name", "_k.age"},
		Headers: []string{"NAME", "AGE"}},
	{Name: "rolebindings", GVR: "rbac.authorization.k8s.io/v1/rolebindings", ModelName: "k8srolebinding", ColCount: 4, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.roleref", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "ROLE REF", "AGE"}},
	{Name: "clusterrolebindings", GVR: "rbac.authorization.k8s.io/v1/clusterrolebindings", ModelName: "k8sclusterrolebinding", ColCount: 3, KeyIdx: []int{0},
		Fields:  []string{"metadata.name", "_k.roleref", "_k.age"},
		Headers: []string{"NAME", "ROLE REF", "AGE"}},
	{Name: "resourcequotas", GVR: "v1/resourcequotas", ModelName: "k8sresourcequota", ColCount: 11, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.age", "_k.requestcpu", "_k.requestmemory", "_k.limitcpu", "_k.limitmemory", "_k.usedrequestcpu", "_k.usedrequestmemory", "_k.usedlimitcpu", "_k.usedlimitmemory"},
		Headers: []string{"NAMESPACE", "NAME", "AGE", "REQUEST CPU", "REQUEST MEMORY", "LIMIT CPU", "LIMIT MEMORY", "USED REQUEST CPU", "USED REQUEST MEMORY", "USED LIMIT CPU", "USED LIMIT MEMORY"}},
	{Name: "limitranges", GVR: "v1/limitranges", ModelName: "k8slimitrange", ColCount: 3, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "AGE"}},
	{Name: "poddisruptionbudgets", GVR: "policy/v1/poddisruptionbudgets", ModelName: "k8spoddisruptionbudget", ColCount: 6, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.minavailable", "_k.maxunavailable", "_k.alloweddisruptions", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "MIN AVAILABLE", "MAX UNAVAILABLE", "ALLOWED DISRUPTIONS", "AGE"}},
	{Name: "crds", GVR: "apiextensions.k8s.io/v1/customresourcedefinitions", ModelName: "k8scrd", ColCount: 5, KeyIdx: []int{0},
		Fields:  []string{"metadata.name", "spec.group", "_k.version", "spec.scope", "_k.age"},
		Headers: []string{"NAME", "GROUP", "VERSION", "SCOPE", "AGE"}},
	{Name: "events", GVR: "v1/events", ModelName: "k8sevent", ColCount: 10, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "type", "reason", "_k.object", "message", "_k.source", "count", "_k.firstseen", "_k.lastseen"},
		Headers: []string{"NAMESPACE", "NAME", "TYPE", "REASON", "OBJECT", "MESSAGE", "SOURCE", "COUNT", "FIRST SEEN", "LAST SEEN"}},
}

var k8sIstioClientPolls = []K8sResourcePollDef{
	{Name: "virtualservices", GVR: "networking.istio.io/v1beta1/virtualservices", ModelName: "istiovirtualservice", ColCount: 5, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.gateways", "_k.hosts", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "GATEWAYS", "HOSTS", "AGE"}},
	{Name: "destinationrules", GVR: "networking.istio.io/v1beta1/destinationrules", ModelName: "istiodestinationrule", ColCount: 4, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.host", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "HOST", "AGE"}},
	{Name: "gateways", GVR: "networking.istio.io/v1beta1/gateways", ModelName: "istiogateway", ColCount: 5, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.servers", "_k.selector", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "SERVERS", "SELECTOR", "AGE"}},
	{Name: "serviceentries", GVR: "networking.istio.io/v1beta1/serviceentries", ModelName: "istioserviceentry", ColCount: 7, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.hosts", "_k.location", "_k.resolution", "_k.ports", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "HOSTS", "LOCATION", "RESOLUTION", "PORTS", "AGE"}},
	{Name: "peerauthentications", GVR: "security.istio.io/v1beta1/peerauthentications", ModelName: "istiopeerauthentication", ColCount: 4, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.mode", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "MODE", "AGE"}},
	{Name: "authorizationpolicies", GVR: "security.istio.io/v1beta1/authorizationpolicies", ModelName: "istioauthorizationpolicy", ColCount: 4, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.action", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "ACTION", "AGE"}},
	{Name: "sidecars", GVR: "networking.istio.io/v1beta1/sidecars", ModelName: "istiosidecar", ColCount: 3, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "AGE"}},
	{Name: "envoyfilters", GVR: "networking.istio.io/v1alpha3/envoyfilters", ModelName: "istioenvoyfilter", ColCount: 3, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "AGE"}},
}

var k8sVClusterClientPolls = []K8sResourcePollDef{
	{Name: "vclusters", GVR: "management.loft.sh/v1/virtualclusterinstances", ModelName: "k8svcluster", ColCount: 10, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.status", "_k.k8sversion", "_k.distro", "_k.connected", "_k.syncedpods", "_k.syncedservices", "_k.syncedingresses", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "STATUS", "K8S VERSION", "DISTRO", "CONNECTED", "SYNCED PODS", "SYNCED SERVICES", "SYNCED INGRESSES", "AGE"}},
}

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

import (
	"fmt"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

// K8sResourcePollDef defines a K8s resource type for data-driven poll registration.
type K8sResourcePollDef struct {
	Name      string
	GVR       string
	Fields    []string
	Headers   []string
	ModelName string
	ColCount  int
	KeyIdx    []int
}

// registerClientResourcePoll registers a client-go API poll for a K8s resource.
func registerClientResourcePoll(p *l8tpollaris.L8Pollaris, def K8sResourcePollDef) {
	poll := createBaseK8sClientPoll(def.Name)
	poll.What = createClientTableSpec(def.GVR, def.Fields, def.Headers)
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = map[string]string{def.ModelName: def.ModelName}
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(def.ColCount, def.KeyIdx...))
	attr.Rules = append(attr.Rules, createTableToInstances())
	poll.Attributes = append(poll.Attributes, attr)
	p.Polling[poll.Name] = poll
}

// registerKubectlResourcePoll registers a kubectl-based poll for a K8s resource.
func registerKubectlResourcePoll(p *l8tpollaris.L8Pollaris, name, what, modelName string, colCount int, keyIdx []int) {
	poll := createBaseK8sPoll(name)
	poll.What = what
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = map[string]string{modelName: modelName}
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(colCount, keyIdx...))
	attr.Rules = append(attr.Rules, createTableToInstances())
	poll.Attributes = []*l8tpollaris.L8PAttribute{attr}
	p.Polling[poll.Name] = poll
}

func createTableToInstances() *l8tpollaris.L8PRule {
	rule := &l8tpollaris.L8PRule{}
	rule.Name = "CTableToInstances"
	rule.Params = make(map[string]*l8tpollaris.L8PParameter)
	return rule
}

// Core workload polls (nodes, pods, deployments, statefulsets, daemonsets, services, namespaces, networkpolicies)
var k8sCoreClientPolls = []K8sResourcePollDef{
	{Name: "nodes", GVR: "v1/nodes", ModelName: "k8snode", ColCount: 10, KeyIdx: []int{0},
		Fields:  []string{"metadata.name", "_k.roles", "_k.age", "status.nodeInfo.kubeletVersion", "_k.internalip", "_k.externalip", "status.nodeInfo.osImage", "status.nodeInfo.kernelVersion", "status.nodeInfo.containerRuntimeVersion"},
		Headers: []string{"NAME", "ROLES", "AGE", "VERSION", "INTERNAL-IP", "EXTERNAL-IP", "OS-IMAGE", "KERNEL-VERSION", "CONTAINER-RUNTIME"}},
	{Name: "pods", GVR: "v1/pods", ModelName: "k8spod", ColCount: 10, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.ready", "status.phase", "_k.restarts", "_k.age", "status.podIP", "spec.nodeName", "_k.nominatednode"},
		Headers: []string{"NAMESPACE", "NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP", "NODE", "NOMINATED NODE"}},
	{Name: "deployments", GVR: "apps/v1/deployments", ModelName: "k8sdeployment", ColCount: 9, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.ready", "status.updatedReplicas", "status.availableReplicas", "_k.age", "_k.containers", "_k.images", "_k.selector"},
		Headers: []string{"NAMESPACE", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE", "CONTAINERS", "IMAGES", "SELECTOR"}},
	{Name: "statefulsets", GVR: "apps/v1/statefulsets", ModelName: "k8sstatefulset", ColCount: 6, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.ready", "_k.age", "_k.containers", "_k.images"},
		Headers: []string{"NAMESPACE", "NAME", "READY", "AGE", "CONTAINERS", "IMAGES"}},
	{Name: "daemonsets", GVR: "apps/v1/daemonsets", ModelName: "k8sdaemonset", ColCount: 12, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "status.desiredNumberScheduled", "status.currentNumberScheduled", "status.numberReady", "_k.uptodate", "_k.available", "_k.nodeselector", "_k.age", "_k.containers", "_k.images", "_k.selector"},
		Headers: []string{"NAMESPACE", "NAME", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE", "NODE SELECTOR", "AGE", "CONTAINERS", "IMAGES", "SELECTOR"}},
	{Name: "services", GVR: "v1/services", ModelName: "k8sservice", ColCount: 8, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "spec.type", "spec.clusterIP", "_k.externalip", "_k.ports", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE"}},
	{Name: "namespaces", GVR: "v1/namespaces", ModelName: "k8snamespace", ColCount: 3, KeyIdx: []int{0},
		Fields:  []string{"metadata.name", "status.phase", "_k.age"},
		Headers: []string{"NAME", "STATUS", "AGE"}},
	{Name: "networkpolicies", GVR: "networking.k8s.io/v1/networkpolicies", ModelName: "k8snetworkpolicy", ColCount: 4, KeyIdx: []int{0, 1},
		Fields:  []string{"metadata.namespace", "metadata.name", "_k.podsel", "_k.age"},
		Headers: []string{"NAMESPACE", "NAME", "POD-SELECTOR", "AGE"}},
}

// CreateK8sBootPolls creates the kubectl-based Kubernetes polling configuration.
func CreateK8sBootPolls() *l8tpollaris.L8Pollaris {
	k8sPollaris := &l8tpollaris.L8Pollaris{}
	k8sPollaris.Name = "kubernetes"
	k8sPollaris.Groups = []string{common.BOOT_STAGE_00}
	k8sPollaris.Polling = make(map[string]*l8tpollaris.L8Poll)

	registerKubectlResourcePoll(k8sPollaris, "nodes", "get nodes -o wide", "k8snode", 10, []int{0})
	registerKubectlResourcePoll(k8sPollaris, "pods", "get pods -A -o wide", "k8spod", 10, []int{0, 1})
	registerKubectlResourcePoll(k8sPollaris, "deployments", "get deployments -A -o wide", "k8sdeployment", 9, []int{0, 1})
	registerKubectlResourcePoll(k8sPollaris, "statefulsets", "get statefulsets -A -o wide", "k8sstatefulset", 6, []int{0, 1})
	registerKubectlResourcePoll(k8sPollaris, "daemonsets", "get daemonsets -A -o wide", "k8sdaemonset", 12, []int{0, 1})
	registerKubectlResourcePoll(k8sPollaris, "services", "get services -A -o wide", "k8sservice", 8, []int{0, 1})
	registerKubectlResourcePoll(k8sPollaris, "namespaces", "get namespaces -A -o wide", "k8snamespace", 3, []int{0})
	registerKubectlResourcePoll(k8sPollaris, "networkpolicies", "get netpol -A -o wide", "k8snetworkpolicy", 4, []int{0, 1})

	createLogs(k8sPollaris)

	createNodeDetails(k8sPollaris)
	createPodDetails(k8sPollaris)
	createDeploymentDetails(k8sPollaris)
	createStatefulsetDetails(k8sPollaris)
	createDaemonsetDetails(k8sPollaris)
	createServiceDetails(k8sPollaris)
	createNamespaceDetails(k8sPollaris)
	createNetworkPolicyDetails(k8sPollaris)

	return k8sPollaris
}

// CreateK8sClientBootPolls creates the client-go Kubernetes polling configuration.
func CreateK8sClientBootPolls() *l8tpollaris.L8Pollaris {
	k8sPollaris := &l8tpollaris.L8Pollaris{}
	k8sPollaris.Name = "kubernetesapi"
	k8sPollaris.Groups = []string{common.BOOT_STAGE_00}
	k8sPollaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	for _, def := range k8sCoreClientPolls {
		registerClientResourcePoll(k8sPollaris, def)
	}
	for _, def := range k8sExtendedClientPolls {
		registerClientResourcePoll(k8sPollaris, def)
	}
	for _, def := range k8sIstioClientPolls {
		registerClientResourcePoll(k8sPollaris, def)
	}
	for _, def := range k8sVClusterClientPolls {
		registerClientResourcePoll(k8sPollaris, def)
	}
	return k8sPollaris
}

func createBaseK8sPoll(name string) *l8tpollaris.L8Poll {
	poll := &l8tpollaris.L8Poll{}
	poll.Name = name
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Cadence = EVERY_5_MINUTES_ALWAYS
	poll.Protocol = l8tpollaris.L8PProtocol_L8PKubectl
	return poll
}

func createBaseK8sClientPoll(name string) *l8tpollaris.L8Poll {
	poll := &l8tpollaris.L8Poll{}
	poll.Name = name
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Cadence = EVERY_5_MINUTES_ALWAYS
	poll.Protocol = l8tpollaris.L8PProtocol_L8PKubernetesAPI
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	return poll
}

func createClientTableSpec(gvr string, fields, columnNames []string) string {
	return fmt.Sprintf(`{"result":"table","mode":"list","gvr":"%s","fields":["%s"],"columnNames":["%s"]}`,
		gvr,
		stringsJoin(fields),
		stringsJoin(columnNames))
}

func stringsJoin(values []string) string {
	if len(values) == 0 {
		return ""
	}
	result := values[0]
	for i := 1; i < len(values); i++ {
		result += `","` + values[i]
	}
	return result
}

func createLogs(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("logs")
	poll.What = "logs -n $namespace $podname"
	poll.Cadence = DISABLED
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func createNodeDetails(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("nodedetails")
	poll.What = "get node $nodename -o json"
	poll.Cadence = DISABLED
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func createPodDetails(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("poddetails")
	poll.What = "get pods -o json -n $namespace $podname"
	poll.Cadence = DISABLED
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func createDeploymentDetails(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("deploymentdetails")
	poll.What = "get deployment -o json -n $namespace $deploymentname"
	poll.Cadence = DISABLED
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func createStatefulsetDetails(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("statefulsetdetails")
	poll.What = "get statefulset -o json -n $namespace $statefulsetname"
	poll.Cadence = DISABLED
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func createDaemonsetDetails(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("daemonsetdetails")
	poll.What = "get daemonset -o json -n $namespace $daemonsetname"
	poll.Cadence = DISABLED
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func createServiceDetails(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("servicedetails")
	poll.What = "get service -o json -n $namespace $servicename"
	poll.Cadence = DISABLED
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func createNamespaceDetails(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("namespacedetails")
	poll.What = "get namespace $namespace -o json"
	poll.Cadence = DISABLED
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func createNetworkPolicyDetails(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("networkpolicydetails")
	poll.What = "get netpol -o json -n $namespace $networkpolicyname"
	poll.Cadence = DISABLED
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

// LogsJob creates a collection job for retrieving pod logs from a Kubernetes cluster.
func LogsJob(cluster, context, namespace, podname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "logs"
	job.Arguments = map[string]string{"namespace": namespace, "podname": podname}
	return job
}

// NodeDetailsJob creates a collection job for retrieving detailed node information.
func NodeDetailsJob(cluster, context, nodename string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "nodedetails"
	job.Arguments = map[string]string{"nodename": nodename}
	return job
}

// PodDetailsJob creates a collection job for retrieving detailed pod information.
func PodDetailsJob(cluster, context, namespace, podname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "poddetails"
	job.Arguments = map[string]string{"namespace": namespace, "podname": podname}
	return job
}

// DeploymentDetailsJob creates a collection job for retrieving detailed deployment information.
func DeploymentDetailsJob(cluster, context, namespace, deploymentname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "deploymentdetails"
	job.Arguments = map[string]string{"namespace": namespace, "deploymentname": deploymentname}
	return job
}

// StatefulsetDetailsJob creates a collection job for retrieving detailed statefulset information.
func StatefulsetDetailsJob(cluster, context, namespace, statefulsetname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "statefulsetdetails"
	job.Arguments = map[string]string{"namespace": namespace, "statefulsetname": statefulsetname}
	return job
}

// DaemonsetDetailsJob creates a collection job for retrieving detailed daemonset information.
func DaemonsetDetailsJob(cluster, context, namespace, daemonsetname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "daemonsetdetails"
	job.Arguments = map[string]string{"namespace": namespace, "daemonsetname": daemonsetname}
	return job
}

// ServiceDetailsJob creates a collection job for retrieving detailed service information.
func ServiceDetailsJob(cluster, context, namespace, servicename string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "servicedetails"
	job.Arguments = map[string]string{"namespace": namespace, "servicename": servicename}
	return job
}

// NamespaceDetailsJob creates a collection job for retrieving detailed namespace information.
func NamespaceDetailsJob(cluster, context, namespace string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "namespacedetails"
	job.Arguments = map[string]string{"namespace": namespace}
	return job
}

// NetworkPolicyDetailsJob creates a collection job for retrieving detailed network policy information.
func NetworkPolicyDetailsJob(cluster, context, namespace, networkpolicyname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "networkpolicydetails"
	job.Arguments = map[string]string{"namespace": namespace, "networkpolicyname": networkpolicyname}
	return job
}

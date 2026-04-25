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
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

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

func createBaseK8sPoll(name string) *l8tpollaris.L8Poll {
	poll := &l8tpollaris.L8Poll{}
	poll.Name = name
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Cadence = EVERY_5_MINUTES_ALWAYS
	poll.Protocol = l8tpollaris.L8PProtocol_L8PKubectl
	return poll
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

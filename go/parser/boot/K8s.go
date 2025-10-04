package boot

import (
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

func CreateK8sBootPolls() *l8tpollaris.L8Pollaris {
	k8sPollaris := &l8tpollaris.L8Pollaris{}
	k8sPollaris.Name = "kubernetes"
	k8sPollaris.Groups = []string{common.BOOT_STAGE_00}
	k8sPollaris.Polling = make(map[string]*l8tpollaris.L8Poll)
	createNodesPoll(k8sPollaris)
	createPodsPoll(k8sPollaris)
	createDeploymentsPoll(k8sPollaris)
	createStatefulsetsPoll(k8sPollaris)
	createDaemonsetsPoll(k8sPollaris)
	createServicesPoll(k8sPollaris)
	createNamespacesPoll(k8sPollaris)
	createNetworkPoliciesPoll(k8sPollaris)

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

func createNodesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("nodes")
	poll.What = "get nodes -o wide"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createNodesTable())
	p.Polling[poll.Name] = poll
}

func createPodsPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("pods")
	poll.What = "get pods -A -o wide"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createPodsTable())
	p.Polling[poll.Name] = poll
}

func createDeploymentsPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("deployments")
	poll.What = "get deployments -A -o wide"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createDeplymentsTable())
	p.Polling[poll.Name] = poll
}

func createStatefulsetsPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("statefulsets")
	poll.What = "get statefulsets -A -o wide"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createStatefulsetsTable())
	p.Polling[poll.Name] = poll
}

func createDaemonsetsPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("daemonsets")
	poll.What = "get daemonsets -A -o wide"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createDaemonsetsTable())
	p.Polling[poll.Name] = poll
}

func createServicesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("services")
	poll.What = "get services -A -o wide"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createServicesTable())
	p.Polling[poll.Name] = poll
}

func createNamespacesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("namespaces")
	poll.What = "get namespaces -A -o wide"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createNamespacesTable())
	p.Polling[poll.Name] = poll
}

func createNetworkPoliciesPoll(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("networkpolicies")
	poll.What = "get netpol -A -o wide"
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createNetworkPoliciesTable())
	p.Polling[poll.Name] = poll
}

func createNodesTable() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "k8scluster.nodes"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(10, 0))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createPodsTable() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "k8scluster.pods"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(10, 6))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createDeplymentsTable() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "k8scluster.deployments"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(9, 1))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createStatefulsetsTable() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "k8scluster.statefulsets"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(6, 1))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createDaemonsetsTable() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "k8scluster.daemonsets"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(12, 1))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createServicesTable() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "k8scluster.services"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(8, 1))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createNamespacesTable() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "k8scluster.namespaces"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(3, 0))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createNetworkPoliciesTable() *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = "k8scluster.networkpolicies"
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	attr.Rules = append(attr.Rules, createToTable(4, 1))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
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

func LogsJob(cluster, context, namespace, podname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "logs"
	job.Arguments = map[string]string{"namespace": namespace, "podname": podname}
	return job
}

func NodeDetailsJob(cluster, context, nodename string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "nodedetails"
	job.Arguments = map[string]string{"nodename": nodename}
	return job
}

func PodDetailsJob(cluster, context, namespace, podname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "poddetails"
	job.Arguments = map[string]string{"namespace": namespace, "podname": podname}
	return job
}

func DeploymentDetailsJob(cluster, context, namespace, deploymentname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "deploymentdetails"
	job.Arguments = map[string]string{"namespace": namespace, "deploymentname": deploymentname}
	return job
}

func StatefulsetDetailsJob(cluster, context, namespace, statefulsetname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "statefulsetdetails"
	job.Arguments = map[string]string{"namespace": namespace, "statefulsetname": statefulsetname}
	return job
}

func DaemonsetDetailsJob(cluster, context, namespace, daemonsetname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "daemonsetdetails"
	job.Arguments = map[string]string{"namespace": namespace, "daemonsetname": daemonsetname}
	return job
}

func ServiceDetailsJob(cluster, context, namespace, servicename string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "servicedetails"
	job.Arguments = map[string]string{"namespace": namespace, "servicename": servicename}
	return job
}

func NamespaceDetailsJob(cluster, context, namespace string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "namespacedetails"
	job.Arguments = map[string]string{"namespace": namespace}
	return job
}

func NetworkPolicyDetailsJob(cluster, context, namespace, networkpolicyname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "networkpolicydetails"
	job.Arguments = map[string]string{"namespace": namespace, "networkpolicyname": networkpolicyname}
	return job
}

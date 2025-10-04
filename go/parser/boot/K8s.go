package boot

import (
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
)

func CreateK8sBootPolls() *l8poll.L8Pollaris {
	k8sPollaris := &l8poll.L8Pollaris{}
	k8sPollaris.Name = "kubernetes"
	k8sPollaris.Groups = []string{common.BOOT_STAGE_00}
	k8sPollaris.Polling = make(map[string]*l8poll.L8Poll)
	createNodesPoll(k8sPollaris)
	createPodsPoll(k8sPollaris)
	createLogs(k8sPollaris)
	createDetails(k8sPollaris)
	return k8sPollaris
}

func createNodesPoll(p *l8poll.L8Pollaris) {
	poll := createBaseK8sPoll("nodes")
	poll.What = "get nodes -o wide"
	poll.Operation = l8poll.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createNodesTable())
	p.Polling[poll.Name] = poll
}

func createPodsPoll(p *l8poll.L8Pollaris) {
	poll := createBaseK8sPoll("pods")
	poll.What = "get pods -A -o wide"
	poll.Operation = l8poll.L8C_Operation_L8C_Table
	poll.Attributes = make([]*l8poll.L8P_Attribute, 0)
	poll.Attributes = append(poll.Attributes, createPodsTable())
	p.Polling[poll.Name] = poll
}

func createNodesTable() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "k8scluster.nodes"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createToTable(10, 0))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createPodsTable() *l8poll.L8P_Attribute {
	attr := &l8poll.L8P_Attribute{}
	attr.PropertyId = "k8scluster.pods"
	attr.Rules = make([]*l8poll.L8P_Rule, 0)
	attr.Rules = append(attr.Rules, createToTable(10, 6))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createBaseK8sPoll(name string) *l8poll.L8Poll {
	poll := &l8poll.L8Poll{}
	poll.Name = name
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Cadence = EVERY_5_MINUTES_ALWAYS
	poll.Protocol = l8poll.L8C_Protocol_L8P_Kubectl
	return poll
}

func createLogs(p *l8poll.L8Pollaris) {
	poll := createBaseK8sPoll("logs")
	poll.What = "logs -n $namespace $podname"
	poll.Cadence = DISABLED
	poll.Operation = l8poll.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func createDetails(p *l8poll.L8Pollaris) {
	poll := createBaseK8sPoll("details")
	poll.What = "get pods -o json -n $namespace $podname"
	poll.Cadence = DISABLED
	poll.Operation = l8poll.L8C_Operation_L8C_Get
	p.Polling[poll.Name] = poll
}

func LogsJob(cluster, context, namespace, podname string) *l8poll.CJob {
	job := &l8poll.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "logs"
	job.Arguments = map[string]string{"namespace": namespace, "podname": podname}
	return job
}

func DetailsJob(cluster, context, namespace, podname string) *l8poll.CJob {
	job := &l8poll.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "details"
	job.Arguments = map[string]string{"namespace": namespace, "podname": podname}
	return job
}

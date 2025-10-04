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
	createLogs(k8sPollaris)
	createDetails(k8sPollaris)
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

func createDetails(p *l8tpollaris.L8Pollaris) {
	poll := createBaseK8sPoll("details")
	poll.What = "get pods -o json -n $namespace $podname"
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

func DetailsJob(cluster, context, namespace, podname string) *l8tpollaris.CJob {
	job := &l8tpollaris.CJob{}
	job.TargetId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "details"
	job.Arguments = map[string]string{"namespace": namespace, "podname": podname}
	return job
}

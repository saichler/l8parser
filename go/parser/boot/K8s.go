package boot

import (
	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8pollaris/go/types"
)

func CreateK8sBootPolls() *types.Pollaris {
	k8sPollaris := &types.Pollaris{}
	k8sPollaris.Name = "kubernetes"
	k8sPollaris.Groups = []string{common.BOOT_GROUP}
	k8sPollaris.Polling = make(map[string]*types.Poll)
	createNodesPoll(k8sPollaris)
	createPodsPoll(k8sPollaris)
	createLogs(k8sPollaris)
	return k8sPollaris
}

func createNodesPoll(p *types.Pollaris) {
	poll := createBaseK8sPoll("nodes")
	poll.What = "get nodes -o wide"
	poll.Operation = types.Operation_OTable
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createNodesTable())
	p.Polling[poll.Name] = poll
}

func createPodsPoll(p *types.Pollaris) {
	poll := createBaseK8sPoll("pods")
	poll.What = "get pods -A -o wide"
	poll.Operation = types.Operation_OTable
	poll.Attributes = make([]*types.Attribute, 0)
	poll.Attributes = append(poll.Attributes, createPodsTable())
	p.Polling[poll.Name] = poll
}

func createNodesTable() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "k8scluster.nodes"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createToTable(10, 0))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createPodsTable() *types.Attribute {
	attr := &types.Attribute{}
	attr.PropertyId = "k8scluster.pods"
	attr.Rules = make([]*types.Rule, 0)
	attr.Rules = append(attr.Rules, createToTable(10, 6))
	attr.Rules = append(attr.Rules, createTableToMap())
	return attr
}

func createBaseK8sPoll(name string) *types.Poll {
	poll := &types.Poll{}
	poll.Name = name
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Cadence = DEFAULT_CADENCE
	poll.Protocol = types.Protocol_PK8s
	return poll
}

func createLogs(p *types.Pollaris) {
	poll := createBaseK8sPoll("logs")
	poll.What = "logs -n |namespace| |podname|"
	poll.Cadence = -1
	poll.Operation = types.Operation_OGet
	p.Polling[poll.Name] = poll
}

func LogsJob(cluster, context, namespace, podname string) *types.CJob {
	job := &types.CJob{}
	job.DeviceId = cluster
	job.HostId = context
	job.PollarisName = "kubernetes"
	job.JobName = "logs"
	job.Arguments = map[string]string{"namespace": namespace, "podname": podname}
	return job
}

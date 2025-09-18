package service

import (
	"errors"

	"github.com/saichler/l8parser/go/parser/rules"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/types/l8poll"

	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
)

type _Parser struct {
	rules map[string]rules.ParsingRule
}

var Parser = newParser()

func newParser() *_Parser {
	p := &_Parser{}
	p.rules = make(map[string]rules.ParsingRule)
	con := &rules.Contains{}
	p.rules[con.Name()] = con
	set := &rules.Set{}
	p.rules[set.Name()] = set
	totable := &rules.StringToCTable{}
	p.rules[totable.Name()] = totable
	tableToMap := &rules.CTableToMapProperty{}
	p.rules[tableToMap.Name()] = tableToMap
	ifTableToPhysicals := &rules.IfTableToPhysicals{}
	p.rules[ifTableToPhysicals.Name()] = ifTableToPhysicals
	entityMibToPhysicals := &rules.EntityMibToPhysicals{}
	p.rules[entityMibToPhysicals.Name()] = entityMibToPhysicals
	inferDeviceType := &rules.InferDeviceType{}
	p.rules[inferDeviceType.Name()] = inferDeviceType
	mapToDeviceStatus := &rules.MapToDeviceStatus{}
	p.rules[mapToDeviceStatus.Name()] = mapToDeviceStatus
	return p
}

func (this *_Parser) Parse(job *l8poll.CJob, any interface{}, resources ifs.IResources) error {
	if job.Error != "" {
		return errors.New(job.Error)
	}

	if job.Result == nil || len(job.Result) < 4 {
		return resources.Logger().Error("Invalid job result ", job.TargetId, " - ", job.PollarisName,
			" - ", job.JobName, " - ", string(job.Result))
	}

	workSpace := make(map[string]interface{})
	enc := object.NewDecode(job.Result, 0, resources.Registry())
	data, err := enc.Get()
	if err != nil {
		return resources.Logger().Error(err)
	}
	poll, err := pollaris.Poll(job.PollarisName, job.JobName, resources)
	if err != nil {
		return resources.Logger().Error("cannot find poll for polaris ", job.PollarisName, ":", job.JobName)
	}
	workSpace[rules.Input] = data
	if poll.Attributes == nil {
		return resources.Logger().Error("No attributes are defined on pollaris "+job.PollarisName, ":", job.JobName)
	}

	for _, attr := range poll.Attributes {
		workSpace[rules.PropertyId] = attr.PropertyId
		for _, rData := range attr.Rules {
			if rData.Params != nil {
				for p, v := range rData.Params {
					workSpace[p] = v.Value
				}
			}
			ruleImpl, ok := this.rules[rData.Name]
			if !ok {
				return resources.Logger().Error("Cannot find parsing rule ", rData.Name)
			}
			err = ruleImpl.Parse(resources, workSpace, rData.Params, any, poll.What)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

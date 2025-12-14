package service

import (
	"fmt"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"reflect"
)

func (this *ParsingService) createElementInstance(job *l8tpollaris.CJob) interface{} {
	newElem := reflect.New(reflect.ValueOf(this.elem).Elem().Type())
	field := newElem.Elem().FieldByName(this.primaryKey)
	if !field.CanSet() {
		panic("cannot set field " + this.primaryKey)
	}
	fmt.Println(job.TargetId)
	field.Set(reflect.ValueOf(job.TargetId))
	return newElem.Interface()
}

func (this *ParsingService) JobComplete(job *l8tpollaris.CJob, resources ifs.IResources) {
	poll, err := pollaris.Poll(job.PollarisName, job.JobName, resources)
	if err != nil {
		resources.Logger().Error("ParsingCenter:" + err.Error())
		return
	}

	if job.Error != "" {
		resources.Logger().Error("ParsingCenter: job error = ", job.Error)
		return
	}

	if job.Error == "" && poll.Attributes != nil {
		elem := this.createElementInstance(job)
		err = Parser.Parse(job, elem, resources)
		if err != nil {
			resources.Logger().Error("ParsingCenter.JobComplete: ", job.TargetId, " - ", job.PollarisName, " - ", job.JobName, " - ", err.Error())
			return
		}
		if this.vnic == nil {
			resources.Logger().Error("No Vnic to notify inventory")
			return
		}

		cacheServiceName, cacheServiceArea := targets.Links.Cache(job.LinksId)
		fmt.Println(cacheServiceName, ":", cacheServiceArea, ":", elem)
		this.vnic.Leader(cacheServiceName, cacheServiceArea, ifs.PATCH, elem)
	}
}

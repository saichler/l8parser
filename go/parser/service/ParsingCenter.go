package service

import (
	"reflect"

	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8types/go/ifs"
)

func (this *ParsingService) createElementInstance(job *types.CJob) interface{} {
	newElem := reflect.New(reflect.ValueOf(this.elem).Elem().Type())
	field := newElem.Elem().FieldByName(this.primaryKey)
	field.Set(reflect.ValueOf(job.DeviceId))
	return newElem.Interface()
}

func (this *ParsingService) JobComplete(job *types.CJob, resources ifs.IResources) {
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
			resources.Logger().Error("ParsingCenter.JobComplete: ", job.DeviceId, " - ", job.PollarisName, " - ", job.JobName, " - ", err.Error())
			return
		}
		if this.vnic == nil {
			resources.Logger().Error("No Vnic to notify inventory")
			return
		}
		/*
			err := this.vnic.Proximity(job.IService.ServiceName, byte(job.IService.ServiceArea), ifs.PATCH, elem)
			if err != nil {
				this.vnic.Resources().Logger().Error(err.Error())
			} else {
				resources.Logger().Info("ParsingCenter.JobComplete: ", job.DeviceId, " - ", job.PollarisName, " - ", job.JobName, " Patch Sent/Receive")
			}*/
		go func() {
			this.vnic.ProximityRequest(job.IService.ServiceName, byte(job.IService.ServiceArea), ifs.PATCH, elem)
		}()
	}
}

package service

import (
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8types/go/ifs"
	"reflect"
)

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
		newElem := reflect.New(reflect.ValueOf(this.elem).Elem().Type())
		field := newElem.Elem().FieldByName(this.primaryKey)
		field.Set(reflect.ValueOf(job.DeviceId))
		elem := newElem.Interface()
		err := Parser.Parse(job, elem, resources)
		if err != nil {
			panic(err)
		}
		if this.vnic == nil {
			resources.Logger().Error("No Vnic to notify inventory")
			return
		}
		_, err = this.vnic.Proximity(job.IService.ServiceName, byte(job.IService.ServiceArea),
			ifs.PATCH, elem)
		if err != nil {
			this.vnic.Resources().Logger().Error(err.Error())
		}
		this.vnic.Resources().Logger().Info("Sent model to ", job.IService.ServiceName,
			" area ", job.IService.ServiceArea)
	}
}

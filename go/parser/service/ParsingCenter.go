package service

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/types/l8services"

	"github.com/saichler/l8types/go/ifs"
)

func (this *ParsingService) createElementInstance(job *l8tpollaris.CJob) interface{} {
	newElem := reflect.New(reflect.ValueOf(this.elem).Elem().Type())
	field := newElem.Elem().FieldByName(this.primaryKey)
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

		key := linkKey(job.LinkData)
		_, ok := this.registeredLinks.Load(key)
		if !ok {
			job.LinkData.Mode = int32(ifs.M_Leader)
			job.LinkData.Interval = 5
			this.vnic.RegisterServiceLink(job.LinkData)
			this.registeredLinks.Store(key, true)
		}

		this.vnic.Leader(job.LinkData.ZsideServiceName, byte(job.LinkData.ZsideServiceArea), ifs.PATCH, elem)
	}
}

func linkKey(link *l8services.L8ServiceLink) string {
	buff := bytes.Buffer{}
	buff.WriteString(link.ZsideServiceName)
	buff.WriteString(strconv.Itoa(int(link.ZsideServiceArea)))
	return buff.String()
}

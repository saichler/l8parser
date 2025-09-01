package service

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
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
		err = Parser.Parse(job, elem, resources)
		if err != nil {
			resources.Logger().Error("ParsingCenter: ", job.DeviceId, " - ", job.PollarisName, " - ", job.JobName, " - ", err.Error())
			return
		}
		if this.vnic == nil {
			resources.Logger().Error("No Vnic to notify inventory")
			return
		}
		err = this.vnic.Proximity(job.IService.ServiceName, byte(job.IService.ServiceArea),
			ifs.PATCH, elem)
		if err != nil {
			this.vnic.Resources().Logger().Error(err.Error())
		} else {
			name := jobFileName(job)
			if strings.Contains(name, "ifTable") {
				fmt.Println("*************************")
				nd := elem.(*types2.NetworkDevice)
				fmt.Println(nd)
				fmt.Println("**********************")
				time.Sleep(time.Second)
			}
			this.vnic.Resources().Logger().Info("Patch Job ", jobFileName(job))
		}
	}
}

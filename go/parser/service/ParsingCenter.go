package service

import (
	"bytes"
	"reflect"
	"strconv"
	"time"

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
		itemsQueueKey := this.itemsQueueKey(job)
		this.itemsQueueMtx.Lock()
		defer this.itemsQueueMtx.Unlock()
		q, ok := this.itemsQueue[itemsQueueKey]
		if !ok {
			q = NewInventoryQueue(job.IService.ServiceName, byte(job.IService.ServiceArea))
			this.itemsQueue[itemsQueueKey] = q
		}
		q.add(elem)
	}
}

func (this *ParsingService) itemsQueueKey(job *types.CJob) string {
	buff := bytes.Buffer{}
	buff.WriteString(job.IService.ServiceName)
	buff.WriteString(strconv.Itoa(int(job.IService.ServiceArea)))
	return buff.String()
}

func (this *ParsingService) watchItemsQueue() {
	for this.active {
		this.flushItemQueue()
		time.Sleep(time.Second * 5)
	}
}

func (this *ParsingService) flushItemQueue() {
	this.itemsQueueMtx.Lock()
	defer this.itemsQueueMtx.Unlock()
	if this.active {
		for _, q := range this.itemsQueue {
			q.flush(this.vnic)
		}
	}
}

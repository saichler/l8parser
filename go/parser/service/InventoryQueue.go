package service

import "github.com/saichler/l8types/go/ifs"

type InventoryQueue struct {
	queue       []interface{}
	serviceName string
	serviceArea byte
}

func NewInventoryQueue(serviceName string, serviceArea byte) *InventoryQueue {
	iq := &InventoryQueue{}
	iq.queue = make([]interface{}, 0)
	iq.serviceName = serviceName
	iq.serviceArea = serviceArea
	return iq
}

func (this *InventoryQueue) add(item interface{}) {
	this.queue = append(this.queue, item)
}

func (this *InventoryQueue) items() []interface{} {
	items := this.queue
	this.queue = make([]interface{}, 0)
	return items
}

func (this *InventoryQueue) flush(vnic ifs.IVNic) {
	if len(this.queue) > 0 {
		vnic.Proximity(this.serviceName, this.serviceArea, ifs.PATCH, this.items())
	}
}

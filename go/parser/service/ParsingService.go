package service

import (
	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceType = "ParsingService"
)

type ParsingService struct {
	resources  ifs.IResources
	elem       interface{}
	primaryKey string
	vnic       ifs.IVNic
}

func (this *ParsingService) Activate(serviceName string, serviceArea byte,
	r ifs.IResources, l ifs.IServiceCacheListener, args ...interface{}) error {

	this.resources = r
	this.resources.Registry().Register(&types.CMap{})
	this.resources.Registry().Register(&types.CTable{})
	this.resources.Registry().Register(&types.Job{})
	this.elem = args[0]
	this.primaryKey = args[1].(string)
	vnic, ok := l.(ifs.IVNic)
	if ok {
		this.vnic = vnic
	}
	this.resources.Introspector().Inspect(this.elem)
	return nil
}

func (this *ParsingService) DeActivate() error {
	this.vnic = nil
	this.resources = nil
	this.elem = nil
	return nil
}

func (this *ParsingService) Post(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	job := pb.Element().(*types.Job)
	vnic.Resources().Logger().Info("Received Job ", job.PollarisName, ":", job.JobName, " completed!")
	this.JobComplete(job, this.resources)
	return nil
}
func (this *ParsingService) Put(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}
func (this *ParsingService) Patch(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}
func (this *ParsingService) Delete(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}
func (this *ParsingService) Get(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}
func (this *ParsingService) GetCopy(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	return nil
}
func (this *ParsingService) Failed(pb ifs.IElements, vnic ifs.IVNic, msg *ifs.Message) ifs.IElements {
	return nil
}
func (this *ParsingService) TransactionMethod() ifs.ITransactionMethod {
	return nil
}
func (this *ParsingService) WebService() ifs.IWebService {
	return nil
}

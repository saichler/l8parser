package service

import (
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"os"
	"sync"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/strings"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	JobFileLocation = "./jobsPersistency/"
)

type ParsingService struct {
	resources   ifs.IResources
	elem        interface{}
	primaryKey  string
	vnic        ifs.IVNic
	persistJobs bool
	//itemsQueue    map[string]*InventoryQueue
	//itemsQueueMtx *sync.Mutex
	active          bool
	registeredLinks *sync.Map
}

func Activate(linksID string, serviceItem interface{}, persist bool, vnic ifs.IVNic, primaryKeys ...string) {
	parserServiceName, parserServiceArea := targets.Links.Parser(linksID)
	sla := ifs.NewServiceLevelAgreement(&ParsingService{}, parserServiceName, parserServiceArea, true, nil)
	sla.SetServiceItem(serviceItem)
	sla.SetPrimaryKeys(primaryKeys...)
	sla.SetArgs(persist)
	vnic.Resources().Services().Activate(sla, vnic)
}

func (this *ParsingService) Activate(sla *ifs.ServiceLevelAgreement, vnic ifs.IVNic) error {
	this.vnic = vnic
	this.registeredLinks = &sync.Map{}
	this.resources = vnic.Resources()
	this.resources.Registry().Register(&l8tpollaris.CMap{})
	this.resources.Registry().Register(&l8tpollaris.CTable{})
	this.resources.Registry().Register(&l8tpollaris.CJob{})
	this.elem = sla.ServiceItem()
	this.primaryKey = sla.PrimaryKeys()[0]
	this.persistJobs = sla.Args()[0].(bool)
	vnic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(this.elem, sla.PrimaryKeys()...)
	//this.itemsQueueMtx = &sync.Mutex{}
	//this.itemsQueue = make(map[string]*InventoryQueue)
	this.active = true

	this.resources.Introspector().Inspect(this.elem)
	if this.persistJobs {
		os.Mkdir(JobFileLocation, 0777)
	}
	//go this.watchItemsQueue()
	return nil
}

func (this *ParsingService) DeActivate() error {
	//this.itemsQueueMtx.Lock()
	//defer this.itemsQueueMtx.Unlock()
	this.active = false
	this.vnic = nil
	this.resources = nil
	this.elem = nil
	//this.itemsQueue = nil
	return nil
}

func (this *ParsingService) Post(pbs ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	for _, pb := range pbs.Elements() {
		job := pb.(*l8tpollaris.CJob)
		if this.persistJobs {
			data, err := protojson.Marshal(job)
			if err != nil {
				vnic.Resources().Logger().Error("Failed to marshal job to JSON", "error", err)
			} else {
				err = os.WriteFile(jobFileName(job), data, 0777)
				if err != nil {
					vnic.Resources().Logger().Error("Failed to save job to file", "error", err)
				}
			}
		}
		vnic.Resources().Logger().Info("Received Job ", job.TargetId, " - ", job.HostId, " - ", job.PollarisName, " - ", job.JobName, " response")
		this.JobComplete(job, this.resources)
	}
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
func (this *ParsingService) TransactionConfig() ifs.ITransactionConfig {
	return nil
}
func (this *ParsingService) WebService() ifs.IWebService {
	return nil
}

func jobFileName(job *l8tpollaris.CJob) string {
	return strings.New(JobFileLocation, job.PollarisName, ".", job.JobName, ".", job.TargetId, ".", job.HostId).String()
}

func LoadJob(pollarisName, jobName, deviceId, hostId string) (*l8tpollaris.CJob, error) {
	filename := strings.New(JobFileLocation, pollarisName, ".", jobName, ".", deviceId, ".", hostId).String()
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	job := &l8tpollaris.CJob{}
	err = protojson.Unmarshal(data, job)
	return job, err
}

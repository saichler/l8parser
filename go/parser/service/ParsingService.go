package service

import (
	"os"

	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/strings"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	ServiceType     = "ParsingService"
	JobFileLocation = "./jobsPersistency/"
)

type ParsingService struct {
	resources   ifs.IResources
	elem        interface{}
	primaryKey  string
	vnic        ifs.IVNic
	persistJobs bool
}

func (this *ParsingService) Activate(serviceName string, serviceArea byte,
	r ifs.IResources, l ifs.IServiceCacheListener, args ...interface{}) error {

	this.resources = r
	this.resources.Registry().Register(&types.CMap{})
	this.resources.Registry().Register(&types.CTable{})
	this.resources.Registry().Register(&types.CJob{})
	this.elem = args[0]
	this.primaryKey = args[1].(string)
	this.persistJobs = args[2].(bool)

	vnic, ok := l.(ifs.IVNic)
	if ok {
		this.vnic = vnic
	}
	this.resources.Introspector().Inspect(this.elem)
	if this.persistJobs {
		os.Mkdir(JobFileLocation, 0777)
	}
	return nil
}

func (this *ParsingService) DeActivate() error {
	this.vnic = nil
	this.resources = nil
	this.elem = nil
	return nil
}

func (this *ParsingService) Post(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
	job := pb.Element().(*types.CJob)
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

func jobFileName(job *types.CJob) string {
	return strings.New(JobFileLocation, job.PollarisName, ".", job.JobName, ".", job.DeviceId, ".", job.HostId).String()
}

func LoadJob(pollarisName, jobName, deviceId, hostId string) (*types.CJob, error) {
	filename := strings.New(JobFileLocation, pollarisName, ".", jobName, ".", deviceId, ".", hostId).String()
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	job := &types.CJob{}
	err = protojson.Unmarshal(data, job)
	return job, err
}

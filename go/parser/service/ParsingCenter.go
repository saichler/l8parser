/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"fmt"
	"reflect"

	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"google.golang.org/protobuf/proto"
)

// createElementInstance creates a new instance of the configured element type
// and initializes its primary key field with the job's host ID.
//
// HostId is used (not TargetId) so that callers can post multiple targets per
// host without the host's identifying primary key (e.g. ClusterName for K8s
// resources) inheriting per-target uniqueness suffixes. For setups where one
// target == one host, HostId == TargetId and behavior is unchanged.
func (this *ParsingService) createElementInstance(job *l8tpollaris.CJob) interface{} {
	newElem := reflect.New(reflect.ValueOf(this.elem).Elem().Type())
	field := newElem.Elem().FieldByName(this.primaryKey)
	if !field.CanSet() {
		panic("cannot set field " + this.primaryKey)
	}
	field.Set(reflect.ValueOf(job.HostId))
	return newElem.Interface()
}

// JobComplete is called when a collection job completes. It parses the job results
// using the Parser, creates an element instance, and sends the parsed data to the
// inventory cache service via PATCH operation.
// For polls using CTableToInstances, it sends each created instance individually.
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
		instances, err := Parser.ParseMulti(job, elem, resources)
		if err != nil {
			resources.Logger().Error("ParsingCenter.JobComplete: ", job.TargetId, " - ", job.PollarisName, " - ", job.JobName, " - ", err.Error())
			return
		}
		if this.vnic == nil {
			resources.Logger().Error("No Vnic to notify inventory")
			return
		}

		cacheServiceName, cacheServiceArea := targets.Links.Cache(job.LinksId)
		if len(instances) > 0 {
			for _, inst := range instances {
				this.agg.AddElement(inst, ifs.Leader, "", cacheServiceName, cacheServiceArea, ifs.PATCH)
			}
		} else {
			this.agg.AddElement(elem, ifs.Leader, "", cacheServiceName, cacheServiceArea, ifs.PATCH)
		}
	}
}

// HandleDelete processes a delete CJob from the collector. It decodes the
// resource keys from CJob.Result (a serialized CMap), constructs a minimal
// proto instance with primary key fields set, and forwards it to the
// inventory cache with ifs.DELETE action.
func (this *ParsingService) HandleDelete(job *l8tpollaris.CJob) {
	cmap := &l8tpollaris.CMap{}
	if err := proto.Unmarshal(job.Result, cmap); err != nil {
		this.resources.Logger().Error("HandleDelete unmarshal: ", err.Error())
		return
	}

	namespace := string(cmap.Data["namespace"])
	name := string(cmap.Data["name"])

	key := name
	if namespace != "" {
		key = namespace + "/" + name
	}

	cacheServiceName, cacheServiceArea := targets.Links.Cache(job.LinksId)

	elem := this.createElementInstance(job)
	// Set the Key field for the composite primary key
	v := reflect.ValueOf(elem).Elem()
	keyField := v.FieldByName("Key")
	if keyField.IsValid() && keyField.CanSet() {
		keyField.Set(reflect.ValueOf(key))
	}

	fmt.Printf("[PARSER-FWD-DELETE] linksId=%s cluster=%s key=%s -> cache=(%s,%d)\n",
		job.LinksId, job.HostId, key, cacheServiceName, cacheServiceArea)

	this.agg.AddElement(elem, ifs.Leader, "", cacheServiceName, cacheServiceArea, ifs.DELETE)
}

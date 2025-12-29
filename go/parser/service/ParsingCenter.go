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
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"reflect"
)

// createElementInstance creates a new instance of the configured element type
// and initializes its primary key field with the job's target ID.
func (this *ParsingService) createElementInstance(job *l8tpollaris.CJob) interface{} {
	newElem := reflect.New(reflect.ValueOf(this.elem).Elem().Type())
	field := newElem.Elem().FieldByName(this.primaryKey)
	if !field.CanSet() {
		panic("cannot set field " + this.primaryKey)
	}
	fmt.Println(job.TargetId)
	field.Set(reflect.ValueOf(job.TargetId))
	return newElem.Interface()
}

// JobComplete is called when a collection job completes. It parses the job results
// using the Parser, creates an element instance, and sends the parsed data to the
// inventory cache service via PATCH operation.
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

		cacheServiceName, cacheServiceArea := targets.Links.Cache(job.LinksId)
		fmt.Println(cacheServiceName, ":", cacheServiceArea, ":", elem)
		this.agg.AddElement(elem, ifs.Leader, "", cacheServiceName, cacheServiceArea, ifs.PATCH)
		//this.vnic.Leader(cacheServiceName, cacheServiceArea, ifs.PATCH, elem)
	}
}

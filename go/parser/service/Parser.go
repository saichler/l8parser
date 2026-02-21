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

// Package service provides the core parsing service for L8Parser.
// It contains the Parser singleton that orchestrates rule execution for transforming
// collected data into structured inventory objects, and the ParsingService that
// integrates with the L8 ecosystem for processing collection jobs.
package service

import (
	"errors"
	"fmt"

	"github.com/saichler/l8parser/go/parser/rules"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"

	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
)

// _Parser is the main parser engine that manages parsing rules and executes them on job results.
// It maintains a registry of all available parsing rules and coordinates their execution.
type _Parser struct {
	rules map[string]rules.ParsingRule
}

// Parser is the singleton instance of the parser engine, initialized with all available rules.
var Parser = newParser()

// newParser creates and initializes a new Parser instance with all registered parsing rules.
func newParser() *_Parser {
	p := &_Parser{}
	p.rules = make(map[string]rules.ParsingRule)
	con := &rules.Contains{}
	p.rules[con.Name()] = con
	set := &rules.Set{}
	p.rules[set.Name()] = set
	totable := &rules.StringToCTable{}
	p.rules[totable.Name()] = totable
	tableToMap := &rules.CTableToMapProperty{}
	p.rules[tableToMap.Name()] = tableToMap
	ifTableToPhysicals := &rules.IfTableToPhysicals{}
	p.rules[ifTableToPhysicals.Name()] = ifTableToPhysicals
	entityMibToPhysicals := &rules.EntityMibToPhysicals{}
	p.rules[entityMibToPhysicals.Name()] = entityMibToPhysicals
	inferDeviceType := &rules.InferDeviceType{}
	p.rules[inferDeviceType.Name()] = inferDeviceType
	mapToDeviceStatus := &rules.MapToDeviceStatus{}
	p.rules[mapToDeviceStatus.Name()] = mapToDeviceStatus
	setTimeSeries := &rules.SetTimeSeries{}
	p.rules[setTimeSeries.Name()] = setTimeSeries
	return p
}

// Parse executes the parsing rules for a completed collection job.
// It deserializes the job result, looks up the corresponding poll configuration,
// and executes each rule defined in the poll's attributes to transform the data.
// Parameters: job (completed collection job), any (target object to populate), resources (system resources).
func (this *_Parser) Parse(job *l8tpollaris.CJob, any interface{}, resources ifs.IResources) error {
	fmt.Println("Parser.Parse: pollarisName=", job.PollarisName, "jobName=", job.JobName, "targetId=", job.TargetId)

	if job.Error != "" {
		fmt.Println("Parser.Parse: job has error=", job.Error)
		return errors.New(job.Error)
	}

	if job.Result == nil || len(job.Result) < 4 {
		fmt.Println("Parser.Parse: invalid job result, len=", len(job.Result))
		return resources.Logger().Error("Invalid job result ", job.TargetId, " - ", job.PollarisName,
			" - ", job.JobName, " - ", string(job.Result))
	}

	workSpace := make(map[string]interface{})
	enc := object.NewDecode(job.Result, 0, resources.Registry())
	data, err := enc.Get()
	if err != nil {
		fmt.Println("Parser.Parse: decode error=", err)
		return resources.Logger().Error(err)
	}
	fmt.Println("Parser.Parse: decoded data type=", fmt.Sprintf("%T", data))

	poll, err := pollaris.Poll(job.PollarisName, job.JobName, resources)
	if err != nil {
		fmt.Println("Parser.Parse: cannot find poll, pollarisName=", job.PollarisName, "jobName=", job.JobName)
		return resources.Logger().Error("cannot find poll for polaris ", job.PollarisName, ":", job.JobName)
	}
	workSpace[rules.Input] = data
	workSpace[rules.JobEnded] = job.Ended
	if poll.Attributes == nil {
		fmt.Println("Parser.Parse: no attributes defined")
		return resources.Logger().Error("No attributes are defined on pollaris "+job.PollarisName, ":", job.JobName)
	}

	fmt.Println("Parser.Parse: processing", len(poll.Attributes), "attributes")
	for i, attr := range poll.Attributes {
		fmt.Println("Parser.Parse: attr[", i, "] propertyId=", attr.PropertyId, "rules count=", len(attr.Rules))
		workSpace[rules.PropertyId] = attr.PropertyId
		for j, rData := range attr.Rules {
			fmt.Println("Parser.Parse: attr[", i, "] rule[", j, "] name=", rData.Name)
			if rData.Params != nil {
				for p, v := range rData.Params {
					workSpace[p] = v.Value
				}
			}
			ruleImpl, ok := this.rules[rData.Name]
			if !ok {
				fmt.Println("Parser.Parse: rule not found:", rData.Name)
				return resources.Logger().Error("Cannot find parsing rule ", rData.Name)
			}
			err = ruleImpl.Parse(resources, workSpace, rData.Params, any, poll.What)
			if err != nil {
				fmt.Println("Parser.Parse: rule error:", err)
				return err
			}
		}
	}
	fmt.Println("Parser.Parse: completed successfully for", job.PollarisName, "/", job.JobName)
	return nil
}

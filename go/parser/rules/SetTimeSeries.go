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

package rules

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
)

// SetTimeSeries is a parsing rule that creates an L8TimeSeriesPoint from the input value
// and the job's end timestamp, then sets it on the target property.
// The property framework handles appending to the repeated field.
// This is used for time-series fields like CpuUsagePercent and MemoryUsagePercent.
type SetTimeSeries struct{}

// Name returns the rule identifier "SetTimeSeries".
func (this *SetTimeSeries) Name() string {
	return "SetTimeSeries"
}

// ParamNames returns the required parameter names for this rule.
func (this *SetTimeSeries) ParamNames() []string {
	return []string{From}
}

// Parse executes the SetTimeSeries rule logic.
// It extracts a value from input, converts it to float64, gets the job end timestamp
// from the workspace, creates an L8TimeSeriesPoint, and sets it on the target property.
func (this *SetTimeSeries) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*l8tpollaris.L8PParameter, any interface{}, pollWhat string) error {
	input := workSpace[Input]
	_propertyId := workSpace[PropertyId]
	propertyId := _propertyId.(string)

	fmt.Println("SetTimeSeries: propertyId=", propertyId)
	fmt.Println("SetTimeSeries: input type=", fmt.Sprintf("%T", input))
	fmt.Println("SetTimeSeries: pollWhat=", pollWhat)

	if input == nil {
		fmt.Println("SetTimeSeries: input is nil")
		return resources.Logger().Error("nil input for SetTimeSeries")
	}

	value, kind, err := GetValueInput(resources, input, params, pollWhat)
	fmt.Println("SetTimeSeries: value=", value, "kind=", kind, "err=", err)
	if err != nil {
		return err
	}

	if value == nil {
		fmt.Println("SetTimeSeries: value is nil for propertyId=", propertyId)
		return resources.Logger().Error("nil value for property id", propertyId)
	}

	floatValue, err := convertToFloat64(value, kind)
	fmt.Println("SetTimeSeries: floatValue=", floatValue, "convertErr=", err)
	if err != nil {
		return resources.Logger().Error("SetTimeSeries: cannot convert value to float64:", err.Error())
	}

	var stamp int64
	if ended, ok := workSpace[JobEnded]; ok {
		if s, ok := ended.(int64); ok {
			stamp = s
		}
	}
	fmt.Println("SetTimeSeries: stamp=", stamp)

	point := &l8api.L8TimeSeriesPoint{
		Stamp: stamp,
		Value: floatValue,
	}
	fmt.Println("SetTimeSeries: created point=", point)

	if _propertyId != nil {
		modifiedPropertyId := injectIndexOrKey(propertyId, workSpace)
		fmt.Println("SetTimeSeries: modifiedPropertyId=", modifiedPropertyId)

		instance, err := properties.PropertyOf(modifiedPropertyId, resources)
		fmt.Println("SetTimeSeries: PropertyOf instance=", instance, "err=", err)
		if err != nil {
			return resources.Logger().Error("error parsing instance path", err.Error())
		}
		if instance != nil {
			_, _, err := instance.Set(any, point)
			fmt.Println("SetTimeSeries: Set result err=", err)
			if err != nil {
				return resources.Logger().Error("error setting time series value:", err.Error())
			}
		} else {
			fmt.Println("SetTimeSeries: instance is nil for modifiedPropertyId=", modifiedPropertyId)
		}
	} else {
		fmt.Println("SetTimeSeries: _propertyId is nil")
	}
	workSpace[Output] = point
	fmt.Println("SetTimeSeries: completed successfully for propertyId=", propertyId)
	return nil
}

// convertToFloat64 converts various numeric types to float64.
func convertToFloat64(value interface{}, kind reflect.Kind) (float64, error) {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(value).Float(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(reflect.ValueOf(value).Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(reflect.ValueOf(value).Uint()), nil
	case reflect.String:
		return strconv.ParseFloat(value.(string), 64)
	default:
		return 0, errors.New("unsupported type for float64 conversion: " + kind.String())
	}
}

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
	"reflect"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
)

// Set is a parsing rule that directly sets a value from input to a target property.
// It supports PropertyId path injection for nested collections and handles type-safe value assignment.
// Parameters: "from" (source field in the input data).
type Set struct{}

// Name returns the rule identifier "Set".
func (this *Set) Name() string {
	return "Set"
}

// ParamNames returns the required parameter names for this rule.
func (this *Set) ParamNames() []string {
	return []string{}
}

// Parse executes the Set rule logic, extracting a value and setting it to the target property.
func (this *Set) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*l8tpollaris.L8PParameter, any interface{}, pollWhat string) error {
	input := workSpace[Input]
	_propertyId := workSpace[PropertyId]
	propertyId := _propertyId.(string)

	if input == nil {
		return resources.Logger().Error("nil input for job")
	}

	value, _, err := GetValueInput(resources, input, params, pollWhat)
	if err != nil || value == nil {
		// Missing/blank OID data is expected for some devices — skip gracefully
		return nil
	}

	// Skip SNMP error strings gracefully - device doesn't support this OID
	if strVal, ok := value.(string); ok {
		if isSnmpErrorString(strVal) {
			return nil
		}
	}

	if _propertyId != nil {
		// Inject slice index or map key into PropertyId before creating property instance
		modifiedPropertyId := injectIndexOrKey(propertyId, workSpace)

		instance, err := properties.PropertyOf(modifiedPropertyId, resources)
		if err != nil {
			return resources.Logger().Error("error parsing instance path", err.Error())
		}
		if instance != nil {
			value = coerceValue(resources, value, instance, workSpace)
			_, _, err = instance.Set(any, value)
			if err != nil {
				return resources.Logger().Error("error setting property value:", err.Error())
			}
		}
	}
	workSpace[Output] = value
	return nil
}

// coerceValue converts the input value to match the target property's type when
// a direct assignment would fail. Handles SNMP int64→bool (TruthValue) and
// int64→*L8TimeSeriesPoint conversions.
func coerceValue(resources ifs.IResources, value interface{}, instance *properties.Property, workSpace map[string]interface{}) interface{} {
	node := instance.Node()
	if node == nil {
		return value
	}
	typeName := node.TypeName

	valueKind := reflect.TypeOf(value).Kind()

	// int64 → bool (SNMP TruthValue: 1=true, 2=false; ifAdminStatus: 1=up, 2=down)
	if typeName == "bool" && (valueKind == reflect.Int64 || valueKind == reflect.Int) {
		return reflect.ValueOf(value).Int() == 1
	}

	// int64/uint64/float64 → *L8TimeSeriesPoint (wrap numeric into time series point)
	if typeName == "L8TimeSeriesPoint" {
		var floatVal float64
		switch valueKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			floatVal = float64(reflect.ValueOf(value).Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			floatVal = float64(reflect.ValueOf(value).Uint())
		case reflect.Float32, reflect.Float64:
			floatVal = reflect.ValueOf(value).Float()
		default:
			return value
		}
		var stamp int64
		if ended, ok := workSpace[JobEnded]; ok {
			if s, ok := ended.(int64); ok {
				stamp = s
			}
		}
		return &l8api.L8TimeSeriesPoint{Stamp: stamp, Value: floatVal}
	}

	// Check if the value type matches the node type
	valueType := reflect.TypeOf(value).String()
	if valueType != typeName {
		propertyId, _ := instance.PropertyId()
		resources.Logger().Info("coerceValue type mismatch: property=%s, nodeType=%s, valueType=%s, value=%v",
			propertyId, typeName, valueType, value)
		info, err := resources.Registry().Info(typeName)
		if err != nil {
			return defaultValueForType(typeName)
		}
		value, err = info.NewInstance()
		if err != nil {
			return defaultValueForType(typeName)
		}
	}

	return value
}

// defaultValueForType returns the zero/default value for a given type name.
func defaultValueForType(typeName string) interface{} {
	switch typeName {
	case "string":
		return ""
	case "bool":
		return false
	case "int32":
		return int32(0)
	case "int64":
		return int64(0)
	case "uint32":
		return uint32(0)
	case "uint64":
		return uint64(0)
	case "float32":
		return float32(0)
	case "float64":
		return float64(0)
	default:
		return nil
	}
}

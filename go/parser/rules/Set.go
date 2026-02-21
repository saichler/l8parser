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
	"fmt"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8types/go/ifs"
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

		fmt.Println("[DEBUG Set] propertyId=", propertyId, " modifiedPropertyId=", modifiedPropertyId, " valueType=", fmt.Sprintf("%T", value))

		instance, err := properties.PropertyOf(modifiedPropertyId, resources)
		if err != nil {
			return resources.Logger().Error("error parsing instance path", err.Error())
		}
		if instance != nil {
			fmt.Println("[DEBUG Set] calling instance.Set for modifiedPropertyId=", modifiedPropertyId)
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("[DEBUG Set] PANIC on propertyId=", propertyId, " modifiedPropertyId=", modifiedPropertyId, " value=", value, " recover=", r)
					}
				}()
				_, _, err = instance.Set(any, value)
			}()
			if err != nil {
				return resources.Logger().Error("error setting property value:", err.Error())
			}
			fmt.Println("[DEBUG Set] instance.Set succeeded for modifiedPropertyId=", modifiedPropertyId)
		}
	}
	workSpace[Output] = value
	return nil
}

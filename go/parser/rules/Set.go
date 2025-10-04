package rules

import (
	"fmt"
	"strings"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8reflect/go/reflect/properties"
)

type Set struct{}

func (this *Set) Name() string {
	return "Set"
}

func (this *Set) ParamNames() []string {
	return []string{}
}

func (this *Set) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*l8tpollaris.L8P_Parameter, any interface{}, pollWhat string) error {
	input := workSpace[Input]
	_propertyId := workSpace[PropertyId]
	propertyId := _propertyId.(string)

	if input == nil {
		return resources.Logger().Error("nil input for job")
	}

	value, _, err := GetValueInput(resources, input, params, pollWhat)
	if err != nil {
		return err
	}

	if value == nil {
		return resources.Logger().Error("nil value for property id", propertyId)
	}

	if _propertyId != nil {
		// Inject slice index or map key into PropertyId before creating property instance
		modifiedPropertyId := injectIndexOrKey(propertyId, workSpace)

		instance, err := properties.PropertyOf(modifiedPropertyId, resources)
		if err != nil {
			return resources.Logger().Error("error parsing instance path", err.Error())
		}
		if instance != nil {
			_, _, err := instance.Set(any, value)
			if err != nil {
				fmt.Println(value)
				return resources.Logger().Error("error setting property value:", err.Error())
			}
		}
	}
	workSpace[Output] = value
	return nil
}

// injectIndexOrKey injects slice indices or map keys into PropertyId paths
// Format: <{reflect.Kind}value> before the attribute that needs indexing
func injectIndexOrKey(propertyId string, workSpace map[string]interface{}) string {
	// Map of collection attributes that need indexing/keying
	collectionMappings := map[string]string{
		"physicals":      "{24}physical-0", // map<string, Physical> - use string key
		"logicals":       "{24}logical-0",  // map<string, Logical> - use string key
		"networklinks":   "{2}0",           // repeated NetworkLink - use int index (alt name)
		"network_links":  "{2}0",           // repeated NetworkLink - use int index
		"chassis":        "{2}0",           // repeated Chassis - use int index
		"ports":          "{2}0",           // repeated Port - use int index
		"power_supplies": "{2}0",           // repeated PowerSupply - use int index
		"powersupplies":  "{2}0",           // repeated PowerSupply - use int index (alt name)
		"fans":           "{2}0",           // repeated Fan - use int index
		"modules":        "{2}0",           // repeated Module - use int index
		"cpus":           "{2}0",           // repeated Cpu - use int index
		"interfaces":     "{2}0",           // repeated Interface - use int index
		"processes":      "{2}0",           // repeated ProcessInfo - use int index
	}

	// Field name mappings for proto compatibility
	fieldMappings := map[string]string{
		"powersupplies": "powersupplies", // powersupplies -> power_supplies
		"networklinks":  "networklinks",  // networklinks -> network_links
		"networkhealth": "networkhealth", // networkhealth -> network_health
	}

	parts := strings.Split(propertyId, ".")
	result := make([]string, 0, len(parts))

	for i, part := range parts {
		// Apply field name mapping first
		mappedPart := part
		if mapped, exists := fieldMappings[part]; exists {
			mappedPart = mapped
		}

		// Apply collection indexing
		if indexKey, exists := collectionMappings[part]; exists {
			// Check if this is not the last part (we need a following attribute)
			if i < len(parts)-1 {
				// Inject the index/key before the next attribute
				result = append(result, mappedPart+"<"+indexKey+">")
			} else {
				result = append(result, mappedPart)
			}
		} else {
			result = append(result, mappedPart)
		}
	}

	modifiedId := strings.Join(result, ".")

	return modifiedId
}

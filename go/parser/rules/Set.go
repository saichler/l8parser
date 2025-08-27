package rules

import (
	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/reflect/go/reflect/properties"
)

type Set struct{}

func (this *Set) Name() string {
	return "Set"
}

func (this *Set) ParamNames() []string {
	return []string{}
}

func (this *Set) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*types.Parameter, any interface{}) error {
	input := workSpace[Input]
	path := workSpace[PropertyId]

	if input == nil {
		return resources.Logger().Error("nil input for job")
	}

	value, _, err := GetValueInput(resources, input, params)
	if err != nil {
		return err
	}

	if path != nil {
		instance, err := properties.PropertyOf(path.(string), resources)
		if err != nil {
			return resources.Logger().Error("error parsing instance path", err.Error())
		}
		if instance != nil {
			if instance.Node().TypeName == "string" {
				value = string(value.([]byte))
			}
			_, _, err := instance.Set(any, value)
			if err != nil {
				return resources.Logger().Error("error setting property value:", err.Error())
			}
		}
	}
	workSpace[Output] = value
	return nil
}

package rules

import (
	"errors"
	"strconv"

	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
)

type ParsingRule interface {
	Name() string
	ParamNames() []string
	Parse(ifs.IResources, map[string]interface{}, map[string]*types.Parameter, interface{}) error
}

func getStringInput(resources ifs.IResources, input interface{}, params map[string]*types.Parameter) (string, error) {
	m, ok := input.(*types.CMap)
	if ok {
		from := params[From]
		if from == nil {
			return "", resources.Logger().Error("missing 'from' key in map input")
		}
		strData := m.Data[from.Value]
		if strData == nil || len(strData) == 0 {
			resources.Logger().Error("Value for From ", from.Name, " is blank")
			return "", errors.New("Value for From " + from.Name + " is blank")
		}
		enc := object.NewDecode(strData, 0, resources.Registry())
		strInt, _ := enc.Get()
		str, ok := strInt.(string)
		if ok {
			return str, nil
		}
		byts, ok := strInt.([]byte)
		if ok {
			return string(byts), nil
		}
		return "", resources.Logger().Error("'from' key not a string")
	}
	byts, ok := input.([]byte)
	if ok {
		return string(byts), nil
	}
	return "", resources.Logger().Error("'from' key not a []byte")
}

func getIntInput(workSpace map[string]interface{}, paramName string) (int, error) {
	v, ok := workSpace[paramName].(string)
	if !ok {
		return -1, errors.New("'" + paramName + "' does not exist")
	}
	i, e := strconv.Atoi(v)
	if e != nil {
		return -1, e
	}
	return i, nil
}

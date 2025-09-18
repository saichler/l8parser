package rules

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/saichler/l8pollaris/go/types/l8poll"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
)

type ParsingRule interface {
	Name() string
	ParamNames() []string
	Parse(ifs.IResources, map[string]interface{}, map[string]*l8poll.L8P_Parameter, interface{}, string) error
}

func convertToString(value interface{}, kind reflect.Kind) (string, error) {
	switch kind {
	case reflect.String:
		return value.(string), nil
	case reflect.Slice:
		if byts, ok := value.([]byte); ok {
			return string(byts), nil
		}
		return "", errors.New("slice is not []byte")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(reflect.ValueOf(value).Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(reflect.ValueOf(value).Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(reflect.ValueOf(value).Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(value.(bool)), nil
	default:
		return "", errors.New("unsupported type for string conversion: " + kind.String())
	}
}

// GetValueInput extracts any value type from input data and returns the value, its reflect.Kind, and any error
func GetValueInput(resources ifs.IResources, input interface{}, params map[string]*l8poll.L8P_Parameter, pollWhat string) (interface{}, reflect.Kind, error) {
	m, ok := input.(*l8poll.CMap)
	if ok {
		if len(m.Data) == 0 {
			return nil, reflect.Invalid, errors.New("no data found in map:" + pollWhat)
		}
		from := params[From]
		if from == nil {
			return nil, reflect.Invalid, errors.New("missing 'from' key in map input")
		}
		rawData := m.Data[from.Value]
		if rawData == nil || len(rawData) == 0 {
			return nil, reflect.Invalid, errors.New("Value for From " + from.Value + " is blank")
		}
		enc := object.NewDecode(rawData, 0, resources.Registry())
		value, err := enc.Get()
		if err != nil {
			return nil, reflect.Invalid, errors.New("failed to decode value: " + err.Error())
		}
		if value == nil {
			return nil, reflect.Invalid, errors.New("failed to decode value")
		}
		return value, reflect.TypeOf(value).Kind(), nil
	}

	byts, ok := input.([]byte)
	if ok {
		return byts, reflect.Slice, nil
	}

	// Handle direct values
	if input != nil {
		return input, reflect.TypeOf(input).Kind(), nil
	}

	return nil, reflect.Invalid, errors.New("unsupported input type")
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

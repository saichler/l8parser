package rules

import (
	"reflect"

	"github.com/saichler/l8pollaris/go/types/l8poll"
	"github.com/saichler/l8types/go/ifs"
	problerTypes "github.com/saichler/probler/go/types"
	"github.com/saichler/l8reflect/go/reflect/properties"
)

type MapToDeviceStatus struct{}

func (this *MapToDeviceStatus) Name() string {
	return "MapToDeviceStatus"
}

func (this *MapToDeviceStatus) ParamNames() []string {
	return []string{From}
}

func (this *MapToDeviceStatus) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*l8poll.L8P_Parameter, any interface{}, pollWhat string) error {
	input := workSpace[Input]
	_propertyId := workSpace[PropertyId]
	propertyId := _propertyId.(string)

	if input == nil {
		return resources.Logger().Error("nil input for MapToDeviceStatus")
	}

	value, kind, err := GetValueInput(resources, input, params, pollWhat)
	if err != nil {
		return err
	}

	if value == nil {
		return resources.Logger().Error("nil value for property id", propertyId)
	}

	// Convert map[int32]bool to device status enum
	deviceStatus := problerTypes.DeviceStatus_DEVICE_STATUS_UNKNOWN
	if kind == reflect.Map {
		statusMap, ok := value.(map[int32]bool)
		if ok {
			// Analyze the map to determine overall device status
			// If any protocol is online (true), consider device ONLINE
			// If all protocols are offline (false), consider device OFFLINE
			hasOnline := false
			hasOffline := false

			for _, isOnline := range statusMap {
				if isOnline {
					hasOnline = true
				} else {
					hasOffline = true
				}
			}

			if hasOnline && !hasOffline {
				deviceStatus = problerTypes.DeviceStatus_DEVICE_STATUS_ONLINE
			} else if !hasOnline && hasOffline {
				deviceStatus = problerTypes.DeviceStatus_DEVICE_STATUS_OFFLINE
			} else if hasOnline && hasOffline {
				deviceStatus = problerTypes.DeviceStatus_DEVICE_STATUS_PARTIAL
			} else {
				deviceStatus = problerTypes.DeviceStatus_DEVICE_STATUS_UNKNOWN
			}
		} else {
			// Try to handle other map types if needed
			resources.Logger().Error("Expected map[int32]bool, got different map type")
			deviceStatus = problerTypes.DeviceStatus_DEVICE_STATUS_UNKNOWN
		}
	} else {
		// If not a map, assume unknown
		resources.Logger().Error("Expected map[int32]bool for device status, got:", kind.String())
		deviceStatus = problerTypes.DeviceStatus_DEVICE_STATUS_UNKNOWN
	}

	if _propertyId != nil {
		instance, err := properties.PropertyOf(propertyId, resources)
		if err != nil {
			return resources.Logger().Error("error parsing instance path", err.Error())
		}
		if instance != nil {
			_, _, err := instance.Set(any, deviceStatus)
			if err != nil {
				return resources.Logger().Error("error setting device status value:", err.Error())
			}
		}
	}

	workSpace[Output] = deviceStatus
	return nil
}

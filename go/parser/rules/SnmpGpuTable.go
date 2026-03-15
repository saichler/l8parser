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
	"reflect"
	"strconv"
	"strings"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
)

// SnmpGpuTable is a parsing rule that transforms NVIDIA GPU SNMP table data
// (indexed OIDs: {base}.{metric_id}.{gpu_index}) into GpuDevice.Gpus repeated fields.
// It groups entries by GPU index and maps each metric OID suffix to the corresponding
// property using either Set (static) or SetTimeSeries (dynamic) semantics.
//
// Parameters:
//   - "oid_base": base OID prefix (e.g., "1.3.6.1.4.1.53246.1.1.1.1")
//   - "mapping": comma-separated "oidSuffix:propertyName:type" triples
//     type is "set" for static or "ts" for time series
//     Example: "1:devicename:set,2:deviceuuid:set,5:gpuutilizationpercent:ts"
type SnmpGpuTable struct{}

// Name returns the rule identifier "SnmpGpuTable".
func (this *SnmpGpuTable) Name() string {
	return "SnmpGpuTable"
}

// ParamNames returns the required parameter names for this rule.
func (this *SnmpGpuTable) ParamNames() []string {
	return []string{"oid_base", "mapping"}
}

// gpuFieldMapping holds a parsed mapping entry.
type gpuFieldMapping struct {
	oidSuffix    int
	propertyName string
	isTimeSeries bool
}

// Parse executes the SnmpGpuTable rule logic.
func (this *SnmpGpuTable) Parse(resources ifs.IResources, workSpace map[string]interface{},
	params map[string]*l8tpollaris.L8PParameter, any interface{}, pollWhat string) error {

	input := workSpace[Input]
	_propertyId := workSpace[PropertyId]
	if input == nil || _propertyId == nil {
		return nil
	}

	propertyId := _propertyId.(string)

	cmap, ok := input.(*l8tpollaris.CMap)
	if !ok || cmap == nil || len(cmap.Data) == 0 {
		return nil
	}

	// Parse parameters
	oidBaseParam := params["oid_base"]
	mappingParam := params["mapping"]
	if oidBaseParam == nil || mappingParam == nil {
		return resources.Logger().Error("SnmpGpuTable: missing oid_base or mapping parameter")
	}

	oidBase := oidBaseParam.Value
	// Normalize: ensure oid_base starts with a dot
	if !strings.HasPrefix(oidBase, ".") {
		oidBase = "." + oidBase
	}

	mappings := parseMappings(mappingParam.Value)
	if len(mappings) == 0 {
		return resources.Logger().Error("SnmpGpuTable: no valid mappings parsed from: ", mappingParam.Value)
	}

	// Build a lookup from OID suffix to mapping
	suffixMap := make(map[int]*gpuFieldMapping)
	for i := range mappings {
		suffixMap[mappings[i].oidSuffix] = &mappings[i]
	}

	// Get job timestamp for time series
	var stamp int64
	if ended, ok := workSpace[JobEnded]; ok {
		if s, ok := ended.(int64); ok {
			stamp = s
		}
	}

	// Iterate CMap entries and group by GPU index
	for oidKey, rawData := range cmap.Data {
		if len(rawData) == 0 {
			continue
		}

		// Extract metric_id and gpu_index from OID key
		// OID format: {oidBase}.{metric_id}.{gpu_index}
		if !strings.HasPrefix(oidKey, oidBase) {
			continue
		}

		suffix := oidKey[len(oidBase):]
		if strings.HasPrefix(suffix, ".") {
			suffix = suffix[1:]
		}

		parts := strings.SplitN(suffix, ".", 2)
		if len(parts) != 2 {
			continue
		}

		metricId, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		gpuIndex, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		mapping, exists := suffixMap[metricId]
		if !exists {
			continue
		}

		// Decode the value
		enc := object.NewDecode(rawData, 0, resources.Registry())
		value, err := enc.Get()
		if err != nil || value == nil {
			continue
		}

		// Skip SNMP error strings
		if strVal, ok := value.(string); ok {
			if isSnmpErrorString(strVal) {
				continue
			}
		}

		// Build the full property path with GPU index injection
		// e.g., "gpudevice.gpus" -> "gpudevice.gpus<{2}INDEX>.devicename"
		fullPropertyId := fmt.Sprintf("%s<{2}%d>.%s", propertyId, gpuIndex, mapping.propertyName)
		fullPropertyId = injectIndexOrKey(fullPropertyId, workSpace)

		if mapping.isTimeSeries {
			// Convert to time series point
			floatVal, err := toFloat64(value)
			if err != nil {
				continue
			}
			point := &l8api.L8TimeSeriesPoint{Stamp: stamp, Value: floatVal}
			instance, err := properties.PropertyOf(fullPropertyId, resources)
			if err != nil || instance == nil {
				continue
			}
			_, _, err = instance.Set(any, point)
			if err != nil {
				resources.Logger().Error("SnmpGpuTable: error setting time series for GPU ", gpuIndex, ":", err.Error())
			}
		} else {
			// Set static value
			instance, err := properties.PropertyOf(fullPropertyId, resources)
			if err != nil || instance == nil {
				continue
			}
			value = coerceValue(resources, value, instance, workSpace)
			_, _, err = instance.Set(any, value)
			if err != nil {
				resources.Logger().Error("SnmpGpuTable: error setting value for GPU ", gpuIndex, ":", err.Error())
			}
		}
	}

	return nil
}

// parseMappings parses a comma-separated mapping string into gpuFieldMapping entries.
// Format: "oidSuffix:propertyName:type,oidSuffix:propertyName:type,..."
// type is "set" for static or "ts" for time series.
func parseMappings(mappingStr string) []gpuFieldMapping {
	result := make([]gpuFieldMapping, 0)
	entries := strings.Split(mappingStr, ",")
	for _, entry := range entries {
		parts := strings.SplitN(strings.TrimSpace(entry), ":", 3)
		if len(parts) != 3 {
			continue
		}
		oidSuffix, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		result = append(result, gpuFieldMapping{
			oidSuffix:    oidSuffix,
			propertyName: parts[1],
			isTimeSeries: parts[2] == "ts",
		})
	}
	return result
}

// toFloat64 converts a value to float64 for time series points.
func toFloat64(value interface{}) (float64, error) {
	kind := reflect.TypeOf(value).Kind()
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
		return 0, fmt.Errorf("unsupported type for float64: %s", kind.String())
	}
}

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
	"reflect"
	"strings"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
)

// CTableToInstances converts a CTable into individual model instances, one per row.
// Each row becomes a separate proto instance with fields mapped from table columns.
// ClusterName is set from the job's target ID (stored in workspace as TargetId),
// and Key is built from the key columns (e.g., "namespace/name").
type CTableToInstances struct{}

func (this *CTableToInstances) Name() string {
	return "CTableToInstances"
}

func (this *CTableToInstances) ParamNames() []string {
	return []string{""}
}

func (this *CTableToInstances) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*l8tpollaris.L8PParameter, any interface{}, pollWhat string) error {
	table, ok := workSpace[Output].(*l8tpollaris.CTable)
	if !ok {
		outVal := workSpace[Output]
		actualType := "nil"
		if outVal != nil {
			actualType = reflect.TypeOf(outVal).String()
		}
		return errors.New("CTableToInstances: expected *CTable but got " + actualType)
	}

	keyColumns, e := getIntArrInput(workSpace, KeyColumn)
	if e != nil {
		return e
	}

	targetId, _ := workSpace[TargetId].(string)

	elemType := reflect.ValueOf(any).Elem().Type()
	toString := strings2.New()
	toString.TypesPrefix = true

	instances := make([]interface{}, 0, len(table.Rows))
	for _, row := range table.Rows {
		inst := reflect.New(elemType)
		instElem := inst.Elem()

		keyBuilder := strings2.New()
		for ki, kc := range keyColumns {
			val := getValue(row.Data[int32(kc)], resources)
			if val == nil {
				continue
			}
			if ki > 0 {
				keyBuilder.Add("/")
			}
			keyBuilder.Add(toString.ToString(reflect.ValueOf(val)))
		}

		for i := 0; i < len(table.Columns); i++ {
			val := getValue(row.Data[int32(i)], resources)
			if val == nil {
				continue
			}
			attrName := getAttributeNameFromColumn(table.Columns[int32(i)])
			field := findFieldByJsonName(instElem, attrName)
			if !field.IsValid() || !field.CanSet() {
				continue
			}
			setFieldValue(field, val)
		}

		clusterField := instElem.FieldByName("ClusterName")
		if clusterField.IsValid() && clusterField.CanSet() {
			clusterField.Set(reflect.ValueOf(targetId))
		}
		keyField := instElem.FieldByName("Key")
		if keyField.IsValid() && keyField.CanSet() {
			keyField.Set(reflect.ValueOf(keyBuilder.String()))
		}

		instances = append(instances, inst.Interface())
	}

	workSpace[Instances] = instances
	return nil
}

func findFieldByJsonName(v reflect.Value, jsonName string) reflect.Value {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		tag := sf.Tag.Get("protobuf")
		if tag == "" {
			continue
		}
		name := extractJsonName(tag)
		if strings.EqualFold(name, jsonName) {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}

func extractJsonName(tag string) string {
	const prefix = "json="
	idx := 0
	for idx < len(tag) {
		pos := indexOf(tag[idx:], prefix)
		if pos == -1 {
			return ""
		}
		start := idx + pos + len(prefix)
		end := indexOf(tag[start:], ",")
		if end == -1 {
			return tag[start:]
		}
		return tag[start : start+end]
	}
	return ""
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func setFieldValue(field reflect.Value, val interface{}) {
	valRef := reflect.ValueOf(val)
	if valRef.Type().AssignableTo(field.Type()) {
		field.Set(valRef)
		return
	}
	if valRef.Type().ConvertibleTo(field.Type()) {
		field.Set(valRef.Convert(field.Type()))
		return
	}
	if field.Kind() == reflect.String {
		str := strings2.New()
		str.TypesPrefix = true
		field.Set(reflect.ValueOf(str.ToString(valRef)))
	}
}

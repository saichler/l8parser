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
	"strings"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
)

// CTableToInstances converts a CTable into individual model instances, one per row.
// Each row becomes a separate proto instance with fields mapped from table columns.
// ClusterName is set from the job's host ID (stored in workspace under the
// TargetId key for historical reasons — see Parser.go), and Key is built from
// the key columns (e.g., "namespace/name").
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
	// Plain string conversion (no type prefix). The Key field is just a flat
	// identifier (e.g. "namespace/name"). Earlier versions left TypesPrefix=true
	// here, which produced "{24}namespace/{24}name" — visible in the UI and
	// breaking any consumer that expected a clean composite key.
	toString := strings2.New()
	toString.TypesPrefix = false

	fmt.Printf("[CTABLE->INSTANCES] elemType=%s rows=%d cols=%d keyCols=%v targetId=%q\n",
		elemType.Name(), len(table.Rows), len(table.Columns), keyColumns, targetId)

	instances := make([]interface{}, 0, len(table.Rows))
	for rowIdx, row := range table.Rows {
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
			setFieldValue(field, val, resources)
		}

		clusterField := instElem.FieldByName("ClusterName")
		if clusterField.IsValid() && clusterField.CanSet() {
			clusterField.Set(reflect.ValueOf(targetId))
		}
		keyField := instElem.FieldByName("Key")
		if keyField.IsValid() && keyField.CanSet() {
			keyField.Set(reflect.ValueOf(keyBuilder.String()))
		}

		if rowIdx < 3 {
			fmt.Printf("[CTABLE->INSTANCE] row=%d cluster=%q key=%q\n",
				rowIdx, targetId, keyBuilder.String())
		}

		instances = append(instances, inst.Interface())
	}

	fmt.Printf("[CTABLE->INSTANCES-DONE] elemType=%s instances=%d\n",
		elemType.Name(), len(instances))

	workSpace[Instances] = instances
	return nil
}

// findFieldByJsonName resolves a CTable column name (already lowered by the
// collector's getAttributeNameFromColumn — dashes, spaces and parens stripped,
// but underscores preserved) to a field on the target proto struct.
//
// Two facts make this trickier than it looks:
//
//  1. The collector's column normalization strips dashes (e.g. INTERNAL-IP →
//     "internalip") but keeps underscores (CONTAINERS_JSON → "containers_json").
//     So a single column name can be either the camelCase proto JSON name
//     ("internalip") OR the snake_case Go field tag ("containers_json").
//  2. protoc emits TWO json identifiers per field — the camelCase one inside
//     the `protobuf:"...,json=foo,proto3"` tag, and the snake_case one in the
//     plain Go `json:"foo,omitempty"` tag — and it OMITS the `json=` clause
//     entirely for single-word fields. Either tag alone is insufficient:
//     - protobuf json= alone misses single-word fields (name, age, roles…)
//     - go json: alone misses snake_case (internal_ip ≠ internalip)
//
// The fix: try BOTH candidate names per field, plus a third comparison that
// strips underscores from both sides (so "internal_ip" matches "internalip").
// Order doesn't matter — EqualFold is symmetric.
func findFieldByJsonName(v reflect.Value, jsonName string) reflect.Value {
	t := v.Type()
	target := stripUnderscores(jsonName)
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}

		// Candidate 1: name from the protobuf tag's `json=` clause (camelCase
		// when present, e.g. internal_ip → internalIp).
		if pbName := protobufJSONName(sf.Tag.Get("protobuf")); pbName != "" {
			if strings.EqualFold(pbName, jsonName) || strings.EqualFold(stripUnderscores(pbName), target) {
				return v.Field(i)
			}
		}

		// Candidate 2: name from the Go `json:` struct tag (snake_case for
		// underscore fields, plain word otherwise — covers single-word fields
		// the protobuf clause omits).
		if goJSON := sf.Tag.Get("json"); goJSON != "" {
			name := goJSON
			if comma := indexOf(name, ","); comma != -1 {
				name = name[:comma]
			}
			if name != "" && (strings.EqualFold(name, jsonName) || strings.EqualFold(stripUnderscores(name), target)) {
				return v.Field(i)
			}
		}
	}
	return reflect.Value{}
}

// protobufJSONName extracts the value of `json=...` from a protoc-emitted
// protobuf struct tag. Returns "" when the clause is absent (which happens
// for single-word fields whose JSON name equals the proto field name).
func protobufJSONName(protobufTag string) string {
	const marker = "json="
	idx := indexOfSubstr(protobufTag, marker)
	if idx == -1 {
		return ""
	}
	rest := protobufTag[idx+len(marker):]
	if comma := indexOf(rest, ","); comma != -1 {
		rest = rest[:comma]
	}
	return rest
}

func stripUnderscores(s string) string {
	if !containsByte(s, '_') {
		return s
	}
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != '_' {
			out = append(out, s[i])
		}
	}
	return string(out)
}

func containsByte(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// indexOfSubstr is a separate name to avoid a self-shadow with indexOf.
// (Both call sites currently want the same behavior; the duplication is
// kept for symmetry / readability.)
func indexOfSubstr(s, substr string) int { return indexOf(s, substr) }

func setFieldValue(field reflect.Value, val interface{}, resources ifs.IResources) {
	valRef := reflect.ValueOf(val)

	// String → typed-int32 enum lookup MUST come before AssignableTo /
	// ConvertibleTo. Go allows `string → int32` via Convert (it produces a
	// rune-valued int32), which would silently store garbage in the field.
	// The enum registry resolves "Running" → 1 etc. for registered enum
	// types and short-circuits before that path is taken. Unregistered
	// types (or unmapped raw values) fall through to the existing logic.
	if valRef.Kind() == reflect.String {
		if v, ok := enumValueForField(field, valRef.String()); ok {
			field.SetInt(int64(v))
			return
		}
	}

	// String → struct via registered STRING serializer (e.g. "1/2" → *K8SReadyState).
	// Mirrors Setter.go's serializer lookup so per-row CTable parsing handles
	// struct fields the same way the Set rule does for property-path setters.
	if valRef.Kind() == reflect.String && resources != nil {
		ftype := field.Type()
		if ftype.Kind() == reflect.Ptr {
			ftype = ftype.Elem()
		}
		if ftype.Kind() == reflect.Struct && ftype.Name() != "" {
			if info, err := resources.Registry().Info(ftype.Name()); err == nil && info != nil {
				if ser := info.Serializer(ifs.STRING); ser != nil {
					if inst, sErr := ser.Unmarshal([]byte(valRef.String()), resources); sErr == nil && inst != nil {
						v := reflect.ValueOf(inst)
						if v.Kind() == reflect.Ptr && field.Kind() != reflect.Ptr {
							v = v.Elem()
						}
						if v.Type().AssignableTo(field.Type()) {
							field.Set(v)
							return
						}
					}
				}
			}
		}
	}

	if valRef.Type().AssignableTo(field.Type()) {
		field.Set(valRef)
		return
	}
	if valRef.Type().ConvertibleTo(field.Type()) {
		// Skip the rune-producing string→int conversion explicitly. If the
		// target was a typed enum, we already tried the registry above; if
		// the target is a raw int32 the parser had no way to interpret a
		// string anyway, and storing a rune would mask the missing mapping.
		if !(valRef.Kind() == reflect.String && field.Kind() == reflect.Int32) {
			field.Set(valRef.Convert(field.Type()))
			return
		}
	}
	if field.Kind() == reflect.String {
		// Fallback: stringify the source value without a type prefix. Including
		// the prefix (e.g. "{2}123") would corrupt user-visible string fields.
		str := strings2.New()
		str.TypesPrefix = false
		field.Set(reflect.ValueOf(str.ToString(valRef)))
	}
}


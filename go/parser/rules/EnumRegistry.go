/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rules

import (
	"reflect"
	"sync"
)

// enumRegistryMu guards enumRegistry. RegisterEnum is typically called once
// per enum at boot, but reads happen on every parser job — a RWMutex keeps
// the hot path lock-free relative to other readers.
var enumRegistryMu sync.RWMutex

// enumRegistry maps a Go enum type name (e.g. "K8SPodStatus") to a
// string→int32 lookup table for converting K8s/external string values
// (e.g. "Running") to the numeric protobuf enum values.
//
// Callers populate the registry at parser activation. The registry is
// purely additive: a missing type name or a missing key falls back to the
// existing setFieldValue behavior, so unrelated parsers are unaffected.
var enumRegistry = map[string]map[string]int32{}

// RegisterEnum registers a string→int32 lookup table for a Go enum type
// identified by its short type name (the value of reflect.Type.Name() —
// not the fully-qualified path). Registration is safe to call multiple
// times; the latest call wins.
//
// Example:
//
//	rules.RegisterEnum("K8SPodStatus", map[string]int32{
//	    "Running":   1,
//	    "Pending":   2,
//	    "Succeeded": 3,
//	})
func RegisterEnum(typeName string, values map[string]int32) {
	if typeName == "" || values == nil {
		return
	}
	enumRegistryMu.Lock()
	defer enumRegistryMu.Unlock()
	// Defensive copy so the caller can mutate its map without affecting
	// the registry.
	copyMap := make(map[string]int32, len(values))
	for k, v := range values {
		copyMap[k] = v
	}
	enumRegistry[typeName] = copyMap
}

// LookupEnum returns the int32 value for a (typeName, key) pair and whether
// the lookup succeeded. Both an unregistered type and a missing key produce
// (0, false). This function is safe for concurrent use.
func LookupEnum(typeName, key string) (int32, bool) {
	enumRegistryMu.RLock()
	defer enumRegistryMu.RUnlock()
	values, ok := enumRegistry[typeName]
	if !ok {
		return 0, false
	}
	v, ok := values[key]
	return v, ok
}

// enumValueForField resolves a string raw value into the int32 enum value
// expected by a target field. It returns (value, true) on success or
// (0, false) when the field's type is not a registered enum or the raw
// string is not in the registered map.
//
// Only int-kinded fields with a non-empty type name (i.e. typed enums like
// `type K8SPodStatus int32`, not bare `int32`) are eligible — a bare int32
// field never matches because no caller would register a map under "int32".
func enumValueForField(field reflect.Value, raw string) (int32, bool) {
	if field.Kind() != reflect.Int32 {
		return 0, false
	}
	typeName := field.Type().Name()
	if typeName == "" || typeName == "int32" {
		return 0, false
	}
	return LookupEnum(typeName, raw)
}

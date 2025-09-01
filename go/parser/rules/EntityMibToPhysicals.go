package rules

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	types2 "github.com/saichler/probler/go/types"
)

type EntityMibToPhysicals struct{}

func (this *EntityMibToPhysicals) Name() string {
	return "EntityMibToPhysicals"
}

func (this *EntityMibToPhysicals) ParamNames() []string {
	return []string{""}
}

// Entity MIB OID structure:
// .1.3.6.1.2.1.47.1.1.1.1.1 - entPhysicalIndex (not in table - used as row key)
// .1.3.6.1.2.1.47.1.1.1.1.2 - entPhysicalDescr
// .1.3.6.1.2.1.47.1.1.1.1.3 - entPhysicalVendorType
// .1.3.6.1.2.1.47.1.1.1.1.4 - entPhysicalContainedIn
// .1.3.6.1.2.1.47.1.1.1.1.5 - entPhysicalClass
// .1.3.6.1.2.1.47.1.1.1.1.6 - entPhysicalParentRelPos
// .1.3.6.1.2.1.47.1.1.1.1.7 - entPhysicalName
// .1.3.6.1.2.1.47.1.1.1.1.8 - entPhysicalHardwareRev
// .1.3.6.1.2.1.47.1.1.1.1.9 - entPhysicalFirmwareRev
// .1.3.6.1.2.1.47.1.1.1.1.10 - entPhysicalSoftwareRev
// .1.3.6.1.2.1.47.1.1.1.1.11 - entPhysicalSerialNum
// .1.3.6.1.2.1.47.1.1.1.1.12 - entPhysicalMfgName
// .1.3.6.1.2.1.47.1.1.1.1.13 - entPhysicalModelName
// .1.3.6.1.2.1.47.1.1.1.1.14 - entPhysicalAlias
// .1.3.6.1.2.1.47.1.1.1.1.15 - entPhysicalAssetID
// .1.3.6.1.2.1.47.1.1.1.1.16 - entPhysicalIsFRU

// Entity Physical Class enum values:
const (
	EntPhysicalClassOther         = 1
	EntPhysicalClassUnknown       = 2
	EntPhysicalClassChassis       = 3
	EntPhysicalClassBackplane     = 4
	EntPhysicalClassContainer     = 5
	EntPhysicalClassPowerSupply   = 6
	EntPhysicalClassFan           = 7
	EntPhysicalClassSensor        = 8
	EntPhysicalClassModule        = 9
	EntPhysicalClassPort          = 10
	EntPhysicalClassStack         = 11
	EntPhysicalClassCpu           = 12
)

func (this *EntityMibToPhysicals) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*types.Parameter, any interface{}, pollWhat string) error {
	// Get the input CTable from workspace
	input := workSpace[Input]
	if input == nil {
		return errors.New("No input data found in workspace")
	}

	// Try to get CTable from input
	table, ok := input.(*types.CTable)
	if !ok {
		return errors.New("Input is not a CTable: " + fmt.Sprintf("%T", input))
	}

	fmt.Printf("DEBUG EntityMibToPhysicals: Found CTable with %d columns and %d rows\n", len(table.Columns), len(table.Rows))
	
	// Print column information for debugging
	for colKey, colName := range table.Columns {
		fmt.Printf("DEBUG EntityMibToPhysicals: Column %d = '%s'\n", colKey, colName)
	}
	
	// Print first row data for debugging
	if len(table.Rows) > 0 {
		for rowKey, row := range table.Rows {
			if rowKey > 2 { // Only show first 3 rows
				break
			}
			fmt.Printf("DEBUG EntityMibToPhysicals: Row %d data:\n", rowKey)
			for colKey, data := range row.Data {
				val := getEntityValue(data, resources)
				fmt.Printf("  Column %d: %v (type %T)\n", colKey, val, val)
			}
		}
	}

	// Get the NetworkDevice
	networkDevice, ok := any.(*types2.NetworkDevice)
	if !ok {
		return errors.New("Target object is not a NetworkDevice")
	}

	// Ensure physicals map exists
	if networkDevice.Physicals == nil {
		networkDevice.Physicals = make(map[string]*types2.Physical)
	}

	// Create the main physical component
	physicalKey := "physical-0"
	physical, exists := networkDevice.Physicals[physicalKey]
	if !exists {
		physical = &types2.Physical{
			Id: physicalKey,
		}
		networkDevice.Physicals[physicalKey] = physical
	}

	// Initialize component maps
	portMap := make(map[string]*types2.Port)

	// Process each row in the Entity MIB table
	for rowKey, row := range table.Rows {
		entityIndex := fmt.Sprintf("%d", rowKey)
		
		// Determine which Entity MIB column this CTable represents
		var entityColumn int
		if len(table.Columns) == 1 {
			// Single column case - get the column number from column name
			for colKey, colName := range table.Columns {
				if colKey == 0 {
					if colNum, err := strconv.Atoi(colName); err == nil {
						entityColumn = colNum
					}
					break
				}
			}
		}
		
		// For now, we'll focus on the key columns that define entity structure
		// Skip processing if this isn't one of the key columns we need
		if entityColumn != 2 && entityColumn != 5 && entityColumn != 7 && entityColumn != 11 && entityColumn != 13 { // descr, class, name, serial, model
			continue
		}
		
		// Get the value for this column
		data, ok := row.Data[0] // Data is always in column 0 for single-column CTable
		if !ok {
			continue
		}
		value := getEntityValue(data, resources)
		if value == nil {
			continue
		}
		
		// Store entity data by index and column for later processing
		// For now, create a basic port for each entity index
		if entityColumn == 5 { // entPhysicalClass
			entityClassInt := 0
			if classStr := fmt.Sprintf("%v", value); classStr != "" {
				if val, err := strconv.Atoi(classStr); err == nil {
					entityClassInt = val
				}
			}
			
			fmt.Printf("DEBUG EntityMibToPhysicals: Entity %s has class %d\n", entityIndex, entityClassInt)
			if entityClassInt == EntPhysicalClassPort {
				port := &types2.Port{
					Id: entityIndex,
				}
				portMap[entityIndex] = port
				fmt.Printf("DEBUG EntityMibToPhysicals: Found port entity %s with class %d\n", entityIndex, entityClassInt)
			}
		}
		
	}

	// Convert maps to slices and assign to physical component
	if len(portMap) > 0 {
		physical.Ports = make([]*types2.Port, 0, len(portMap))
		for _, port := range portMap {
			physical.Ports = append(physical.Ports, port)
		}
		fmt.Printf("DEBUG EntityMibToPhysicals: Created %d ports\n", len(portMap))
	}

	return nil
}

func getEntityValue(data []byte, resources ifs.IResources) interface{} {
	if len(data) == 0 {
		return nil
	}
	obj := object.NewDecode(data, 0, resources.Registry())
	val, _ := obj.Get()
	return val
}

func getEntityStringValue(data []byte, resources ifs.IResources) string {
	val := getEntityValue(data, resources)
	if val == nil {
		return ""
	}
	
	// Handle byte array to string conversion
	if byteArray, ok := val.([]uint8); ok {
		return strings.TrimSpace(string(byteArray))
	}
	
	return strings.TrimSpace(fmt.Sprintf("%v", val))
}
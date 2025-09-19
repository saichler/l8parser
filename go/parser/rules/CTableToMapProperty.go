package rules

import (
	"errors"
	"reflect"
	"strings"

	"github.com/saichler/l8pollaris/go/types/l8poll"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
	"github.com/saichler/l8reflect/go/reflect/properties"
)

type CTableToMapProperty struct{}

func (this *CTableToMapProperty) Name() string {
	return "CTableToMapProperty"
}

func (this *CTableToMapProperty) ParamNames() []string {
	return []string{""}
}

func (this *CTableToMapProperty) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*l8poll.L8P_Parameter, any interface{}, pollWhat string) error {
	table, ok := workSpace[Output].(*l8poll.CTable)
	if !ok {
		return errors.New("Workspace had an invalid output object")
	}

	keyColumn, e := getIntInput(workSpace, KeyColumn)
	if e != nil {
		return e
	}

	propertyId := workSpace[PropertyId].(string)
	toString := strings2.New()
	toString.TypesPrefix = true

	for _, row := range table.Rows {
		pid := strings2.New(propertyId)
		for i := 0; i < len(table.Columns); i++ {
			if i == 0 {
				val := getValue(row.Data[int32(keyColumn)], resources)
				if val == nil {
					break
				}
				pid.Add("<")
				pid.Add(toString.ToString(reflect.ValueOf(val)))
				pid.Add(">.")
			}

			key := strings2.New(pid.String())
			attrName := getAttributeNameFromColumn(table.Columns[int32(i)])
			key.Add(attrName)

			prop, err := properties.PropertyOf(key.String(), resources)
			if err != nil {
				resources.Logger().Error(err.Error())
				continue
			}

			val := getValue(row.Data[int32(i)], resources)
			_, _, err = prop.Set(any, val)
			if err != nil {
				resources.Logger().Error(err.Error())
				continue
			}
		}
	}
	return nil
}

func getAttributeNameFromColumn(value interface{}) string {
	colName := strings.TrimSpace(value.(string))
	colName = strings.ToLower(colName)
	colName = removeChar(colName, "-")
	colName = removeChar(colName, " ")
	return colName
}

func removeChar(colName, c string) string {
	index := strings.LastIndex(colName, c)
	if index == -1 {
		return colName
	}
	return strings2.New(colName[0:index], colName[index+1:]).String()
}

func getValue(data []byte, resources ifs.IResources) interface{} {
	if len(data) == 0 {
		return nil
	}
	obj := object.NewDecode(data, 0, resources.Registry())
	val, _ := obj.Get()
	return val
}

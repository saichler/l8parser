package rules

import (
	"strings"

	"github.com/saichler/l8pollaris/go/types"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
)

type StringToCTable struct{}

func (this *StringToCTable) Name() string {
	return "StringToCTable"
}

func (this *StringToCTable) ParamNames() []string {
	return []string{"columns", "keycolumn"}
}

func (this *StringToCTable) Parse(resources ifs.IResources, workSpace map[string]interface{}, params map[string]*types.Parameter, any interface{}) error {
	input, ok := workSpace[Input].(string)
	if !ok {
		return nil
	}
	colmns, err := getIntInput(workSpace, Columns)
	if err != nil {
		return err
	}

	lines := strings.Split(input, "\n")
	table := &types.CTable{}
	table.Rows = make(map[int32]*types.CRow)
	for i, line := range lines {
		if table.Columns == nil {
			table.Columns = getColumns(line, colmns)
			if len(table.Columns) != colmns {
				return resources.Logger().Error("Number of columns mismatch, expected:", colmns, ", actual:", len(table.Columns))
			}
			continue
		}
		row := &types.CRow{}
		row.Data = getValues(line, table.Columns)
		table.Rows[int32(i)] = row
	}
	workSpace[Output] = table
	return nil
}

func getValues(line string, columns map[int32]string) map[int32][]byte {
	line = strings.TrimSpace(line)
	result := make(map[int32][]byte, 0)
	begin := 0
	for i := 0; i < len(columns); i++ {
		col := columns[int32(i)]
		if begin+len(col) > len(line) {
			result[int32(i)] = []byte{}
		} else {
			value := strings.TrimSpace(line[begin : begin+len(col)])
			obj := object.NewEncode()
			obj.Add(value)
			result[int32(i)] = obj.Data()
			begin = begin + len(col)
		}
	}
	return result
}

func getColumns(line string, expected int) map[int32]string {
	columns := make([]string, 0)
	begin := 0
	open := false

	for i := 0; i < len(line); i++ {
		if line[i] == ' ' && !open {
			open = true
		} else if open && line[i] != ' ' {
			open = false
			col := line[begin:i]
			columns = append(columns, col)
			begin = i
		}
	}

	if begin < len(line)-1 {
		col := line[begin:len(line)]
		columns = append(columns, col)
	}
	if len(columns) > expected {
		newColumns := make([]string, 0)
		for i := 0; i < len(columns); i++ {
			l := len(columns[i])
			t := len(strings.TrimSpace(columns[i]))
			if l == t+1 {
				newColumns = append(newColumns, strings2.New(columns[i], columns[i+1]).String())
				i++
			} else {
				newColumns = append(newColumns, columns[i])
			}
		}
		columns = newColumns
	}
	result := make(map[int32]string)
	for i, col := range columns {
		result[int32(i)] = col
	}
	return result
}

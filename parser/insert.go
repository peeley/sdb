package parser

import (
	"fmt"
	"sdb/types"
)

func ParseInsertStatement(input string) (types.Statement, error){
	trimmed, ok := HasPrefix(input, "insert into")
	if !ok {
		return nil, nil
	}

	tableName := ParseIdentifier(trimmed)
	trimmed, _ = HasPrefix(trimmed, tableName)

	trimmed, ok = HasPrefix(trimmed, "values(")
	if !ok {
		return nil, fmt.Errorf("Expected 'values' after table to insert into")
	}

	var valueList []types.Value
	var err error
	valueList, trimmed, err = ParseValueList(trimmed)
	if err != nil {
		return nil, err
	}

	trimmed, ok = HasPrefix(trimmed, ")")
	if !ok {
		return nil, fmt.Errorf("Expected list of values to end in ')'")
	}

	statement := types.InsertStatement{
		TableName: tableName,
		Values: valueList,
	}

	return statement, nil
}

package parser

import (
	"fmt"
	"sdb/db"
	"sdb/statements"
	"sdb/utils"
)

func ParseInsertStatement(input string) (db.Executable, error) {
	trimmed, ok := utils.HasPrefix(input, "insert into")
	if !ok {
		return nil, nil
	}

	tableName := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, tableName)

	trimmed, _ = utils.HasPrefix(trimmed, "values")
	trimmed, ok = utils.HasPrefix(trimmed, "(")
	if !ok {
		return nil, fmt.Errorf("Expected 'values' after table to insert into")
	}

	var valueList []db.Value
	var err error
	valueList, trimmed, err = utils.ParseValueList(trimmed)
	if err != nil {
		return nil, err
	}

	_, ok = utils.HasPrefix(trimmed, ")")
	if !ok {
		return nil, fmt.Errorf("Expected list of values to end in ')'")
	}

	statement := statements.InsertStatement{
		TableName: tableName,
		Values:    valueList,
	}

	return statement, nil
}

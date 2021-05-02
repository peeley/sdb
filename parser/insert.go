package parser

import (
	"fmt"
	"sdb/statements"
	"sdb/db"
	"sdb/utils"
)

func ParseInsertStatement(input string) (statements.Executable, error){
	trimmed, ok := utils.HasPrefix(input, "insert into")
	if !ok {
		return nil, nil
	}

	tableName := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, tableName)

	trimmed, ok = utils.HasPrefix(trimmed, "values(")
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
		Values: valueList,
	}

	return statement, nil
}

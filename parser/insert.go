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

	trimmed, ok = HasPrefix(trimmed, "values(")
	if !ok {
		return nil, fmt.Errorf("Expected 'values' after table to insert into")
	}

	return nil, nil
}

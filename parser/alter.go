// Noah Snelson
// February 25, 2021
// sdb/parser/alter.go
//
// Contains function to parse `ALTER` queries.

package parser

import (
	"fmt"
	"sdb/types"
	"strings"
)

// Parses `ALTER TABLE` input.
func ParseAlterStatement(input string) (types.Statement, error) {
	trimmed, ok := HasPrefix(input, "alter table")
	if !ok {
		return nil, nil
	}

	tableName := ParseIdentifier(trimmed)
	trimmed = strings.TrimPrefix(trimmed, tableName)
	trimmed = strings.TrimSpace(trimmed)

	trimmed, ok = HasPrefix(input, "add")
	if !ok {
		return nil, fmt.Errorf(
			"Expected `ADD` after table name in `ALTER` statement.",
		)
	}

	newColName := ParseIdentifier(trimmed)
	if newColName == "" {
		return nil, fmt.Errorf(
			"Missing column name after `ADD` in `ALTER` statement.",
		)
	}
	trimmed = strings.TrimPrefix(trimmed, newColName)
	trimmed = strings.TrimSpace(trimmed)

	typeName, err := ParseType(trimmed)
	if err != nil {
		return nil, err
	}

	alterStatement := types.AlterStatement{
		TableName: tableName,
		ColumnName: newColName,
		ColumnType: typeName,
	}

	return alterStatement, nil
}

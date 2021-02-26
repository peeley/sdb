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
	prefix := "alter table"
	if !strings.HasPrefix(input, prefix) {
		return nil, nil
	}

	trimmed := strings.TrimPrefix(input, prefix)
	trimmed = strings.TrimSpace(trimmed)

	tableName := ParseIdentifier(trimmed)
	trimmed = strings.TrimPrefix(trimmed, tableName)
	trimmed = strings.TrimSpace(trimmed)

	addString := "add"
	if !strings.HasPrefix(trimmed, addString) {
		return nil, fmt.Errorf(
			"Expected `ADD` after table name in `ALTER` statement.",
		)
	}
	trimmed = strings.TrimPrefix(trimmed, addString)
	trimmed = strings.TrimSpace(trimmed)

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

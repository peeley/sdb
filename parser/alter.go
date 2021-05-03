// Noah Snelson
// February 25, 2021
// sdb/parser/alter.go
//
// Contains function to parse `ALTER` queries.

package parser

import (
	"fmt"
	"sdb/db"
	"sdb/statements"
	"sdb/utils"
	"strings"
)

// Parses `ALTER TABLE` input.
func ParseAlterStatement(input string) (db.Executable, error) {
	trimmed, ok := utils.HasPrefix(input, "alter table")
	if !ok {
		return nil, nil
	}

	tableName := utils.ParseIdentifier(trimmed)
	trimmed = strings.TrimPrefix(trimmed, tableName)
	trimmed = strings.TrimSpace(trimmed)

	trimmed, ok = utils.HasPrefix(trimmed, "add")
	if !ok {
		return nil, fmt.Errorf(
			"Expected `ADD` after table name in `ALTER` statement.",
		)
	}

	newColName := utils.ParseIdentifier(trimmed)
	if newColName == "" {
		return nil, fmt.Errorf(
			"Missing column name after `ADD` in `ALTER` statement.",
		)
	}
	trimmed = strings.TrimPrefix(trimmed, newColName)
	trimmed = strings.TrimSpace(trimmed)

	typeName, err := utils.ParseType(trimmed)
	if err != nil {
		return nil, err
	}

	alterStatement := statements.AlterStatement{
		TableName:  tableName,
		ColumnName: newColName,
		ColumnType: typeName,
	}

	return alterStatement, nil
}

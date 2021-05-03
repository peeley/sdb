// Noah Snelson
// February 25, 2021
// sdb/parser/select.go
//
// Contains functions for parsing `SELECT` queries.

package parser

import (
	"fmt"
	"sdb/db"
	"sdb/statements"
	"sdb/utils"
)

// Parses `SELECT` input.
func ParseSelectStatement(input string) (db.Executable, error) {
	trimmed, ok := utils.HasPrefix(input, "select")
	if !ok {
		return nil, nil
	}

	colNames := []string{}
	for {
		ident := utils.ParseIdentifier(trimmed)
		trimmed, _ = utils.HasPrefix(trimmed, ident)

		colNames = append(colNames, ident)
		if ident == "*" {
			break
		}

		trimmed, ok = utils.HasPrefix(trimmed, ",")
		if !ok {
			break
		}
	}

	trimmed, ok = utils.HasPrefix(trimmed, "from")
	if !ok {
		return nil, fmt.Errorf("Expected `FROM` after columns in `SELECT`.")
	}

	tableName := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, tableName)

	where, _ := ParseWhereClause(trimmed)

	var joinClause *statements.JoinClause
	if where == nil {
		joinClause, _ = ParseJoinClause(trimmed, tableName)
	}

	statement := statements.SelectStatement{
		TableName:   tableName,
		ColumnNames: colNames,
		WhereClause: where,
		JoinClause:  joinClause,
	}

	return statement, nil
}

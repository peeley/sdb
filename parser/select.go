// Noah Snelson
// February 25, 2021
// sdb/parser/select.go
//
// Contains functions for parsing `SELECT` queries.

package parser

import (
	"fmt"
	"sdb/types"
	"strings"
)

// Parses `SELECT` input.
func ParseSelectStatement(input string) (types.Statement, error) {
	trimmed, ok := HasPrefix(input, "select")
	if !ok {
		return nil, nil
	}

	// FIXME add support to select cols by name
	if trimmed[0] != '*' {
		return nil, fmt.Errorf("Expected `SELECT` followed by `*`")
	}

	trimmed = strings.TrimPrefix(trimmed, "*")
	trimmed = strings.TrimSpace(trimmed)

	trimmed, ok = HasPrefix(trimmed, "from")
	if !ok {
		return nil, fmt.Errorf("Expected `FROM` after columns in `SELECT`.")
	}

	tableName := ParseIdentifier(trimmed)

	statement := types.SelectStatement {
		TableName: tableName,
		Columns: "*",
	}

	return statement, nil
}

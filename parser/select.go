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

	prefix := "select"
	if !strings.HasPrefix(input, prefix) {
		return nil, nil
	}

	trimmed := strings.TrimPrefix(input, prefix)
	trimmed = strings.TrimSpace(trimmed)

	// Currently, `SELECT` only supports querying every column at once via the
	// `*` shortcut.
	if trimmed[0] != '*' {
		return nil, fmt.Errorf("Expected `SELECT` followed by `*`")
	}
	trimmed = strings.TrimPrefix(trimmed, "*")
	trimmed = strings.TrimSpace(trimmed)

	if !strings.HasPrefix(trimmed, "from") {
		return nil, fmt.Errorf("Expected `FROM` after columns in `SELECT`.")
	}
	trimmed = strings.TrimPrefix(trimmed, "from")
	trimmed = strings.TrimSpace(trimmed)

	tableName := ParseIdentifier(trimmed)

	statement := types.SelectStatement {
		TableName: tableName,
		Columns: "*",
	}

	return statement, nil
}

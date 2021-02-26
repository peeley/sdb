// Noah Snelson
// February 25, 2021
// sdb/parser/use.go
//
// Contains functions for parsing `USE` queries.

package parser

import (
	"sdb/types"
	"strings"
)

// Parses `USE <db_name>;` input.
func ParseUseDBStatement(input string) (types.Statement, error) {
	prefix := "use"
	if len(input) < len(prefix) || !strings.HasPrefix(input, prefix) {
		return nil, nil
	}

	trimmed := strings.TrimPrefix(input, prefix)
	trimmed = strings.TrimSpace(trimmed)
	ident := ParseIdentifier(trimmed)

	dropDB := types.UseDBStatement{
		DBName: ident,
	}

	return dropDB, nil
}

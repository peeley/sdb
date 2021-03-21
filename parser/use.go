// Noah Snelson
// February 25, 2021
// sdb/parser/use.go
//
// Contains functions for parsing `USE` queries.

package parser

import (
	"sdb/types"
)

// Parses `USE <db_name>;` input.
func ParseUseDBStatement(input string) (types.Statement, error) {
	trimmed, ok := HasPrefix(input, "use")
	if !ok {
		return nil, nil
	}

	ident := ParseIdentifier(trimmed)

	dropDB := types.UseDBStatement{
		DBName: ident,
	}

	return dropDB, nil
}

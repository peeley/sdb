// Noah Snelson
// February 25, 2021
// sdb/parser/use.go
//
// Contains functions for parsing `USE` queries.

package parser

import (
	"sdb/db"
	"sdb/statements"
	"sdb/utils"
)

// Parses `USE <db_name>;` input.
func ParseUseDBStatement(input string) (db.Executable, error) {
	trimmed, ok := utils.HasPrefix(input, "use")
	if !ok {
		return nil, nil
	}

	ident := utils.ParseIdentifier(trimmed)

	dropDB := statements.UseDBStatement{
		DBName: ident,
	}

	return dropDB, nil
}

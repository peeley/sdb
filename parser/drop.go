// Noah Snelson
// February 25, 2021
// sdb/parser/drop.go
//
// Contains functions for parsing `DROP TABLE` and `DROP DATABASE` queries.

package parser

import (
	"sdb/types"
	"sdb/utils"
)

// Parses `DROP DATABASE <table_name>;` input.
func ParseDropDBStatement(input string) (types.Statement, error) {
	trimmed, ok := utils.HasPrefix(input, "drop database")
	if !ok {
		return nil, nil
	}

	ident := utils.ParseIdentifier(trimmed)

	dropDB := types.DropDBStatement{
		DBName: ident,
	}

	return dropDB, nil
}

// Parses `DROP TABLE <table_name>;` input.
func ParseDropTableStatement(input string) (types.Statement, error) {
	trimmed, ok := utils.HasPrefix(input, "drop table")
	if !ok {
		return nil, nil
	}

	ident := utils.ParseIdentifier(trimmed)

	dropDB := types.DropTableStatement{
		TableName: ident,
	}

	return dropDB, nil
}

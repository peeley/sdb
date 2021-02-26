// Noah Snelson
// February 25, 2021
// sdb/parser/drop.go
//
// Contains functions for parsing `DROP TABLE` and `DROP DATABASE` queries.

package parser

import (
	"sdb/types"
	"strings"
)

// Parses `DROP DATABASE <table_name>;` input.
func ParseDropDBStatement(input string) (types.Statement, error) {
	prefix := "drop database"
	if len(input) < len(prefix) || !strings.HasPrefix(input, prefix) {
		return nil, nil
	}

	trimmed := strings.TrimPrefix(input, prefix)
	trimmed = strings.TrimSpace(trimmed)
	ident := ParseIdentifier(trimmed)

	dropDB := types.DropDBStatement{
		DBName: ident,
	}

	return dropDB, nil
}

// Parses `DROP TABLE <table_name>;` input.
func ParseDropTableStatement(input string) (types.Statement, error) {
	prefix := "drop table"
	if len(input) < len(prefix) || !strings.HasPrefix(input, prefix) {
		return nil, nil
	}

	trimmed := strings.TrimPrefix(input, prefix)
	trimmed = strings.TrimSpace(trimmed)
	ident := ParseIdentifier(trimmed)

	dropDB := types.DropTableStatement{
		TableName: ident,
	}

	return dropDB, nil
}

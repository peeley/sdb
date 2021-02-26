package parser

import (
	"sdb/types"
	"strings"
)

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

func ParseDropTableStatement(input string) (types.Statement, error) {
	prefix := "drop table"
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

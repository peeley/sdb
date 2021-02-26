package parser

import (
	"sdb/types"
	"strings"
)

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

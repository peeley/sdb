package parser

import (
	"sdb/types"
	"sdb/utils"
)

func ParseWhereClause(input string) (*types.WhereClause, error) {
	trimmed, ok := utils.HasPrefix(input, "where")
	if !ok {
		return nil, nil
	}

	colName := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, colName)

	var comparison string
	if (trimmed[0] == '=' || trimmed[0] == '<' || trimmed[0] == '>') {
		comparison = string(trimmed[0])
	} else if( trimmed[:2] == "!=" || trimmed[:2] == "<=" || trimmed[:2] == ">=" ) {
		comparison = string(trimmed[:2])
	}

	trimmed, _ = utils.HasPrefix(trimmed, comparison)

	value, _ := utils.ParseValue(trimmed)

	where := types.WhereClause {
		colName,
		comparison,
		value,
	}

	return &where, nil
}

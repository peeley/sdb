package parser

import (
	"sdb/statements"
	"sdb/utils"
)

func ParseWhereClause(input string) (*statements.WhereClause, error) {
	trimmed, ok := utils.HasPrefix(input, "where")
	if !ok {
		return nil, nil
	}

	colName := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, colName)

	var comparison string
	if trimmed[0] == '=' || trimmed[0] == '<' || trimmed[0] == '>' {
		comparison = string(trimmed[0])
	} else if trimmed[:2] == "!=" || trimmed[:2] == "<=" || trimmed[:2] == ">=" {
		comparison = string(trimmed[:2])
	}

	trimmed, _ = utils.HasPrefix(trimmed, comparison)

	value, _ := utils.ParseValue(trimmed)

	where := statements.WhereClause{
		ColName:         colName,
		Comparison:      comparison,
		ComparisonValue: value,
	}

	return &where, nil
}

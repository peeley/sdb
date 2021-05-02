package parser

import (
	"sdb/utils"
	"sdb/statements"
)

func ParseUpdateStatement(input string) (statements.Executable, error){
	trimmed, ok := utils.HasPrefix(input, "update")
	if !ok {
		return nil, nil
	}

	tableName := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, tableName)
	trimmed, _ = utils.HasPrefix(trimmed, "set")

	colName := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, colName)

	trimmed, _ = utils.HasPrefix(trimmed, "=")

	value, _ := utils.ParseValue(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, value.ToString())

	where, _ := ParseWhereClause(trimmed)

	update := statements.UpdateStatement{
		TableName: tableName,
		UpdatedCol: colName,
		UpdatedValue: value,
		WhereClause: where,
	}

	return update, nil
}

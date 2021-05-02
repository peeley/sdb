package parser

import (
	"errors"
	"sdb/statements"
	"sdb/utils"
)

func ParseDeleteStatement(input string) (statements.Executable, error) {
	trimmed, ok := utils.HasPrefix(input, "delete from")
	if !ok {
		return nil, nil
	}

	tableName := utils.ParseIdentifier(trimmed)
	trimmed, ok = utils.HasPrefix(trimmed, tableName)
	if len(tableName) == 0 || !ok {
		return nil, errors.New("!Expected table name after DELETE FROM.")
	}

	where, err := ParseWhereClause(trimmed)
	if err != nil {
		return nil, err
	}

	delete := statements.DeleteStatment{
		TableName:   tableName,
		WhereClause: where,
	}

	return delete, nil
}

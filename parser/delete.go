package parser

import (
	"errors"
	"sdb/types"
	"sdb/utils"
)

func ParseDeleteStatement(input string) (types.Statement, error) {
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

	delete := types.DeleteStatment {
		TableName: tableName,
		WhereClause: where,
	}

	return delete, nil
}

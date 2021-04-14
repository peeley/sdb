package parser

import (
	"sdb/types"
	"sdb/utils"
	"strings"
)

func ParseJoinClause(input, leftTableName string) (*types.JoinClause, error) {
	leftTableAlias := utils.ParseIdentifier(input)
	trimmed, _ := utils.HasPrefix(input, leftTableAlias)
	if leftTableAlias == "" || trimmed == "" {
		return nil, nil
	}

	var joinType types.JoinType
	if strings.HasPrefix(trimmed, ",") {
		joinType = types.InnerJoin
		trimmed, _ = utils.HasPrefix(trimmed, ",")
	} else if strings.HasPrefix(trimmed, "inner join") {
		joinType = types.InnerJoin
		trimmed, _ = utils.HasPrefix(trimmed, "inner join")
	} else if strings.HasPrefix(trimmed, "left outer join") {
		joinType = types.LeftOuterJoin
		trimmed, _ = utils.HasPrefix(trimmed, "left outer join")
	} else if strings.HasPrefix(trimmed, "right outer join") {
		joinType = types.RightOuterJoin
		trimmed, _ = utils.HasPrefix(trimmed, "right outer join")
	} else {
		return nil, nil
	}

	rightTableName := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, rightTableName)
	rightTableAlias := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, rightTableAlias)

	if joinType == types.InnerJoin {
		trimmed, _ = utils.HasPrefix(trimmed, "where")
	} else {
		trimmed, _ = utils.HasPrefix(trimmed, "on")
	}

	trimmed, _ = utils.HasPrefix(trimmed, leftTableAlias)
	trimmed, _ = utils.HasPrefix(trimmed, ".")
	leftTableColumn := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, leftTableColumn)

	trimmed, _ = utils.HasPrefix(trimmed, "=")

	trimmed, _ = utils.HasPrefix(trimmed, rightTableAlias)
	trimmed, _ = utils.HasPrefix(trimmed, ".")
	rightTableColumn := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, rightTableColumn)

	joinClause := &types.JoinClause{
		JoinType: joinType,
		LeftTable: leftTableName,
		LeftTableAlias: leftTableAlias,
		LeftTableColumn: leftTableColumn,
		RightTable: rightTableName,
		RightTableAlias: rightTableAlias,
		RightTableColumn: rightTableColumn,
	}

	return joinClause, nil
}

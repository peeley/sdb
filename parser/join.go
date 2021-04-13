package parser

import (
	"sdb/types"
	"sdb/utils"
	"strings"
)

func ParseJoinClause(input, leftTableName string) (*types.JoinClause, error) {
	leftTableAlias := utils.ParseIdentifier(input)
	trimmed, _ := utils.HasPrefix(input, leftTableAlias)

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
	}

	rightTableName := utils.ParseIdentifier(trimmed)
	trimmed, _ = utils.HasPrefix(trimmed, rightTableName)
	rightTableAlias := utils.ParseIdentifier(trimmed)

	joinClause := &types.JoinClause{
		JoinType: joinType,
		LeftTable: leftTableName,
		LeftTableAlias: leftTableAlias,
		RightTable: rightTableName,
		RightTableAlias: rightTableAlias,
		LeftTableColumn: "",
		RightTableColumn: "",
	}

	return joinClause, nil
}

package parser

import (
	"sdb/types"
	"sdb/utils"
)

func ParseBeginTransaction(input string) (types.Statement, error) {
	trimmed, ok := utils.HasPrefix(input, "begin transaction")

	if ok && trimmed == "" {
		return types.BeginTransaction{}, nil
	} else {
		return nil, nil
	}
}

func ParseCommit(input string) (types.Statement, error) {
	trimmed, ok := utils.HasPrefix(input, "commit")

	if ok && trimmed == "" {
		return types.Commit{}, nil
	} else {
		return nil, nil
	}
}

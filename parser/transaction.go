package parser

import (
	"sdb/statements"
	"sdb/utils"
)

func ParseBeginTransaction(input string) (statements.Executable, error) {
	trimmed, ok := utils.HasPrefix(input, "begin transaction")

	if ok && trimmed == "" {
		return statements.BeginTransaction{}, nil
	} else {
		return nil, nil
	}
}

func ParseCommit(input string) (statements.Executable, error) {
	trimmed, ok := utils.HasPrefix(input, "commit")

	if ok && trimmed == "" {
		return statements.Commit{}, nil
	} else {
		return nil, nil
	}
}

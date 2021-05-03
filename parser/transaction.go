package parser

import (
	"sdb/statements"
	"sdb/utils"
)

func ParseBeginTransaction(input string) (statements.Executable, error) {
	_, ok := utils.HasPrefix(input, "begin transaction")

	if ok {
		return statements.BeginTransaction{}, nil
	} else {
		return nil, nil
	}
}

func ParseCommit(input string) (statements.Executable, error) {
	_, ok := utils.HasPrefix(input, "commit")

	if ok {
		return statements.Commit{}, nil
	} else {
		return nil, nil
	}
}

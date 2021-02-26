package parser

import (
	"errors"
	"sdb/types"
	"strings"
)

func Parse(input string) (types.Statement, error) {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	if IsComment(input) {
		return types.Comment{}, nil
	}

	if input[len(input)-1] != ';' {
		return nil, errors.New("Missing ';' at end of statement.")
	}

	createTable, err := ParseCreateTableStatement(input)

	if err != nil {
		return nil, err
	} else if createTable != nil{
		return createTable, nil
	}

	createDB, err := ParseCreateDBStatement(input)

	if err != nil {
		return nil, err
	} else if createDB != nil {
		return createDB, nil
	}

	return nil, errors.New("Syntax error.")
}

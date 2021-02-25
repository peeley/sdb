package parser

import (
	"errors"
	"sdb/types"
	"strings"
)

func Parse(input string) (types.Statement, error) {
	input = strings.TrimSpace(input)

	if isComment(input) {
		return types.Comment{}, nil
	}

	create, err := parseCreateStatement(input)

	if err != nil {
		return nil, err
	} else if create != nil{
		return create, nil
	}

	return nil, errors.New("Syntax error.")
}

func isComment(input string) bool {
	return len(input) == 0 || (input[0] == '-' && input[1] == '-')
}

func parseCreateStatement(input string) (types.Statement, error){
	if strings.ToLower(input[:12]) != "create table" {
		return nil, nil
	}

	trimmed := strings.TrimSpace(input[12:])

	colList, err := parseColumnlList(input[12:])

	if colList == nil {
		return
	}
	statement := types.CreateStatement{
		TableName: "table",
		ColumnNames: make(map[string]string),
	}

	return &statement, nil
}

func parseColumnlList(input string) (map[string]string, error) {

	return nil, nil
}

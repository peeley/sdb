package parser

import (
	"fmt"
	"sdb/types"
	"strings"
	"errors"
)

func ParseCreateTableStatement(input string) (types.Statement, error){
	prefix := "create table"
	if len(input) < len(prefix) || !strings.HasPrefix(input, prefix) {
		return nil, nil
	}

	trimmed := strings.TrimPrefix(input, prefix)
	trimmed = strings.TrimSpace(trimmed)

	tableName := ParseIdentifier(trimmed)

	if tableName == "" {
		return nil, errors.New("Missing table name.")
	}

	trimmed = strings.TrimPrefix(trimmed, tableName)
	trimmed = strings.TrimSpace(trimmed)

	colList, err := parseColumnlList(trimmed)

	if err != nil {
		return nil, err
	}

	statement := types.CreateTableStatement{
		TableName: tableName,
		ColumnNames: colList,
	}

	return &statement, nil
}

func ParseCreateDBStatement(input string) (types.Statement, error) {
	prefix := "create database"
	if len(input) < len(prefix) || !strings.HasPrefix(input, prefix) {
		return nil, nil
	}

	trimmed := strings.TrimPrefix(input, prefix)
	trimmed = strings.TrimSpace(trimmed)
	ident := ParseIdentifier(trimmed)

	createDB := types.CreateDBStatement{
		DBName: ident,
	}

	return createDB, nil
}

func parseColumnlList(input string) (map[string]string, error) {
	fmt.Println("parsing column list:", input)
	if len(input) < 1 || input[0] != '(' {
		return nil, errors.New(
			"Expected '(' after table name in CREATE statement.",
		)
	}

	trimmed := strings.TrimPrefix(input, "(")
	cols := make(map[string]string)

	for {
		fmt.Printf("parsing column: '%v'\n", trimmed)
		trimmed = strings.TrimSpace(trimmed)
		ident := ParseIdentifier(trimmed)
		trimmed = strings.TrimPrefix(trimmed, ident)

		trimmed = strings.TrimSpace(trimmed)
		typeName, err := ParseTypename(trimmed)
		if err != nil {
			return nil, err
		}
		trimmed = strings.TrimPrefix(trimmed, typeName)

		cols[ident] = typeName
		trimmed = strings.TrimSpace(trimmed)

		if trimmed[0] == ')' {
			break
		} else if trimmed[0] != ',' {
			return nil, errors.New("Expected ',' in columns list.")
		}

		trimmed = strings.TrimPrefix(trimmed, ",")
	}

	if len(cols) < 1 {
		return nil, errors.New("Empty column list for CREATE statement.")
	}

	return cols, nil
}

package parser

import (
	"errors"
	"fmt"
	"sdb/types"
	"strings"
	"unicode"
)

func Parse(input string) (types.Statement, error) {
	input = strings.TrimSpace(input)

	if isComment(input) {
		return types.Comment{}, nil
	}

	if input[len(input)-1] != ';' {
		return nil, errors.New("Missing ';' at end of statement.")
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
	return len(input) < 2 || (input[0] == '-' && input[1] == '-')
}

func parseCreateStatement(input string) (types.Statement, error){
	prefix := "create table"
	if strings.ToLower(input[:len(prefix)]) != prefix {
		return nil, nil
	}

	trimmed := strings.TrimPrefix(input, prefix)
	trimmed = strings.TrimSpace(trimmed)

	tableName := parseIdentifier(trimmed)

	if tableName == "" {
		return nil, errors.New("Missing table name.")
	}

	trimmed = strings.TrimPrefix(trimmed, tableName)
	trimmed = strings.TrimSpace(trimmed)

	colList, err := parseColumnlList(trimmed)

	if err != nil {
		return nil, err
	}

	statement := types.CreateStatement{
		TableName: tableName,
		ColumnNames: colList,
	}

	return &statement, nil
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
		ident := parseIdentifier(trimmed)
		trimmed = strings.TrimPrefix(trimmed, ident)

		trimmed = strings.TrimSpace(trimmed)
		typeName, err := parseTypename(trimmed)
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

func parseIdentifier(input string) string {
	var builder strings.Builder

	for idx := 0; idx < len(input); idx++ {
		char := rune(input[idx])
		if unicode.IsSpace(char) ||
			!(unicode.IsLetter(char) || unicode.IsNumber(char)) {
			break
		}
		builder.WriteByte(input[idx])
	}

	return builder.String()
}

func parseTypename(input string) (string, error) {
	return "", nil
}

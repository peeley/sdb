// Noah Snelson
// February 25, 2021
// sdb/parser/create.go
//
// Contains parsing functions for `CREATE TABLE` and `CREATE DATABASE` functions.

package parser

import (
	"errors"
	"fmt"
	"sdb/types"
	"strings"
)

// Parses `CREATE TABLE <table_name> (<table_columns>);` input.
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

	colList, err := parseColumnList(trimmed)

	if err != nil {
		return nil, err
	}

	statement := types.CreateTableStatement{
		TableName: tableName,
		Columns: colList,
	}

	return &statement, nil
}

// Parses `CREATE DATABASE <db_name>;` input.
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

// Private utility function to parse <table_columns> into map of
// column name -> column type.
func parseColumnList(input string) (map[string]types.Type, error) {
	if len(input) < 1 || input[0] != '(' {
		return nil, errors.New(
			"Expected '(' after table name in CREATE statement.",
		)
	}

	trimmed := strings.TrimPrefix(input, "(")
	cols := make(map[string]types.Type)

	for {
		trimmed = strings.TrimSpace(trimmed)
		ident := ParseIdentifier(trimmed)
		trimmed = strings.TrimPrefix(trimmed, ident)

		trimmed = strings.TrimSpace(trimmed)
		typeName, err := ParseType(trimmed)
		if err != nil {
			return nil, err
		}
		trimmed = strings.TrimPrefix(trimmed, typeName.ToString())

		cols[ident] = typeName
		trimmed = strings.TrimSpace(trimmed)

		if trimmed[0] == ')' {
			break
		} else if trimmed[0] != ',' {
			return nil, fmt.Errorf("Expected ',' at end of column %v.", trimmed)
		}

		trimmed = strings.TrimPrefix(trimmed, ",")
	}

	if len(cols) < 1 {
		return nil, errors.New("Empty column list for CREATE statement.")
	}

	return cols, nil
}

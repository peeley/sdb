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
	"sdb/types/metatypes"
	"sdb/utils"
	"strings"
)

// Parses `CREATE TABLE <table_name> (<table_columns>);` input.
func ParseCreateTableStatement(input string) (types.Statement, error){
	trimmed, ok := utils.HasPrefix(input, "create table")
	if !ok {
		return nil, nil
	}

	tableName := utils.ParseIdentifier(trimmed)

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
	trimmed, ok := utils.HasPrefix(input, "create database")
	if !ok {
		return nil, nil
	}

	ident := utils.ParseIdentifier(trimmed)

	createDB := types.CreateDBStatement{
		DBName: ident,
	}

	return createDB, nil
}


// Private utility function to parse <table_columns> into map of
// column name -> column type.
func parseColumnList(input string) ([]metatypes.Column, error) {
	trimmed, ok := utils.HasPrefix(input, "(")
	if !ok {
		return nil, errors.New(
			"Expected '(' after table name in CREATE statement.",
		)
	}

	var cols []metatypes.Column

	for {
		trimmed = strings.TrimSpace(trimmed)
		ident := utils.ParseIdentifier(trimmed)
		trimmed = strings.TrimPrefix(trimmed, ident)

		trimmed = strings.TrimSpace(trimmed)
		colType, err := utils.ParseType(trimmed)
		if err != nil {
			return nil, err
		}
		trimmed = strings.TrimPrefix(trimmed, colType.ToString())

		cols = append(cols, metatypes.Column{ ident, colType })
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

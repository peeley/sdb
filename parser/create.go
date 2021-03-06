// Noah Snelson
// February 25, 2021
// sdb/parser/create.go
//
// Contains parsing functions for `CREATE TABLE` and `CREATE DATABASE` functions.

package parser

import (
	"errors"
	"sdb/db"
	"sdb/statements"
	"sdb/utils"
)

// Parses `CREATE TABLE <table_name> (<table_columns>);` input.
func ParseCreateTableStatement(input string) (db.Executable, error) {
	trimmed, ok := utils.HasPrefix(input, "create table")
	if !ok {
		return nil, nil
	}

	tableName := utils.ParseIdentifier(trimmed)

	if tableName == "" {
		return nil, errors.New("Missing table name.")
	}

	trimmed, _ = utils.HasPrefix(trimmed, tableName)
	trimmed, ok = utils.HasPrefix(trimmed, "(")
	if !ok {
		return nil, errors.New(
			"Expected '(' after table name in CREATE statement.",
		)
	}

	colList, err := utils.ParseColumnList(trimmed)

	if err != nil {
		return nil, err
	}

	if len(colList) < 1 {
		return nil, errors.New("Empty column list for CREATE statement.")
	}

	trimmed, _ = utils.HasPrefix(trimmed, utils.ColumnsToString(colList))

	trimmed, ok = utils.HasPrefix(trimmed, ")")
	if !ok {
		return nil, errors.New(
			"Expected ')' after column types in CREATE statement.",
		)
	}

	statement := statements.CreateTableStatement{
		TableName: tableName,
		Columns:   colList,
	}

	return &statement, nil
}

// Parses `CREATE DATABASE <db_name>;` input.
func ParseCreateDBStatement(input string) (db.Executable, error) {
	trimmed, ok := utils.HasPrefix(input, "create database")
	if !ok {
		return nil, nil
	}

	ident := utils.ParseIdentifier(trimmed)

	createDB := statements.CreateDBStatement{
		DBName: ident,
	}

	return createDB, nil
}

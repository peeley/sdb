// Noah Snelson
// May 1, 2021
// sdb/statements/create.go
//
// Contains logic for statements that create databases and tables.

package statements

import (
	"fmt"
	"os"
	"sdb/db"
	"sdb/utils"
)

type CreateDBStatement struct {
	DBName string
}

type CreateTableStatement struct {
	TableName string
	Columns   []db.Column
}

// Executes `CREATE DATABASE <db_name>;` query.
func (statement CreateDBStatement) Execute(state *db.DBState) error {
	err := os.Mkdir(statement.DBName, os.ModeDir|os.ModePerm)

	if err != nil {
		return fmt.Errorf(
			"!Failed to create database %v because it already exists.",
			statement.DBName,
		)
	}

	fmt.Printf("Database %v created.\n", statement.DBName)
	return nil
}

// Executes `CREATE TABLE <table_name> (<table_columns>);` queries.
func (statement CreateTableStatement) Execute(state *db.DBState) error {
	tablePath, exists := utils.TableExists(state, statement.TableName)

	if exists {
		return fmt.Errorf("!Failed to create table %v because it already exists.", statement.TableName)
	}

	tableFile, err := os.Create(tablePath)
	if err != nil {
		return fmt.Errorf("!Failed to create table %v because it already exists.", statement.TableName)
	}

	tableTypesString := utils.ColumnsToString(statement.Columns)
	tableFile.WriteString(tableTypesString)
	tableFile.WriteString("\n")

	fmt.Printf("Table %v created.\n", statement.TableName)
	return nil
}

// Noah Snelson
// February 25, 2021
// sdb/types/statement.go
//
// Contains type declarations for all SQL statements, as well as interfaces for
// executing SQL statements. Core logic of the database can be found in here.

package types

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// All SQL statement types implement this interface. The `Execute` function
// contains the core logic of the query, which is executed in the REPL at the
// `sdb/main.go` main function.
type Statement interface {
	Execute(*DBState) error
}

// These statement structs are what is output by their respective parsing
// functions. Every field of the struct represents a dynamic variable of the
// query, and the `Execute` function they implement uses these fields to
// implement the statement's functionality.
type CreateDBStatement struct {
	DBName string
}

type DropDBStatement struct {
	DBName string
}

type UseDBStatement struct {
	DBName string
}

type CreateTableStatement struct {
	TableName string
	Columns map[string]Type
}

type DropTableStatement struct {
	TableName string
}

type SelectStatement struct {
	TableName string
	Columns string
}

type AlterStatement struct {
	TableName string
	ColumnName string
	ColumnType Type
}

// Comments are essentially no-ops, but still parsed and as such need to
// implement the `Statement interface`
type Comment struct{}

// The `.EXIT` command is the only command currently implemented, so it can fit
// into the `Statement` interface as well.
type ExitCommand struct{}



// --- `Statement` interface implementations -----------------------------------

// Executes `CREATE DATABASE <db_name>;` query.
func (statement CreateDBStatement) Execute(state *DBState) error {
	err := os.Mkdir(statement.DBName, os.ModeDir)

	if err != nil {
		return fmt.Errorf(
			"!Failed to create database %v because it already exists.",
			statement.DBName,
		)
	}

	fmt.Printf("Database %v created.\n", statement.DBName)
	return nil
}

// Executes `DROP DATABASE <db_name>;` query.
func (statement DropDBStatement) Execute(state *DBState) error {
	_, err := os.Stat(statement.DBName)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.DBName)
	}

	os.RemoveAll(statement.DBName)

	fmt.Printf("Database %v deleted.\n", statement.DBName)
	return nil
}

// Executes `DROP TABLE <table_name>;` query. Assumes that the table being
// deleted is in the current database stored in DBState.
func (statement DropTableStatement) Execute(state *DBState) error {
	tablePath, exists := tableExists(state, statement.TableName)

	if !exists {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.TableName)
	}

	err := os.Remove(tablePath)

	if err != nil {
		return err
	}

	fmt.Printf("Deleted table %v.\n", statement.TableName)
	return nil
}

// Executes `USE <db_name>;` queries. Changes the current DB in DBState.
func (statement UseDBStatement) Execute(state *DBState) error {
	_, err := os.Stat(statement.DBName)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.DBName)
	}

	state.CurrentDB = statement.DBName
	fmt.Printf("Using database %v.\n", statement.DBName)
	return nil
}

// Executes `CREATE TABLE <table_name> (<table_columns>);` queries.
func (statement CreateTableStatement) Execute(state *DBState) error {
	tablePath, exists := tableExists(state, statement.TableName)

	if exists {
		return fmt.Errorf("!Failed to create table %v because it already exists.", statement.TableName)
	}

	tableFile, err := os.Create(tablePath)
	if err != nil {
		return fmt.Errorf("!Failed to create table %v because it already exists.", statement.TableName)
	}

	tableTypesString := columnsToString(statement.Columns)
	tableFile.WriteString(tableTypesString)
	tableFile.WriteString("\n")

	fmt.Printf("Table %v created.\n", statement.TableName)
	return nil
}

// Executes `SELECT <columns> FROM <table_name>;` queries. Currently only
// supports querying from every column via <columns> = `*`.
func (statement SelectStatement) Execute(state *DBState) error {
	tableFile, err := openTable(state, statement.TableName)
	if err != nil {
		return fmt.Errorf("!Failed to select from table %v because it does not exist.", statement.TableName)
	}

	// Tables are guaranteed to be empty, only need to read columns as header
	tableReader := bufio.NewReader(tableFile)
	columns, err := tableReader.ReadString('\n')
	if err != nil {
		return err
	}

	fmt.Println(columns)
	return nil
}

// Executes comments - comments are essentially no-ops.
func (statement Comment) Execute(state *DBState) error {
	return nil
}

// Executes the `.EXIT` command, exits from program.
func (statement ExitCommand) Execute(state *DBState) error {
	fmt.Println("\nGoodbye!")
	os.Exit(0)

	// Unreachable, but necessary for return type
	return nil
}

// Executes `ALTER TABLE <table_name> ADD <column_name> <column_type>;`
// statements.
func (statement AlterStatement) Execute(state *DBState) error {
	tableFile, err := openTable(state, statement.TableName)
	if err != nil {
		return fmt.Errorf(
			"!Failed to alter table %v because it does not exist.",
			statement.TableName,
		)
	}

	// read current header from table file
	reader := bufio.NewReader(tableFile)
	currentCols, err := reader.ReadString('\n')
	currentCols = currentCols[:len(currentCols)-1] // chop off last `\n` char

	if err != nil {
		return err
	}

	// create new header string based off current header
	var builder strings.Builder
	builder.WriteString(currentCols)
	builder.WriteString(
		fmt.Sprintf(", %v %v\n",
			statement.ColumnName,
			statement.ColumnType.ToString(),
		),
	)

	// overwrite header in table file with new header
	_, err = tableFile.WriteAt([]byte(builder.String()), 0)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Table %v modified, added column %v.\n",
		statement.TableName,
		statement.ColumnName,
	)

	return nil
}

// Private utility function to convert map representing column names to column
// types to a formatted string.
func columnsToString(columns map[string]Type) string {
	var tableTypesStringBuilder strings.Builder
	idx := 0

	for columnName, columnType := range columns {
		columnString := fmt.Sprintf("%v %v", columnName, columnType.ToString())

		if idx < len(columns) - 1 {
			columnString = fmt.Sprintf("%v, ", columnString)
		}
		idx += 1

		tableTypesStringBuilder.WriteString(columnString)
	}

	return tableTypesStringBuilder.String()
}

// Private utility function, opens table file based on current DBState and given
// table name.
func openTable(state *DBState, tableName string) (*os.File, error) {
	var tablePathBuilder strings.Builder
	tablePathBuilder.WriteString(state.CurrentDB)
	tablePathBuilder.WriteString("/")
	tablePathBuilder.WriteString(tableName)

	tablePath := tablePathBuilder.String()

	// open file with read & write permissions, unix perm bits set to 0777
	tableFile, err := os.OpenFile(tablePath, os.O_RDWR, 0777)
	if err != nil {
		return nil, fmt.Errorf("!Failed to select from table %v because it does not exist.", tableName)
	}

	return tableFile, nil
}

// Private utility function, determines if table exists given current DBState
// and given table name. Return table path and boolean representing existence of
// table.
func tableExists(state *DBState, tableName string) (string, bool) {
	var tablePathBuilder strings.Builder
	tablePathBuilder.WriteString(state.CurrentDB)
	tablePathBuilder.WriteString("/")
	tablePathBuilder.WriteString(tableName)

	tablePath := tablePathBuilder.String()

	_, err := os.Stat(tablePath)

	return tablePath, err == nil
}

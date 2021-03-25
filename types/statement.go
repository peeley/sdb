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
	"sdb/types/metatypes"
	"sdb/utils"
	"strings"
)

// All SQL statement types implement this interface. The `Execute` function
// contains the core logic of the query, which is executed in the REPL at the
// `sdb/main.go` main function.
type Statement interface {
	Execute(*metatypes.DBState) error
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
	Columns map[string]metatypes.Type
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
	ColumnType metatypes.Type
}

type InsertStatement struct {
	TableName string
	Values []metatypes.Value
}

// Comments are essentially no-ops, but still parsed and as such need to
// implement the `Statement interface`
type Comment struct{}


// --- `Statement` interface implementations -----------------------------------

// Executes `CREATE DATABASE <db_name>;` query.
func (statement CreateDBStatement) Execute(state *metatypes.DBState) error {
	err := os.Mkdir(statement.DBName, os.ModeDir | os.ModePerm)

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
func (statement DropDBStatement) Execute(state *metatypes.DBState) error {
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
func (statement DropTableStatement) Execute(state *metatypes.DBState) error {
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
func (statement UseDBStatement) Execute(state *metatypes.DBState) error {
	_, err := os.Stat(statement.DBName)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.DBName)
	}

	state.CurrentDB = statement.DBName
	fmt.Printf("Using database %v.\n", statement.DBName)
	return nil
}

// Executes `CREATE TABLE <table_name> (<table_columns>);` queries.
func (statement CreateTableStatement) Execute(state *metatypes.DBState) error {
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
func (statement SelectStatement) Execute(state *metatypes.DBState) error {
	tableFile, err := openTable(state, statement.TableName, os.O_RDONLY)
	defer tableFile.Close()
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
func (statement Comment) Execute(state *metatypes.DBState) error {
	return nil
}

// Executes `ALTER TABLE <table_name> ADD <column_name> <column_type>;`
// statements.
func (statement AlterStatement) Execute(state *metatypes.DBState) error {
	tableFile, err := openTable(state, statement.TableName, os.O_RDWR)
	defer tableFile.Close()
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

func (statement InsertStatement) Execute(state *metatypes.DBState) error {
	tableFile, err := openTable(state, statement.TableName, os.O_APPEND|os.O_RDWR)
	defer tableFile.Close()
	if err != nil {
		return fmt.Errorf("!Failed to insert into table %v because it does not exist.", statement.TableName)
	}
	defer tableFile.Close()

	reader := bufio.NewReader(tableFile)
	tableHeader, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("!Failed to read from table file %v.", statement.TableName)
	}

	fmt.Println("table header:", tableHeader)
	var tableTypes []metatypes.Type
	var ok bool
	for {
		if tableHeader == "" {
			break
		}

		ident := utils.ParseIdentifier(tableHeader)
		tableHeader, _ = utils.HasPrefix(tableHeader, ident)

		typeName, err := utils.ParseType(tableHeader)
		if err != nil {
			return err
		}
		tableTypes = append(tableTypes, typeName)

		tableHeader, _ = utils.HasPrefix(tableHeader, typeName.ToString())
		tableHeader, ok = utils.HasPrefix(tableHeader, ",")
		if !ok {
			break
		}

	}

	if len(tableTypes) != len(statement.Values) {
		return fmt.Errorf("!Failed, list of values to insert does not match table arity.")
	}
	// check types match
	for statementIdx, tableColType := range tableTypes {
		if !statement.Values[statementIdx].TypeMatches(&tableColType) {
			return fmt.Errorf("!Value %v is not of type %v", statement.Values[statementIdx], tableColType.ToString())
		}
	}

	var rowBuilder strings.Builder
	for idx, val := range statement.Values {
		rowBuilder.WriteString(val.ToString())
		if idx < len(statement.Values)-1 {
			rowBuilder.WriteString(", ")
		}
	}

	writer := bufio.NewWriter(tableFile)
	rowString := rowBuilder.String()

	_, err = writer.WriteString(rowString)
	if err != nil {
		return err
	}
	fmt.Printf("Inserted `%v` into %v\n", rowString, statement.TableName)

	return nil
}

// Private utility function to convert map representing column names to column
// types to a formatted string.
func columnsToString(columns map[string]metatypes.Type) string {
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
func openTable(state *metatypes.DBState, tableName string, flags int) (*os.File, error) {
	var tablePathBuilder strings.Builder
	tablePathBuilder.WriteString(state.CurrentDB)
	tablePathBuilder.WriteString("/")
	tablePathBuilder.WriteString(tableName)

	tablePath := tablePathBuilder.String()

	// open file with mode flags, unix perm bits set to 0777
	tableFile, err := os.OpenFile(tablePath, flags, 0777)
	if err != nil {
		return nil, fmt.Errorf("!Failed to select from table %v because it does not exist.", tableName)
	}

	return tableFile, nil
}

// Private utility function, determines if table exists given current DBState
// and given table name. Return table path and boolean representing existence of
// table.
func tableExists(state *metatypes.DBState, tableName string) (string, bool) {
	var tablePathBuilder strings.Builder
	tablePathBuilder.WriteString(state.CurrentDB)
	tablePathBuilder.WriteString("/")
	tablePathBuilder.WriteString(tableName)

	tablePath := tablePathBuilder.String()

	_, err := os.Stat(tablePath)

	return tablePath, err == nil
}

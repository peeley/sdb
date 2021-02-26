package types

import (
	"fmt"
	"os"
	"strings"
)

type Statement interface {
	Execute(*DBState) error
}

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

type Comment struct{}

type ExitCommand struct{}

func (statement CreateDBStatement) Execute(state *DBState) error {
	err := os.Mkdir(statement.DBName, os.ModeDir)

	if err != nil {
		return fmt.Errorf("!Failed to create database %v because it already exists.", statement.DBName)
	}

	fmt.Printf("Database %v created.\n", statement.DBName)
	return nil
}

func (statement DropDBStatement) Execute(state *DBState) error {
	_, err := os.Stat(statement.DBName)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.DBName)
	}

	os.RemoveAll(statement.DBName)

	fmt.Printf("Database %v deleted.\n", statement.DBName)
	return nil
}

func (statement DropTableStatement) Execute(state *DBState) error {
	var tablePathBuilder strings.Builder
	tablePathBuilder.WriteString(state.CurrentDB)
	tablePathBuilder.WriteString("/")
	tablePathBuilder.WriteString(statement.TableName)

	tablePath := tablePathBuilder.String()

	_, err := os.Stat(tablePath)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.TableName)
	}

	err = os.Remove(tablePath)

	if err != nil {
		return err
	}

	fmt.Printf("Deleted table %v.\n", statement.TableName)
	return nil
}

func (statement UseDBStatement) Execute(state *DBState) error {
	_, err := os.Stat(statement.DBName)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.DBName)
	}

	state.CurrentDB = statement.DBName
	fmt.Printf("Using database %v.\n", statement.DBName)
	return nil
}

func (statement CreateTableStatement) Execute(state *DBState) error {

	var tablePathBuilder strings.Builder
	tablePathBuilder.WriteString(state.CurrentDB)
	tablePathBuilder.WriteString("/")
	tablePathBuilder.WriteString(statement.TableName)

	tablePath := tablePathBuilder.String()

	_, err := os.Stat(tablePath)

	if err == nil {
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

func (statement Comment) Execute(state *DBState) error {
	// no-op
	return nil
}

func (statement ExitCommand) Execute(state *DBState) error{
	fmt.Println("\nGoodbye!")
	os.Exit(0)

	// unreachable, but necessary for return type
	return nil
}

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

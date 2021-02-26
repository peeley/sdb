package types

import (
	"os"
	"fmt"
)

var ConstWidthTypes = []string{"float", "int"}
var VariableWidthTypes = []string{"char", "varchar"}

type DBState struct {
	CurrentDB string
}

func NewState() DBState {
	return DBState{
		CurrentDB: "",
	}
}

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
	ColumnNames map[string]string
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
	// TODO
	fmt.Printf("creating table %v with cols %v\n",
		statement.TableName,
		statement.ColumnNames,
	)

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

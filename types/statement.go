package types

import (
	"os"
	"fmt"
)

var ConstWidthTypes = []string{"float", "int"}
var VariableWidthTypes = []string{"char", "varchar"}

type DBState struct {
	CurrentDB *os.File
}

func NewState() DBState {
	return DBState{
		CurrentDB: nil,
	}
}

type Statement interface {
	Execute(*DBState)
}

type CreateDBStatement struct {
	DBName string
}

type CreateTableStatement struct {
	TableName string
	ColumnNames map[string]string
}

type Comment struct{}

func (statement CreateDBStatement) Execute(state *DBState){
	fmt.Println("creating db", statement.DBName)
}

func (statement CreateTableStatement) Execute(state *DBState){
	// TODO
	fmt.Printf("creating table %v with cols %v\n",
		statement.TableName,
		statement.ColumnNames,
	)
}

func (statement Comment) Execute(state *DBState){}

package types

import (
	"os"
	"fmt"
)

const ConstWidthTypes = ["float", "int"]
const VariableWidthTypes = ["char", "varchar"]

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

type CreateStatement struct {
	TableName string
	ColumnNames map[string]string
}

type Comment struct{}

func (statement CreateStatement) Execute(state *DBState){
	// TODO
	fmt.Printf("creating table %v with cols %v\n",
		statement.TableName,
		statement.ColumnNames,
	)
}

func (statement Comment) Execute(state *DBState){}

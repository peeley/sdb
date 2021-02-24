package types

import "os"

type DBState struct {
	currentDB *os.File
}

func NewState() DBState {
	return DBState{
		currentDB: nil,
	}
}

type Statement interface {
	Execute(*DBState)
}

type CreateStatement struct {
	tableName string
	colNames map[string]string
}

type Comment struct{}

func (statement CreateStatement) Execute(state *DBState){}

func (statement Comment) Execute(state *DBState){}

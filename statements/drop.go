// Noah Snelson
// May 1, 2021
// sdb/statements/drop.go
//
// Contains logic for statements that drop databases and tables.

package statements

import (
	"fmt"
	"os"
	"sdb/db"
	"sdb/utils"
)

type DropDBStatement struct {
	DBName string
}

type DropTableStatement struct {
	TableName string
}

// Executes `DROP TABLE <table_name>;` query. Assumes that the table being
// deleted is in the current database stored in DBState.
func (statement DropTableStatement) Execute(state *db.DBState) error {
	tablePath, exists := utils.TableExists(state, statement.TableName)

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

// Executes `DROP DATABASE <db_name>;` query.
func (statement DropDBStatement) Execute(state *db.DBState) error {
	_, err := os.Stat(statement.DBName)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.DBName)
	}

	os.RemoveAll(statement.DBName)

	fmt.Printf("Database %v deleted.\n", statement.DBName)
	return nil
}

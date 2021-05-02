// Noah Snelson
// May 1, 2021
// sdb/statements/use.go

package statements

import (
	"fmt"
	"os"
	"sdb/db"
)

type UseDBStatement struct {
	DBName string
}

// Executes `USE <db_name>;` queries. Changes the current DB in DBState.
func (statement UseDBStatement) Execute(state *db.DBState) error {
	_, err := os.Stat(statement.DBName)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.DBName)
	}

	state.CurrentDB = statement.DBName
	fmt.Printf("Using database %v.\n", statement.DBName)
	return nil
}

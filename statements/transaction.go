// Noah Snelson
// May 1, 2021
// sdb/statements/transation.go
//
// Contains logic for transaction/commits.

package statements

import (
	"fmt"
	"os"
	"sdb/db"
)

// Transactions don't necessarily hold any information, but are Executables
// nonetheless
type BeginTransaction struct{}
type Commit struct{}

func (statement BeginTransaction) Execute(state *db.DBState) error {
	state.BeginTransaction()
	fmt.Printf("Transaction: %v\n", state.Transaction)
	fmt.Printf("Transaction started.\n")
	return nil
}

func (statement Commit) Execute(state *db.DBState) error {
	if len(state.Transaction.Statements) == 0 {
		return fmt.Errorf("Transaction abort.")
	}

	statements := state.Transaction.Statements
	lockFileNames := state.Transaction.LockFiles
	state.Transaction = nil

	for _, statement := range statements {
		statement.Execute(state)
	}

	for _, lockFileName := range lockFileNames {
		err := os.Remove(lockFileName)
		if err != nil {
			return fmt.Errorf(
				"!Failed to remove transaction lock %v",
				lockFileName,
			)
		}
	}
	state.Transaction = &db.Transaction{}

	fmt.Printf("Transaction committed.\n")
	return nil
}

// Creates lock file to signify table is undergoing transaction. Returns tuple
// of (string, error) signifying the name of the lock file created or any errors
// during creation.
func createTableLock(dbName, tableName string) (string, error) {
	if tableLockExists(dbName, tableName) {
		return "", fmt.Errorf("!Table %v is locked.", tableName)
	}

	lockFileName := dbName + "/." + tableName + "_lock"
	_, err := os.Create(lockFileName)
	if err != nil {
		return "", fmt.Errorf("!Failed to create transaction lock: %v", err)
	}

	return lockFileName, nil
}

func tableLockExists(dbName, tableName string) bool {
	lockFileName := dbName + "/." + tableName + "_lock"
	_, err := os.Stat(lockFileName)

	if err != nil && os.IsExist(err) {
		return true
	}

	return false
}

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
	fmt.Printf("Transaction started.\n")
	return nil
}

func (statement Commit) Execute(state *db.DBState) error {
	if state.Transaction == nil {
		fmt.Println("No current transaction.")
	}

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

	fmt.Printf("Transaction committed.\n")
	return nil
}

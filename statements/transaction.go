// Noah Snelson
// May 1, 2021
// sdb/statements/transation.go
//
// Contains logic for transaction/commits.

package statements

import (
	"fmt"
	"sdb/db"
)

// Transactions don't necessarily hold any information, but are Executables
// nonetheless
type BeginTransaction struct{}
type Commit struct{}

func (statement BeginTransaction) Execute(state *db.DBState) error {
	fmt.Printf("Beginning transaction!\n")
	return nil
}

func (statement Commit) Execute(state *db.DBState) error {
	fmt.Printf("Committing transaction!\n")
	return nil
}

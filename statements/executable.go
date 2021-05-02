// Noah Snelson
// February 25, 2021
// sdb/statements/statement.go

package statements

import (
	"sdb/db"
)

// All SQL statement types implement this interface. The `Execute` function
// contains the core logic of the query, which is executed in the REPL at the
// `sdb/main.go` main function. See each statement's respective file in
// `./<statement>.go`.
type Executable interface {
	Execute(*db.DBState) error
}

// Comments are basically no-ops, but are still Executables
type Comment struct{}

func (c Comment) Execute(_ *db.DBState) error {
	return nil
}

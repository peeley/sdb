// Noah Snelson
// May 2, 2021
// sdb/statements/comment.go

package statements

import (
	"sdb/db"
)

// Comments are basically no-ops, but are still Executable.
type Comment struct{}

func (c Comment) Execute(_ *db.DBState) error {
	return nil
}

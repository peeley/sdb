// Noah Snelson
// May 1, 2021
// sdb/statements/alter.go
//
// Implements logic for ALTER statement

package statements

import (
	"bufio"
	"fmt"
	"os"
	"sdb/db"
	"sdb/utils"
	"strings"
)

type AlterStatement struct {
	TableName  string
	ColumnName string
	ColumnType db.Type
}

// Executes `ALTER TABLE <table_name> ADD <column_name> <column_type>;`
// statements.
func (statement AlterStatement) Execute(state *db.DBState) error {
	tableFile, err := utils.OpenTable(state, statement.TableName, os.O_RDWR)
	if err != nil {
		return fmt.Errorf(
			"!Failed to alter table %v because it does not exist.",
			statement.TableName,
		)
	}
	defer tableFile.Close()

	// read current header from table file
	reader := bufio.NewReader(tableFile)
	currentCols, err := reader.ReadString('\n')
	currentCols = currentCols[:len(currentCols)-1] // chop off last `\n` char

	if err != nil {
		return err
	}

	// create new header string based off current header
	var builder strings.Builder
	builder.WriteString(currentCols)
	builder.WriteString(
		fmt.Sprintf(", %v %v\n",
			statement.ColumnName,
			statement.ColumnType.ToString(),
		),
	)

	// overwrite header in table file with new header
	_, err = tableFile.WriteAt([]byte(builder.String()), 0)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Table %v modified, added column %v.\n",
		statement.TableName,
		statement.ColumnName,
	)

	return nil
}

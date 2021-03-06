// Noah Snelson
// May 1, 2021
// sdb/statements/update.go
//
// Contains logic for UPDATE statements
package statements

import (
	"bufio"
	"fmt"
	"os"
	"sdb/db"
	"sdb/utils"
	"strings"
)

type UpdateStatement struct {
	TableName    string
	UpdatedCol   string
	UpdatedValue *db.Value
	WhereClause  *WhereClause
}

func (statement UpdateStatement) Execute(state *db.DBState) error {
	if state.IsTransacting() {
		// this process is transacting, add this statement to transaction
		lockFileName, err := state.AcquireTableLock(statement.TableName)
		if err != nil {
			return err
		}

		state.Transaction.LockFiles = append(
			state.Transaction.LockFiles,
			lockFileName,
		)
		state.Transaction.Statements = append(
			state.Transaction.Statements,
			statement,
		)

		fmt.Println("Added update to transaction.")

		return nil
	} else if state.TableLockExists(statement.TableName) {
		// another process' transaction has locked the table, can't do anything
		fmt.Printf("!Table %v is locked.\n", statement.TableName)
		return nil
	}

	tableFile, err := utils.OpenTable(state, statement.TableName, os.O_RDONLY)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("!Failed to insert into table %v because it does not exist.", statement.TableName)
	}
	defer tableFile.Close()

	reader := bufio.NewReader(tableFile)
	tableHeader, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("!Failed to read from table file %v.", statement.TableName)
	}

	colNames := utils.TableHeaderToColMap(tableHeader)

	var replaceStringBuilder strings.Builder
	replaceStringBuilder.WriteString(tableHeader)

	updated := 0

	for {
		row, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		rowValues, _, _ := utils.ParseValueList(row)
		if whereApplies(statement.WhereClause, colNames, rowValues) {

			rowValues[colNames[statement.UpdatedCol]] = *statement.UpdatedValue
			updatedRowString := utils.ValueListToString(rowValues)
			replaceStringBuilder.WriteString(updatedRowString)

			updated += 1
		} else {
			replaceStringBuilder.WriteString(row)
		}
	}

	// need to close file before reopening to truncate
	tableFile.Close()
	tableFile, err = utils.OpenTable(state, statement.TableName, os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer tableFile.Close()

	replacedTable := replaceStringBuilder.String()
	tableFile.WriteString(replacedTable)

	fmt.Printf("Updated %v rows.\n", updated)

	return nil
}

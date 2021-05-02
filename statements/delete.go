// Noah Snelson
// May 1, 2021
// sdb/statements/delete.go
//
// Contains logic for DELETE statement

package statements

import (
	"bufio"
	"fmt"
	"os"
	"sdb/db"
	"sdb/utils"
	"strings"
)

type DeleteStatment struct {
	TableName   string
	WhereClause *WhereClause
}

func (statement DeleteStatment) Execute(state *db.DBState) error {
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

	deleted := 0

	for {
		row, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		rowValues, _, _ := utils.ParseValueList(row)
		if !whereApplies(statement.WhereClause, colNames, rowValues) {
			replaceStringBuilder.WriteString(row)
		} else {
			deleted += 1
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

	fmt.Printf("Deleted %v rows.\n", deleted)
	return nil
}

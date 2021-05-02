// Noah Snelson
// May 1, 2021
// sdb/statements/insert.go
//
// Contains logic for INSERT statement.

package statements

import (
	"bufio"
	"fmt"
	"os"
	"sdb/db"
	"sdb/utils"
	"strings"
)

type InsertStatement struct {
	TableName string
	Values    []db.Value
}

func (statement InsertStatement) Execute(state *db.DBState) error {
	tableFile, err := utils.OpenTable(state, statement.TableName, os.O_APPEND|os.O_RDWR)
	if err != nil {
		return fmt.Errorf("!Failed to insert into table %v because it does not exist.", statement.TableName)
	}
	defer tableFile.Close()

	reader := bufio.NewReader(tableFile)
	tableHeader, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("!Failed to read from table file %v.", statement.TableName)
	}

	var tableTypes []db.Type
	var ok bool
	for {
		if tableHeader == "" {
			break
		}

		ident := utils.ParseIdentifier(tableHeader)
		tableHeader, _ = utils.HasPrefix(tableHeader, ident)

		typeName, err := utils.ParseType(tableHeader)
		if err != nil {
			return err
		}
		tableTypes = append(tableTypes, typeName)

		tableHeader, _ = utils.HasPrefix(tableHeader, typeName.ToString())
		tableHeader, ok = utils.HasPrefix(tableHeader, ",")
		if !ok {
			break
		}

	}

	if len(tableTypes) != len(statement.Values) {
		return fmt.Errorf("!Failed, list of values to insert does not match table arity.")
	}
	// check types match
	for statementIdx, tableColType := range tableTypes {
		if !statement.Values[statementIdx].TypeMatches(&tableColType) {
			return fmt.Errorf("!Value %v is not of type %v", statement.Values[statementIdx], tableColType.ToString())
		}
	}

	var rowBuilder strings.Builder
	for idx, val := range statement.Values {
		rowBuilder.WriteString(val.ToString())
		if idx < len(statement.Values)-1 {
			rowBuilder.WriteString(", ")
		}
	}
	rowBuilder.WriteRune('\n')

	rowString := rowBuilder.String()

	_, err = tableFile.WriteString(rowString)
	if err != nil {
		return err
	}
	fmt.Printf("Inserted {%v} into %v\n", strings.TrimSpace(rowString), statement.TableName)

	return nil
}

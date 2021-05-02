// Noah Snelson
// May 1, 2021
// sdb/statements/select.go
//
// Implements logic for SELECT statement.
// Also includes types & logic for joins.

package statements

import (
	"bufio"
	"fmt"
	"os"
	"sdb/db"
	"sdb/utils"
	"strings"
)

type SelectStatement struct {
	TableName   string
	ColumnNames []string
	JoinClause  *JoinClause
	WhereClause *WhereClause
}

// Executes `SELECT <columns> FROM <table_name> [WHERE <condition>];` queries.
func (statement SelectStatement) Execute(state *db.DBState) error {
	tableFile, err := utils.OpenTable(state, statement.TableName, os.O_RDONLY)
	if err != nil {
		return fmt.Errorf("!Failed to select from table %v because it does not exist.", statement.TableName)
	}
	defer tableFile.Close()

	var joinTableFile *os.File
	if statement.JoinClause != nil {
		joinTableFile, err = utils.OpenTable(state, statement.JoinClause.RightTable, os.O_RDONLY)
		if err != nil {
			return fmt.Errorf("!Failed to select from table %v because it does not exist.", statement.JoinClause.RightTable)
		}
		defer joinTableFile.Close()
	}

	var outputBuilder strings.Builder

	tableReader := bufio.NewReader(tableFile)
	tableHeader, err := tableReader.ReadString('\n')
	if err != nil {
		return err
	}

	// write names of selected columns into output
	tableColumns, err := utils.ParseColumnList(tableHeader)
	if err != nil {
		return err
	}

	// need to forward declare variables used in JOINs
	var joinTableColMap map[string]int
	var joinTableRows [][]db.Value
	if statement.JoinClause != nil {
		joinTableReader := bufio.NewReader(joinTableFile)
		joinTableHeader, _ := joinTableReader.ReadString('\n')
		joinTableColMap = utils.TableHeaderToColMap(joinTableHeader)

		for {
			joinRow, err := joinTableReader.ReadString('\n')
			if err != nil {
				break
			}
			joinRowValues, _, _ := utils.ParseValueList(joinRow)
			joinTableRows = append(joinTableRows, joinRowValues)
		}

		// add joined columns to header
		var builder strings.Builder
		builder.WriteString(strings.TrimSpace(tableHeader))
		builder.WriteString(", ")
		builder.WriteString(joinTableHeader)

		tableHeader = builder.String()
	}

	if statement.ColumnNames[0] == "*" {
		outputBuilder.WriteString(tableHeader)
	} else {
		for statementColumnsIdx, statementColumnName := range statement.ColumnNames {
			for _, tableColumn := range tableColumns {
				if tableColumn.Name == statementColumnName {

					outputBuilder.WriteString(tableColumn.Name)
					outputBuilder.WriteString(" ")
					outputBuilder.WriteString(tableColumn.Type.ToString())

					if statementColumnsIdx < len(statement.ColumnNames)-1 {
						outputBuilder.WriteString(", ")
					}
				}
			}
		}
		outputBuilder.WriteString("\n")
	}

	colMap := utils.TableHeaderToColMap(tableHeader)
	var rowStringBuilder strings.Builder

	// iterate through all rows/lines of the table file and process as necessary
	for {
		row, err := tableReader.ReadString('\n')
		if err != nil {
			break
		}

		rowValues, _, _ := utils.ParseValueList(row)

		if statement.JoinClause != nil {
			// this join row to the joining table's matching row
			joined := applyJoin(*statement.JoinClause, colMap, joinTableColMap, joinTableRows, rowValues)
			if joined == "" {
				continue
			} else {
				row = joined
			}
		}

		// filter out rows according to `where`
		if whereApplies(statement.WhereClause, colMap, rowValues) {
			if statement.ColumnNames[0] == "*" {
				outputBuilder.WriteString(row)
			} else {
				// filter out selected columns
				for colNameIdx, colName := range statement.ColumnNames {
					selectedValue := rowValues[colMap[colName]]
					rowStringBuilder.WriteString(selectedValue.ToString())
					if colNameIdx < len(statement.ColumnNames)-1 {
						rowStringBuilder.WriteString(", ")
					}
				}
				rowStringBuilder.WriteString("\n")
			}

			outputBuilder.WriteString(rowStringBuilder.String())
		}
		rowStringBuilder.Reset()
	}

	fmt.Println(outputBuilder.String())
	return nil
}

type JoinType string

const (
	InnerJoin      = "inner join"
	LeftOuterJoin  = "left outer join"
	RightOuterJoin = "right outer join"
)

type JoinClause struct {
	JoinType         JoinType
	LeftTable        string
	LeftTableAlias   string
	LeftTableColumn  string
	RightTable       string
	RightTableAlias  string
	RightTableColumn string
}

// Determines if a row in the 'left' table should be joined to any rows from the
// 'right' table. Assumes that tables are being joined on an equality
// comparison between the joining columns.
func applyJoin(
	joinClause JoinClause,
	colNames map[string]int,
	joinColNames map[string]int,
	joinRows [][]db.Value,
	row []db.Value,
) string {
	var joinedRowsBuilder strings.Builder
	for _, joinRow := range joinRows {
		leftColIdx := colNames[joinClause.LeftTableColumn]
		rightColIdx := joinColNames[joinClause.RightTableColumn]

		var matchingRowBuilder strings.Builder
		if row[leftColIdx] == joinRow[rightColIdx] {
			matchingRowBuilder.WriteString(
				strings.TrimSpace(utils.ValueListToString(row)),
			)
			matchingRowBuilder.WriteString(", ")
			matchingRowBuilder.WriteString(utils.ValueListToString(joinRow))

			joinedRowsBuilder.WriteString(matchingRowBuilder.String())
			matchingRowBuilder.Reset()
		}
	}

	joinedRows := joinedRowsBuilder.String()

	if joinedRows == "" && joinClause.JoinType == LeftOuterJoin {
		return utils.ValueListToString(row)
	}

	return joinedRows
}

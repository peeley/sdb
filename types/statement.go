// Noah Snelson
// February 25, 2021
// sdb/types/statement.go
//
// Contains type declarations for all SQL statements, as well as interfaces for
// executing SQL statements. Core logic of the database can be found in here.

package types

import (
	"bufio"
	"fmt"
	"os"
	"sdb/types/metatypes"
	"sdb/utils"
	"strings"
)

// All SQL statement types implement this interface. The `Execute` function
// contains the core logic of the query, which is executed in the REPL at the
// `sdb/main.go` main function.
type Statement interface {
	Execute(*metatypes.DBState) error
}

// These statement structs are what is output by their respective parsing
// functions. Every field of the struct represents a dynamic variable of the
// query, and the `Execute` function they implement uses these fields to
// implement the statement's functionality.
type CreateDBStatement struct {
	DBName string
}

type DropDBStatement struct {
	DBName string
}

type UseDBStatement struct {
	DBName string
}

type CreateTableStatement struct {
	TableName string
	Columns []metatypes.Column
}

type DropTableStatement struct {
	TableName string
}

type SelectStatement struct {
	TableName string
	ColumnNames []string
	JoinClause *JoinClause
	WhereClause *WhereClause
}

type AlterStatement struct {
	TableName string
	ColumnName string
	ColumnType metatypes.Type
}

type InsertStatement struct {
	TableName string
	Values []metatypes.Value
}

type UpdateStatement struct {
	TableName string
	UpdatedCol string
	UpdatedValue *metatypes.Value
	WhereClause *WhereClause
}

type DeleteStatment struct {
	TableName string
	WhereClause *WhereClause
}

type WhereClause struct {
	ColName string
	Comparison string
	ComparisonValue *metatypes.Value
}

type JoinType string
const (
	InnerJoin = "inner join"
	LeftOuterJoin = "left outer join"
	RightOuterJoin = "right outer join"
)

type JoinClause struct {
	JoinType JoinType
	LeftTable string
	LeftTableAlias string
	LeftTableColumn string
	RightTable string
	RightTableAlias string
	RightTableColumn string
}

// Comments are essentially no-ops, but still parsed and as such need to
// implement the `Statement interface`
type Comment struct{}

// Transactions don't necessarily hold any information, but are a Statement
// nonetheless
type BeginTransaction struct{}
type Commit struct{}


// --- `Statement` interface implementations -----------------------------------

// Executes `CREATE DATABASE <db_name>;` query.
func (statement CreateDBStatement) Execute(state *metatypes.DBState) error {
	err := os.Mkdir(statement.DBName, os.ModeDir | os.ModePerm)

	if err != nil {
		return fmt.Errorf(
			"!Failed to create database %v because it already exists.",
			statement.DBName,
		)
	}

	fmt.Printf("Database %v created.\n", statement.DBName)
	return nil
}

// Executes `DROP DATABASE <db_name>;` query.
func (statement DropDBStatement) Execute(state *metatypes.DBState) error {
	_, err := os.Stat(statement.DBName)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.DBName)
	}

	os.RemoveAll(statement.DBName)

	fmt.Printf("Database %v deleted.\n", statement.DBName)
	return nil
}

// Executes `DROP TABLE <table_name>;` query. Assumes that the table being
// deleted is in the current database stored in DBState.
func (statement DropTableStatement) Execute(state *metatypes.DBState) error {
	tablePath, exists := utils.TableExists(state, statement.TableName)

	if !exists {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.TableName)
	}

	err := os.Remove(tablePath)

	if err != nil {
		return err
	}

	fmt.Printf("Deleted table %v.\n", statement.TableName)
	return nil
}

// Executes `USE <db_name>;` queries. Changes the current DB in DBState.
func (statement UseDBStatement) Execute(state *metatypes.DBState) error {
	_, err := os.Stat(statement.DBName)

	if err != nil {
		return fmt.Errorf("!Failed to delete %v because it does not exist.", statement.DBName)
	}

	state.CurrentDB = statement.DBName
	fmt.Printf("Using database %v.\n", statement.DBName)
	return nil
}

// Executes `CREATE TABLE <table_name> (<table_columns>);` queries.
func (statement CreateTableStatement) Execute(state *metatypes.DBState) error {
	tablePath, exists := utils.TableExists(state, statement.TableName)

	if exists {
		return fmt.Errorf("!Failed to create table %v because it already exists.", statement.TableName)
	}

	tableFile, err := os.Create(tablePath)
	if err != nil {
		return fmt.Errorf("!Failed to create table %v because it already exists.", statement.TableName)
	}

	tableTypesString := utils.ColumnsToString(statement.Columns)
	tableFile.WriteString(tableTypesString)
	tableFile.WriteString("\n")

	fmt.Printf("Table %v created.\n", statement.TableName)
	return nil
}

// Executes `SELECT <columns> FROM <table_name> [WHERE <condition>];` queries.
func (statement SelectStatement) Execute(state *metatypes.DBState) error {
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
	var joinTableRows [][]metatypes.Value
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

					if statementColumnsIdx < len(statement.ColumnNames) - 1 {
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
					if colNameIdx < len(statement.ColumnNames) - 1 {
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

// Executes comments - comments are essentially no-ops.
func (statement Comment) Execute(state *metatypes.DBState) error {
	return nil
}

// Executes `ALTER TABLE <table_name> ADD <column_name> <column_type>;`
// statements.
func (statement AlterStatement) Execute(state *metatypes.DBState) error {
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

func (statement InsertStatement) Execute(state *metatypes.DBState) error {
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

	var tableTypes []metatypes.Type
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

func (statement UpdateStatement) Execute(state *metatypes.DBState) error {
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
		return err;
	}
	defer tableFile.Close()

	replacedTable := replaceStringBuilder.String()
	tableFile.WriteString(replacedTable)

	fmt.Printf("Updated %v rows.\n", updated)

	return nil
}

func (statement DeleteStatment) Execute(state *metatypes.DBState) error {
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
		return err;
	}
	defer tableFile.Close()

	replacedTable := replaceStringBuilder.String()
	tableFile.WriteString(replacedTable)

	fmt.Printf("Deleted %v rows.\n", deleted)
	return nil
}

func (statement BeginTransaction) Execute(state *metatypes.DBState) error {
	fmt.Printf("Beginning transaction!\n")
	return nil;
}

func (statement Commit) Execute(state *metatypes.DBState) error {
	fmt.Printf("Committing transaction!\n")
	return nil;
}

// Determines if `where` clause applies to row
func whereApplies(where *WhereClause, colNames map[string]int, row []metatypes.Value) bool {
	if where == nil {
		return true
	}
	colIndex := colNames[where.ColName]

	rowValue := row[colIndex]

	// FIXME might want to check if types match before comparison
	if where.Comparison == "=" {
		return rowValue.GetValue() == where.ComparisonValue.GetValue()
	} else if where.Comparison == "!=" {
		return rowValue.GetValue() != where.ComparisonValue.GetValue()
	} else if where.Comparison == "<" { // assuming numerical types for less/greater than
		return rowValue.GetValue().(float64) < where.ComparisonValue.GetValue().(float64)
	} else if where.Comparison == "<=" {
		return rowValue.GetValue().(float64) <= where.ComparisonValue.GetValue().(float64)
	} else if where.Comparison == ">" {
		return rowValue.GetValue().(float64) > where.ComparisonValue.GetValue().(float64)
	} else if where.Comparison == ">=" {
		return rowValue.GetValue().(float64) >= where.ComparisonValue.GetValue().(float64)
	}

	return false
}

// Determines if a row in the 'left' table should be joined to any rows from the
// 'right' table. Assumes that tables are being joined on an equality
// comparison between the joining columns.
func applyJoin(joinClause JoinClause, colNames map[string]int, joinColNames map[string]int, joinRows [][]metatypes.Value, row []metatypes.Value) string {
	var joinedRowsBuilder strings.Builder
	for _, joinRow := range joinRows {
		leftColIdx := colNames[joinClause.LeftTableColumn]
		rightColIdx := joinColNames[joinClause.RightTableColumn]

		var matchingRowBuilder strings.Builder
		if row[leftColIdx] == joinRow[rightColIdx] {

			matchingRowBuilder.WriteString(strings.TrimSpace(utils.ValueListToString(row)))
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

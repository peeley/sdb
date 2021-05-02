// Noah Snelson
// May 1, 2021
// sdb/statements/where.go
//
// Implements logic for WHERE clauses in SELECT/UPDATE statements.

package statements

import "sdb/db"

type WhereClause struct {
	ColName         string
	Comparison      string
	ComparisonValue *db.Value
}

// Determines if `where` clause applies to row
func whereApplies(where *WhereClause, colNames map[string]int, row []db.Value) bool {
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

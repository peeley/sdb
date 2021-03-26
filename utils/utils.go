// Noah Snelson
// February 25, 2021
// sdb/parser/go
//
// Common parsing utilities used across several top-level parsing functions.

package utils

import (
	"fmt"
	"sdb/types/metatypes"
	"strconv"
	"strings"
	"unicode"
	"os"
)

// Detects if input is comment by checking if it begins with "--"
func IsComment(input string) bool {
	return len(input) < 2 || strings.HasPrefix(input, "--")
}

func HasPrefix(input string, prefix string) (string, bool) {
	if len(input) < len(prefix) || !strings.HasPrefix(input, prefix) {
		return input, false
	}

	trimmed := strings.TrimPrefix(input, prefix)
	trimmed = strings.TrimSpace(trimmed)

	return trimmed, true
}

// Parses identifiers, which is any sequence of letters, numbers, and `_`
func ParseIdentifier(input string) string {
	var builder strings.Builder

	for idx := 0; idx < len(input); idx++ {
		char := rune(input[idx])
		if !(unicode.IsLetter(char) || unicode.IsNumber(char) || char == rune('_') || char == rune('*')) {
			break
		}
		builder.WriteByte(input[idx])
	}

	return builder.String()
}

// Parses the various types the database supports, like `float`, `int`,
// `char(X)`, and `varchar(X)`.
func ParseType(input string) (metatypes.Type, error) {
	baseType := ParseIdentifier(input)

	for _, typeName := range metatypes.ConstWidthTypes {
		if typeName == baseType {
			return metatypes.NewType(typeName, 0), nil
		}
	}

	trimmed := strings.TrimPrefix(input, baseType)
	if len(trimmed) < 1 || trimmed[0] != '(' {
		return nil, fmt.Errorf("Expected '(' after typename %v.", baseType)
	}
	trimmed = strings.TrimPrefix(trimmed, "(")

	var numberBuilder strings.Builder
	for _, digit := range trimmed {
		if !unicode.IsNumber(digit) {
			break
		}
		numberBuilder.WriteRune(digit)
	}

	numberString := numberBuilder.String()

	trimmed = strings.TrimPrefix(trimmed, numberString)
	if len(trimmed) < 1 || trimmed[0] != ')' {
		return nil, fmt.Errorf("Expected ')' after parameters of type %v.", baseType)
	}
	trimmed = strings.TrimPrefix(trimmed, ")")

	size, err := strconv.Atoi(numberString)
	if err != nil {
		return nil, err
	}

	return metatypes.NewType(baseType, size), nil
}

// Parses type in tuple e.g. 123, 3.14, "hello"
func ParseValue(input string) (*metatypes.Value, error) {

	float, err := ParseFloat(input)
	if float != nil {
		return float, nil
	} else if err != nil {
		return nil, err
	}

	int, err := ParseInt(input)
	if int != nil {
		return int, nil
	} else if err != nil {
		return nil, err
	}

	string, err := ParseString(input)
	if string != nil {
		return string, nil
	} else if err != nil {
		return nil, err
	}

	return &metatypes.Value{ Value: nil, Type: metatypes.Null{}}, nil
}

// Parse floating point numeric of arbitrary precision
func ParseInt(input string) (*metatypes.Value, error) {

	var integerBuilder strings.Builder
	for _, digit := range input {
		if !unicode.IsNumber(digit) {
			break
		}
		integerBuilder.WriteRune(digit)
	}
	integerString := integerBuilder.String()

	if integerString == "" {
		return nil, nil
	}

	integer, err := strconv.Atoi(integerString)

	if err != nil {
		return nil, err
	}

	val := metatypes.Value {
		Value: integer,
		Type: metatypes.Int{},
	}
	return &val, nil
}

// Parse integer
func ParseFloat(input string) (*metatypes.Value, error) {
	var integerBuilder strings.Builder
	for _, digit := range input {
		if !unicode.IsNumber(digit) {
			break
		}
		integerBuilder.WriteRune(digit)
	}
	integerString := integerBuilder.String()
	integer, _ := strconv.Atoi(integerString)

	trimmed, _ := HasPrefix(input, integerString)
	trimmed, ok := HasPrefix(trimmed, ".")
	if !ok {
		return nil, nil
	}

	var decimalBuilder strings.Builder
	decimalBuilder.WriteString("0.")
	for _, digit := range trimmed {
		if !unicode.IsNumber(digit) {
			break
		}
		decimalBuilder.WriteRune(digit)
	}
	decimalString := decimalBuilder.String()
	decimal, _ := strconv.ParseFloat(decimalString, 64)

	val := &metatypes.Value{
		Value: float64(integer) + decimal,
		Type: metatypes.Float{},
	}
	return val, nil
}

// Parse string.
// Always returns a value of varchar(length of string)
// This is checked against the column var/varchar(length) later
func ParseString(input string) (*metatypes.Value, error) {
	trimmed, ok := HasPrefix(input, "'")
	if !ok {
		return nil, fmt.Errorf("Expected string to start with `'`")
	}

	var stringBuilder strings.Builder
	for _, letter := range trimmed {
		if letter == rune('\'') {
			break
		}
		stringBuilder.WriteRune(letter)
	}
	string := stringBuilder.String()

	trimmed, _ = HasPrefix(trimmed, string)
	trimmed, ok = HasPrefix(trimmed, "'")
	if !ok {
		return nil, fmt.Errorf("Expected string to end with `'`")
	}

	val := &metatypes.Value{
		Value: string,
		Type: metatypes.VarChar{ Size: len(string) },
	}

	return val, nil
}

func ParseValueList(input string) ([]metatypes.Value, string, error) {
	var valueList []metatypes.Value

	trimmed := input
	var ok bool
	for {
		value, err := ParseValue(trimmed)
		if err != nil {
			return nil, input, err
		}
		trimmed, _ = HasPrefix(trimmed, value.ToString())
		valueList = append(valueList, *value)
		trimmed, ok = HasPrefix(trimmed, ",")
		if !ok {
			return valueList, trimmed, nil
		}
	}
}

func ValueListToString(list []metatypes.Value) string {
	var stringBuilder strings.Builder
	for idx, val := range list {
		stringBuilder.WriteString(val.ToString())
		if idx < len(list) - 1 {
			stringBuilder.WriteString(", ")
		}
	}
	stringBuilder.WriteString("\n")

	return stringBuilder.String()
}

// Determines if table exists given current DBState and given table name. Return
// table path and boolean representing existence of table.
func TableExists(state *metatypes.DBState, tableName string) (string, bool) {
	var tablePathBuilder strings.Builder
	tablePathBuilder.WriteString(state.CurrentDB)
	tablePathBuilder.WriteString("/")
	tablePathBuilder.WriteString(tableName)

	tablePath := tablePathBuilder.String()

	_, err := os.Stat(tablePath)

	return tablePath, err == nil
}

// Opens table file based on current DBState and given table name.
func OpenTable(state *metatypes.DBState, tableName string, flags int) (*os.File, error) {
	var tablePathBuilder strings.Builder
	tablePathBuilder.WriteString(state.CurrentDB)
	tablePathBuilder.WriteString("/")
	tablePathBuilder.WriteString(tableName)

	tablePath := tablePathBuilder.String()

	// open file with mode flags, unix perm bits set to 0777
	tableFile, err := os.OpenFile(tablePath, flags, 0777)
	if err != nil {
		return nil, fmt.Errorf("!Failed to select from table %v because it does not exist.", tableName)
	}

	return tableFile, nil
}

// Convert mapping of column names -> column types to a formatted string.
func ColumnsToString(columns []metatypes.Column) string {
	var tableTypesStringBuilder strings.Builder
	var columnString string

	for idx, column := range columns {
		columnString = fmt.Sprintf("%v %v", column.Name, column.Type.ToString())
		tableTypesStringBuilder.WriteString(columnString)

		if idx < len(columns) - 1 {
			tableTypesStringBuilder.WriteString(", ")
		}
	}

	return tableTypesStringBuilder.String()
}

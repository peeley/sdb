// Noah Snelson
// February 25, 2021
// sdb/parser/utils.go
//
// Common parsing utilities used across several top-level parsing functions.

package parser

import (
	"fmt"
	"sdb/types"
	"strconv"
	"strings"
	"unicode"
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
		if unicode.IsSpace(char) ||
			!(unicode.IsLetter(char) || unicode.IsNumber(char) || char == rune('_')) {
			break
		}
		builder.WriteByte(input[idx])
	}

	return builder.String()
}

// Parses the various types the database supports, like `float`, `int`,
// `char(X)`, and `varchar(X)`.
func ParseType(input string) (types.Type, error) {
	baseType := ParseIdentifier(input)

	for _, typeName := range types.ConstWidthTypes {
		if typeName == baseType {
			return types.NewType(typeName, 0), nil
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

	return types.NewType(baseType, size), nil
}

// Parses type in tuple e.g. 123, 3.14, "hello"
func ParseValue(input string) (*types.Value, error) {

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

	return &types.Value{ Value: nil, Type: types.Null{}}, nil
}

// Parse floating point numeric of arbitrary precision
func ParseInt(input string) (*types.Value, error) {

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

	val := types.Value {
		Value: integer,
		Type: types.Int{},
	}
	return &val, nil
}

// Parse integer
func ParseFloat(input string) (*types.Value, error) {
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

	val := &types.Value{
		Value: float64(integer) + decimal,
		Type: types.Float{},
	}
	return val, nil
}

// Parse string.
// Always returns a value of varchar(length of string)
// This is checked against the column var/varchar(length) later
func ParseString(input string) (*types.Value, error) {
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

	val := &types.Value{
		Value: string,
		Type: types.VarChar{ Size: len(string) },
	}

	return val, nil
}

func ParseValueList(input string) ([]types.Value, string, error) {
	var valueList []types.Value

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

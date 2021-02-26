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

// Detects if input is `.EXIT` command
func IsExitCommand(input string) bool {
	return strings.HasPrefix(".exit", input)
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

package parser

import (
	"errors"
	"sdb/types"
	"strings"
	"fmt"
)

func Parse(input string) (types.Statement, error) {
	fmt.Println("parsing: ", input)
	input = strings.TrimSpace(input)

	if isComment(input) {
		return types.Comment{}, nil
	}

	return nil, errors.New("Syntax error.")
}

func isComment(input string) bool {
	return len(input) == 0 || (input[0] == '-' && input[1] == '-')
}

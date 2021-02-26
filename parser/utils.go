package parser

import (
	"strings"
	"unicode"
)

func IsComment(input string) bool {
	return len(input) < 2 || (input[0] == '-' && input[1] == '-')
}


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

func ParseTypename(input string) (string, error) {
	return "", nil
}

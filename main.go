// Noah Snelson
// February 25, 2021
// sdb/main.go
//
// This file is the main script that is run when a user executes the `sdb`
// binary. It represents a REPL wherein users can type and SQL query, have the
// query parsed, and then executed. The command also reads from stdin, so files
// can be piped/redirected into the REPL.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sdb/parser"
	"sdb/types/metatypes"
	"strings"
)

// Main REPL loop, runs until user terminates the loop or EOF is detected.
func main(){
	reader := bufio.NewReader(os.Stdin)
	dbstate := metatypes.NewState()

	var inputBuilder strings.Builder
	var input string
	var err error
	for {
		fmt.Print("\n> ")

		for {
			input, err = reader.ReadString('\n')
			inputBuilder.WriteString(input)

			fmt.Printf("%v\n", strings.TrimSpace(input))
			if strings.HasPrefix(input, "--") {
				break
			}

			if strings.HasSuffix(strings.TrimSpace(input), ";") {
				fmt.Println()
				break
			}

			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(input)), ".exit") {
				fmt.Println("\nGoodbye!")
				os.Exit(0)
			}

			fmt.Print(">> ")
		}

		if err != nil {
			fmt.Println("\n\nGoodbye!")
			break
		}

		inputStatement := inputBuilder.String()
		inputBuilder.Reset()

		statement, err := parser.Parse(inputStatement)

		if err != nil {
			fmt.Println(err)
			continue
		}

		err = statement.Execute(&dbstate)

		if err != nil {
			fmt.Println(err)
		}
	}
}

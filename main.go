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
	"sdb/types"
)

// Main REPL loop, runs until user terminates the loop or EOF is detected.
func main(){
	reader := bufio.NewReader(os.Stdin)
	dbstate := types.NewState()

	testStr := "\"123.56\""
	parsed, err := parser.ParseValue(testStr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(parsed.GetValue(), parsed.GetType().ToString())
	}

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\n\nGoodbye!")
			break
		}

		fmt.Println(input)

		statement, err := parser.Parse(input)

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

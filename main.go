package main

import (
	"bufio"
	"fmt"
	"os"
	"sdb/parser"
	"sdb/types"
)

func main(){
	reader := bufio.NewReader(os.Stdin)
	dbstate := types.NewState()

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
		} else {
			statement.Execute(&dbstate)
		}
	}
}
